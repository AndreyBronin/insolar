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

package thread

import (
	"sync"
)

type Controller struct {
	cancelMu sync.Mutex
	cancel   chan struct{}
}

func NewController() *Controller {
	return &Controller{cancel: make(chan struct{})}
}

func (c *Controller) Cancel() <-chan struct{} {
	c.cancelMu.Lock()
	defer c.cancelMu.Unlock()

	return c.cancel
}

func (c *Controller) Pulse() {
	c.cancelMu.Lock()
	defer c.cancelMu.Unlock()

	toClose := c.cancel
	c.cancel = make(chan struct{})
	close(toClose)
}
