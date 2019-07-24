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

package executor_test

import (
	"context"
	"math/rand"
	"testing"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/insolar/insolar/ledger/object"
	"github.com/insolar/insolar/testutils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilamentModifierDefault_SetResult(t *testing.T) {
	t.Parallel()
	mc := minimock.NewController(t)
	ctx := inslogger.TestContext(t)

	var (
		pcs        insolar.PlatformCryptographyScheme
		indexes    object.IndexStorage
		records    object.RecordStorage
		calculator *executor.FilamentCalculatorMock
		manager    *executor.FilamentModifierDefault
		sender     *bus.SenderMock
		pulses     *pulse.CalculatorMock
	)
	resetComponents := func() {
		pcs = testutils.NewPlatformCryptographyScheme()
		indexes = object.NewIndexStorageMemory()
		records = object.NewRecordMemory()
		calculator = executor.NewFilamentCalculatorMock(mc)
		sender = bus.NewSenderMock(mc)
		manager = executor.NewFilamentModifier(indexes, records, pcs, calculator, pulses, sender)
	}

	validResult := record.Result{Object: gen.ID()}

	resetComponents()
	t.Run("object id is empty", func(t *testing.T) {
		_, err := manager.SetResult(ctx, insolar.ID{}, gen.JetID(), validResult)
		assert.Error(t, err)

		mc.Finish()
	})

	resetComponents()
	t.Run("jet is not valid", func(t *testing.T) {
		_, err := manager.SetResult(ctx, gen.ID(), insolar.JetID{}, validResult)
		assert.Error(t, err)

		mc.Finish()
	})

	resetComponents()
	t.Run("index does not exist", func(t *testing.T) {
		_, err := manager.SetResult(ctx, gen.ID(), gen.JetID(), validResult)
		assert.Error(t, err)

		mc.Finish()
	})

	resetComponents()
	t.Run("problems during duplicates searching", func(t *testing.T) {
		calculator.ResultDuplicateMock.Return(nil, errors.New("smth terrible happened"))
		resultID := gen.ID()
		err := indexes.SetIndex(ctx, resultID.Pulse(), record.Index{
			ObjID:    validResult.Object,
			Lifeline: record.Lifeline{},
		})
		require.NoError(t, err)

		_, err = manager.SetResult(ctx, resultID, gen.JetID(), validResult)
		require.Error(t, err)
		require.Contains(t, err.Error(), "smth terrible happened")

		mc.Finish()
	})

	resetComponents()
	t.Run("happy basic", func(t *testing.T) {
		reqID := gen.ID()
		validResult := record.Result{Object: gen.ID(), Request: *insolar.NewReference(reqID)}

		resultID := gen.ID()
		resultID.SetPulse(insolar.FirstPulseNumber + 2)
		latestPendingID := gen.ID()
		latestPendingID.SetPulse(insolar.FirstPulseNumber + 1)
		jetID := gen.JetID()

		expectedFilamentRecord := record.PendingFilament{
			RecordID:       resultID,
			PreviousRecord: &latestPendingID,
		}
		virtual := record.Wrap(expectedFilamentRecord)
		hash := record.HashVirtual(pcs.ReferenceHasher(), virtual)
		expectedFilamentRecordID := *insolar.NewID(resultID.Pulse(), hash)

		calculator.PendingRequestsFunc = func(_ context.Context, pn insolar.PulseNumber, id insolar.ID) ([]record.CompositeFilamentRecord, error) {
			require.Equal(t, resultID.Pulse(), pn)
			require.Equal(t, validResult.Object, id)

			return []record.CompositeFilamentRecord{{RecordID: expectedFilamentRecordID}}, nil
		}
		calculator.ResultDuplicateFunc = func(_ context.Context, inPN insolar.PulseNumber, inObjID insolar.ID, inResID insolar.ID, inRes record.Result) (*record.CompositeFilamentRecord, error) {
			require.Equal(t, inPN, resultID.Pulse())
			require.Equal(t, validResult.Object, inObjID)
			require.Equal(t, resultID, inResID)
			return nil, nil
		}
		calculator.FindRecordFunc = func(_ context.Context, startFrom insolar.ID, objID insolar.ID, reqID insolar.ID) (r record.CompositeFilamentRecord, r1 error) {
			require.Equal(t, validResult.Object, objID)
			require.Equal(t, *validResult.Request.Record(), reqID)
			require.Equal(t, expectedFilamentRecordID, startFrom)

			req := record.IncomingRequest{}
			reqVirt := record.Wrap(req)
			return record.CompositeFilamentRecord{Record: record.Material{Virtual: &reqVirt}}, nil
		}

		latestPendingPulse := latestPendingID.Pulse()
		err := indexes.SetIndex(ctx, resultID.Pulse(), record.Index{
			ObjID: validResult.Object,
			Lifeline: record.Lifeline{
				PendingPointer:      &latestPendingID,
				EarliestOpenRequest: &latestPendingPulse,
			},
		})
		require.NoError(t, err)

		_, err = manager.SetResult(ctx, resultID, jetID, validResult)
		assert.NoError(t, err)

		idx, err := indexes.ForID(ctx, resultID.Pulse(), validResult.Object)
		require.NoError(t, err)

		require.Equal(t, expectedFilamentRecordID, *idx.Lifeline.PendingPointer)
		require.Equal(t, resultID.Pulse(), *idx.Lifeline.EarliestOpenRequest)

		rec, err := records.ForID(ctx, expectedFilamentRecordID)
		require.NoError(t, err)
		virtual = record.Wrap(expectedFilamentRecord)
		require.Equal(t, record.Material{Virtual: &virtual, JetID: jetID}, rec)

		rec, err = records.ForID(ctx, resultID)
		require.NoError(t, err)
		virtual = record.Wrap(validResult)
		require.Equal(t, record.Material{Virtual: &virtual, JetID: jetID}, rec)

		mc.Finish()
	})

	resetComponents()
	t.Run("notification about outgoing", func(t *testing.T) {
		reqID := gen.ID()
		outReqID := gen.ID()
		validResult := record.Result{Object: gen.ID(), Request: *insolar.NewReference(reqID)}

		resultID := gen.ID()
		resultID.SetPulse(insolar.FirstPulseNumber + 2)
		latestPendingID := gen.ID()
		latestPendingID.SetPulse(insolar.FirstPulseNumber + 1)
		jetID := gen.JetID()

		expectedFilamentRecord := record.PendingFilament{
			RecordID:       resultID,
			PreviousRecord: &latestPendingID,
		}
		virtual := record.Wrap(expectedFilamentRecord)
		hash := record.HashVirtual(pcs.ReferenceHasher(), virtual)
		expectedFilamentRecordID := *insolar.NewID(resultID.Pulse(), hash)

		caller := gen.Reference()
		calculator.PendingRequestsFunc = func(_ context.Context, pn insolar.PulseNumber, id insolar.ID) ([]record.CompositeFilamentRecord, error) {
			require.Equal(t, resultID.Pulse(), pn)
			require.Equal(t, validResult.Object, id)

			req := record.OutgoingRequest{
				ReturnMode: record.ReturnSaga,
				Caller:     caller,
				Reason:     *insolar.NewReference(reqID),
			}
			reqVirt := record.Wrap(req)
			return []record.CompositeFilamentRecord{{RecordID: outReqID, Record: record.Material{Virtual: &reqVirt}}}, nil
		}
		calculator.ResultDuplicateFunc = func(_ context.Context, inPN insolar.PulseNumber, inObjID insolar.ID, inResID insolar.ID, inRes record.Result) (*record.CompositeFilamentRecord, error) {
			require.Equal(t, inPN, resultID.Pulse())
			require.Equal(t, validResult.Object, inObjID)
			require.Equal(t, resultID, inResID)
			return nil, nil
		}
		calculator.FindRecordFunc = func(_ context.Context, startFrom insolar.ID, objID insolar.ID, reqID insolar.ID) (r record.CompositeFilamentRecord, r1 error) {
			require.Equal(t, validResult.Object, objID)
			require.Equal(t, *validResult.Request.Record(), reqID)
			require.Equal(t, expectedFilamentRecordID, startFrom)

			req := record.IncomingRequest{}
			reqVirt := record.Wrap(req)
			return record.CompositeFilamentRecord{RecordID: reqID, Record: record.Material{Virtual: &reqVirt}}, nil
		}

		sender.SendRoleFunc = func(_ context.Context, msg *message.Message, role insolar.DynamicRole, objRef insolar.Reference) (r <-chan *message.Message, r1 func()) {
			require.Equal(t, insolar.DynamicRoleVirtualExecutor, role)
			require.Equal(t, validResult.Object, *objRef.Record())

			notification := payload.SagaCallAcceptNotification{}
			err := notification.Unmarshal(msg.Payload)
			require.NoError(t, err)

			require.Equal(t, validResult.Object, notification.ObjectID)
			require.Equal(t, outReqID, notification.OutgoingReqID)

			return nil, func() {}
		}

		latestPendingPulse := latestPendingID.Pulse()
		err := indexes.SetIndex(ctx, resultID.Pulse(), record.Index{
			ObjID: validResult.Object,
			Lifeline: record.Lifeline{
				PendingPointer:      &latestPendingID,
				EarliestOpenRequest: &latestPendingPulse,
			},
		})
		require.NoError(t, err)

		_, err = manager.SetResult(ctx, resultID, jetID, validResult)
		assert.NoError(t, err)

		mc.Finish()
	})

	resetComponents()
	t.Run("happy resets pending pointers in index", func(t *testing.T) {
		resultID := gen.ID()
		resultID.SetPulse(insolar.FirstPulseNumber + 2)
		latestPendingID := gen.ID()
		latestPendingID.SetPulse(insolar.FirstPulseNumber + 1)
		jetID := gen.JetID()

		calculator.PendingRequestsFunc = func(_ context.Context, pn insolar.PulseNumber, id insolar.ID) ([]record.CompositeFilamentRecord, error) {
			require.Equal(t, resultID.Pulse(), pn)
			require.Equal(t, validResult.Object, id)

			return []record.CompositeFilamentRecord{}, nil
		}
		calculator.ResultDuplicateFunc = func(_ context.Context, pn insolar.PulseNumber, objID insolar.ID, inResID insolar.ID, _ record.Result) (*record.CompositeFilamentRecord, error) {
			require.Equal(t, resultID, inResID)
			require.Equal(t, validResult.Object, objID)
			require.Equal(t, resultID.Pulse(), pn)

			return nil, nil
		}

		latestPendingPulse := latestPendingID.Pulse()
		err := indexes.SetIndex(ctx, resultID.Pulse(), record.Index{
			ObjID: validResult.Object,
			Lifeline: record.Lifeline{
				PendingPointer:      &latestPendingID,
				EarliestOpenRequest: &latestPendingPulse,
			},
		})
		require.NoError(t, err)

		_, err = manager.SetResult(ctx, resultID, jetID, validResult)
		assert.NoError(t, err)

		idx, err := indexes.ForID(ctx, resultID.Pulse(), validResult.Object)
		require.NoError(t, err)

		assert.Nil(t, idx.Lifeline.EarliestOpenRequest)

		mc.Finish()
	})
}

