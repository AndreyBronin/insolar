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

package store

import (
	"io/ioutil"
	"math/rand"
	"os"
	"testing"

	fuzz "github.com/google/gofuzz"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/instrumentation/inslogger"
)

type testKey struct {
	id    []byte
	scope Scope
}

func (k testKey) Scope() Scope {
	return k.scope
}

func (k testKey) ID() []byte {
	return k.id
}

func TestDB_Components(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)

	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	assert.NoError(t, err)
	badger, err := NewBadgerDB(tmpdir)
	defer badger.Stop(ctx)
	require.NoError(t, err)

	mock := NewMemoryMockDB()

	type data struct {
		key   testKey
		value []byte
	}
	var datas []data

	f := fuzz.New().NilChance(0).NumElements(5, 10)
	f = f.Funcs(func(d *data, c fuzz.Continue) {
		id := make([]byte, 10)
		rand.Read(id)
		d.key = testKey{
			scope: Scope(rand.Int31()),
			id:    id,
		}
		d.value = make([]byte, 10)
		rand.Read(d.value)
	})
	f.Fuzz(&datas)

	for _, d := range datas {
		{
			err := badger.Set(d.key, d.value)
			assert.NoError(t, err)
		}
		{
			err := mock.Set(d.key, d.value)
			assert.NoError(t, err)
		}
	}
	for _, d := range datas {
		{
			val, err := badger.Get(d.key)
			assert.NoError(t, err)
			assert.Equal(t, d.value, val)
		}
		{
			val, err := mock.Get(d.key)
			assert.NoError(t, err)
			assert.Equal(t, d.value, val)
		}
	}
}
