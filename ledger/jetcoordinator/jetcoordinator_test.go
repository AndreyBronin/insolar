/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package jetcoordinator_test

import (
	"bytes"
	"context"
	"sort"
	"testing"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/ledgertestutils"
	"github.com/insolar/insolar/logicrunner"
	"github.com/insolar/insolar/network/nodekeeper"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/assert"
)

func TestJetCoordinator_QueryRole(t *testing.T) {
	lr, err := logicrunner.NewLogicRunner(&configuration.LogicRunner{
		BuiltIn: &configuration.BuiltIn{},
	})
	assert.NoError(t, err)
	keeper := nodekeeper.NewNodeKeeper(testutils.TestNode(core.RecordRef{}))
	c := core.Components{LogicRunner: lr, NodeNetwork: keeper}
	ledger, cleaner := ledgertestutils.TmpLedger(t, "", c)
	defer cleaner()

	am := ledger.GetArtifactManager()
	pm := ledger.GetPulseManager()
	jc := ledger.GetJetCoordinator()

	ctx, _ := inslogger.WithField(context.Background(), "testname", t.Name())
	pulse, err := pm.Current(ctx)
	assert.NoError(t, err)

	ref := func(r string) core.RecordRef { return core.NewRefFromBase58(r) }

	keeper.AddActiveNodes([]*core.Node{
		{NodeID: ref("53jNWvey7Nzyh4ZaLdJDf3SRgoD4GpWuwHgrgvVVGLbDkk3A7cwStSmBU2X7s4fm6cZtemEyJbce9dM9SwNxbsxf"), Roles: []core.NodeRole{core.RoleVirtual}},
		{NodeID: ref("4gU79K6woTZDvn4YUFHauNKfcHW69X42uyk8ZvRevCiMv3PLS24eM1vcA9mhKPv8b2jWj9J5RgGN9CB7PUzCtBsj"), Roles: []core.NodeRole{core.RoleLightMaterial}},
	})

	sorted := func(list []core.RecordRef) []core.RecordRef {
		sort.Slice(list, func(i, j int) bool {
			return bytes.Compare(list[i][:], list[j][:]) < 0
		})
		return list
	}

	selected, err := jc.QueryRole(ctx, core.RoleVirtualExecutor, *am.GenesisRef(), pulse.PulseNumber)
	assert.NoError(t, err)
	assert.Equal(t, []core.RecordRef{ref("53jNWvey7Nzyh4ZaLdJDf3SRgoD4GpWuwHgrgvVVGLbDkk3A7cwStSmBU2X7s4fm6cZtemEyJbce9dM9SwNxbsxf")}, selected)

	selected, err = jc.QueryRole(ctx, core.RoleLightValidator, *am.GenesisRef(), pulse.PulseNumber)
	assert.NoError(t, err)
	assert.Equal(t, sorted([]core.RecordRef{ref("4gU79K6woTZDvn4YUFHauNKfcHW69X42uyk8ZvRevCiMv3PLS24eM1vcA9mhKPv8b2jWj9J5RgGN9CB7PUzCtBsj")}), sorted(selected))

	selected, err = jc.QueryRole(ctx, core.RoleHeavyExecutor, *am.GenesisRef(), pulse.PulseNumber)
	assert.Error(t, err)
}
