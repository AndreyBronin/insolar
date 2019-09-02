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
//
// +build apitests

package smoketests

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/insolar/insolar/apitests/tests/insolarapi"
	"github.com/insolar/insolar/apitests/tests/internalapi"

	"github.com/insolar/insolar/apitests/apihelper"
	"github.com/insolar/insolar/apitests/apihelper/apilogger"
	"github.com/stretchr/testify/require"
)

// Information api

func TestGetSeed(t *testing.T) {
	seed := insolarapi.GetSeed(t)
	require.NotEmpty(t, seed)
}

// Member api

func TestCreateMember(t *testing.T) {
	member := insolarapi.CreateMember(t)
	require.NotEmpty(t, member.MemberReference, "MemberReference")
}

func TestMemberTransfer(t *testing.T) {
	member1 := insolarapi.CreateMember(t)
	member2 := insolarapi.CreateMember(t)
	transfer := member1.Transfer(t, member2.MemberReference, "1")
	apihelper.CheckResponseHasNoError(t, transfer)
	apilogger.Println("Transfer OK. Fee: " + transfer.Result.CallResult.Fee)
	require.NotEmpty(t, transfer.Result.CallResult.Fee, "Fee")
}

func TestGetMember(t *testing.T) {
	member1 := insolarapi.CreateMember(t)
	resp := member1.GetMember(t)
	apihelper.CheckResponseHasNoError(t, resp)
	require.Equal(t, member1.MemberReference, resp.Result.CallResult.Reference, "Reference")
	require.Empty(t, resp.Result.CallResult.MigrationAddress, "MigrationAddress")
}

// Migration api

func TestMemberMigrationCreate(t *testing.T) {
	// TODO everything related to migration.addAddresses move to 'before' function
	absPath, _ := filepath.Abs("../resources/migration-addresses.json")
	b, err := ioutil.ReadFile(absPath)
	require.NoError(t, err)

	type addressesList struct {
		Addresses []string `json:"Addresses"`
	}
	list := addressesList{}
	json.Unmarshal(b, &list)

	response := internalapi.AddMigrationAddresses(t, list.Addresses)
	require.NotNil(t, response)
	member := insolarapi.MemberMigrationCreate(t)
	require.NotEmpty(t, member)
	require.NotEmpty(t, member.MemberResponseResult)
}

func TestDepositTransfer(t *testing.T) {
	member := insolarapi.MemberMigrationCreate(t)
	response := member.DepositTransfer(t)
	require.NotEmpty(t, response.Result.CallResult)
}
