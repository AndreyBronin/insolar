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

package node

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/packets"
	"github.com/insolar/insolar/testutils"
)

func TestNode_Version(t *testing.T) {
	n := NewNode(testutils.RandomRef(), insolar.StaticRoleVirtual, nil, "127.0.0.1", "123")
	assert.Equal(t, "123", n.Version())
	n.(MutableNode).SetVersion("234")
	assert.Equal(t, "234", n.Version())
}

func TestNode_GetState(t *testing.T) {
	n := NewNode(testutils.RandomRef(), insolar.StaticRoleVirtual, nil, "127.0.0.1", "123")
	assert.Equal(t, insolar.NodeReady, n.GetState())
	n.(MutableNode).SetState(insolar.NodeUndefined)
	assert.Equal(t, insolar.NodeUndefined, n.GetState())
	n.(MutableNode).ChangeState()
	assert.Equal(t, insolar.NodePending, n.GetState())
	n.(MutableNode).ChangeState()
	assert.Equal(t, insolar.NodeReady, n.GetState())
	n.(MutableNode).ChangeState()
	assert.Equal(t, insolar.NodeReady, n.GetState())
}

func TestNode_GetGlobuleID(t *testing.T) {
	n := NewNode(testutils.RandomRef(), insolar.StaticRoleVirtual, nil, "127.0.0.1", "123")
	assert.EqualValues(t, 0, n.GetGlobuleID())
}

func TestNode_LeavingETA(t *testing.T) {
	n := NewNode(testutils.RandomRef(), insolar.StaticRoleVirtual, nil, "127.0.0.1", "123")
	assert.Equal(t, insolar.NodeReady, n.GetState())
	n.(MutableNode).SetLeavingETA(25)
	assert.Equal(t, insolar.NodeLeaving, n.GetState())
	assert.EqualValues(t, 25, n.LeavingETA())
}

func TestNode_ShortID(t *testing.T) {
	n := NewNode(testutils.RandomRef(), insolar.StaticRoleVirtual, nil, "127.0.0.1", "123")
	assert.EqualValues(t, GenerateUintShortID(n.ID()), n.ShortID())
	n.(MutableNode).SetShortID(11)
	assert.EqualValues(t, 11, n.ShortID())
}

func TestClaimToNode(t *testing.T) {
	address, err := packets.NewNodeAddress("123.234.55.66:12345")
	require.NoError(t, err)

	claim := packets.NodeJoinClaim{
		NodeRef:     testutils.RandomRef(),
		NodePK:      testutils.BrokenPK(),
		ShortNodeID: 10,
		NodeAddress: address,
	}

	_, err = ClaimToNode("", &claim)
	assert.Error(t, err)
	claim.NodePK = [packets.PublicKeyLength]byte{}
	n, err := ClaimToNode("", &claim)
	assert.NoError(t, err)
	assert.Equal(t, claim.NodeRef, n.ID())
	assert.EqualValues(t, 10, n.ShortID())
	assert.Equal(t, claim.NodeAddress.String(), n.Address())
}
