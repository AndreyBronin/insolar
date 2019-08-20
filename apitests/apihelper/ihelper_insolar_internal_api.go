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

package apihelper

import (
	"testing"

	"github.com/insolar/insolar/apitests/apihelper/apilogger"
	"github.com/insolar/insolar/apitests/scripts/insolar_internal_api"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

const (

	//information_api
	GetStatusMethod = "node.getStatus"
	//migration_api
	MigrationAddAddresses = "migration.addAddresses"
	MigrationGetInfo      = "migration.getInfo"
	DepositMigration      = "deposit.migration"
	DeactivateDaemon      = "migration.deactivateDaemon"
	ActivateDaemon        = "migration.activateDaemon"

	//member api
	MemberGetBalance = "member.getBalance"
)

var internalMemberApi = GetInternalClient().MemberApi
var internalMigrationApi = GetInternalClient().MigrationApi
var internalObserverApi = GetInternalClient().ObserverApi
var internalInformationApi = GetInternalClient().InformationApi

func GetInternalClient() *insolar_internal_api.APIClient {
	c := insolar_internal_api.Configuration{
		BasePath: url,
	}
	return insolar_internal_api.NewAPIClient(&c)
}

func GetStatus(t *testing.T) insolar_internal_api.NodeGetStatusResponse200Result {
	body := insolar_internal_api.NodeGetStatusRequest{
		Jsonrpc: JSONRPCVersion,
		Id:      GetRequestId(),
		Method:  GetStatusMethod,
		Params:  nil,
	}
	apilogger.LogApiRequest(GetStatusMethod, body, nil)
	response, http, err := internalInformationApi.GetStatus(nil, body)
	require.Nil(t, err)
	apilogger.LogApiResponse(http, response)
	CheckResponseHasNoError(t, response)

	return response.Result
}

func AddMigrationAddresses(t *testing.T) insolar_internal_api.MigrationDeactivateDaemonResponse200 {
	ms, _ := NewMemberSignature()
	uuids, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	adminPub, _ := LoadAdminMemberKeys() //todo getinfo

	body := insolar_internal_api.MigrationAddAddressesRequest{
		Jsonrpc: JSONRPCVersion,
		Id:      GetRequestId(),
		Method:  ContractCall,
		Params: insolar_internal_api.MigrationAddAddressesRequestParams{
			Seed:     GetSeed(t),
			CallSite: MigrationAddAddresses,
			CallParams: insolar_internal_api.MigrationAddAddressesRequestParamsCallParams{
				MigrationAddresses: []string{uuids.String()},
			},
			PublicKey: adminPub,
			Reference: "",
		},
	}
	d, s, m := Sign(body, ms.PrivateKey)
	apilogger.LogApiRequest(MigrationAddAddresses, body, m)
	response, http, err := internalMigrationApi.AddMigrationAddresses(nil, d, s, body)
	require.Nil(t, err)
	apilogger.LogApiResponse(http, response)
	CheckResponseHasNoError(t, response)
	apilogger.Printf("response id: %d", response.Id)
	return response
}

func GetMigrationInfo(t *testing.T) insolar_internal_api.MigrationGetInfoResponse200 {
	body := insolar_internal_api.MigrationGetInfoRequest{
		Jsonrpc: JSONRPCVersion,
		Id:      GetRequestId(),
		Method:  MigrationGetInfo,
		Params:  nil,
	}
	apilogger.LogApiRequest(MigrationGetInfo, body, nil)
	response, http, err := internalMigrationApi.GetInfo(nil, body)
	require.Nil(t, err)
	apilogger.LogApiResponse(http, response)
	require.NotEmpty(t, response.Result.MigrationAdminMember)
	return response
}

func MigrationDeposit(t *testing.T) insolar_internal_api.DepositMigrationResponse200 {
	body := insolar_internal_api.DepositMigrationRequest{
		Jsonrpc: JSONRPCVersion,
		Id:      GetRequestId(),
		Method:  ContractCall,
		Params: insolar_internal_api.DepositMigrationRequestParams{
			Seed:     GetSeed(t),
			CallSite: DepositMigration,
			CallParams: insolar_internal_api.DepositMigrationRequestParamsCallParams{
				Amount:           "1000",
				EthTxHash:        "Eth_TxHash_test",
				MigrationAddress: "", //todo getinfo
			},
			PublicKey: "", //migrationDaemonMember
			Reference: "", //migrationDaemonMember
		},
	}
	apilogger.LogApiRequest(MigrationGetInfo, body, nil)
	response, http, err := internalMigrationApi.DepositMigration(nil, "", "", body) //migrationDaemonMember
	require.Nil(t, err)
	apilogger.LogApiResponse(http, response)
	require.NotEmpty(t, response.Result.CallResult.MemberReference)
	return response
}

func ObserverToken(t *testing.T) map[string]interface{} {
	response, http, err := internalObserverApi.TokenGetInfo(nil)
	require.Nil(t, err)
	apilogger.LogApiResponse(http, response)
	return response
}

func GetBalance(t *testing.T, member MemberObject) insolar_internal_api.MemberGetBalanceResponse200 {
	body := insolar_internal_api.MemberGetBalanceRequest{
		Jsonrpc: JSONRPCVersion,
		Id:      GetRequestId(),
		Method:  ContractCall,
		Params: insolar_internal_api.MemberGetBalanceRequestParams{
			Seed:     GetSeed(t),
			CallSite: MemberGetBalance,
			CallParams: insolar_internal_api.MemberGetBalanceRequestParamsCallParams{
				Reference: member.MemberReference,
			},
			PublicKey: string(member.Signature.PemPublicKey),
			Reference: member.MemberReference,
		},
	}
	d, s, m := Sign(body, member.Signature.PrivateKey)
	apilogger.LogApiRequest(MigrationGetInfo, body, m)
	response, http, err := internalMemberApi.GetBalance(nil, d, s, body)
	require.Nil(t, err)
	apilogger.LogApiResponse(http, response)
	require.NotEmpty(t, response.Result.CallResult.Balance)
	return response
}

func MigrationDeactivateDaemon(t *testing.T, migrationDaemonReference string) insolar_internal_api.MigrationDeactivateDaemonResponse200 {

	body := insolar_internal_api.MigrationDeactivateDaemonRequest{
		Jsonrpc: JSONRPCVersion,
		Id:      GetRequestId(),
		Method:  ContractCall,
		Params: insolar_internal_api.MigrationDeactivateDaemonRequestParams{
			Seed:     GetSeed(t),
			CallSite: DeactivateDaemon,
			CallParams: insolar_internal_api.MigrationDeactivateDaemonRequestParamsCallParams{
				Reference: migrationDaemonReference, //migrationdaemon
			},
			PublicKey: "", //admin
			Reference: "", //admin
		},
	}
	//d, s, m := Sign(body, admin.PrivateKey)
	apilogger.LogApiRequest(MigrationAddAddresses, body, nil)
	response, http, err := internalMigrationApi.MigrationDeactivateDaemon(nil, "", "", body)
	require.Nil(t, err)
	CheckResponseHasNoError(t, response)
	apilogger.LogApiResponse(http, response)
	apilogger.Printf("response id: %d", response.Id)
	return response
}

func MigrationActivateDaemon(t *testing.T, migrationDaemonReference string) insolar_internal_api.MigrationDeactivateDaemonResponse200 {

	body := insolar_internal_api.MigrationActivateDaemonRequest{
		Jsonrpc: JSONRPCVersion,
		Id:      GetRequestId(),
		Method:  ContractCall,
		Params: insolar_internal_api.MigrationActivateDaemonRequestParams{
			Seed:     GetSeed(t),
			CallSite: ActivateDaemon,
			CallParams: insolar_internal_api.MigrationActivateDaemonRequestParamsCallParams{
				Reference: migrationDaemonReference, //migrationdaemon
			},
			PublicKey: "", //admin
			Reference: "", //admin
		},
	}
	//d, s, m := Sign(body, admin.PrivateKey)
	apilogger.LogApiRequest(MigrationAddAddresses, body, nil)
	response, http, err := internalMigrationApi.MigrationChangeDaemon(nil, "", "", body)
	require.Nil(t, err)
	CheckResponseHasNoError(t, response)
	apilogger.LogApiResponse(http, response)
	apilogger.Printf("response id: %d", response.Id)
	return response
}

func getMigrationAdmin(t *testing.T) string {
	return GetMigrationInfo(t).Result.MigrationAdminMember
}