func TestFilamentCalculatorDefault_Requests(t *testing.T) {
	t.Parallel()
	mc := minimock.NewController(t)
	ctx := inslogger.TestContext(t)

	var (
		indexes    object.IndexStorage
		records    *object.RecordMemory
		pcs        insolar.PlatformCryptographyScheme
		calculator *executor.FilamentCalculatorDefault
	)
	resetComponents := func() {
		indexes = object.NewIndexStorageMemory()
		records = object.NewRecordMemory()
		pcs = testutils.NewPlatformCryptographyScheme()
		calculator = executor.NewFilamentCalculator(indexes, records, nil, nil, nil, nil)
	}

	resetComponents()
	t.Run("returns error if object does not exist", func(t *testing.T) {
		_, err := calculator.Requests(ctx, gen.ID(), gen.ID(), gen.PulseNumber(), gen.PulseNumber())
		assert.Error(t, err)

		mc.Finish()
	})

	resetComponents()
	t.Run("empty response", func(t *testing.T) {
		objectID := gen.ID()
		fromID := gen.ID()
		err := indexes.SetIndex(ctx, fromID.Pulse(), record.Index{
			ObjID: objectID,
		})
		require.NoError(t, err)

		recs, err := calculator.Requests(ctx, objectID, fromID, gen.PulseNumber(), gen.PulseNumber())
		assert.NoError(t, err)
		assert.Equal(t, 0, len(recs))

		mc.Finish()
	})

	resetComponents()
	t.Run("happy basic", func(t *testing.T) {
		b := newFilamentBuilder(ctx, pcs, records)
		storageRecs := make([]record.CompositeFilamentRecord, 5)
		storageRecs[0] = b.Append(insolar.FirstPulseNumber+1, record.IncomingRequest{Nonce: rand.Uint64(), CallType: record.CTMethod})
		storageRecs[1] = b.Append(insolar.FirstPulseNumber+2, record.IncomingRequest{Nonce: rand.Uint64(), CallType: record.CTMethod})
		storageRecs[2] = b.Append(insolar.FirstPulseNumber+2, record.IncomingRequest{Nonce: rand.Uint64(), CallType: record.CTMethod})
		storageRecs[3] = b.Append(insolar.FirstPulseNumber+3, record.IncomingRequest{Nonce: rand.Uint64(), CallType: record.CTMethod})
		storageRecs[4] = b.Append(insolar.FirstPulseNumber+4, record.IncomingRequest{Nonce: rand.Uint64(), CallType: record.CTMethod})

		objectID := gen.ID()
		fromID := storageRecs[3].MetaID
		earliestPending := storageRecs[0].MetaID.Pulse()
		err := indexes.SetIndex(ctx, fromID.Pulse(), record.Index{
			ObjID: objectID,
			Lifeline: record.Lifeline{
				PendingPointer:      &storageRecs[3].MetaID,
				EarliestOpenRequest: &earliestPending,
			},
		})
		require.NoError(t, err)

		// First time, records accessed from storage.
		recs, err := calculator.Requests(ctx, objectID, fromID, storageRecs[1].MetaID.Pulse(), storageRecs[3].MetaID.Pulse())
		assert.NoError(t, err)
		require.Equal(t, 3, len(recs))
		assert.Equal(t, []record.CompositeFilamentRecord{storageRecs[3], storageRecs[2], storageRecs[1]}, recs)

		// Second time storage is cleared. Records are accessed from cache.
		for _, rec := range storageRecs {
			records.DeleteForPN(ctx, rec.MetaID.Pulse())
		}
		recs, err = calculator.Requests(ctx, objectID, fromID, storageRecs[1].MetaID.Pulse(), storageRecs[3].MetaID.Pulse())
		assert.NoError(t, err)
		require.Equal(t, 3, len(recs))
		assert.Equal(t, []record.CompositeFilamentRecord{storageRecs[3], storageRecs[2], storageRecs[1]}, recs)

		mc.Finish()
	})
}

