//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
//

package transport

import (
	"context"
	"io"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/network/hostnetwork/resolver"

	"github.com/pkg/errors"
)

// StreamProcessor interface provides callback method to process data stream
type StreamProcessor interface {
	ProcessStream(address string, reader io.Reader) error
}

// DatagramProcessor interface provides callback method to process received datagrams
type DatagramProcessor interface {
	ProcessDatagram(address string, buf []byte) error
}

// DatagramTransport interface provides methods to send and receive datagrams
type DatagramTransport interface {
	component.Starter
	component.Stopper

	SendDatagram(ctx context.Context, address string, data []byte) error
	SetDatagramProcessor(processor DatagramProcessor)
}

// StreamTransport interface provides methods to send and receive data streams
type StreamTransport interface {
	component.Starter
	component.Stopper

	SendBuffer(ctx context.Context, address string, data []byte) error
	SetStreamProcessor(processor StreamProcessor)
}

// NewDatagramTransport creates new UDP DatagramTransport
func NewDatagramTransport(cfg configuration.Transport) (DatagramTransport, string, error) {
	return newUDPTransport(cfg.Address, cfg.FixedPublicAddress)
}

// NewTransport creates new Network with particular configuration
func NewTransport(cfg configuration.Transport) (StreamTransport, string, error) {
	switch cfg.Protocol {
	case "TCP":
		return newTCPTransport(cfg.Address, cfg.FixedPublicAddress)
	default:
		return nil, "", errors.New("invalid transport configuration")
	}
}
