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

package object_test

import (
	"math/rand"
	"testing"

	fuzz "github.com/google/gofuzz"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/internal/ledger/store"
	"github.com/insolar/insolar/ledger/object"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIndex_Components(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)

	indexMemory := object.NewInMemoryIndex()
	indexDB := object.NewIndexDB(store.NewMemoryMockDB())

	type tempIndex struct {
		id  insolar.ID
		idx object.Lifeline
	}

	var indices []tempIndex

	f := fuzz.New().Funcs(func(t *tempIndex, c fuzz.Continue) {
		t.id = gen.ID()
		ls := gen.ID()
		pn := gen.PulseNumber()
		t.idx = object.Lifeline{
			LatestState:  &ls,
			LatestUpdate: pn,
			JetID:        gen.JetID(),
			Delegates:    []object.LifelineDelegate{},
		}
	})
	f.NumElements(5, 10).NilChance(0).Fuzz(&indices)
	pn := gen.PulseNumber()

	t.Run("saves correct index-value", func(t *testing.T) {
		for _, i := range indices {
			memErr := indexMemory.Set(ctx, pn, i.id, i.idx)
			dbErr := indexDB.Set(ctx, pn, i.id, i.idx)
			require.NoError(t, memErr)
			require.NoError(t, dbErr)
		}

		for _, i := range indices {
			resIndexMem, memErr := indexMemory.ForID(ctx, pn, i.id)
			resIndexDB, dbErr := indexDB.ForID(ctx, pn, i.id)
			require.NoError(t, memErr)
			require.NoError(t, dbErr)

			idxBuf, _ := i.idx.Marshal()

			memBuf, _ := resIndexMem.Marshal()
			assert.Equal(t, idxBuf, memBuf)
			assert.Equal(t, i.idx.JetID, resIndexMem.JetID)
			assert.Equal(t, i.idx.LatestState, resIndexMem.LatestState)
			assert.Equal(t, i.idx.LatestUpdate, resIndexMem.LatestUpdate)

			dbBuf, _ := resIndexDB.Marshal()
			assert.Equal(t, idxBuf, dbBuf)
			assert.Equal(t, i.idx.JetID, resIndexDB.JetID)
			assert.Equal(t, i.idx.LatestState, resIndexDB.LatestState)
			assert.Equal(t, i.idx.LatestUpdate, resIndexDB.LatestUpdate)
		}
	})

	t.Run("returns error when no index-value for id", func(t *testing.T) {
		t.Parallel()

		for i := int32(0); i < rand.Int31n(10); i++ {
			_, memErr := indexMemory.ForID(ctx, pn, gen.ID())
			_, dbErr := indexDB.ForID(ctx, pn, gen.ID())
			require.Error(t, memErr)
			require.Error(t, dbErr)
			assert.Equal(t, object.ErrLifelineNotFound, memErr)
			assert.Equal(t, object.ErrLifelineNotFound, dbErr)
		}
	})

	t.Run("override indices is ok", func(t *testing.T) {
		t.Parallel()

		indexMemory := object.NewInMemoryIndex()
		indexDB := object.NewIndexDB(store.NewMemoryMockDB())

		for _, i := range indices {
			memErr := indexMemory.Set(ctx, pn, i.id, i.idx)
			dbErr := indexDB.Set(ctx, pn, i.id, i.idx)
			require.NoError(t, memErr)
			require.NoError(t, dbErr)
		}

		for _, i := range indices {
			memErr := indexMemory.Set(ctx, pn, i.id, i.idx)
			dbErr := indexDB.Set(ctx, pn, i.id, i.idx)
			assert.NoError(t, memErr)
			assert.NoError(t, dbErr)
		}
	})
}