func TestFilamentCalculatorDefault_PendingRequests(t *testing.T) {
	t.Parallel()
	mc := minimock.NewController(t)
	ctx := inslogger.TestContext(t)

	var (
		indexes     object.IndexStorage
		records     object.RecordStorage
		coordinator *jet.CoordinatorMock
		jetFetcher  *executor.JetFetcherMock
		sender      *bus.SenderMock
		pcs         insolar.PlatformCryptographyScheme
		calculator  *executor.FilamentCalculatorDefault
	)
	resetComponents := func() {
		indexes = object.NewIndexStorageMemory()
		records = object.NewRecordMemory()
		coordinator = jet.NewCoordinatorMock(mc)
		jetFetcher = executor.NewJetFetcherMock(mc)
		sender = bus.NewSenderMock(mc)
		pcs = testutils.NewPlatformCryptographyScheme()
		calculator = executor.NewFilamentCalculator(indexes, records, coordinator, jetFetcher, sender, nil)
	}

	resetComponents()
	t.Run("returns error if object does not exist", func(t *testing.T) {
		_, err := calculator.PendingRequests(ctx, gen.PulseNumber(), gen.ID())
		assert.Error(t, err)

		mc.Finish()
	})

	resetComponents()
	t.Run("empty response", func(t *testing.T) {
		objectID := gen.ID()
		fromPulse := gen.PulseNumber()
		err := indexes.SetIndex(ctx, fromPulse, record.Index{
			ObjID: objectID,
		})
		require.NoError(t, err)

		recs, err := calculator.PendingRequests(ctx, fromPulse, objectID)
		require.NoError(t, err)
		require.Equal(t, 0, len(recs))

		mc.Finish()
	})

	resetComponents()
	t.Run("happy basic", func(t *testing.T) {
		b := newFilamentBuilder(ctx, pcs, records)
		rec1 := b.Append(insolar.FirstPulseNumber+1, record.IncomingRequest{Nonce: rand.Uint64(), CallType: record.CTMethod})
		rec2 := b.Append(insolar.FirstPulseNumber+2, record.IncomingRequest{Nonce: rand.Uint64(), CallType: record.CTMethod})
		b.Append(insolar.FirstPulseNumber+3, record.Result{Request: *insolar.NewReference(rec1.RecordID)})
		rec4 := b.Append(insolar.FirstPulseNumber+3, record.IncomingRequest{Nonce: rand.Uint64(), CallType: record.CTMethod})
		b.Append(insolar.FirstPulseNumber+4, record.IncomingRequest{Nonce: rand.Uint64(), CallType: record.CTMethod})

		objectID := gen.ID()
		fromPulse := rec4.MetaID.Pulse()
		earliestPending := rec1.MetaID.Pulse()
		err := indexes.SetIndex(ctx, fromPulse, record.Index{
			ObjID: objectID,
			Lifeline: record.Lifeline{
				PendingPointer:      &rec4.MetaID,
				EarliestOpenRequest: &earliestPending,
			},
		})
		require.NoError(t, err)

		recs, err := calculator.PendingRequests(ctx, fromPulse, objectID)
		require.NoError(t, err)
		require.Equal(t, 2, len(recs))
		require.Equal(t, []record.CompositeFilamentRecord{rec2, rec4}, recs)

		mc.Finish()
	})

	resetComponents()
	t.Run("happy fetches from light", func(t *testing.T) {
		b := newFilamentBuilder(ctx, pcs, records)
		rec1 := b.Append(insolar.FirstPulseNumber+1, record.IncomingRequest{Nonce: rand.Uint64()})
		rec2 := b.Append(insolar.FirstPulseNumber+2, record.IncomingRequest{Nonce: rand.Uint64()})
		// This result is not in the storage.
		missingRec := b.AppendNoPersist(insolar.FirstPulseNumber+3, record.Result{Request: *insolar.NewReference(rec1.RecordID)})
		rec4 := b.Append(insolar.FirstPulseNumber+4, record.IncomingRequest{Nonce: rand.Uint64()})
		b.Append(insolar.FirstPulseNumber+5, record.IncomingRequest{Nonce: rand.Uint64()})

		objectID := gen.ID()
		fromPulse := rec4.MetaID.Pulse()
		earliestPending := rec1.MetaID.Pulse()
		err := indexes.SetIndex(ctx, fromPulse, record.Index{
			ObjID: objectID,
			Lifeline: record.Lifeline{
				PendingPointer:      &rec4.MetaID,
				EarliestOpenRequest: &earliestPending,
			},
		})
		require.NoError(t, err)

		coordinator.IsBeyondLimitFunc = func(_ context.Context, target insolar.PulseNumber) (bool, error) {
			require.Equal(t, missingRec.MetaID.Pulse(), target)
			return false, nil
		}

		jetID := gen.JetID()
		jetFetcher.FetchFunc = func(_ context.Context, targetID insolar.ID, pn insolar.PulseNumber) (*insolar.ID, error) {
			require.Equal(t, objectID, targetID)
			require.Equal(t, missingRec.MetaID.Pulse(), pn)
			id := insolar.ID(jetID)
			return &id, nil
		}

		node := gen.Reference()
		coordinator.NodeForJetFunc = func(_ context.Context, jet insolar.ID, target insolar.PulseNumber) (*insolar.Reference, error) {
			require.Equal(t, insolar.ID(jetID), jet)
			require.Equal(t, missingRec.MetaID.Pulse(), target)
			return &node, nil
		}

		coordinator.MeMock.Return(node)

		recs, err := calculator.PendingRequests(ctx, fromPulse, objectID)
		require.Error(t, err, "returns error if trying to fetch from self")

		coordinator.MeMock.Return(gen.Reference())

		sender.SendTargetFunc = func(_ context.Context, msg *message.Message, target insolar.Reference) (<-chan *message.Message, func()) {
			pl, err := payload.Unmarshal(msg.Payload)
			require.NoError(t, err)

			getFilament, ok := pl.(*payload.GetFilament)
			require.True(t, ok)

			require.Equal(t, objectID, getFilament.ObjectID)
			require.Equal(t, missingRec.MetaID, getFilament.StartFrom)
			require.Equal(t, earliestPending, getFilament.ReadUntil)

			require.NoError(t, err)
			respMsg, err := payload.NewMessage(&payload.FilamentSegment{
				ObjectID: objectID,
				Records:  []record.CompositeFilamentRecord{missingRec},
			})
			require.NoError(t, err)
			meta := payload.Meta{Payload: respMsg.Payload}
			buf, err := meta.Marshal()
			require.NoError(t, err)
			respMsg.Payload = buf
			ch := make(chan *message.Message, 1)
			ch <- respMsg
			return ch, func() {}
		}

		recs, err = calculator.PendingRequests(ctx, fromPulse, objectID)
		require.NoError(t, err)
		require.Equal(t, 2, len(recs))
		require.Equal(t, []record.CompositeFilamentRecord{rec2, rec4}, recs)

		mc.Finish()
	})

	resetComponents()
	t.Run("happy fetches from heavy", func(t *testing.T) {
		b := newFilamentBuilder(ctx, pcs, records)
		rec1 := b.Append(insolar.FirstPulseNumber+1, record.IncomingRequest{Nonce: rand.Uint64()})
		rec2 := b.Append(insolar.FirstPulseNumber+2, record.IncomingRequest{Nonce: rand.Uint64()})
		// This result is not in the storage.
		missingRec := b.AppendNoPersist(insolar.FirstPulseNumber+3, record.Result{Request: *insolar.NewReference(rec1.RecordID)})
		rec4 := b.Append(insolar.FirstPulseNumber+4, record.IncomingRequest{Nonce: rand.Uint64()})
		b.Append(insolar.FirstPulseNumber+5, record.IncomingRequest{Nonce: rand.Uint64()})

		objectID := gen.ID()
		fromPulse := rec4.MetaID.Pulse()
		earliestPending := rec1.MetaID.Pulse()
		err := indexes.SetIndex(ctx, fromPulse, record.Index{
			ObjID: objectID,
			Lifeline: record.Lifeline{
				PendingPointer:      &rec4.MetaID,
				EarliestOpenRequest: &earliestPending,
			},
		})
		require.NoError(t, err)

		coordinator.IsBeyondLimitFunc = func(_ context.Context, target insolar.PulseNumber) (bool, error) {
			require.Equal(t, missingRec.MetaID.Pulse(), target)
			return true, nil
		}

		node := gen.Reference()
		coordinator.HeavyFunc = func(_ context.Context) (*insolar.Reference, error) {
			return &node, nil
		}
		coordinator.MeMock.Return(node)

		recs, err := calculator.PendingRequests(ctx, fromPulse, objectID)
		assert.Error(t, err, "returns error if trying to fetch from self")

		coordinator.MeMock.Return(gen.Reference())

		sender.SendTargetFunc = func(_ context.Context, msg *message.Message, target insolar.Reference) (<-chan *message.Message, func()) {
			pl, err := payload.Unmarshal(msg.Payload)
			require.NoError(t, err)

			getFilament, ok := pl.(*payload.GetFilament)
			require.True(t, ok)

			require.Equal(t, objectID, getFilament.ObjectID)
			require.Equal(t, missingRec.MetaID, getFilament.StartFrom)
			require.Equal(t, earliestPending, getFilament.ReadUntil)

			require.NoError(t, err)
			respMsg, err := payload.NewMessage(&payload.FilamentSegment{
				ObjectID: objectID,
				Records:  []record.CompositeFilamentRecord{missingRec},
			})
			require.NoError(t, err)
			meta := payload.Meta{Payload: respMsg.Payload}
			buf, err := meta.Marshal()
			require.NoError(t, err)
			respMsg.Payload = buf
			ch := make(chan *message.Message, 1)
			ch <- respMsg
			return ch, func() {}
		}

		recs, err = calculator.PendingRequests(ctx, fromPulse, objectID)
		require.NoError(t, err)
		require.Equal(t, 2, len(recs))
		require.Equal(t, []record.CompositeFilamentRecord{rec2, rec4}, recs)

		mc.Finish()
	})

	resetComponents()
	t.Run("ignore not detached outgoings", func(t *testing.T) {
		b := newFilamentBuilder(ctx, pcs, records)
		rec1 := b.Append(insolar.FirstPulseNumber+1, record.OutgoingRequest{
			Nonce:      rand.Uint64(),
			CallType:   record.CTMethod,
			ReturnMode: record.ReturnResult,
		})

		objectID := gen.ID()
		fromPulse := rec1.MetaID.Pulse()
		earliestPending := rec1.MetaID.Pulse()
		err := indexes.SetIndex(ctx, fromPulse, record.Index{
			ObjID: objectID,
			Lifeline: record.Lifeline{
				PendingPointer:      &rec1.MetaID,
				EarliestOpenRequest: &earliestPending,
			},
		})
		require.NoError(t, err)

		recs, err := calculator.PendingRequests(ctx, fromPulse, objectID)
		require.NoError(t, err)
		require.Equal(t, 0, len(recs))

		mc.Finish()
	})

	// b := newFilamentBuilder(ctx, pcs, records)
	// rec1 := b.Append(insolar.FirstPulseNumber+1, record.IncomingRequest{Nonce: rand.Uint64()})
	// rec2 := b.Append(insolar.FirstPulseNumber+2, record.IncomingRequest{Nonce: rand.Uint64()})
	// // This result is not in the storage.
	// missingRec := b.AppendNoPersist(insolar.FirstPulseNumber+3, record.Result{Request: *insolar.NewReference(rec1.RecordID)})
	// rec4 := b.Append(insolar.FirstPulseNumber+4, record.IncomingRequest{Nonce: rand.Uint64()})
	// b.Append(insolar.FirstPulseNumber+5, record.IncomingRequest{Nonce: rand.Uint64()})
	//

	resetComponents()
	t.Run("ignore closed outgoing", func(t *testing.T) {
		b := newFilamentBuilder(ctx, pcs, records)
		reason := b.Append(insolar.FirstPulseNumber+1, record.IncomingRequest{Nonce: rand.Uint64()})
		outgoing := b.Append(insolar.FirstPulseNumber+1, record.OutgoingRequest{
			Nonce:      rand.Uint64(),
			Reason:     *insolar.NewReference(reason.RecordID),
			CallType:   record.CTMethod,
			ReturnMode: record.ReturnSaga,
		})
		_ = b.Append(insolar.FirstPulseNumber+1, record.Result{Request: *insolar.NewReference(reason.RecordID)})
		outgoingRes := b.Append(insolar.FirstPulseNumber+1, record.Result{Request: *insolar.NewReference(outgoing.RecordID)})

		objectID := gen.ID()
		fromPulse := outgoingRes.MetaID.Pulse()
		earliestPending := outgoingRes.MetaID.Pulse()
		err := indexes.SetIndex(ctx, fromPulse, record.Index{
			ObjID: objectID,
			Lifeline: record.Lifeline{
				PendingPointer:      &outgoingRes.MetaID,
				EarliestOpenRequest: &earliestPending,
			},
		})
		require.NoError(t, err)

		recs, err := calculator.PendingRequests(ctx, fromPulse, objectID)
		require.NoError(t, err)
		require.Equal(t, 0, len(recs))

		mc.Finish()
	})

	resetComponents()
	t.Run("return outgoing with closed reason and no result", func(t *testing.T) {
		b := newFilamentBuilder(ctx, pcs, records)
		reason := b.Append(insolar.FirstPulseNumber+1, record.IncomingRequest{Nonce: rand.Uint64()})
		outgoing := b.Append(insolar.FirstPulseNumber+1, record.OutgoingRequest{
			Nonce:      rand.Uint64(),
			Reason:     *insolar.NewReference(reason.RecordID),
			CallType:   record.CTMethod,
			ReturnMode: record.ReturnSaga,
		})
		reasonRes := b.Append(insolar.FirstPulseNumber+1, record.Result{Request: *insolar.NewReference(reason.RecordID)})

		objectID := gen.ID()
		fromPulse := reasonRes.MetaID.Pulse()
		earliestPending := reasonRes.MetaID.Pulse()
		err := indexes.SetIndex(ctx, fromPulse, record.Index{
			ObjID: objectID,
			Lifeline: record.Lifeline{
				PendingPointer:      &reasonRes.MetaID,
				EarliestOpenRequest: &earliestPending,
			},
		})
		require.NoError(t, err)

		recs, err := calculator.PendingRequests(ctx, fromPulse, objectID)
		require.NoError(t, err)
		require.Equal(t, 1, len(recs))
		require.Equal(t, outgoing, recs[0])

		mc.Finish()
	})

	resetComponents()
	t.Run("returns error, when there is closed outgoing with open reason", func(t *testing.T) {
		b := newFilamentBuilder(ctx, pcs, records)
		reason := b.Append(insolar.FirstPulseNumber+1, record.IncomingRequest{Nonce: rand.Uint64()})
		outgoing := b.Append(insolar.FirstPulseNumber+1, record.OutgoingRequest{
			Nonce:      rand.Uint64(),
			Reason:     *insolar.NewReference(reason.RecordID),
			CallType:   record.CTMethod,
			ReturnMode: record.ReturnSaga,
		})
		outgoingRes := b.Append(insolar.FirstPulseNumber+1, record.Result{Request: *insolar.NewReference(outgoing.RecordID)})

		objectID := gen.ID()
		fromPulse := outgoingRes.MetaID.Pulse()
		earliestPending := outgoingRes.MetaID.Pulse()
		err := indexes.SetIndex(ctx, fromPulse, record.Index{
			ObjID: objectID,
			Lifeline: record.Lifeline{
				PendingPointer:      &outgoingRes.MetaID,
				EarliestOpenRequest: &earliestPending,
			},
		})
		require.NoError(t, err)

		_, err = calculator.PendingRequests(ctx, fromPulse, objectID)
		require.Error(t, err)

		mc.Finish()
	})
}

