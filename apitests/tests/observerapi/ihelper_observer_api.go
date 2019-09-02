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

package observerapi

import (
	"testing"

	"github.com/insolar/insolar/apitests/apihelper"

	"github.com/insolar/insolar/apitests/apiclient/insolar_observer_api"
	"github.com/insolar/insolar/apitests/apihelper/apilogger"
	"github.com/stretchr/testify/require"
)

const (
	//observerUrl    = "http://wallet-api.qa-wallet.k8s-dev.insolar.io"
	ObserverUrl    = "http://wallet-api.qa-gamma.k8s-dev.insolar.io"
	JSONRPCVersion = "2.0"
)

var observerObserverApi = GetObserverClient().ObserverApi
var observerInformationApi = GetObserverClient().InformationApi

func GetObserverClient() *insolar_observer_api.APIClient {
	c := insolar_observer_api.Configuration{
		BasePath: ObserverUrl,
	}
	return insolar_observer_api.NewAPIClient(&c)
}

func Notification(t *testing.T) insolar_observer_api.NotificationResponse200 {
	apilogger.LogApiRequest("get /api/notification", nil, nil)
	response, http, err := observerInformationApi.Notification(nil)
	apilogger.LogApiResponse(http, response)
	require.Nil(t, err)
	apihelper.CheckResponseHasNoError(t, response)
	return response
}

func Balance(t *testing.T, reference string) insolar_observer_api.BalanceResponse200 {
	apilogger.LogApiRequest("get /api/balance", nil, nil)
	apilogger.Println("reference = " + reference)
	response, http, err := observerObserverApi.Balance(nil, reference)
	apilogger.LogApiResponse(http, response)
	require.Nil(t, err)
	apihelper.CheckResponseHasNoError(t, response)
	return response
}

func Fee(t *testing.T, amount string) insolar_observer_api.FeeResponse200 {
	apilogger.LogApiRequest("get /api/fee", nil, nil)
	apilogger.Println("amount = " + amount)
	response, http, err := observerObserverApi.Fee(nil, amount)
	apilogger.LogApiResponse(http, response)
	require.Nil(t, err)
	apihelper.CheckResponseHasNoError(t, response)
	return response
}

func Member(t *testing.T, reference string) insolar_observer_api.MemberResponse200 {
	apilogger.LogApiRequest("get /api/member", nil, nil)
	apilogger.Println("reference = " + reference)
	response, http, err := observerObserverApi.Member(nil, reference)
	apilogger.LogApiResponse(http, response)
	require.Nil(t, err)
	apihelper.CheckResponseHasNoError(t, response)
	return response
}

func Transaction(t *testing.T, txId string) insolar_observer_api.TransactionResponse200 {
	apilogger.LogApiRequest("get /api/transaction", nil, nil)
	apilogger.Println("txId = " + txId)
	response, http, err := observerObserverApi.Transaction(nil, txId)
	apilogger.LogApiResponse(http, response)
	require.Nil(t, err)
	apihelper.CheckResponseHasNoError(t, response)
	return response
}
func TransactionList(t *testing.T, reference string) []insolar_observer_api.InlineResponse200 {
	apilogger.LogApiRequest("get /api/transactionList", nil, nil)
	apilogger.Println("reference = " + reference)
	response, http, err := observerObserverApi.TransactionList(nil, reference)
	apilogger.LogApiResponse(http, response)
	require.Nil(t, err)
	//apihelper.CheckResponseHasNoError(t, response) //todo not allowed, because array in response
	return response
}
