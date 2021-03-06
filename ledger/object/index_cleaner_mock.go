package object

// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

import (
	"context"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar"
)

// IndexCleanerMock implements IndexCleaner
type IndexCleanerMock struct {
	t minimock.Tester

	funcDeleteForPN          func(ctx context.Context, pn insolar.PulseNumber)
	inspectFuncDeleteForPN   func(ctx context.Context, pn insolar.PulseNumber)
	afterDeleteForPNCounter  uint64
	beforeDeleteForPNCounter uint64
	DeleteForPNMock          mIndexCleanerMockDeleteForPN
}

// NewIndexCleanerMock returns a mock for IndexCleaner
func NewIndexCleanerMock(t minimock.Tester) *IndexCleanerMock {
	m := &IndexCleanerMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.DeleteForPNMock = mIndexCleanerMockDeleteForPN{mock: m}
	m.DeleteForPNMock.callArgs = []*IndexCleanerMockDeleteForPNParams{}

	return m
}

type mIndexCleanerMockDeleteForPN struct {
	mock               *IndexCleanerMock
	defaultExpectation *IndexCleanerMockDeleteForPNExpectation
	expectations       []*IndexCleanerMockDeleteForPNExpectation

	callArgs []*IndexCleanerMockDeleteForPNParams
	mutex    sync.RWMutex
}

// IndexCleanerMockDeleteForPNExpectation specifies expectation struct of the IndexCleaner.DeleteForPN
type IndexCleanerMockDeleteForPNExpectation struct {
	mock   *IndexCleanerMock
	params *IndexCleanerMockDeleteForPNParams

	Counter uint64
}

// IndexCleanerMockDeleteForPNParams contains parameters of the IndexCleaner.DeleteForPN
type IndexCleanerMockDeleteForPNParams struct {
	ctx context.Context
	pn  insolar.PulseNumber
}

// Expect sets up expected params for IndexCleaner.DeleteForPN
func (mmDeleteForPN *mIndexCleanerMockDeleteForPN) Expect(ctx context.Context, pn insolar.PulseNumber) *mIndexCleanerMockDeleteForPN {
	if mmDeleteForPN.mock.funcDeleteForPN != nil {
		mmDeleteForPN.mock.t.Fatalf("IndexCleanerMock.DeleteForPN mock is already set by Set")
	}

	if mmDeleteForPN.defaultExpectation == nil {
		mmDeleteForPN.defaultExpectation = &IndexCleanerMockDeleteForPNExpectation{}
	}

	mmDeleteForPN.defaultExpectation.params = &IndexCleanerMockDeleteForPNParams{ctx, pn}
	for _, e := range mmDeleteForPN.expectations {
		if minimock.Equal(e.params, mmDeleteForPN.defaultExpectation.params) {
			mmDeleteForPN.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmDeleteForPN.defaultExpectation.params)
		}
	}

	return mmDeleteForPN
}

// Inspect accepts an inspector function that has same arguments as the IndexCleaner.DeleteForPN
func (mmDeleteForPN *mIndexCleanerMockDeleteForPN) Inspect(f func(ctx context.Context, pn insolar.PulseNumber)) *mIndexCleanerMockDeleteForPN {
	if mmDeleteForPN.mock.inspectFuncDeleteForPN != nil {
		mmDeleteForPN.mock.t.Fatalf("Inspect function is already set for IndexCleanerMock.DeleteForPN")
	}

	mmDeleteForPN.mock.inspectFuncDeleteForPN = f

	return mmDeleteForPN
}

// Return sets up results that will be returned by IndexCleaner.DeleteForPN
func (mmDeleteForPN *mIndexCleanerMockDeleteForPN) Return() *IndexCleanerMock {
	if mmDeleteForPN.mock.funcDeleteForPN != nil {
		mmDeleteForPN.mock.t.Fatalf("IndexCleanerMock.DeleteForPN mock is already set by Set")
	}

	if mmDeleteForPN.defaultExpectation == nil {
		mmDeleteForPN.defaultExpectation = &IndexCleanerMockDeleteForPNExpectation{mock: mmDeleteForPN.mock}
	}

	return mmDeleteForPN.mock
}

//Set uses given function f to mock the IndexCleaner.DeleteForPN method
func (mmDeleteForPN *mIndexCleanerMockDeleteForPN) Set(f func(ctx context.Context, pn insolar.PulseNumber)) *IndexCleanerMock {
	if mmDeleteForPN.defaultExpectation != nil {
		mmDeleteForPN.mock.t.Fatalf("Default expectation is already set for the IndexCleaner.DeleteForPN method")
	}

	if len(mmDeleteForPN.expectations) > 0 {
		mmDeleteForPN.mock.t.Fatalf("Some expectations are already set for the IndexCleaner.DeleteForPN method")
	}

	mmDeleteForPN.mock.funcDeleteForPN = f
	return mmDeleteForPN.mock
}

