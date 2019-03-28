/*
 *    Copyright 2019 Insolar Technologies
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

package artifactmanager

import (
	"context"

	"github.com/insolar/insolar/conveyor/adapter"
	"github.com/insolar/insolar/conveyor/adapter/adapterid"
	"github.com/insolar/insolar/conveyor/fsm"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/log"
	"github.com/pkg/errors"
)

// GetCodeTask is task for adapter for getting code
type GetCodeTask struct {
	// TODO: don't let adapter and component know about Parcel type, get every needed info in
	Parcel insolar.Parcel
}

// GetCodeResp is response for adapter for getting code
type GetCodeResp struct {
	Reply insolar.Reply
	Err   error
}

// GetCodeProcessor is worker for adapter for getting code
type GetCodeProcessor struct {
	Handlers HandlerStorage `inject:""`
}

// NewGetCodeProcessor returns new instance of processor which get code
func NewGetCodeProcessor() adapter.Processor {
	return &GetCodeProcessor{}
}

// Process implements Processor interface
func (p *GetCodeProcessor) Process(task adapter.AdapterTask, nestedEventHelper adapter.NestedEventHelper, cancelInfo adapter.CancelInfo) interface{} {
	payload, ok := task.TaskPayload.(GetCodeTask)
	var msg GetCodeResp
	if !ok {
		msg.Err = errors.Errorf("[ GetCodeProcessor.Process ] Incorrect payload type: %T", task.TaskPayload)
		return msg
	}

	ctx := context.Background()
	reply, err := p.Handlers.handleGetCode(ctx, payload.Parcel)
	msg = GetCodeResp{reply, err}
	log.Info("[ GetCodeProcessor.Process ] Process was dome successfully")

	return msg
}

// GetCodeHelper is helper for GetCodeProcessor
type GetCodeHelper struct{}

// GetCode makes correct message and send it to adapter
func (r *GetCodeHelper) GetCode(element fsm.SlotElementHelper, parcel insolar.Parcel, respHandlerID uint32) error {
	task := GetCodeTask{
		Parcel: parcel,
	}
	err := element.SendTask(adapterid.GetCode, task, respHandlerID)
	return errors.Wrap(err, "[ GetCodeHelper.SendResponse ] Can't SendTask")
}
