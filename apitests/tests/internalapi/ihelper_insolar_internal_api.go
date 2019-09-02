///
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
///
// +build apitests

package internalapi

import (
	"testing"

	"github.com/insolar/insolar/apitests/apihelper"
	"github.com/insolar/insolar/apitests/tests/insolarapi"
	"github.com/insolar/insolar/apitests/tests/observerapi"

	"github.com/insolar/insolar/apitests/apiclient/insolar_internal_api"
	"github.com/insolar/insolar/apitests/apihelper/apilogger"
	"github.com/stretchr/testify/require"
)

const (
	JSONRPCVersion = "2.0"
	url            = "http://localhost:19003"
	ContractCall   = "contract.call"
	// information_api
	GetStatusMethod = "node.getStatus"
	// migration_api
	AddAddresses     = "migration.addAddresses"
	CheckDaemon      = "migration.checkDaemon"
	DepositMigration = "deposit.migration"
	DeactivateDaemon = "migration.deactivateDaemon"
	ActivateDaemon   = "migration.activateDaemon"
	NetworkGetInfo   = "network.getInfo"

	// member api
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
		Id:      apihelper.GetRequestId(),
		Method:  GetStatusMethod,
		Params:  nil,
	}
	apilogger.LogApiRequest(body.Method, body, nil)
	response, http, err := internalInformationApi.GetStatus(nil, body)
	apilogger.LogApiResponse(http, response)
	require.Nil(t, err)
	apihelper.CheckResponseHasNoError(t, response)

	return response.Result
}

func GetSeedInternal(t *testing.T) string {
	body := insolar_internal_api.NodeGetSeedRequest{
		Jsonrpc: JSONRPCVersion,
		Id:      apihelper.GetRequestId(),
		Method:  insolarapi.GetSeedMethod,
	}
	apilogger.LogApiRequest(body.Method, body, nil)
	response, http, err := internalInformationApi.GetSeed(nil, body)
	apilogger.LogApiResponse(http, response)
	require.Nil(t, err)
	apihelper.CheckResponseHasNoError(t, response)

	return response.Result.Seed
}

func GetInfo(t *testing.T) insolar_internal_api.NetworkGetInfoResponse200 {
	body := insolar_internal_api.NetworkGetInfoRequest{
		Jsonrpc: JSONRPCVersion,
		Id:      apihelper.GetRequestId(),
		Method:  NetworkGetInfo,
		Params:  nil,
	}
	apilogger.LogApiRequest(body.Method, body, nil)
	response, http, _ := internalInformationApi.GetInfo(nil, body)
	apilogger.LogApiResponse(http, response)
	apihelper.CheckResponseHasNoError(t, response)

	return response
}

func AddMigrationAddresses(t *testing.T, addresses []string) insolar_internal_api.MigrationDeactivateDaemonResponse200 {
	privateKey, publicKey := apihelper.LoadAdminMemberKeys()
	signature, err := apihelper.CreateMemberSignature(publicKey, privateKey)
	require.Nil(t, err)
	body := insolar_internal_api.MigrationAddAddressesRequest{
		Jsonrpc: JSONRPCVersion,
		Id:      apihelper.GetRequestId(),
		Method:  ContractCall,
		Params: insolar_internal_api.MigrationAddAddressesRequestParams{
			Seed:     GetSeedInternal(t),
			CallSite: AddAddresses,
			CallParams: insolar_internal_api.MigrationAddAddressesRequestParamsCallParams{
				MigrationAddresses: addresses,
			},
			PublicKey: publicKey,
			Reference: GetInfo(t).Result.MigrationAdminMember,
		},
	}
	d, s, m := apihelper.Sign(body, signature.PrivateKey)
	apilogger.LogApiRequest(body.Params.CallSite, body, m)
	response, http, err := internalMigrationApi.AddMigrationAddresses(nil, d, s, body)
	apilogger.LogApiResponse(http, response)
	require.Nil(t, err)
	apihelper.CheckResponseHasNoError(t, response)
	apilogger.Printf("response id: %d", response.Id)
	return response
}

func MigrationCheckDaemon(t *testing.T, adminReference string, daemonRef string) insolar_internal_api.MigrationCheckDaemonResponse200 {
	privateKey, publicKey := apihelper.LoadAdminMemberKeys()
	signature, err := apihelper.CreateMemberSignature(publicKey, privateKey)
	require.Nil(t, err)
	body := insolar_internal_api.MigrationCheckDaemonRequest{
		Jsonrpc: JSONRPCVersion,
		Id:      apihelper.GetRequestId(),
		Method:  ContractCall,
		Params: insolar_internal_api.MigrationCheckDaemonRequestParams{
			Seed:     GetSeedInternal(t),
			CallSite: CheckDaemon,
			CallParams: insolar_internal_api.MigrationCheckDaemonRequestParamsCallParams{
				Reference: daemonRef,
			},
			PublicKey: publicKey,
			Reference: adminReference,
		},
	}
	d, s, m := apihelper.Sign(body, signature.PrivateKey)
	apilogger.LogApiRequest(body.Params.CallSite, body, m)
	response, http, err := internalMigrationApi.MigrationCheckDaemon(nil, d, s, body)
	apilogger.LogApiResponse(http, response)
	require.Nil(t, err)
	apihelper.CheckResponseHasNoError(t, response)
	apilogger.Printf("response id: %d", response.Id)
	return response
}

