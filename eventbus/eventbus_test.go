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

package eventbus

import (
	"testing"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/eventbus/event"
	"github.com/insolar/insolar/network/servicenetwork"
	"github.com/stretchr/testify/assert"
)

type req struct {
	ref    core.RecordRef
	method string
	args   []byte
}

type runner struct {
	requests  []req
	responses []core.Reaction
}

func (r *runner) Start(components core.Components) error { return nil }
func (r *runner) Stop() error                            { return nil }

func (r *runner) Execute(e core.Event) (core.Reaction, error) {
	if len(r.responses) == 0 {
		panic("no request expected")
	}
	m := e.(*event.CallMethodEvent)
	r.requests = append(r.requests, req{e.GetReference(), m.Method, m.Arguments})
	resp := r.responses[0]
	r.responses = r.responses[1:]

	return resp, nil
}

func TestNew(t *testing.T) {
	t.Skip("need repair")
	r := new(runner)
	r.requests = make([]req, 0)
	r.responses = make([]core.Reaction, 0)
	cfg := configuration.NewConfiguration()
	network, err := servicenetwork.NewServiceNetwork(cfg.Host, cfg.Node)
	assert.NoError(t, err)
	eb, err := New(configuration.Configuration{})
	eb.Start(core.Components{
		"core.LogicRunner": r,
		"core.Network":     network,
	})
	if err != nil {
		t.Fatal(err)
	}
	if eb == nil {
		t.Fatal("no object created")
	}
}

// TODO: fix network interaction
// func TestRoute(t *testing.T) {
// 	r := new(runner)
// 	r.requests = make([]req, 0)
// 	r.responses = make([]core.Reaction, 0)
//
// 	dht, err := NewNode()
// 	assert.NoError(t, err)
// 	ctx := getDefaultCtx(dht)
//
// 	mr, _ := New(r, dht)
// 	reference := dht.GetOriginHost(ctx).ID.String()
//
// 	t.Run("success", func(t *testing.T) {
// 		r.responses = append(r.responses, core.Reaction{Data: []byte("data"), Result: []byte("result"), Error: nil})
// 		resp, err := mr.Route(
// 			ctx, core.Event{Reference: core.NewRefFromBase58(reference), Method: "SomeMethod", Arguments: []byte("args")},
// 		)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		if string(resp.Data) != "data" {
// 			t.Fatal("unexpected data")
// 		}
// 		if string(resp.Result) != "result" {
// 			t.Fatal("unexpected data")
// 		}
// 		if len(r.requests) != 1 {
// 			t.Fatal("unexpected number of requests registered")
// 		}
// 		req := r.requests[0]
// 		r.requests = r.requests[1:]
//
// 		if req.ref != reference {
// 			t.Fatal("unexpected data")
// 		}
// 		if req.method != "SomeMethod" {
// 			t.Fatal("unexpected data")
// 		}
// 		if string(req.args) != "args" {
// 			t.Fatal("unexpected data")
// 		}
// 	})
// 	t.Run("error", func(t *testing.T) {
// 		r.responses = append(r.responses, core.Reaction{Data: []byte{}, Result: []byte{}, Error: errors.New("wtf")})
// 		_, err := mr.Route(
// 			ctx, core.Event{Reference: core.NewRefFromBase58(reference), Method: "SomeMethod", Arguments: []byte("args")},
// 		)
// 		if err == nil {
// 			t.Fatal("error expected")
// 		}
//
// 		if len(r.requests) != 1 {
// 			t.Fatal("unexpected number of requests registered")
// 		}
// 		req := r.requests[0]
// 		r.requests = r.requests[1:]
//
// 		if req.ref != reference {
// 			t.Fatal("unexpected data")
// 		}
// 		if req.method != "SomeMethod" {
// 			t.Fatal("unexpected data")
// 		}
// 		if string(req.args) != "args" {
// 			t.Fatal("unexpected data")
// 		}
// 	})
//
// 	t.Run("referenceNotFound", func(t *testing.T) {
// 		_, err := mr.Route(
// 			ctx,
// 			core.Event{
// 				Reference: core.NewRefFromBase58("refNotFound"),
// 				Method:    "SomeMethod",
// 				Arguments: []byte("args"),
// 			},
// 		)
// 		assert.Error(t, err)
// 	})
// }
