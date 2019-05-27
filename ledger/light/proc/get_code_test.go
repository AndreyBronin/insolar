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

package proc_test

import (
	"context"
	"testing"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/blob"
	"github.com/insolar/insolar/ledger/light/proc"
	"github.com/insolar/insolar/ledger/object/mocks"
	"github.com/stretchr/testify/require"
)

func TestGetCode_Proceed(t *testing.T) {
	mc := minimock.NewController(t)
	a := require.New(t)
	ctx := inslogger.TestContext(t)

	replyTo := make(chan bus.Reply, 1)
	blobValue := blob.Blob{Value: []byte{1, 2, 3}}
	blobID := gen.ID()
	codeRec := codeRecord(blobID)
	codeRef := gen.Reference()
	getCode := proc.NewGetCode(codeRef, replyTo)
	records := mocks.NewRecordAccessorMock(mc)
	records.ForIDFunc = func(c context.Context, id insolar.ID) (record.Material, error) {
		a.Equal(*codeRef.Record(), id)
		return codeRec, nil
	}
	blobs := blob.NewAccessorMock(mc)
	blobs.ForIDFunc = func(c context.Context, id insolar.ID) (blob.Blob, error) {
		a.Equal(blobID, id)
		return blobValue, nil
	}
	getCode.Dep.RecordAccessor = records
	getCode.Dep.BlobAccessor = blobs

	err := getCode.Proceed(ctx)
	a.NoError(err)

	unwrappedCodeRec := record.Unwrap(codeRec.Virtual)

	rep := <-replyTo
	a.Equal(bus.Reply{Reply: &reply.Code{
		Code:        blobValue.Value,
		MachineType: unwrappedCodeRec.(*record.Code).MachineType,
	}}, rep)
}

func codeRecord(codeID insolar.ID) record.Material {
	return record.Material{
		Virtual: &record.Virtual{
			Union: &record.Virtual_Code{
				Code: &record.Code{
					Code:        codeID,
					MachineType: insolar.MachineTypeBuiltin,
				},
			},
		},
	}
}
