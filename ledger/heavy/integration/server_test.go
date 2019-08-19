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

package integration_test

import (
	"context"
	"crypto"
	"sync"

	"github.com/insolar/insolar/network"

	"math"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/cryptography"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/jetcoordinator"
	"github.com/insolar/insolar/insolar/node"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/store"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/keystore"
	"github.com/insolar/insolar/ledger/artifact"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/genesis"
	"github.com/insolar/insolar/ledger/heavy/executor"
	"github.com/insolar/insolar/ledger/heavy/handler"
	"github.com/insolar/insolar/ledger/heavy/pulsemanager"
	"github.com/insolar/insolar/ledger/object"
	"github.com/insolar/insolar/log"
	networknode "github.com/insolar/insolar/network/node"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var (
	light = nodeMock{
		ref:     gen.Reference(),
		shortID: 1,
		role:    insolar.StaticRoleLightMaterial,
	}
	heavy = nodeMock{
		ref:     gen.Reference(),
		shortID: 2,
		role:    insolar.StaticRoleHeavyMaterial,
	}
	virtual = nodeMock{
		ref:     gen.Reference(),
		shortID: 3,
		role:    insolar.StaticRoleVirtual,
	}
)

func NodeHeavy() insolar.Reference {
	return heavy.ref
}

const PulseStep insolar.PulseNumber = 10

type Server struct {
	pm           insolar.PulseManager
	pulse        insolar.Pulse
	lock         sync.RWMutex
	clientSender bus.Sender
}

func DefaultHeavyConfig() configuration.Configuration {
	cfg := configuration.Configuration{}
	cfg.KeysPath = "../../light/integration/testdata/bootstrap_keys.json"
	cfg.Ledger.LightChainLimit = math.MaxInt32
	cfg.Ledger.JetSplit.DepthLimit = math.MaxUint8
	cfg.Ledger.JetSplit.ThresholdOverflowCount = math.MaxInt32
	cfg.Ledger.JetSplit.ThresholdRecordsCount = math.MaxInt32
	cfg.Bus.ReplyTimeout = time.Minute
	cfg.Ledger.Storage = configuration.Storage{
		DataDirectory: "./db",
	}
	return cfg
}

func defaultReceiveCallback(meta payload.Meta, pl payload.Payload) []payload.Payload {
	return nil
}

func Test_test(t *testing.T) {
	s, err := NewServer(context.Background(), DefaultHeavyConfig(), insolar.GenesisHeavyConfig{}, nil)
	assert.NoError(t, err)
	s.Stop()
}

func NewServer(
	ctx context.Context,
	cfg configuration.Configuration,
	genesisCfg insolar.GenesisHeavyConfig,
	receiveCallback func(meta payload.Meta, pl payload.Payload) []payload.Payload,
) (*Server, error) {
	// Cryptography.
	var (
		KeyProcessor  insolar.KeyProcessor
		CryptoScheme  insolar.PlatformCryptographyScheme
		CryptoService insolar.CryptographyService
	)
	{
		var err error
		// Private key storage.
		ks, err := keystore.NewKeyStore(cfg.KeysPath)
		if err != nil {
			return nil, errors.Wrap(err, "failed to load KeyStore")
		}
		// Public key manipulations.
		KeyProcessor = platformpolicy.NewKeyProcessor()
		// Platform cryptography.
		CryptoScheme = platformpolicy.NewPlatformCryptographyScheme()
		// Sign, verify, etc.
		CryptoService = cryptography.NewCryptographyService()

		c := component.Manager{}
		c.Inject(CryptoService, CryptoScheme, KeyProcessor, ks)
	}

	// Network.
	var (
		NodeNetwork network.NodeNetwork
	)
	{
		NodeNetwork = newNodeNetMock(&light)
	}

	// Role calculations.
	var (
		Coordinator jet.Coordinator
		Pulses      *pulse.DB
		Jets        jet.Storage
		Nodes       *node.Storage
		DB          *store.BadgerDB
	)
	{
		var err error
		DB, err = store.NewBadgerDB(cfg.Ledger.Storage.DataDirectory)
		if err != nil {
			panic(errors.Wrap(err, "failed to initialize DB"))
		}
		Nodes = node.NewStorage()
		Pulses = pulse.NewDB(DB)
		Jets = jet.NewStore()

		c := jetcoordinator.NewJetCoordinator(cfg.Ledger.LightChainLimit)
		c.PulseCalculator = Pulses
		c.PulseAccessor = Pulses
		c.JetAccessor = Jets
		c.OriginProvider = NodeNetwork
		c.PlatformCryptographyScheme = CryptoScheme
		c.Nodes = Nodes

		Coordinator = c
	}

	logger := log.NewWatermillLogAdapter(inslogger.FromContext(ctx))
	// Communication.
	var (
		ServerBus, ClientBus       *bus.Bus
		ServerPubSub, ClientPubSub *gochannel.GoChannel
	)
	{
		ServerPubSub = gochannel.NewGoChannel(gochannel.Config{}, logger)
		ClientPubSub = gochannel.NewGoChannel(gochannel.Config{}, logger)
		ServerBus = bus.NewBus(cfg.Bus, ServerPubSub, Pulses, Coordinator, CryptoScheme)

		c := jetcoordinator.NewJetCoordinator(cfg.Ledger.LightChainLimit)
		c.PulseCalculator = Pulses
		c.PulseAccessor = Pulses
		c.JetAccessor = Jets
		c.OriginProvider = newNodeNetMock(&virtual)
		c.PlatformCryptographyScheme = CryptoScheme
		c.Nodes = Nodes
		ClientBus = bus.NewBus(cfg.Bus, ClientPubSub, Pulses, c, CryptoScheme)
	}

	// Heavy components.
	var (
		PulseManager insolar.PulseManager
		Handler      *handler.Handler
		Genesis      *genesis.Genesis
		Records      *object.RecordDB
		JetKeeper    executor.JetKeeper
	)
	{
		Records = object.NewRecordDB(DB)
		indexes := object.NewIndexDB(DB, Records)
		drops := drop.NewDB(DB)
		jets := jet.NewDBStore(DB)
		JetKeeper = executor.NewJetKeeper(jets, DB, Pulses)
		// c.rollback = executor.NewDBRollback(JetKeeper, Pulses, drops, Records, indexes, jets, Pulses)

		sp := pulse.NewStartPulse()

		backupMaker, err := executor.NewBackupMaker(ctx, DB, cfg.Ledger.Backup, JetKeeper.TopSyncPulse())
		if err != nil {
			return nil, errors.Wrap(err, "failed create backuper")
		}

		pm := pulsemanager.NewPulseManager()
		// pm.Bus = Bus
		pm.NodeNet = NodeNetwork
		pm.NodeSetter = Nodes
		pm.Nodes = Nodes
		pm.PulseAppender = Pulses
		pm.PulseAccessor = Pulses
		pm.JetModifier = jets
		pm.StartPulse = sp
		pm.FinalizationKeeper = executor.NewFinalizationKeeperDefault(JetKeeper, Pulses, cfg.Ledger.LightChainLimit)

		h := handler.New(cfg.Ledger)
		h.RecordAccessor = Records
		h.RecordModifier = Records
		h.JetCoordinator = Coordinator
		h.IndexAccessor = indexes
		h.IndexModifier = indexes
		// h.Bus = Bus
		h.DropModifier = drops
		h.PCS = CryptoScheme
		h.PulseAccessor = Pulses
		h.PulseCalculator = Pulses
		h.StartPulse = sp
		h.JetModifier = jets
		h.JetAccessor = jets
		h.JetTree = jets
		h.DropDB = drops
		h.JetKeeper = JetKeeper
		h.BackupMaker = backupMaker
		h.Sender = ClientBus

		PulseManager = pm
		Handler = h

		artifactManager := &artifact.Scope{
			PulseNumber:    insolar.FirstPulseNumber,
			PCS:            CryptoScheme,
			RecordAccessor: Records,
			RecordModifier: Records,
			IndexModifier:  indexes,
			IndexAccessor:  indexes,
		}
		Genesis = &genesis.Genesis{
			ArtifactManager: artifactManager,
			BaseRecord: &genesis.BaseRecord{
				DB:             DB,
				DropModifier:   drops,
				PulseAppender:  Pulses,
				PulseAccessor:  Pulses,
				RecordModifier: Records,
				IndexModifier:  indexes,
			},

			DiscoveryNodes:  genesisCfg.DiscoveryNodes,
			ContractsConfig: genesisCfg.ContractsConfig,
		}

		_ = Genesis
		_ = Handler
	}

	// Start routers with handlers.
	{
		outHandler := func(msg *message.Message) error {
			meta := payload.Meta{}
			err := meta.Unmarshal(msg.Payload)
			if err != nil {
				panic(errors.Wrap(err, "failed to unmarshal meta"))
			}

			pl, err := payload.Unmarshal(meta.Payload)
			if err != nil {
				panic(nil)
			}
			go func() {
				var replies []payload.Payload
				if receiveCallback != nil {
					replies = receiveCallback(meta, pl)
				} else {
					replies = defaultReceiveCallback(meta, pl)
				}

				for _, rep := range replies {
					msg, err := payload.NewMessage(rep)
					if err != nil {
						panic(err)
					}
					ClientBus.Reply(context.Background(), meta, msg)
				}
			}()

			clientHandler := func(msg *message.Message) (messages []*message.Message, e error) {
				return nil, nil
			}
			// Republish as incoming to client.
			_, err = ClientBus.IncomingMessageRouter(clientHandler)(msg)

			if err != nil {
				panic(err)
			}
			return nil
		}

		inRouter, err := message.NewRouter(message.RouterConfig{}, logger)
		if err != nil {
			panic(err)
		}
		outRouter, err := message.NewRouter(message.RouterConfig{}, logger)
		if err != nil {
			panic(err)
		}

		outRouter.AddNoPublisherHandler(
			"Outgoing",
			bus.TopicOutgoing,
			ServerPubSub,
			outHandler,
		)

		inRouter.AddMiddleware(
			middleware.InstantAck,
			ServerBus.IncomingMessageRouter,
		)

		startRouter(ctx, inRouter)
		startRouter(ctx, outRouter)
	}

	inslogger.FromContext(ctx).WithFields(map[string]interface{}{
		"light":   light.ID().String(),
		"virtual": virtual.ID().String(),
		"heavy":   heavy.ID().String(),
	}).Info("started test server")

	if err := Genesis.Start(ctx); err != nil {
		log.Fatalf("genesis failed on heavy with error: %v", err)
	}

	s := &Server{
		pm:           PulseManager,
		pulse:        *insolar.GenesisPulse,
		clientSender: ClientBus,
	}
	return s, nil
}

func startRouter(ctx context.Context, router *message.Router) {
	go func() {
		if err := router.Run(ctx); err != nil {
			inslogger.FromContext(ctx).Error("Error while running router", err)
		}
	}()
	<-router.Running()
}

func (s *Server) SetPulse(ctx context.Context) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.pulse = insolar.Pulse{
		PulseNumber: s.pulse.PulseNumber + PulseStep,
	}
	err := s.pm.Set(ctx, s.pulse)
	if err != nil {
		panic(err)
	}
}

