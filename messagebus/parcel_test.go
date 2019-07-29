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

package messagebus

import (
	"testing"
)

func Test_parcelFactory_Create_CheckLogLevel(t *testing.T) {
	//ctx := inslogger.TestContext(t)
	//
	///* Prepare CryptographyService mock, Parcel factory, DelegationToken factory */
	//mock := testutils.NewCryptographyServiceMock(t)
	//mock.SignFunc = func(p []byte) (r *insolar.Signature, r1 error) {
	//	signature := insolar.SignatureFromBytes(nil)
	//	return &signature, nil
	//}
	//mock.GetPublicKeyFunc = func() (r crypto.PublicKey, r1 error) {
	//	return nil, nil
	//}
	//
	//parcelFactory := NewParcelFactory()
	//
	//cm := &component.Manager{}
	//cm.Register(mock, parcelFactory)
	//cm.Inject(parcelFactory)
	//assert.NoError(t, cm.Init(ctx))
	//assert.NoError(t, cm.Start(ctx))
	//
	//ref := testutils.RandomRef()
	//pulse := insolar.Pulse{PulseNumber: 0}
	//msg := message.CallMethod{}
	//
	//parcel, err := parcelFactory.Create(ctx, &msg, ref, nil, pulse)
	//
	//assert.NoError(t, err)
	//assert.Equal(t, inslogger.GetLoggerLevel(ctx), insolar.NoLevel)
	//
	//ctx_new := inslogger.WithLoggerLevel(ctx, insolar.DebugLevel)
	//assert.NotEqual(t, inslogger.GetLoggerLevel(ctx_new), insolar.NoLevel)
	//assert.NotEqual(t, inslogger.GetLoggerLevel(ctx), insolar.DebugLevel)
	//
	//parcel, err = parcelFactory.Create(ctx_new, &msg, ref, nil, pulse)
	//
	////ctx = parcel.Context(ctx)
	//assert.NoError(t, err)
	//assert.Equal(t, inslogger.GetLoggerLevel(ctx), insolar.DebugLevel)
}
