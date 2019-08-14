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

// +build functest

package functest

import (
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/testutils"
)

// TODO: https://insolar.atlassian.net/browse/WLT-768
func TestDepositTransferToken(t *testing.T) {
	member := fullMigration(t, "Eth_TxHash_test")

	firstBalance := getBalanceNoErr(t, member, member.ref)
	secondBalance := new(big.Int).Add(firstBalance, big.NewInt(100))

	var err error
	for i := 0; i <= 20; i++ {
		time.Sleep(time.Second)
		_, _, err = makeSignedRequest(member, "deposit.transfer", map[string]interface{}{"amount": "100", "ethTxHash": "Eth_TxHash_test"})
		if err != nil {
			require.Error(t, err)
			require.Contains(t, err.Error(), "hold period didn't end")
		} else {
			break
		}
	}
	require.NoError(t, err)

	var finalBalance *big.Int
	for i := 0; i <= 10; i++ {
		time.Sleep(time.Second)
		finalBalance = getBalanceNoErr(t, member, member.ref)

		if secondBalance.Cmp(finalBalance) == 0 {
			break
		}
	}

	require.Equal(t, secondBalance, finalBalance)
}

func TestDepositTransferBeforeUnhold(t *testing.T) {
	member := fullMigration(t, "Eth_TxHash_test")

	_, err := signedRequestWithEmptyRequestRef(t, member, "deposit.transfer", map[string]interface{}{"amount": "100", "ethTxHash": "Eth_TxHash_test"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "hold period didn't end")
}

func TestDepositTransferBiggerAmount(t *testing.T) {
	member := fullMigration(t, "Eth_TxHash_test")

	_, err := signedRequestWithEmptyRequestRef(t, member, "deposit.transfer", map[string]interface{}{"amount": "1001", "ethTxHash": "Eth_TxHash_test"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "not enough balance for transfer")
}

func TestDepositTransferAnotherTx(t *testing.T) {
	member := fullMigration(t, "Eth_TxHash_test")

	_, err := signedRequestWithEmptyRequestRef(t, member, "deposit.transfer", map[string]interface{}{"amount": "100", "ethTxHash": "Eth_TxHash_foo"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "can't find deposit")
}

func TestDepositTransferWrongValueAmount(t *testing.T) {
	member := fullMigration(t, "Eth_TxHash_test")

	_, err := signedRequestWithEmptyRequestRef(t, member, "deposit.transfer", map[string]interface{}{"amount": "foo", "ethTxHash": "Eth_TxHash_test"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "can't parse input amount")
}

func TestDepositTransferNotEnoughConfirms(t *testing.T) {
	migrationAddress := testutils.RandomString()
	member := createMigrationMemberForMA(t, migrationAddress)

	migrate(t, member.ref, "1000", "Eth_TxHash_test", migrationAddress, 2)
	migrate(t, member.ref, "1000", "Eth_TxHash_test", migrationAddress, 0)

	_, err := signedRequestWithEmptyRequestRef(t, member, "deposit.transfer", map[string]interface{}{"amount": "100", "ethTxHash": "Eth_TxHash_test"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "number of confirms is less then 3")
}
