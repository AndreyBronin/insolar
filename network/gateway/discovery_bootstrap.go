package gateway

import (
	"context"
	"errors"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/hostnetwork/host"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/hostnetwork/packet"
)

func newDiscoveryBootstrap(b *Base) *DiscoveryBootstrap {
	return &DiscoveryBootstrap{b}
}

// DiscoveryBootstrap void network state
type DiscoveryBootstrap struct {
	*Base
}

func (g *DiscoveryBootstrap) Run(ctx context.Context) {

	permit, err := g.authorize(ctx)
	if err != nil {
		// log warn
		g.Gatewayer.SwitchState(insolar.NoNetworkState)
	}

	g.NodeKeeper.GetConsensusInfo().SetIsJoiner(false)

	claim, _ := g.NodeKeeper.GetOriginJoinClaim()
	pulse, err := g.PulseAccessor.Latest(ctx)
	if err != nil {
		pulse = insolar.Pulse{PulseNumber: 1}
	}

	resp, _ := g.BootstrapRequester.Bootstrap(ctx, permit, claim, &pulse)

	if resp.Code == packet.Reject {
		g.Gatewayer.SwitchState(insolar.NoNetworkState)
		return
	}

	if resp.Code == packet.Accepted {
		//  ConsensusWaiting, ETA
		g.bootstrapETA = insolar.PulseNumber(resp.ETA)
		g.Gatewayer.SwitchState(insolar.WaitConsensus)
		return
	}

	// var err error

	// cert := g.CertificateManager.GetCertificate()

	// TODO: shaffle discovery nodes

	// ping ?
	// Authorize(utc) permit, check version
	// process response: trueAccept, redirect with permit, posibleAccept(regen shortId, updateScedule, update time utc)
	// check majority
	// handle reconect to other network
	// fake pulse

	// if network.OriginIsDiscovery(cert) {
	// 	_, err = g.Bootstrapper.BootstrapDiscovery(ctx)
	// 	// if the network is up and complete, we return discovery nodes via consensus
	// 	if err == bootstrap.ErrReconnectRequired {
	// 		log.Debugf("[ Bootstrap ] Connecting discovery node %s as joiner", g.NodeKeeper.GetOrigin().ID())
	// 		g.NodeKeeper.GetOrigin().(node.MutableNode).SetState(insolar.NodePending)
	// 		g.bootstrapJoiner(ctx)
	// 		return
	// 	}
	//
	// }
}

func (g *DiscoveryBootstrap) GetState() insolar.NetworkState {
	return insolar.DiscoveryBootstrap
}

func (g *DiscoveryBootstrap) OnPulse(ctx context.Context, pu insolar.Pulse) error {
	return g.Base.OnPulse(ctx, pu)
}

func (g *DiscoveryBootstrap) ShoudIgnorePulse(context.Context, insolar.Pulse) bool {
	return false
}

func (g *DiscoveryBootstrap) authorize(ctx context.Context) (*packet.Permit, error) {
	cert := g.CertificateManager.GetCertificate()
	discoveryNodes := network.ExcludeOrigin(cert.GetDiscoveryNodes(), g.NodeKeeper.GetOrigin().ID())
	// todo: sort discoveryNodes

	for _, n := range discoveryNodes {
		if g.NodeKeeper.GetAccessor().GetActiveNode(*n.GetNodeRef()) != nil {
			inslogger.FromContext(ctx).Info("Skip discovery already in active list: ", n.GetNodeRef().String())
			continue
		}

		h, _ := host.NewHostN(n.GetHost(), *n.GetNodeRef())

		res, err := g.BootstrapRequester.Authorize(ctx, h, cert)
		if err != nil {
			inslogger.FromContext(ctx).Errorf("Error authorizing to discovery node %s: %s", h.String(), err.Error())
			continue
		}

		return res.Permit, nil
	}

	return nil, errors.New("Failed to authorize to any discovery node.")
}