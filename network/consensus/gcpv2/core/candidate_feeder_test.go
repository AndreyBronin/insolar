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

package core

import (
	"github.com/insolar/insolar/network/consensus/gcpv2/api"
	"testing"

	ccommon "github.com/insolar/insolar/network/consensus/common"

	"github.com/insolar/insolar/network/consensus/gcpv2/common"
	"github.com/insolar/insolar/network/consensus/gcpv2/packets"
	"github.com/stretchr/testify/require"
)

func TestPickNextJoinCandidate(t *testing.T) {
	require.Equal(t, (&SequencialCandidateFeeder{}).PickNextJoinCandidate(), nil)

	s := &SequencialCandidateFeeder{buf: make([]api.CandidateProfile, 1)}
	c := common.NewCandidateProfileMock(t)
	s.buf[0] = c
	require.Equal(t, s.PickNextJoinCandidate(), c)
}

func TestRemoveJoinCandidate(t *testing.T) {
	require.False(t, (&SequencialCandidateFeeder{}).RemoveJoinCandidate(false, ccommon.ShortNodeID(0)))

	s := &SequencialCandidateFeeder{buf: make([]api.CandidateProfile, 1)}
	c := common.NewCandidateProfileMock(t)
	s.buf[0] = c
	c.GetNodeIDMock.Set(func() ccommon.ShortNodeID { return ccommon.ShortNodeID(1) })
	require.False(t, s.RemoveJoinCandidate(false, ccommon.ShortNodeID(2)))

	c.GetNodeIDMock.Set(func() ccommon.ShortNodeID { return ccommon.ShortNodeID(1) })
	require.True(t, s.RemoveJoinCandidate(false, ccommon.ShortNodeID(1)))

	require.Equal(t, s.buf, []api.CandidateProfile(nil))

	s.buf = make([]api.CandidateProfile, 2)
	s.buf[0] = c
	c2 := common.NewCandidateProfileMock(t)
	s.buf[1] = c2
	require.True(t, s.RemoveJoinCandidate(false, ccommon.ShortNodeID(1)))

	require.Equal(t, len(s.buf), 1)

	require.True(t, len(s.buf) > 0 && s.buf[0] == c2)
}

func TestAddJoinCandidate(t *testing.T) {
	require.Panics(t, func() { (&SequencialCandidateFeeder{}).AddJoinCandidate(nil) })

	f := packets.NewFullIntroductionReaderMock(t)
	s := &SequencialCandidateFeeder{}
	s.AddJoinCandidate(f)
	require.True(t, len(s.buf) == 1 && s.buf[0] == f)
}
