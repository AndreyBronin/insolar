//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
//

package adapters

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/consensus/common"
	"github.com/insolar/insolar/network/consensus/gcpv2/census"
	common2 "github.com/insolar/insolar/network/consensus/gcpv2/common"
	"github.com/insolar/insolar/network/consensus/gcpv2/errors"
)

type MisbehaviorRegistry struct{}

func NewMisbehaviorRegistry() *MisbehaviorRegistry {
	return &MisbehaviorRegistry{}
}

func (mr *MisbehaviorRegistry) AddReport(report errors.MisbehaviorReport) {
	ctx := context.TODO()

	inslogger.FromContext(ctx).Warnf("Got MisbehaviorReport")
}

type MandateRegistry struct {
	cloudHash common2.CloudStateHash
}

func NewMandateRegistry(cloudHash common2.CloudStateHash) *MandateRegistry {
	return &MandateRegistry{
		cloudHash: cloudHash,
	}
}

func (mr *MandateRegistry) FindRegisteredProfile(host common.HostIdentityHolder) common2.HostProfile {
	panic("implement me")
}

func (mr *MandateRegistry) GetPrimingCloudHash() common2.CloudStateHash {
	return mr.cloudHash
}

type OfflinePopulation struct {
	// TODO: should't use nodekeeper here.
	nodeKeeper   network.NodeKeeper
	manager      insolar.CertificateManager
	keyProcessor insolar.KeyProcessor
}

func NewOfflinePopulation(nodeKeeper network.NodeKeeper, manager insolar.CertificateManager, keyProcessor insolar.KeyProcessor) *OfflinePopulation {
	return &OfflinePopulation{
		nodeKeeper:   nodeKeeper,
		manager:      manager,
		keyProcessor: keyProcessor,
	}
}

func (op *OfflinePopulation) FindRegisteredProfile(identity common.HostIdentityHolder) common2.HostProfile {
	node := op.nodeKeeper.GetAccessor().GetActiveNodeByAddr(identity.GetHostAddress().String())
	cert := op.manager.GetCertificate()

	return NewNodeIntroProfile(node, cert, op.keyProcessor)
}

type VersionedRegistries struct {
	mandateRegistry     census.MandateRegistry
	misbehaviorRegistry census.MisbehaviorRegistry
	offlinePopulation   census.OfflinePopulation

	pulseData common.PulseData
}

func NewVersionedRegistries(
	mandateRegistry census.MandateRegistry,
	misbehaviorRegistry census.MisbehaviorRegistry,
	offlinePopulation census.OfflinePopulation,
) *VersionedRegistries {
	return &VersionedRegistries{
		mandateRegistry:     mandateRegistry,
		misbehaviorRegistry: misbehaviorRegistry,
		offlinePopulation:   offlinePopulation,
	}
}

func (c *VersionedRegistries) CommitNextPulse(pd common.PulseData, population census.OnlinePopulation) census.VersionedRegistries {
	pd.EnsurePulseData()
	cp := *c
	cp.pulseData = pd
	return &cp
}

func (c *VersionedRegistries) GetMisbehaviorRegistry() census.MisbehaviorRegistry {
	return c.misbehaviorRegistry
}

func (c *VersionedRegistries) GetMandateRegistry() census.MandateRegistry {
	return c.mandateRegistry
}

func (c *VersionedRegistries) GetOfflinePopulation() census.OfflinePopulation {
	return c.offlinePopulation
}

func (c *VersionedRegistries) GetVersionPulseData() common.PulseData {
	return c.pulseData
}