func TestFilamentCalculatorDefault_ResultDuplicate(t *testing.T) {
	t.Parallel()
	mc := minimock.NewController(t)
	ctx := inslogger.TestContext(t)

	var (
		indexes     object.IndexStorage
		records     object.RecordStorage
		coordinator *jet.CoordinatorMock
		jetFetcher  *executor.JetFetcherMock
		sender      *bus.SenderMock
		pcs         insolar.PlatformCryptographyScheme
		calculator  *executor.FilamentCalculatorDefault
	)
	resetComponents := func() {
		indexes = object.NewIndexStorageMemory()
		records = object.NewRecordMemory()
		coordinator = jet.NewCoordinatorMock(mc)
		jetFetcher = executor.NewJetFetcherMock(mc)
		sender = bus.NewSenderMock(mc)
		pcs = testutils.NewPlatformCryptographyScheme()
		calculator = executor.NewFilamentCalculator(indexes, records, coordinator, jetFetcher, sender, nil)
	}

	resetComponents()
	t.Run("returns error if reason is empty", func(t *testing.T) {
		_, err := calculator.ResultDuplicate(ctx, gen.PulseNumber(), gen.ID(), gen.ID(), record.Result{})
		assert.Error(t, err)

		mc.Finish()
	})

	resetComponents()
	t.Run("no records", func(t *testing.T) {
		objectID := gen.ID()
		fromPulse := gen.PulseNumber()
		err := indexes.SetIndex(ctx, fromPulse, record.Index{
			ObjID: objectID,
		})
		require.NoError(t, err)

		res, err := calculator.ResultDuplicate(ctx, fromPulse, objectID, gen.ID(), record.Result{Request: gen.Reference()})

		assert.NoError(t, err)
		assert.Nil(t, res)

		mc.Finish()
	})

	resetComponents()
	t.Run("returns result. result duplicate is found", func(t *testing.T) {
		b := newFilamentBuilder(ctx, pcs, records)
		req := record.IncomingRequest{Nonce: rand.Uint64(), Reason: *insolar.NewReference(*insolar.NewID(insolar.FirstPulseNumber, nil))}
		req1 := b.Append(insolar.FirstPulseNumber+1, req)
		res := record.Result{Request: *insolar.NewReference(req1.RecordID)}
		res1 := b.Append(insolar.FirstPulseNumber+2, res)

		objectID := gen.ID()
		fromPulse := res1.MetaID.Pulse()
		err := indexes.SetIndex(ctx, fromPulse, record.Index{
			ObjID: objectID,
			Lifeline: record.Lifeline{
				PendingPointer: &res1.MetaID,
			},
		})
		require.NoError(t, err)

		fRes, err := calculator.ResultDuplicate(ctx, fromPulse, objectID, res1.RecordID, res)
		require.NoError(t, err)
		require.Equal(t, *fRes, res1)

		mc.Finish()
	})

	resetComponents()
	t.Run("returns result. request not found", func(t *testing.T) {
		b := newFilamentBuilder(ctx, pcs, records)
		req := b.Append(
			insolar.FirstPulseNumber+1,
			record.IncomingRequest{Nonce: rand.Uint64(), Reason: *insolar.NewReference(*insolar.NewID(insolar.FirstPulseNumber, nil))},
		)

		objectID := gen.ID()
		fromPulse := req.MetaID.Pulse()
		err := indexes.SetIndex(ctx, fromPulse, record.Index{
			ObjID: objectID,
			Lifeline: record.Lifeline{
				PendingPointer: &req.MetaID,
			},
		})
		require.NoError(t, err)

		_, err = calculator.ResultDuplicate(ctx, fromPulse, objectID, req.RecordID, record.Result{Request: gen.Reference()})
		require.Error(t, err)

		mc.Finish()
	})

	resetComponents()
	t.Run("returns no result. request found", func(t *testing.T) {
		b := newFilamentBuilder(ctx, pcs, records)
		req := record.IncomingRequest{Nonce: rand.Uint64(), Reason: *insolar.NewReference(*insolar.NewID(insolar.FirstPulseNumber, nil))}
		req1 := b.Append(insolar.FirstPulseNumber+1, req)
		res := record.Result{Request: *insolar.NewReference(req1.RecordID)}
		resID := insolar.NewID(insolar.FirstPulseNumber+1, []byte{1})

		objectID := gen.ID()
		fromPulse := req1.MetaID.Pulse()
		err := indexes.SetIndex(ctx, fromPulse, record.Index{
			ObjID: objectID,
			Lifeline: record.Lifeline{
				PendingPointer: &req1.MetaID,
			},
		})
		require.NoError(t, err)

		fRes, err := calculator.ResultDuplicate(ctx, fromPulse, objectID, *resID, res)
		require.NoError(t, err)
		require.Nil(t, fRes)

		mc.Finish()
	})
}

