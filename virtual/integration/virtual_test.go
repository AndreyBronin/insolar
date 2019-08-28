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

package integration

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

func TestVirtual_BasicOperations(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	cfg := DefaultVMConfig()
	s, err := NewServer(ctx, cfg, nil)
	require.NoError(t, err)
	defer s.Stop(ctx)

	// First pulse goes in storage then interrupts.
	s.SetPulse(ctx)

	t.Run("happy path", func(t *testing.T) {
		SendMessage(ctx, s, &payload.CallMethod{
			Request: &record.IncomingRequest{
				Caller:          insolar.Reference{},
				CallerPrototype: insolar.Reference{},
				Reason:          insolar.Reference{},
				APINode:         insolar.Reference{},
			},
			PulseNumber: s.pulse.PulseNumber,
		})
	})
}
