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

package artifactmanager

import (
	"context"
	"fmt"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow/bus"
)

type Dependencies struct {
	FetchJet   func(*FetchJet) *FetchJet
	WaitHot    func(*WaitHot) *WaitHot
	GetIndex   func(*GetIndex) *GetIndex
	SendObject func(p *SendObject) *SendObject
}

type ReturnReply struct {
	ReplyTo chan<- bus.Reply
	Err     error
	Reply   insolar.Reply
	Pub     message.Publisher
}

func (p *ReturnReply) Proceed(context.Context) error {
	fmt.Println("lol here love, Return reply", p.Reply, p.Err)
	p.ReplyTo <- bus.Reply{Reply: p.Reply, Err: p.Err}
	// if p.Pub != nil {
	// 	msg
	// 	p.Pub.Publish("outbound", msg)
	// }
	return nil
}