func TestFilamentCalculatorDefault_RequestDuplicate(t *testing.T) {
	t.Parallel()
	mc := minimock.NewController(t)
	ctx := inslogger.TestContext(t)

	var (
		indexes     object.IndexStorage
		records     object.RecordStorage
		coordinator *jet.CoordinatorMock
		jetFetcher  *executor.JetFetcherMock
		sender      *bus.SenderMock
		pcs         insolar.PlatformCryptographyScheme
		calculator  *executor.FilamentCalculatorDefault
	)
	resetComponents := func() {
		indexes = object.NewIndexStorageMemory()
		records = object.NewRecordMemory()
		coordinator = jet.NewCoordinatorMock(mc)
		jetFetcher = executor.NewJetFetcherMock(mc)
		sender = bus.NewSenderMock(mc)
		pcs = testutils.NewPlatformCryptographyScheme()
		calculator = executor.NewFilamentCalculator(indexes, records, coordinator, jetFetcher, sender, nil)
	}

	resetComponents()
	t.Run("returns error if reason is empty", func(t *testing.T) {
		_, _, err := calculator.RequestDuplicate(ctx, gen.ID(), gen.ID(), &record.IncomingRequest{})
		assert.Error(t, err)

		mc.Finish()
	})

	resetComponents()
	t.Run("no records", func(t *testing.T) {
		objectID := gen.ID()
		fromPulse := gen.PulseNumber()
		err := indexes.SetIndex(ctx, fromPulse, record.Index{
			ObjID: objectID,
		})
		require.NoError(t, err)

		req, res, err := calculator.RequestDuplicate(ctx, objectID, gen.IDWithPulse(fromPulse), &record.IncomingRequest{
			Reason: gen.Reference(),
		})

		assert.NoError(t, err)
		assert.Nil(t, req)
		assert.Nil(t, res)

		mc.Finish()
	})

	resetComponents()
	t.Run("returns request and result", func(t *testing.T) {
		b := newFilamentBuilder(ctx, pcs, records)
		reason := *insolar.NewReference(*insolar.NewID(insolar.FirstPulseNumber, nil))
		req := record.IncomingRequest{Nonce: rand.Uint64(), Reason: reason}
		req1 := b.Append(insolar.FirstPulseNumber+1, req)
		res1 := b.Append(insolar.FirstPulseNumber+2, record.Result{Request: *insolar.NewReference(req1.RecordID)})

		objectID := gen.ID()
		err := indexes.SetIndex(ctx, req1.RecordID.Pulse(), record.Index{
			ObjID: objectID,
			Lifeline: record.Lifeline{
				PendingPointer: &res1.MetaID,
			},
		})
		require.NoError(t, err)

		fReq, fRes, err := calculator.RequestDuplicate(ctx, objectID, req1.RecordID, &req)
		assert.NoError(t, err)
		require.Equal(t, fReq, &req1)
		assert.Equal(t, fRes, &res1)

		mc.Finish()
	})

	resetComponents()
	t.Run("returns only request", func(t *testing.T) {
		b := newFilamentBuilder(ctx, pcs, records)
		reason := *insolar.NewReference(*insolar.NewID(insolar.FirstPulseNumber, nil))
		reqR := record.IncomingRequest{Nonce: rand.Uint64(), Reason: reason}
		req1 := b.Append(insolar.FirstPulseNumber+1, reqR)
		reqR2 := record.IncomingRequest{Nonce: rand.Uint64(), Reason: *insolar.NewReference(*insolar.NewID(insolar.FirstPulseNumber, nil))}
		req2 := b.Append(insolar.FirstPulseNumber+2, reqR2)

		objectID := gen.ID()
		err := indexes.SetIndex(ctx, req1.RecordID.Pulse(), record.Index{
			ObjID: objectID,
			Lifeline: record.Lifeline{
				PendingPointer: &req2.MetaID,
			},
		})
		require.NoError(t, err)

		fReq, fRes, err := calculator.RequestDuplicate(ctx, objectID, req1.RecordID, &reqR)
		require.NoError(t, err)
		require.Equal(t, *fReq, req1)
		require.Nil(t, fRes)

		mc.Finish()
	})

}