// DeleteForPN implements IndexCleaner
func (mmDeleteForPN *IndexCleanerMock) DeleteForPN(ctx context.Context, pn insolar.PulseNumber) {
	mm_atomic.AddUint64(&mmDeleteForPN.beforeDeleteForPNCounter, 1)
	defer mm_atomic.AddUint64(&mmDeleteForPN.afterDeleteForPNCounter, 1)

	if mmDeleteForPN.inspectFuncDeleteForPN != nil {
		mmDeleteForPN.inspectFuncDeleteForPN(ctx, pn)
	}

	params := &IndexCleanerMockDeleteForPNParams{ctx, pn}

	// Record call args
	mmDeleteForPN.DeleteForPNMock.mutex.Lock()
	mmDeleteForPN.DeleteForPNMock.callArgs = append(mmDeleteForPN.DeleteForPNMock.callArgs, params)
	mmDeleteForPN.DeleteForPNMock.mutex.Unlock()

	for _, e := range mmDeleteForPN.DeleteForPNMock.expectations {
		if minimock.Equal(e.params, params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return
		}
	}

	if mmDeleteForPN.DeleteForPNMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmDeleteForPN.DeleteForPNMock.defaultExpectation.Counter, 1)
		want := mmDeleteForPN.DeleteForPNMock.defaultExpectation.params
		got := IndexCleanerMockDeleteForPNParams{ctx, pn}
		if want != nil && !minimock.Equal(*want, got) {
			mmDeleteForPN.t.Errorf("IndexCleanerMock.DeleteForPN got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		return

	}
	if mmDeleteForPN.funcDeleteForPN != nil {
		mmDeleteForPN.funcDeleteForPN(ctx, pn)
		return
	}
	mmDeleteForPN.t.Fatalf("Unexpected call to IndexCleanerMock.DeleteForPN. %v %v", ctx, pn)

}

// DeleteForPNAfterCounter returns a count of finished IndexCleanerMock.DeleteForPN invocations
func (mmDeleteForPN *IndexCleanerMock) DeleteForPNAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmDeleteForPN.afterDeleteForPNCounter)
}

// DeleteForPNBeforeCounter returns a count of IndexCleanerMock.DeleteForPN invocations
func (mmDeleteForPN *IndexCleanerMock) DeleteForPNBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmDeleteForPN.beforeDeleteForPNCounter)
}

// Calls returns a list of arguments used in each call to IndexCleanerMock.DeleteForPN.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmDeleteForPN *mIndexCleanerMockDeleteForPN) Calls() []*IndexCleanerMockDeleteForPNParams {
	mmDeleteForPN.mutex.RLock()

	argCopy := make([]*IndexCleanerMockDeleteForPNParams, len(mmDeleteForPN.callArgs))
	copy(argCopy, mmDeleteForPN.callArgs)

	mmDeleteForPN.mutex.RUnlock()

	return argCopy
}

// MinimockDeleteForPNDone returns true if the count of the DeleteForPN invocations corresponds
// the number of defined expectations
func (m *IndexCleanerMock) MinimockDeleteForPNDone() bool {
	for _, e := range m.DeleteForPNMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.DeleteForPNMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterDeleteForPNCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcDeleteForPN != nil && mm_atomic.LoadUint64(&m.afterDeleteForPNCounter) < 1 {
		return false
	}
	return true
}

// MinimockDeleteForPNInspect logs each unmet expectation
func (m *IndexCleanerMock) MinimockDeleteForPNInspect() {
	for _, e := range m.DeleteForPNMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to IndexCleanerMock.DeleteForPN with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.DeleteForPNMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterDeleteForPNCounter) < 1 {
		if m.DeleteForPNMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to IndexCleanerMock.DeleteForPN")
		} else {
			m.t.Errorf("Expected call to IndexCleanerMock.DeleteForPN with params: %#v", *m.DeleteForPNMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcDeleteForPN != nil && mm_atomic.LoadUint64(&m.afterDeleteForPNCounter) < 1 {
		m.t.Error("Expected call to IndexCleanerMock.DeleteForPN")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *IndexCleanerMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockDeleteForPNInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *IndexCleanerMock) MinimockWait(timeout mm_time.Duration) {
	timeoutCh := mm_time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-mm_time.After(10 * mm_time.Millisecond):
		}
	}
}

func (m *IndexCleanerMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockDeleteForPNDone()
}
