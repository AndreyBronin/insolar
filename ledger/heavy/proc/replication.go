//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package proc

import (
	"context"
	"fmt"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/heavy/executor"
	"github.com/insolar/insolar/ledger/object"
	"github.com/pkg/errors"
	"go.opencensus.io/stats"
)

type Replication struct {
	message payload.Meta
	cfg     configuration.Ledger

	dep struct {
		records  object.RecordModifier
		indexes  object.IndexModifier
		pcs      insolar.PlatformCryptographyScheme
		pulses   pulse.Accessor
		drops    drop.Modifier
		jets     jet.Modifier
		keeper   executor.JetKeeper
		backuper executor.BackupMaker
	}
}

func NewReplication(msg payload.Meta, cfg configuration.Ledger) *Replication {
	return &Replication{
		message: msg,
		cfg:     cfg,
	}
}

func (p *Replication) Dep(
	records object.RecordModifier,
	indexes object.IndexModifier,
	pcs insolar.PlatformCryptographyScheme,
	pulses pulse.Accessor,
	drops drop.Modifier,
	jets jet.Modifier,
	keeper executor.JetKeeper,
	backuper executor.BackupMaker,
) {
	p.dep.records = records
	p.dep.indexes = indexes
	p.dep.pcs = pcs
	p.dep.pulses = pulses
	p.dep.drops = drops
	p.dep.jets = jets
	p.dep.keeper = keeper
	p.dep.backuper = backuper
}

func (p *Replication) Proceed(ctx context.Context) error {
	pl, err := payload.Unmarshal(p.message.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal payload")
	}
	msg, ok := pl.(*payload.Replication)
	if !ok {
		return fmt.Errorf("unexpected payload %T", pl)
	}

	storeRecords(ctx, p.dep.records, p.dep.pcs, msg.Pulse, msg.Records)
	if err := storeIndexes(ctx, p.dep.indexes, msg.Indexes, msg.Pulse); err != nil {
		return errors.Wrap(err, "failed to store indexes")
	}

	dr, err := storeDrop(ctx, p.dep.drops, msg.Drop)
	if err != nil {
		return errors.Wrap(err, "failed to store drop")
	}

	jetKeeper := p.dep.keeper
	topSyncPulse := jetKeeper.TopSyncPulse()
	if err := jetKeeper.AddDropConfirmation(ctx, dr.Pulse, dr.JetID, dr.Split); err != nil {
		return errors.Wrapf(err, "failed to add jet to JetKeeper jet=%v", dr.JetID.DebugString())
	}

	if !p.cfg.Backup.Enabled {
		if err := jetKeeper.AddBackupConfirmation(ctx, dr.Pulse); err != nil {
			inslogger.FromContext(ctx).Fatal("AddBackupConfirmation return error: ", err)
		}
	}
	if topSyncPulse != jetKeeper.TopSyncPulse() {
		FinalizePulse(ctx, p.dep.backuper, jetKeeper, dr.Pulse)
	}

	stats.Record(ctx,
		statReceivedHeavyPayloadCount.M(1),
	)

	return nil
}

func storeIndexes(
	ctx context.Context,
	mod object.IndexModifier,
	indexes []record.Index,
	pn insolar.PulseNumber,
) error {
	for _, idx := range indexes {
		err := mod.SetIndex(ctx, pn, idx)
		if err != nil {
			return errors.Wrapf(err, "heavyserver: index storing failed")
		}
	}

	return nil
}

func storeDrop(
	ctx context.Context,
	drops drop.Modifier,
	rawDrop []byte,
) (*drop.Drop, error) {
	d, err := drop.Decode(rawDrop)
	if err != nil {
		inslogger.FromContext(ctx).Error(err)
		return nil, err
	}
	err = drops.Set(ctx, *d)
	if err != nil {
		return nil, errors.Wrapf(err, "heavyserver: drop storing failed")
	}

	return d, nil
}

func storeRecords(
	ctx context.Context,
	mod object.RecordModifier,
	pcs insolar.PlatformCryptographyScheme,
	pn insolar.PulseNumber,
	records []record.Material,
) {
	inslog := inslogger.FromContext(ctx)

	for _, rec := range records {
		hash := record.HashVirtual(pcs.ReferenceHasher(), rec.Virtual)
		id := insolar.NewID(pn, hash)
		err := mod.Set(ctx, *id, rec)
		if err != nil {
			inslog.Error(err, "heavyserver: store record failed")
			continue
		}
	}
}

func FinalizePulse(ctx context.Context, backuper executor.BackupMaker, jetKeeper executor.JetKeeper, pulse insolar.PulseNumber) {
	go func() {
		logger := inslogger.FromContext(ctx)
		err := backuper.Do(ctx, pulse)
		if err != nil {
			if err == executor.ErrAlreadyDone {
				logger.Warn("BackupMaker says, that work already done")
				return
			}
			logger.Fatalf("Can't do backup: ", err)
		}
		err = jetKeeper.AddBackupConfirmation(ctx, pulse)
		if err != nil {
			logger.Fatalf("Can't add backup confirmation: ", err)
		}
		inslogger.FromContext(ctx).Infof("Pulse %d completely finalized ( drops + hots + backup )", pulse)
	}()
}
