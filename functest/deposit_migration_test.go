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
	"encoding/base64"
	"math/big"
	"testing"

	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/testutils"
)

func TestMigrationToken(t *testing.T) {
	err := activateDaemons(t)
	require.NoError(t, err)
	migrationAddress := testutils.RandomString()
	member := createMigrationMemberForMA(t, migrationAddress)

	deposit := migrate(t, member.ref, "1000", "Test_TxHash", migrationAddress, 0)
	firstMemberBalance := deposit["balance"].(string)

	require.Equal(t, "0", firstMemberBalance)
	firstMABalance := getBalanceNoErr(t, &migrationAdmin, migrationAdmin.ref)

	deposit = migrate(t, member.ref, "1000", "Test_TxHash", migrationAddress, 1)
	deposit = migrate(t, member.ref, "1000", "Test_TxHash", migrationAddress, 2)
	confirmerReferencesMap := deposit["confirmerReferences"].(string)
	sm := make(foundation.StableMap)
	decoded, err := base64.StdEncoding.DecodeString(confirmerReferencesMap)
	require.NoError(t, err)
	err = sm.UnmarshalBinary(decoded)
	require.NoError(t, err)

	for i := 0; i < 3; i++ {
		require.Equal(t, sm[migrationDaemons[i].ref], "1000")
	}
	require.Equal(t, deposit["ethTxHash"], "Test_TxHash")
	require.Equal(t, deposit["amount"], "1000")

	secondMemberBalance := deposit["balance"].(string)
	require.Equal(t, "1000", secondMemberBalance)
	secondMABalance := getBalanceNoErr(t, &migrationAdmin, migrationAdmin.ref)
	dif := new(big.Int).Sub(firstMABalance, secondMABalance)
	require.Equal(t, "1000", dif.String())
}

func TestMigrationTokenOnDifferentDeposits(t *testing.T) {
	migrationAddress := testutils.RandomString()
	member := createMigrationMemberForMA(t, migrationAddress)

	_ = migrate(t, member.ref, "1000", "Test_TxHash", migrationAddress, 0)
	secondDeposit := migrate(t, member.ref, "1000", "Test_TxHash", migrationAddress, 1)
	sm := make(foundation.StableMap)
	confirmerReferencesMap := secondDeposit["confirmerReferences"].(string)
	decoded, err := base64.StdEncoding.DecodeString(confirmerReferencesMap)
	err = sm.UnmarshalBinary(decoded)
	require.NoError(t, err)
	require.Equal(t, sm[migrationDaemons[0].ref], "1000")
	require.Equal(t, sm[migrationDaemons[1].ref], "1000")
}

func TestMigrationTokenNotInTheList(t *testing.T) {
	migrationAddress := generateMigrationAddress()
	_, err := signedRequestWithEmptyRequestRef(t, &migrationAdmin,
		"deposit.migration",
		map[string]interface{}{"amount": "1000", "ethTxHash": "TxHash", "migrationAddress": migrationAddress})
	require.Error(t, err)
	require.Contains(t, err.Error(), "this migration daemon is not in the list")
}

func TestMigrationTokenZeroAmount(t *testing.T) {
	migrationAddress := generateMigrationAddress()
	_ = createMigrationMemberForMA(t, migrationAddress)

	result, err := signedRequestWithEmptyRequestRef(t,
		&migrationDaemons[0],
		"deposit.migration",
		map[string]interface{}{"amount": "0", "ethTxHash": "TxHash", "migrationAddress": migrationAddress})

	require.Error(t, err)
	require.Contains(t, err.Error(), "amount must be greater than zero")
	require.Nil(t, result)

}

func TestMigrationTokenMistakeField(t *testing.T) {
	migrationAddress := generateMigrationAddress()
	_ = createMigrationMemberForMA(t, migrationAddress)

	result, err := signedRequestWithEmptyRequestRef(t,
		&migrationDaemons[0],
		"deposit.migration",
		map[string]interface{}{"amount1": "0", "ethTxHash": "TxHash", "migrationAddress": migrationAddress})
	require.Error(t, err)
	require.Contains(t, err.Error(), " incorect input: failed to get 'amount' param")
	require.Nil(t, result)
}

func TestMigrationTokenNilValue(t *testing.T) {
	migrationAddress := generateMigrationAddress()
	_ = createMigrationMemberForMA(t, migrationAddress)

	result, err := signedRequestWithEmptyRequestRef(t, &migrationDaemons[0], "deposit.migration", map[string]interface{}{"amount": "20", "ethTxHash": nil, "migrationAddress": migrationAddress})
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to get 'ethTxHash' param")
	require.Nil(t, result)

}

func TestMigrationTokenMaxAmount(t *testing.T) {
	migrationAddress := generateMigrationAddress()
	member := createMigrationMemberForMA(t, migrationAddress)

	result, err := signedRequest(t,
		&migrationDaemons[0],
		"deposit.migration",
		map[string]interface{}{"amount": "500000000000000000", "ethTxHash": "ethTxHash", "migrationAddress": migrationAddress})
	require.NoError(t, err)
	require.Equal(t, result.(map[string]interface{})["memberReference"].(string), member.ref)
}

func TestMigrationDoubleMigrationFromSameDaemon(t *testing.T) {
	migrationAddress := generateMigrationAddress()
	member := createMigrationMemberForMA(t, migrationAddress)

	resultMigr1, err := signedRequest(t,
		&migrationDaemons[0], "deposit.migration", map[string]interface{}{"amount": "20", "ethTxHash": "ethTxHash", "migrationAddress": migrationAddress})
	require.NoError(t, err)
	require.Equal(t, resultMigr1.(map[string]interface{})["memberReference"].(string), member.ref)

	_, err = signedRequestWithEmptyRequestRef(t,
		&migrationDaemons[0],
		"deposit.migration",
		map[string]interface{}{"amount": "20", "ethTxHash": "ethTxHash", "migrationAddress": migrationAddress})
	require.Error(t, err)
	require.Contains(t, err.Error(), "confirm from this migration daemon already exists")
}

func TestMigrationAnotherAmountSameTx(t *testing.T) {

	migrationAddress := generateMigrationAddress()
	_ = createMigrationMemberForMA(t, migrationAddress)

	_, err := signedRequest(t,
		&migrationDaemons[0], "deposit.migration", map[string]interface{}{"amount": "20", "ethTxHash": "ethTxHash", "migrationAddress": migrationAddress})
	require.NoError(t, err)

	_, err = signedRequest(t,
		&migrationDaemons[1],
		"deposit.migration",
		map[string]interface{}{"amount": "30", "ethTxHash": "ethTxHash", "migrationAddress": migrationAddress})
	require.NoError(t, err)

	_, _, err = makeSignedRequest(
		&migrationDaemons[2],
		"deposit.migration",
		map[string]interface{}{"amount": "30", "ethTxHash": "ethTxHash", "migrationAddress": migrationAddress})
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to check amount in confirmation from migration daemon")
}
