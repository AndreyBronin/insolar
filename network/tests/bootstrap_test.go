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

// +build networktest

package tests

import (
	"fmt"
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type bootstrapSuite struct {
	testSuite
}

func (s *bootstrapSuite) Setup() {
	var err error
	s.pulsar, err = NewTestPulsar(reqTimeoutMs*10, pulseDelta*10)
	require.NoError(s.t, err)

	inslogger.FromContext(s.ctx).Info("SetupTest")

	for i := 0; i < s.bootstrapCount; i++ {
		role := insolar.StaticRoleVirtual
		if i == 0 {
			role = insolar.StaticRoleHeavyMaterial
		}

		s.bootstrapNodes = append(s.bootstrapNodes, s.newNetworkNodeWithRole(fmt.Sprintf("bootstrap_%d", i), role))
	}

	s.SetupNodesNetwork(s.bootstrapNodes)

	pulseReceivers := make([]string, 0)
	for _, node := range s.bootstrapNodes {
		pulseReceivers = append(pulseReceivers, node.host)
	}

	log.Info("Start test pulsar")
	err = s.pulsar.Start(s.ctx, pulseReceivers)
	require.NoError(s.t, err)
}

func (s *bootstrapSuite) stopBootstrapSuite() {
	inslogger.FromContext(s.ctx).Info("stopNetworkSuite")

	suiteLogger.Info("Stop bootstrap nodes")
	for _, n := range s.bootstrapNodes {
		err := n.componentManager.Stop(n.ctx)
		assert.NoError(s.t, err)
	}
}

func (s *bootstrapSuite) waitForConsensus(consensusCount int) {
	for i := 0; i < consensusCount; i++ {
		for _, n := range s.bootstrapNodes {
			<-n.consensusResult
		}
	}
}

func newBootstraptSuite(t *testing.T, bootstrapCount int) *bootstrapSuite {
	return &bootstrapSuite{
		testSuite: newTestSuite(t, bootstrapCount, 0),
	}
}

func startBootstrapSuite(t *testing.T) *bootstrapSuite {
	t.Skip("Skip until fix consensus bugs")

	s := newBootstraptSuite(t, 11)
	s.Setup()
	return s
}

func TestBootstrap(t *testing.T) {
	s := startBootstrapSuite(t)
	defer s.stopBootstrapSuite()

	s.StartNodesNetwork(s.bootstrapNodes)

	s.waitForConsensus(2)
	s.AssertActiveNodesCountDelta(0)

	s.waitForConsensus(1)
	s.AssertActiveNodesCountDelta(0)
	s.AssertWorkingNodesCountDelta(0)
}