func MigrationDeposit(t *testing.T) insolar_internal_api.DepositMigrationResponse200 {
	body := insolar_internal_api.DepositMigrationRequest{
		Jsonrpc: JSONRPCVersion,
		Id:      apihelper.GetRequestId(),
		Method:  ContractCall,
		Params: insolar_internal_api.DepositMigrationRequestParams{
			Seed:     GetSeedInternal(t),
			CallSite: DepositMigration,
			CallParams: insolar_internal_api.DepositMigrationRequestParamsCallParams{
				Amount:           "1000",
				EthTxHash:        "Eth_TxHash_test",
				MigrationAddress: "",
			},
			PublicKey: "", //migrationDaemonMember
			Reference: "", //migrationDaemonMember
		},
	}
	apilogger.LogApiRequest(body.Params.CallSite, body, nil)
	response, http, err := internalMigrationApi.DepositMigration(nil, "", "", body) //migrationDaemonMember
	apilogger.LogApiResponse(http, response)
	require.Nil(t, err)
	require.NotEmpty(t, response.Result.CallResult.MemberReference)
	return response
}

func ObserverToken(t *testing.T) insolar_internal_api.TokenResponse200 {
	apilogger.Println(observerapi.ObserverUrl + "/api/token")
	response, http, err := internalObserverApi.TokenGetInfo(nil)
	apilogger.LogApiResponse(http, response)
	require.Nil(t, err)
	return response
}

func ObserverAddressesCount(t *testing.T) insolar_internal_api.AddressesCountResponse200 {
	apilogger.Println(observerapi.ObserverUrl + "/admin/migration/addresses/count")
	response, http, err := internalObserverApi.GetMigrationAddressCount(nil)
	apilogger.LogApiResponse(http, response)
	require.Nil(t, err)
	return response
}

func ObserverGetMigrationAddresses(t *testing.T) []string { //todo query params not generate - bug
	apilogger.Println(observerapi.ObserverUrl + "/admin/migration/addresses")
	response, http, err := internalObserverApi.GetMigrationAddresses(nil)
	apilogger.LogApiResponse(http, response)
	require.Nil(t, err)
	return response
}

func GetBalance(t *testing.T, member insolarapi.MemberObject) insolar_internal_api.MemberGetBalanceResponse200 {
	body := insolar_internal_api.MemberGetBalanceRequest{
		Jsonrpc: JSONRPCVersion,
		Id:      apihelper.GetRequestId(),
		Method:  ContractCall,
		Params: insolar_internal_api.MemberGetBalanceRequestParams{
			Seed:     GetSeedInternal(t),
			CallSite: MemberGetBalance,
			CallParams: insolar_internal_api.MemberGetBalanceRequestParamsCallParams{
				Reference: member.MemberReference,
			},
			PublicKey: string(member.Signature.PemPublicKey),
			Reference: member.MemberReference,
		},
	}
	d, s, m := apihelper.Sign(body, member.Signature.PrivateKey)
	apilogger.LogApiRequest(body.Params.CallSite, body, m)
	response, http, err := internalMemberApi.GetBalance(nil, d, s, body)
	apilogger.LogApiResponse(http, response)
	apilogger.Println(string(err.Error()))
	require.Nil(t, err.Error())

	require.NotEmpty(t, response.Result.CallResult.Balance)
	return response
}

func MigrationDeactivateDaemon(t *testing.T, migrationDaemonReference string) insolar_internal_api.MigrationDeactivateDaemonResponse200 {

	body := insolar_internal_api.MigrationDeactivateDaemonRequest{
		Jsonrpc: JSONRPCVersion,
		Id:      apihelper.GetRequestId(),
		Method:  ContractCall,
		Params: insolar_internal_api.MigrationDeactivateDaemonRequestParams{
			Seed:     GetSeedInternal(t),
			CallSite: DeactivateDaemon,
			CallParams: insolar_internal_api.MigrationDeactivateDaemonRequestParamsCallParams{
				Reference: migrationDaemonReference, // migrationdaemon
			},
			PublicKey: "", // admin
			Reference: "", // admin
		},
	}
	// d, s, m := Sign(body, admin.PrivateKey)
	apilogger.LogApiRequest(body.Params.CallSite, body, nil)
	response, http, err := internalMigrationApi.MigrationDeactivateDaemon(nil, "", "", body)
	apilogger.LogApiResponse(http, response)
	require.Nil(t, err)
	apihelper.CheckResponseHasNoError(t, response)
	apilogger.Printf("response id: %d", response.Id)
	return response
}

func MigrationActivateDaemon(t *testing.T, migrationDaemonReference string) insolar_internal_api.MigrationDeactivateDaemonResponse200 {

	body := insolar_internal_api.MigrationActivateDaemonRequest{
		Jsonrpc: JSONRPCVersion,
		Id:      apihelper.GetRequestId(),
		Method:  ContractCall,
		Params: insolar_internal_api.MigrationActivateDaemonRequestParams{
			Seed:     GetSeedInternal(t),
			CallSite: ActivateDaemon,
			CallParams: insolar_internal_api.MigrationActivateDaemonRequestParamsCallParams{
				Reference: migrationDaemonReference, // migrationdaemon
			},
			PublicKey: "", // admin
			Reference: "", // admin
		},
	}
	// d, s, m := Sign(body, admin.PrivateKey)
	apilogger.LogApiRequest(body.Params.CallSite, body, nil)
	response, http, err := internalMigrationApi.MigrationChangeDaemon(nil, "", "", body)
	apilogger.LogApiResponse(http, response)
	require.Nil(t, err)
	apihelper.CheckResponseHasNoError(t, response)
	apilogger.Printf("response id: %d", response.Id)
	return response
}
