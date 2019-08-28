///
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
///

package integration

import (
	"context"
	"fmt"
	"io"
	"math"
	"sync"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/contractrequester"
	"github.com/insolar/insolar/cryptography"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/jetcoordinator"
	"github.com/insolar/insolar/insolar/node"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/keystore"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/logicrunner"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/logicexecutor"
	"github.com/insolar/insolar/logicrunner/machinesmanager"
	"github.com/insolar/insolar/logicrunner/pulsemanager"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/platformpolicy"
)

func NodeLight() insolar.Reference {
	return light.ref
}

const PulseStep insolar.PulseNumber = 10

type Server struct {
	pm                insolar.PulseManager
	componentManager  *component.Manager
	stopper           func()
	pulse             insolar.Pulse
	lock              sync.RWMutex
	clientSender      bus.Sender
	logicRunner       *logicrunner.LogicRunner
	contractRequester *contractrequester.ContractRequester

	ExternalPubSub, IncomingPubSub *gochannel.GoChannel
}

func DefaultVMConfig() configuration.Configuration {
	cfg := configuration.Configuration{}
	cfg.KeysPath = "testdata/bootstrap_keys.json"
	cfg.Ledger.LightChainLimit = math.MaxInt32
	cfg.LogicRunner = configuration.NewLogicRunner()
	cfg.Bus.ReplyTimeout = time.Minute
	return cfg
}

func DefaultLightResponse(pl payload.Payload) []payload.Payload {
	switch pl.(type) {
	// getters
	case *payload.GetFilament, *payload.GetCode, *payload.GetRequest, *payload.GetRequestInfo:
		return nil
	// setters
	case *payload.SetResult, *payload.SetCode, *payload.SetIncomingRequest, *payload.SetOutgoingRequest:
		return nil
	}

	panic(fmt.Sprintf("unexpected message to lightt %T", pl))
}

func defaultReceiveCallback(meta payload.Meta, pl payload.Payload) []payload.Payload {
	if meta.Receiver == NodeLight() {
		return DefaultLightResponse(pl)
	}
	return nil
}

func checkError(ctx context.Context, err error, message string) {
	if err == nil {
		return
	}
	inslogger.FromContext(ctx).Fatalf("%v: %v", message, err.Error())
}