type filamentBuilder struct {
	records   object.RecordModifier
	currentID insolar.ID
	ctx       context.Context
	pcs       insolar.PlatformCryptographyScheme
}

func newFilamentBuilder(
	ctx context.Context,
	pcs insolar.PlatformCryptographyScheme,
	records object.RecordModifier,
) *filamentBuilder {
	return &filamentBuilder{
		ctx:     ctx,
		records: records,
		pcs:     pcs,
	}
}

func (b *filamentBuilder) Append(pn insolar.PulseNumber, rec record.Record) record.CompositeFilamentRecord {
	return b.append(pn, rec, true)
}

func (b *filamentBuilder) AppendNoPersist(pn insolar.PulseNumber, rec record.Record) record.CompositeFilamentRecord {
	return b.append(pn, rec, false)
}

func (b *filamentBuilder) append(pn insolar.PulseNumber, rec record.Record, persist bool) record.CompositeFilamentRecord {
	var composite record.CompositeFilamentRecord
	{
		virtual := record.Wrap(rec)
		hash := record.HashVirtual(b.pcs.ReferenceHasher(), virtual)
		id := *insolar.NewID(pn, hash)
		material := record.Material{Virtual: &virtual, JetID: insolar.ZeroJetID}
		if persist {
			err := b.records.Set(b.ctx, id, material)
			if err != nil {
				panic(err)
			}
		}
		composite.RecordID = id
		composite.Record = material
	}

	{
		rec := record.PendingFilament{RecordID: composite.RecordID}
		if !b.currentID.IsEmpty() {
			curr := b.currentID
			rec.PreviousRecord = &curr
		}
		virtual := record.Wrap(rec)
		hash := record.HashVirtual(b.pcs.ReferenceHasher(), virtual)
		id := *insolar.NewID(pn, hash)
		material := record.Material{Virtual: &virtual, JetID: insolar.ZeroJetID}
		if persist {
			err := b.records.Set(b.ctx, id, material)
			if err != nil {
				panic(err)
			}
		}
		composite.MetaID = id
		composite.Meta = material
	}

	b.currentID = composite.MetaID

	return composite
}
