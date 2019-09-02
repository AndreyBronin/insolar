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
	"time"

	"github.com/insolar/insolar/apitests/apiclient/insolar_observer_api"
	"github.com/insolar/insolar/apitests/apihelper"
	"github.com/insolar/insolar/apitests/apihelper/apilogger"
	"github.com/insolar/insolar/apitests/tests/insolarapi"
	"github.com/insolar/insolar/apitests/tests/observerapi"
	"github.com/stretchr/testify/require"
)

func TestNotification(t *testing.T) {
	response := observerapi.Notification(t) //todo how this working?
	require.NotEmpty(t, response.Notification)
}

func TestObserverSmoke(t *testing.T) {
	//precondition
	amount := "1"
	member1 := insolarapi.CreateMember(t)
	member2 := insolarapi.CreateMember(t)
	transfer := member1.Transfer(t, member2.MemberReference, amount)
	apihelper.CheckResponseHasNoError(t, transfer)
	fee := transfer.Result.CallResult.Fee
	apilogger.Println("Transfer OK. Fee: " + fee)
	require.NotEmpty(t, fee, "Fee")

	apilogger.Println("Sleep 60 sek ")
	time.Sleep(60 * time.Second)

	//getMember
	getMember(t, member1)
	getMember(t, member2)

	//getBalance
	getBalance(t, member1, "9989999999")
	getBalance(t, member2, "10000000001")

	//getTransaction
	response := observerapi.Transaction(t, transfer.Result.RequestReference)

	require.Equal(t, amount, response.Amount)
	require.Equal(t, fee, response.Fee)
	require.Equal(t, member1.MemberReference, response.FromMemberReference)
	require.Equal(t, member2.MemberReference, response.ToMemberReference)
	require.Equal(t, "SUCCESS", response.Status)

	require.Empty(t, response.Error)

	//transactionsList
	transactions1 := getTransactionList(t, member1)
	assertTransactionsContains(transactions1, response, t)
	transactions2 := getTransactionList(t, member2)
	assertTransactionsContains(transactions2, response, t)
}

func getTransactionList(t *testing.T, member insolarapi.MemberObject) []insolar_observer_api.InlineResponse200 {
	transactions := observerapi.TransactionList(t, member.MemberReference+"?direction=all&limit=10")
	require.NotEmpty(t, transactions)
	return transactions
}

func assertTransactionsContains(transactions []insolar_observer_api.InlineResponse200, tx insolar_observer_api.TransactionResponse200, t *testing.T) {
	for _, v := range transactions {
		if v.Index == tx.Index {
			apilogger.Println("transaction=" + string(v.TxID))
			require.Equal(t, tx.Index, v.Index)
			require.Equal(t, tx.Status, v.Status)
			require.Equal(t, tx.ToMemberReference, v.ToMemberReference)
			require.Equal(t, tx.FromMemberReference, v.FromMemberReference)
			require.Equal(t, tx.Fee, v.Fee)
			require.Equal(t, tx.Amount, v.Amount)
		}
	}
}

func getBalance(t *testing.T, member insolarapi.MemberObject, expBalance string) {
	balance := observerapi.Balance(t, member.MemberReference)
	require.Empty(t, balance.Error)
	require.Equal(t, expBalance, balance.Balance)
}

func getMember(t *testing.T, member insolarapi.MemberObject) {
	get := observerapi.Member(t, member.MemberReference)
	require.Empty(t, get.Error)
	require.NotEmpty(t, get.Balance)
	//require.NotEmpty(t, get.MigrationAddress)
}
