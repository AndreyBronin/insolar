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

package claimhandler

import (
	"bytes"
	"crypto/rand"
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensusv1/packets"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueue_PushClaim(t *testing.T) {
	queue := Queue{}
	elemCount := 20
	entr := insolar.Entropy{}
	_, err := rand.Read(entr[:])
	assert.NoError(t, err)
	for i := 0; i < elemCount; i++ {
		claim := getJoinClaim(t, testutils.RandomRef())
		queue.PushClaim(claim, getPriority(claim.NodeRef, entr))
	}
	assert.Equal(t, queue.Len(), elemCount)
}

func TestQueue_Pop(t *testing.T) {
	queue := Queue{}
	elemCount := 20
	entr := insolar.Entropy{}
	_, err := rand.Read(entr[:])
	assert.NoError(t, err)
	for i := 0; i < elemCount; i++ {
		claim := getJoinClaim(t, testutils.RandomRef())
		queue.PushClaim(claim, getPriority(claim.NodeRef, entr))
	}
	assert.Equal(t, queue.Len(), elemCount)

	claim := queue.PopClaim().(*packets.NodeJoinClaim)
	refLen := len(claim.NodeRef.Bytes())
	prevPriority := make([]byte, refLen)
	priority := make([]byte, refLen)
	copy(prevPriority, getPriority(claim.NodeRef, entr))
	for i := 1; i < elemCount; i++ {
		claim := queue.PopClaim().(*packets.NodeJoinClaim)
		copy(priority, getPriority(claim.NodeRef, entr))
		assert.True(t, bytes.Compare(prevPriority, priority) > 0)
		copy(prevPriority, priority)
	}
}

func getJoinClaim(t *testing.T, ref insolar.Reference) *packets.NodeJoinClaim {
	nodeJoinClaim := &packets.NodeJoinClaim{}
	nodeJoinClaim.ShortNodeID = insolar.ShortNodeID(77)
	nodeJoinClaim.RelayNodeID = insolar.ShortNodeID(26)
	nodeJoinClaim.ProtocolVersionAndFlags = uint32(99)
	nodeJoinClaim.JoinsAfter = uint32(67)
	nodeJoinClaim.NodeRoleRecID = 32
	nodeJoinClaim.NodeRef = ref
	_, err := rand.Read(nodeJoinClaim.NodePK[:])
	require.NoError(t, err)
	addr, err := packets.NewNodeAddress("127.0.0.1:5566")
	require.NoError(t, err)
	nodeJoinClaim.NodeAddress = addr

	return nodeJoinClaim
}
