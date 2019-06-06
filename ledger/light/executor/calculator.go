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

package executor

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
)

//go:generate minimock -i github.com/insolar/insolar/ledger/light/executor.Calculator -o ./ -s _mock.go

// Calculator provides methods for calculating jets
type Calculator interface {
	// MineForPulse returns current node's jets for a provided pulse
	MineForPulse(ctx context.Context, pn insolar.PulseNumber) []insolar.JetID
}

// CalculatorDefault is a struct, that implements jet.Calculator
type CalculatorDefault struct {
	coordinator jet.Coordinator
	jetAccessor jet.Accessor
}

// NewCalculator returns a new instance of a jet calculator.
func NewCalculator(jetCoordinator jet.Coordinator, jetAccessor jet.Accessor) *CalculatorDefault {
	return &CalculatorDefault{coordinator: jetCoordinator, jetAccessor: jetAccessor}
}

// MineForPulse returns current node's jets for a provided pulse
func (c *CalculatorDefault) MineForPulse(ctx context.Context, pn insolar.PulseNumber) []insolar.JetID {
	var res []insolar.JetID

	jetIDs := c.jetAccessor.All(ctx, pn)
	me := c.coordinator.Me()

	for _, jetID := range jetIDs {
		executor, err := c.coordinator.LightExecutorForJet(ctx, insolar.ID(jetID), pn)
		if err != nil {
			continue
		}
		if *executor == me {
			res = append(res, jetID)
		}
	}

	return res
}
