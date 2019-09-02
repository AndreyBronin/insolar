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
// +build apitests

package smoketests

import (
	"testing"

	"github.com/insolar/insolar/apitests/tests/insolarapi"

	"github.com/insolar/insolar/apitests/tests/internalapi"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestMigrationAddAddresses(t *testing.T) {
	uuids, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	response := internalapi.AddMigrationAddresses(t, []string{uuids.String()})
	require.NotEmpty(t, response.Result)
	require.Empty(t, response.Error)
}

func TestMigrationDeposit(t *testing.T) {
	response := internalapi.MigrationDeposit(t)
	require.NotEmpty(t, response.Result)
	require.Empty(t, response.Error)
}

func TestObserverGetToken(t *testing.T) {
	response := internalapi.ObserverToken(t) //not worked https://insolar.atlassian.net/browse/INS-3401
	require.NotEmpty(t, response)
}

func TestObserverAddressesCount(t *testing.T) { // https://insolar.atlassian.net/browse/INS-3401
	response := internalapi.ObserverAddressesCount(t)
	require.NotEmpty(t, response)
}

func TestObserverGetMigrationAddresses(t *testing.T) { //https://insolar.atlassian.net/browse/INS-3401
	response := internalapi.ObserverGetMigrationAddresses(t) //qa bug autogenerate https://insolar.atlassian.net/browse/INS-3399
	require.NotEmpty(t, response)
}

func TestMemberGetBalance(t *testing.T) {
	member := insolarapi.CreateMember(t)
	response := internalapi.GetBalance(t, member)
	//require.NotEmpty(t, response.Result.CallResult.Deposits)
	require.NotEmpty(t, response.Result.CallResult.Balance)
	require.NotEmpty(t, response.Result.RequestReference)
	require.NotEmpty(t, response.Result.TraceID)
}

func TestMigrationDeactivateDaemon(t *testing.T) {
	response := internalapi.MigrationDeactivateDaemon(t, "")
	require.NotEmpty(t, response.Result)
	require.Empty(t, response.Error)
}

func TestMigrationActivateDaemon(t *testing.T) {
	response := internalapi.MigrationActivateDaemon(t, "")
	require.NotEmpty(t, response.Result)
	require.Empty(t, response.Error)
}

func TestMigrationCheckDaemon(t *testing.T) {
	info := internalapi.GetInfo(t)
	response := internalapi.MigrationCheckDaemon(t, info.Result.MigrationAdminMember, info.Result.MigrationDaemonMembers[0])
	require.NotEmpty(t, response.Result)
	require.Empty(t, response.Error)
}

func TestGetStatus(t *testing.T) {
	response := internalapi.GetStatus(t)
	require.Equal(t, "CompleteNetworkState", response.NetworkState)
	require.NotEmpty(t, response.ActiveListSize)
	require.NotEmpty(t, response.Entropy)
	for _, v := range response.Nodes {
		require.Equal(t, true, v.IsWorking)
	}
	require.Equal(t, false, response.Origin.IsWorking) //bug https://insolar.atlassian.net/browse/INS-3213
	require.NotEmpty(t, response.PulseNumber)
	require.NotEmpty(t, response.Version) //bug https://insolar.atlassian.net/browse/INS-3404
}

func TestGetInfo(t *testing.T) {
	response := internalapi.GetInfo(t)
	require.NotEmpty(t, response)
	require.NotEmpty(t, response.Result.RootMember)
	require.NotEmpty(t, response.Result.RootDomain)
	require.NotEmpty(t, response.Result.NodeDomain)
	require.NotEmpty(t, response.Result.TraceID)
	require.NotEmpty(t, response.Result.MigrationAdminMember)
	require.NotEmpty(t, response.Result.MigrationDaemonMembers)
}

func TestGetSeedInternal(t *testing.T) {
	response := internalapi.GetSeedInternal(t)
	require.NotEmpty(t, response)
}
