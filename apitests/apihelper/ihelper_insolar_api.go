package apihelper

import (
	"github.com/insolar/insolar/apitests/apiclient/insolar_api"
	"github.com/insolar/insolar/apitests/apihelper/apilogger"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

var id int32 = 0

const (
	url            = "http://localhost:19102"
	JSONRPCVersion = "2.0"
	ApiCall        = "api.call"
	//ContractCall   = "contract.call"
	//information_api
	GetSeedMethod = "node.getSeed"
	GetInfoMethod = "network.getInfo"

	//member_api
	MemberCreateMethod   = "member.create"
	MemberTransferMethod = "member.transfer"
	MemberGetMethod      = "member.get"
	//migration_api
	MemberMigrationCreateMethod = "member.migrationCreate"
	DepositTransferMethod       = "deposit.transfer"
)

var informationApi = GetClient().InformationApi
var memberApi = GetClient().MemberApi
var migrationApi = GetClient().MigrationApi

func GetClient() *insolar_api.APIClient {
	c := insolar_api.Configuration{
		BasePath: url,
	}
	return insolar_api.NewAPIClient(&c)
}

func GetRequestId() int32 {
	id++
	return id
}

func GetSeed(t *testing.T) string {
	r := insolar_api.NodeGetSeedRequest{
		Jsonrpc: JSONRPCVersion,
		Id:      GetRequestId(),
		Method:  GetSeedMethod,
	}
	return GetSeedRequest(t, r)
}

func GetSeedRequest(t *testing.T, r insolar_api.NodeGetSeedRequest) string {
	apilogger.LogApiRequest(GetSeedMethod, r, nil)
	response, http, err := informationApi.GetSeed(nil, r)
	require.Nil(t, err)
	apilogger.LogApiResponse(http, response)
	CheckResponseHasNoError(t, response)
	return response.Result.Seed
}

func GetInfo(t *testing.T) insolar_api.NetworkGetInfoResponseResult {
	infoBody := insolar_api.NetworkGetInfoRequest{
		Jsonrpc: JSONRPCVersion,
		Id:      GetRequestId(),
		Method:  GetInfoMethod,
		Params:  nil,
	}
	apilogger.LogApiRequest(GetInfoMethod, infoBody, nil)
	response, http, err := informationApi.GetInfo(nil, infoBody)
	require.Nil(t, err)
	apilogger.LogApiResponse(http, response)
	CheckResponseHasNoError(t, response)

	return response.Result
}

func GetRootMember(t *testing.T) string {
	return GetInfo(t).RootMember
}

func CreateMember(t *testing.T) MemberObject {
	var err error
	ms, _ := NewMemberSignature()
	seed := GetSeed(t)

	request := insolar_api.MemberCreateRequest{
		Jsonrpc: JSONRPCVersion,
		Id:      GetRequestId(),
		Method:  ApiCall,
		Params: insolar_api.MemberCreateRequestParams{
			Seed:      seed,
			CallSite:  MemberCreateMethod,
			PublicKey: string(ms.PemPublicKey),
		},
	}
	d, s, m := sign(request, ms.PrivateKey)
	apilogger.LogApiRequest(MemberCreateMethod, request, m)
	response, http, err := memberApi.MemberCreate(nil, d, s, request)
	require.Nil(t, err)
	apilogger.LogApiResponse(http, response)
	CheckResponseHasNoError(t, response)
	apilogger.Println("Member created: " + response.Result.CallResult.Reference)
	return MemberObject{
		MemberReference:      response.Result.CallResult.Reference,
		Signature:            ms,
		MemberResponseResult: response,
	}
}

func (member *MemberObject) GetMember(t *testing.T) insolar_api.MemberGetResponse {
	seed := GetSeed(t)
	request := insolar_api.MemberGetRequest{
		Jsonrpc: JSONRPCVersion,
		Id:      GetRequestId(),
		Method:  ApiCall,
		Params: insolar_api.MemberGetRequestParams{
			Seed:       seed,
			CallSite:   MemberGetMethod,
			CallParams: nil,
			PublicKey:  string(member.Signature.PemPublicKey),
		},
	}
	d, s, m := sign(request, member.Signature.PrivateKey)
	apilogger.LogApiRequest(MemberGetMethod, request, m)
	response, http, err := memberApi.MemberGet(nil, d, s, request)
	require.Nil(t, err)
	apilogger.LogApiResponse(http, response)
	CheckResponseHasNoError(t, response)
	return response
}

func (member *MemberObject) Transfer(t *testing.T, toMemberRef string, amount string) insolar_api.MemberTransferResponse {
	seed := GetSeed(t)
	request := insolar_api.MemberTransferRequest{
		Jsonrpc: JSONRPCVersion,
		Id:      GetRequestId(),
		Method:  ApiCall,
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
	d, s, m := sign(request, member.Signature.PrivateKey)
	apilogger.LogApiRequest(MemberTransferMethod, request, m)
	response, http, err := memberApi.MemberTransfer(nil, d, s, request)
	require.Nil(t, err)
	apilogger.LogApiResponse(http, response)
	CheckResponseHasNoError(t, response)
	apilogger.Println("Transfer OK. Fee: " + response.Result.CallResult.Fee)
	return response
}

func MemberMigrationCreate(t *testing.T) MemberObject {
	var err error
	ms, _ := NewMemberSignature()
	seed := GetSeed(t)

	request := insolar_api.MemberMigrationCreateRequest{
		Jsonrpc: JSONRPCVersion,
		Id:      1,
		Method:  ApiCall,
		Params: insolar_api.MemberMigrationCreateRequestParams{
			Seed:       seed,
			CallSite:   MemberMigrationCreateMethod,
			CallParams: nil,
			PublicKey:  string(ms.PemPublicKey),
		},
	}
	d, s, m := sign(request, ms.PrivateKey)
	apilogger.LogApiRequest(MemberMigrationCreateMethod, request, m)
	response, http, err := migrationApi.MemberMigrationCreate(nil, d, s, request)
	apilogger.LogApiResponse(http, response)
	CheckResponseHasNoError(t, response)
	if err != nil {
		log.Fatalln(err)
	}

	return MemberObject{
		MemberReference: response.Result.CallResult.Reference,
		Signature:       ms,
		MemberResponseResult: insolar_api.MemberCreateResponse{
			Jsonrpc: response.Jsonrpc,
			Id:      response.Id,
			Result: insolar_api.MemberCreateResponseResult{
				CallResult: insolar_api.MemberCreateResponseResultCallResult{
					Reference: response.Result.CallResult.Reference,
				},
				RequestReference: response.Result.RequestReference,
				TraceID:          response.Result.TraceID,
			},
			Error: insolar_api.MemberCreateResponseError{
				Data: insolar_api.MemberCreateResponseErrorData{
					RequestReference: response.Error.Data.RequestReference,
					TraceID:          response.Error.Data.TraceID,
				},
				Code:    response.Error.Code,
				Message: response.Error.Message,
			},
		},
	}
}

func DepositTransfer(t *testing.T) insolar_api.DepositTransferResponse {
	var err error
	ms, _ := NewMemberSignature()

	request := insolar_api.DepositTransferRequest{
		Jsonrpc: JSONRPCVersion,
		Id:      1,
		Method:  ApiCall,
		Params: insolar_api.DepositTransferRequestParams{
			Seed:     GetSeed(t),
			CallSite: DepositTransferMethod,
			CallParams: insolar_api.DepositTransferRequestParamsCallParams{
				Amount:    "1000",
				EthTxHash: "",
			},
			PublicKey: string(ms.PemPublicKey),
		},
	}

	d, s, m := sign(request, ms.PrivateKey)
	apilogger.LogApiRequest(DepositTransferMethod, request, m)
	response, http, err := migrationApi.DepositTransfer(nil, d, s, request)
	apilogger.LogApiResponse(http, response)
	CheckResponseHasNoError(t, response)
	if err != nil {
		log.Fatalln(err)
	}
	return response
}
