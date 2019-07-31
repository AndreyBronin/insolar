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

package nodenetwork

import (
	"context"
	"net"
	"sync"

	"github.com/insolar/insolar/network/hostnetwork/resolver"
	"github.com/insolar/insolar/network/node"

	"github.com/insolar/insolar/instrumentation/inslogger"

	"github.com/pkg/errors"
	"go.opencensus.io/stats"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/version"
)

// NewNodeNetwork create active node component
func NewNodeNetwork(configuration configuration.Transport, certificate insolar.Certificate) (insolar.NodeNetwork, error) {
	origin, err := createOrigin(configuration, certificate)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create origin node")
	}
	nodeKeeper := NewNodeKeeper(origin)
	if !network.OriginIsDiscovery(certificate) {
		origin.(node.MutableNode).SetState(insolar.NodePending)
	}
	return nodeKeeper, nil
}

func createOrigin(configuration configuration.Transport, certificate insolar.Certificate) (insolar.NetworkNode, error) {
	publicAddress, err := resolveAddress(configuration)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to resolve public address")
	}

	role := certificate.GetRole()
	if role == insolar.StaticRoleUnknown {
		log.Info("[ createOrigin ] Use insolar.StaticRoleLightMaterial, since no role in certificate")
		role = insolar.StaticRoleLightMaterial
	}

	return node.NewNode(
		*certificate.GetNodeRef(),
		role,
		certificate.GetPublicKey(),
		publicAddress,
		version.Version,
	), nil
}

func resolveAddress(configuration configuration.Transport) (string, error) {
	addr, err := net.ResolveTCPAddr("tcp", configuration.Address)
	if err != nil {
		return "", err
	}
	address, err := resolver.Resolve(configuration.FixedPublicAddress, addr.String())
	if err != nil {
		return "", err
	}
	return address, nil
}

// NewNodeKeeper create new NodeKeeper
func NewNodeKeeper(origin insolar.NetworkNode) network.NodeKeeper {
	nk := &nodekeeper{
		cloudHash: make([]byte, 64),
		origin:    origin,
		syncNodes: make([]insolar.NetworkNode, 0),
	}
	nk.SetInitialSnapshot([]insolar.NetworkNode{})
	return nk
}

type nodekeeper struct {
	origin insolar.NetworkNode

	cloudHashLock sync.RWMutex
	cloudHash     []byte

	activeLock sync.RWMutex
	snapshot   *node.Snapshot
	accessor   *node.Accessor

	syncLock  sync.Mutex
	syncNodes []insolar.NetworkNode

	TerminationHandler insolar.TerminationHandler `inject:""`
	//CryptographyService insolar.CryptographyService `inject:""`
}

func (nk *nodekeeper) GetSnapshotCopy() *node.Snapshot {
	nk.activeLock.RLock()
	defer nk.activeLock.RUnlock()

	return nk.snapshot.Copy()
}

func (nk *nodekeeper) SetInitialSnapshot(nodes []insolar.NetworkNode) {
	nk.activeLock.Lock()
	defer nk.activeLock.Unlock()

	nodesMap := make(map[insolar.Reference]insolar.NetworkNode)
	for _, n := range nodes {
		nodesMap[n.ID()] = n
	}
	nk.snapshot = node.NewSnapshot(insolar.FirstPulseNumber, nodesMap)
	nk.accessor = node.NewAccessor(nk.snapshot)

	nk.syncLock.Lock()
	nk.syncNodes = nk.accessor.GetActiveNodes()
	nk.syncLock.Unlock()
}

func (nk *nodekeeper) GetAccessor() network.Accessor {
	nk.activeLock.RLock()
	defer nk.activeLock.RUnlock()

	return nk.accessor
}

func (nk *nodekeeper) GetWorkingNode(ref insolar.Reference) insolar.NetworkNode {
	return nk.GetAccessor().GetWorkingNode(ref)
}

func (nk *nodekeeper) GetWorkingNodesByRole(role insolar.DynamicRole) []insolar.Reference {
	return nk.GetAccessor().GetWorkingNodesByRole(role)
}

func (nk *nodekeeper) GetOrigin() insolar.NetworkNode {
	return nk.origin
}

func (nk *nodekeeper) GetCloudHash() []byte {
	nk.cloudHashLock.RLock()
	defer nk.cloudHashLock.RUnlock()

	return nk.cloudHash
}

func (nk *nodekeeper) SetCloudHash(cloudHash []byte) {
	nk.cloudHashLock.Lock()
	defer nk.cloudHashLock.Unlock()

	nk.cloudHash = cloudHash
}

func (nk *nodekeeper) GetWorkingNodes() []insolar.NetworkNode {
	return nk.GetAccessor().GetWorkingNodes()
}

func (nk *nodekeeper) Sync(ctx context.Context, nodes []insolar.NetworkNode) {
	nk.syncLock.Lock()
	defer nk.syncLock.Unlock()

	inslogger.FromContext(ctx).Debugf("Sync, nodes: %d", len(nodes))
	nk.syncNodes = nodes
}

// syncOrigin synchronize data in origin node with node from active list in case when they are different objects
func (nk *nodekeeper) syncOrigin(n insolar.NetworkNode) {
	if nk.origin == n {
		return
	}
	mutableOrigin := nk.origin.(node.MutableNode)
	mutableOrigin.SetState(n.GetState())
	if n.GetState() == insolar.NodeLeaving {
		mutableOrigin.SetLeavingETA(n.LeavingETA())
	}
	mutableOrigin.SetShortID(n.ShortID())
}

func (nk *nodekeeper) MoveSyncToActive(ctx context.Context, number insolar.PulseNumber) {
	nk.activeLock.Lock()
	nk.syncLock.Lock()
	defer func() {
		nk.syncLock.Unlock()
		nk.activeLock.Unlock()
	}()

	if len(nk.syncNodes) == 0 {
		inslogger.FromContext(ctx).Info("[ MoveSyncToActive ] New active list confirmed and unchanged")
		return
	}

	mergeResult, _ := GetMergedCopy(nk.syncNodes)
	inslogger.FromContext(ctx).Infof("[ MoveSyncToActive ] New active list confirmed. Active list size: %d -> %d",
		len(nk.accessor.GetActiveNodes()),
		len(mergeResult.ActiveList),
	)

	nk.snapshot = node.NewSnapshot(number, mergeResult.ActiveList)
	nk.accessor = node.NewAccessor(nk.snapshot)
	stats.Record(ctx, network.ActiveNodes.M(int64(len(nk.accessor.GetActiveNodes()))))

	nk.syncNodes = nil
}