func NewServer(
	ctx context.Context,
	cfg configuration.Configuration,
	receiveCallback func(meta payload.Meta, pl payload.Payload) []payload.Payload) (*Server, error) {

	cm := component.Manager{}

	// Cryptography.
	var (
		KeyProcessor  insolar.KeyProcessor
		CryptoScheme  insolar.PlatformCryptographyScheme
		CryptoService insolar.CryptographyService
		KeyStore      insolar.KeyStore
	)
	{
		var err error
		// Private key storage.
		KeyStore, err = keystore.NewKeyStore(cfg.KeysPath)
		if err != nil {
			return nil, errors.Wrap(err, "failed to load KeyStore")
		}
		// Public key manipulations.
		KeyProcessor = platformpolicy.NewKeyProcessor()
		// Platform cryptography.
		CryptoScheme = platformpolicy.NewPlatformCryptographyScheme()
		// Sign, verify, etc.
		CryptoService = cryptography.NewCryptographyService()
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
		Pulses      *pulse.StorageMem
		Jets        jet.Storage
		Nodes       *node.Storage
	)
	{
		Nodes = node.NewStorage()
		Pulses = pulse.NewStorageMem()
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

	// PulseManager
	var (
		PulseManager *pulsemanager.PulseManager
	)
	{
		PulseManager = pulsemanager.NewPulseManager()
	}

	wmLogger := log.NewWatermillLogAdapter(inslogger.FromContext(ctx))

	// Communication.
	var (
		ClientBus                      *bus.Bus
		ExternalPubSub, IncomingPubSub *gochannel.GoChannel
	)
	{
		ExternalPubSub = gochannel.NewGoChannel(gochannel.Config{}, wmLogger)
		IncomingPubSub = gochannel.NewGoChannel(gochannel.Config{}, wmLogger)

		c := jetcoordinator.NewJetCoordinator(cfg.Ledger.LightChainLimit)
		c.PulseCalculator = Pulses
		c.PulseAccessor = Pulses
		c.JetAccessor = Jets
		c.OriginProvider = newNodeNetMock(&virtual)
		c.PlatformCryptographyScheme = CryptoScheme
		c.Nodes = Nodes
		ClientBus = bus.NewBus(cfg.Bus, IncomingPubSub, Pulses, c, CryptoScheme)
	}

	logicRunner, err := logicrunner.NewLogicRunner(&cfg.LogicRunner, IncomingPubSub, ClientBus)
	checkError(ctx, err, "failed to start LogicRunner")

	contractRequester, err := contractrequester.New()
	checkError(ctx, err, "failed to start ContractRequester")

	// TODO: remove this hack in INS-3341
	contractRequester.LR = logicRunner

	pm := pulsemanager.NewPulseManager()

	cm.Register(
		CryptoScheme,
		KeyStore,
		CryptoService,
		KeyProcessor,
		logicRunner,
		logicexecutor.NewLogicExecutor(),
		logicrunner.NewRequestsExecutor(),
		machinesmanager.NewMachinesManager(),
		NodeNetwork,
		pm,
		PulseManager,
	)

	components := []interface{}{
		ClientBus,
		IncomingPubSub,
		contractRequester,
		artifacts.NewClient(ClientBus),
		artifacts.NewDescriptorsCache(),
		Coordinator,
		Pulses,
		jet.NewStore(),
		node.NewStorage(),
	}
	components = append(components, []interface{}{
		CryptoService,
		KeyProcessor,
	}...)

	cm.Inject(components...)

	err = cm.Init(ctx)
	checkError(ctx, err, "failed to init components")

	// Start routers with handlers.
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

		// Republish as incoming to self.
		if meta.Receiver == virtual.ID() {
			err = ExternalPubSub.Publish(bus.TopicIncoming, msg)
			if err != nil {
				panic(err)
			}
			return nil
		}

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

	stopper := startWatermill(
		ctx, wmLogger, IncomingPubSub, ClientBus,
		outHandler,
		logicRunner.FlowDispatcher.Process,
		contractRequester.ReceiveResult,
	)

	PulseManager.FlowDispatcher = logicRunner.FlowDispatcher

	inslogger.FromContext(ctx).WithFields(map[string]interface{}{
		"light":   light.ID().String(),
		"virtual": virtual.ID().String(),
	}).Info("started test server")

	s := &Server{
		pm:               PulseManager,
		componentManager: &cm,
		stopper:          stopper,
		pulse:            *insolar.GenesisPulse,
		clientSender:     ClientBus,
	}
	return s, nil
}

func (s *Server) Stop(ctx context.Context) {
	panicIfErr(s.componentManager.Stop(ctx))
	s.stopper()
	panicIfErr(s.ExternalPubSub.Close())
	panicIfErr(s.IncomingPubSub.Close())
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

func (s *Server) SendToSelf(ctx context.Context, pl payload.Payload) (<-chan *message.Message, func()) {
	msg, err := payload.NewMessage(pl)
	if err != nil {
		panic(err)
	}
	inslogger.FromContext(ctx).Info(msg)
	return s.clientSender.SendTarget(ctx, msg, virtual.ID())
}

func (s *Server) Pulse() insolar.PulseNumber {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.pulse.PulseNumber
}

func startWatermill(
	ctx context.Context,
	logger watermill.LoggerAdapter,
	sub message.Subscriber,
	b *bus.Bus,
	outHandler, inHandler, resultsHandler message.NoPublishHandlerFunc,
) func() {
	inRouter, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		panic(err)
	}
	outRouter, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		panic(err)
	}

	outRouter.AddNoPublisherHandler(
		"OutgoingHandler",
		bus.TopicOutgoing,
		sub,
		outHandler,
	)

	inRouter.AddMiddleware(
		b.IncomingMessageRouter,
	)

	inRouter.AddNoPublisherHandler(
		"IncomingHandler",
		bus.TopicIncoming,
		sub,
		inHandler,
	)

	inRouter.AddNoPublisherHandler(
		"IncomingRequestResultHandler",
		bus.TopicIncomingRequestResults,
		sub,
		resultsHandler)

	startRouter(ctx, inRouter)
	startRouter(ctx, outRouter)

	return stopWatermill(ctx, inRouter, outRouter)
}

func stopWatermill(ctx context.Context, routers ...io.Closer) func() {
	return func() {
		for _, r := range routers {
			err := r.Close()
			if err != nil {
				inslogger.FromContext(ctx).Error("Error while closing router", err)
			}
		}
	}
}

func startRouter(ctx context.Context, router *message.Router) {
	go func() {
		if err := router.Run(ctx); err != nil {
			inslogger.FromContext(ctx).Error("Error while running router", err)
		}
	}()
	<-router.Running()
}
