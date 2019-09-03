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

package insolarapi

import (
	"log"
	"testing"

	"github.com/insolar/insolar/apitests/apiclient/insolar_api"
	"github.com/insolar/insolar/apitests/apihelper"
	"github.com/insolar/insolar/apitests/apihelper/apilogger"
	"github.com/stretchr/testify/require"
)

const (
	//url = "http://localhost:19103"
	url            = "http://wallet-api.qa-wallet.k8s-dev.insolar.io"
	JSONRPCVersion = "2.0"
	ContractCall   = "contract.call"
	// information_api
	GetSeedMethod = "node.getSeed"

	// member_api
	MemberCreateMethod   = "member.create"
	MemberTransferMethod = "member.transfer"
	MemberGetMethod      = "member.get"
	// migration_api
	MemberMigrationCreateMethod = "member.migrationCreate"
	DepositTransferMethod       = "deposit.transfer"
)

var informationApi = GetClient().InformationApi
var memberApi = GetClient().MemberApi
var migrationApi = GetClient().MigrationApi

type MemberObject apihelper.MemberObject

func GetClient() *insolar_api.APIClient {
	c := insolar_api.Configuration{
		BasePath: url,
	}
	return insolar_api.NewAPIClient(&c)
}

func GetSeed(t *testing.T) string {
	r := insolar_api.NodeGetSeedRequest{
		Jsonrpc: JSONRPCVersion,
		Id:      apihelper.GetRequestId(),
		Method:  GetSeedMethod,
	}
	return GetSeedRequest(t, r)
}

func GetSeedRequest(t *testing.T, r insolar_api.NodeGetSeedRequest) string {
	apilogger.LogApiRequest(r.Method, r, nil)
	response, http, err := informationApi.GetSeed(nil, r)
	apilogger.LogApiResponse(http, response)
	require.Nil(t, err)
	apihelper.CheckResponseHasNoError(t, response)
	return response.Result.Seed
}

func CreateMember(t *testing.T) MemberObject {
	var err error
	ms, _ := apihelper.NewMemberSignature()
	seed := GetSeed(t)

	request := insolar_api.MemberCreateRequest{
		Jsonrpc: JSONRPCVersion,
		Id:      apihelper.GetRequestId(),
		Method:  ContractCall,
		Params: insolar_api.MemberCreateRequestParams{
			Seed:      seed,
			CallSite:  MemberCreateMethod,
			PublicKey: string(ms.PemPublicKey),
		},
	}
	d, s, m := apihelper.Sign(request, ms.PrivateKey)
	apilogger.LogApiRequest(request.Params.CallSite, request, m)
	response, http, err := memberApi.MemberCreate(nil, d, s, request)
	require.Nil(t, err)
	apilogger.LogApiResponse(http, response)
	apihelper.CheckResponseHasNoError(t, response)
	apilogger.Println("Member created: " + response.Result.CallResult.Reference)
	return MemberObject{
		MemberReference:      response.Result.CallResult.Reference,
		Signature:            ms,
		MemberResponseResult: response,
	}
}

func (member *MemberObject) GetMember(t *testing.T) insolar_api.MemberGetResponse200 {
	seed := GetSeed(t)
	request := insolar_api.MemberGetRequest{
		Jsonrpc: JSONRPCVersion,
		Id:      apihelper.GetRequestId(),
		Method:  ContractCall,
		Params: insolar_api.MemberGetRequestParams{
			Seed:       seed,
			CallSite:   MemberGetMethod,
			CallParams: nil,
			PublicKey:  string(member.Signature.PemPublicKey),
		},
	}
	d, s, m := apihelper.Sign(request, member.Signature.PrivateKey)
	apilogger.LogApiRequest(request.Params.CallSite, request, m)
	response, http, err := memberApi.MemberGet(nil, d, s, request)
	apilogger.LogApiResponse(http, response)
	require.Nil(t, err)
	return response
}

