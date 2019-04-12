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

package pool

import (
	"context"
	"io"
	"sync"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/metrics"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/transport"
)

type connectionPool struct {
	transport transport.StreamTransport

	entryHolder entryHolder
	mutex       sync.RWMutex
}

func newConnectionPool(t transport.StreamTransport) *connectionPool {
	return &connectionPool{
		transport:   t,
		entryHolder: newEntryHolder(),
	}
}

func (cp *connectionPool) GetConnection(ctx context.Context, host *host.Host) (io.ReadWriteCloser, error) {
	logger := inslogger.FromContext(ctx)

	entry, ok := cp.getEntry(host)

	logger.Debugf("[ GetConnection ] Finding entry for connection to %s in pool: %t", host, ok)

	if ok {
		return entry.Open(ctx)
	}

	logger.Debugf("[ GetConnection ] Missing entry for connection to %s in pool ", host)
	entry = cp.getOrCreateEntry(ctx, host)

	return entry.Open(ctx)
}

func (cp *connectionPool) CloseConnection(ctx context.Context, host *host.Host) {
	cp.mutex.Lock()
	defer cp.mutex.Unlock()

	logger := inslogger.FromContext(ctx)

	entry, ok := cp.entryHolder.Get(host)
	logger.Debugf("[ CloseConnection ] Finding entry for connection to %s in pool: %t", host, ok)

	if ok {
		entry.Close()

		logger.Debugf("[ CloseConnection ] Delete entry for connection to %s from pool", host)
		cp.entryHolder.Delete(host)
		metrics.NetworkConnections.Dec()
	}
}

func (cp *connectionPool) HandleConnection(host *host.Host, conn io.ReadWriteCloser) error {
	panic("implement me")
}

func (cp *connectionPool) getEntry(host *host.Host) (entry, bool) {
	cp.mutex.RLock()
	defer cp.mutex.RUnlock()

	return cp.entryHolder.Get(host)
}

func (cp *connectionPool) getOrCreateEntry(ctx context.Context, host *host.Host) entry {
	logger := inslogger.FromContext(ctx)

	cp.mutex.Lock()
	defer cp.mutex.Unlock()

	entry, ok := cp.entryHolder.Get(host)
	logger.Debugf("[ getOrCreateEntry ] Finding entry for connection to %s in pool: %t", host, ok)

	if ok {
		return entry
	}

	logger.Debugf("[ getOrCreateEntry ] Failed to retrieve entry for connection to %s, creating it", host)

	entry = newEntry(cp.transport, host, cp.CloseConnection)

	cp.entryHolder.Add(host, entry)
	size := cp.entryHolder.Size()
	logger.Debugf(
		"[ getOrCreateEntry ] Added entry for connection to %s. Current pool size: %d",
		host,
		size,
	)
	metrics.NetworkConnections.Inc()

	return entry
}

func (cp *connectionPool) Reset() {
	cp.mutex.Lock()
	defer cp.mutex.Unlock()

	cp.entryHolder.Iterate(func(entry entry) {
		entry.Close()
	})
	cp.entryHolder.Clear()
	metrics.NetworkConnections.Set(float64(cp.entryHolder.Size()))
}