func (s *Server) Pulse() insolar.PulseNumber {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.pulse.PulseNumber
}

func (s *Server) Send(ctx context.Context, pl payload.Payload) (<-chan *message.Message, func()) {
	msg, err := payload.NewMessage(pl)
	if err != nil {
		panic(err)
	}
	return s.clientSender.SendTarget(ctx, msg, gen.Reference())
}

func (s *Server) Stop() {
}

type nodeMock struct {
	ref     insolar.Reference
	shortID insolar.ShortNodeID
	role    insolar.StaticRole
}

func (n *nodeMock) ID() insolar.Reference {
	return n.ref
}

func (n *nodeMock) ShortID() insolar.ShortNodeID {
	return n.shortID
}

func (n *nodeMock) Role() insolar.StaticRole {
	return n.role
}

func (n *nodeMock) PublicKey() crypto.PublicKey {
	panic("implement me")
}

func (n *nodeMock) Address() string {
	return ""
}

func (n *nodeMock) GetGlobuleID() insolar.GlobuleID {
	panic("implement me")
}

func (n *nodeMock) Version() string {
	panic("implement me")
}

func (n *nodeMock) LeavingETA() insolar.PulseNumber {
	panic("implement me")
}

func (n *nodeMock) GetState() insolar.NodeState {
	return insolar.NodeReady
}

func (n *nodeMock) GetPower() insolar.Power {
	return 1
}

type nodeNetMock struct {
	me insolar.NetworkNode
}

func (n *nodeNetMock) GetAccessor(insolar.PulseNumber) network.Accessor {
	return networknode.NewAccessor(networknode.NewSnapshot(insolar.GenesisPulse.PulseNumber, []insolar.NetworkNode{&virtual, &heavy, &light}))
}

func newNodeNetMock(me insolar.NetworkNode) *nodeNetMock {
	return &nodeNetMock{me: me}
}

func (n *nodeNetMock) GetOrigin() insolar.NetworkNode {
	return n.me
}