func (member *MemberObject) Transfer(t *testing.T, toMemberRef string, amount string) insolar_api.MemberTransferResponse200 {
	seed := GetSeed(t)
	request := insolar_api.MemberTransferRequest{
		Jsonrpc: JSONRPCVersion,
		Id:      apihelper.GetRequestId(),
		Method:  ContractCall,
		Params: insolar_api.MemberTransferRequestParams{
			Seed:     seed,
			CallSite: MemberTransferMethod,
			CallParams: insolar_api.MemberTransferRequestParamsCallParams{
				Amount:            amount,
				ToMemberReference: toMemberRef,
			},
			PublicKey: string(member.Signature.PemPublicKey),
			Reference: member.MemberResponseResult.Result.CallResult.Reference,
		},
	}
	d, s, m := apihelper.Sign(request, member.Signature.PrivateKey)
	apilogger.LogApiRequest(request.Params.CallSite, request, m)
	response, http, err := memberApi.MemberTransfer(nil, d, s, request)
	require.Nil(t, err)
	apilogger.LogApiResponse(http, response)
	return response
}

func MemberMigrationCreate(t *testing.T) MemberObject {
	var err error
	ms, _ := apihelper.NewMemberSignature()
	seed := GetSeed(t)

	request := insolar_api.MemberMigrationCreateRequest{
		Jsonrpc: JSONRPCVersion,
		Id:      apihelper.GetRequestId(),
		Method:  ContractCall,
		Params: insolar_api.MemberMigrationCreateRequestParams{
			Seed:       seed,
			CallSite:   MemberMigrationCreateMethod,
			CallParams: nil,
			PublicKey:  string(ms.PemPublicKey),
		},
	}
	d, s, m := apihelper.Sign(request, ms.PrivateKey)
	apilogger.LogApiRequest(request.Params.CallSite, request, m)
	response, http, err := migrationApi.MemberMigrationCreate(nil, d, s, request)
	apilogger.LogApiResponse(http, response)
	apihelper.CheckResponseHasNoError(t, response)
	if err != nil {
		log.Fatalln(err)
	}

	return MemberObject{
		MemberReference: response.Result.CallResult.Reference,
		Signature:       ms,
		MemberResponseResult: insolar_api.MemberCreateResponse200{
			Jsonrpc: response.Jsonrpc,
			Id:      response.Id,
			Result: insolar_api.MemberCreateResponse200Result{
				CallResult: insolar_api.MemberCreateResponse200ResultCallResult{
					Reference: response.Result.CallResult.Reference,
				},
				RequestReference: response.Result.RequestReference,
				TraceID:          response.Result.TraceID,
			},
			Error: insolar_api.MemberCreateResponse200Error{
				Data: insolar_api.MemberCreateResponse200ErrorData{
					RequestReference: response.Error.Data.RequestReference,
					TraceID:          response.Error.Data.TraceID,
				},
				Code:    response.Error.Code,
				Message: response.Error.Message,
			},
		},
	}
}

func (member *MemberObject) DepositTransfer(t *testing.T) insolar_api.DepositTransferResponse200 {
	var err error
	ms, _ := apihelper.NewMemberSignature()

	request := insolar_api.DepositTransferRequest{
		Jsonrpc: JSONRPCVersion,
		Id:      apihelper.GetRequestId(),
		Method:  ContractCall,
		Params: insolar_api.DepositTransferRequestParams{
			Seed:     GetSeed(t),
			CallSite: DepositTransferMethod,
			CallParams: insolar_api.DepositTransferRequestParamsCallParams{
				Amount:    "1000",
				EthTxHash: "",
			},
			PublicKey: string(ms.PemPublicKey),
			Reference: member.MemberReference,
		},
	}

	d, s, m := apihelper.Sign(request, ms.PrivateKey)
	apilogger.LogApiRequest(request.Params.CallSite, request, m)
	response, http, err := migrationApi.DepositTransfer(nil, d, s, request)
	apilogger.LogApiResponse(http, response)
	apihelper.CheckResponseHasNoError(t, response)
	if err != nil {
		log.Fatalln(err)
	}
	return response
}
