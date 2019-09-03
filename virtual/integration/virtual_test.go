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
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

func TestVirtual_BasicOperations(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	cfg := DefaultVMConfig()

	t.Run("happy path", func(t *testing.T) {
		objRef := gen.Reference()
		objID := objRef.GetLocal()
		s, err := NewServer(ctx, cfg, func(meta payload.Meta, pl payload.Payload) []payload.Payload {
			if meta.Receiver != NodeLight() {
				return nil
			}

			switch pl.(type) {
			// getters
			case *payload.SetIncomingRequest, *payload.SetOutgoingRequest:
				return []payload.Payload{&payload.RequestInfo{
					ObjectID:  *objID,
					RequestID: *objID,
					Result:    record.NewResultFromFace(&record.Result{}).Payload,
				}}
			// setters
			case *payload.SetResult, *payload.SetCode:
				return []payload.Payload{&payload.ID{}}
			case *payload.HasPendings:
				return []payload.Payload{&payload.PendingsInfo{HasPendings: true}}
			case *payload.GetObject:
				return []payload.Payload{&payload.State{Record: objRef.Bytes()}}
			}

			panic(fmt.Sprintf("unexpected message to light %T", pl))
		})

		require.NoError(t, err)
		defer s.Stop(ctx)

		// First pulse goes in storage then interrupts.
		s.SetPulse(ctx)

		protoRef := gen.Reference()

		res := SendMessage(ctx, s, &payload.CallMethod{
			Request: &record.IncomingRequest{
				CallType:        record.CTSaveAsChild,
				Caller:          gen.Reference(),
				CallerPrototype: gen.Reference(),
				Reason:          gen.Reference(),
				APINode:         gen.Reference(),
				Object:          &objRef,
				Prototype:       &protoRef,
			},
			PulseNumber: s.pulse.PulseNumber,
		})

		_, isError := res.(*payload.Error)

		require.False(t, isError, "result expected not to be error %v", res)
	})
}
