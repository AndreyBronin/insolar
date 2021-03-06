package writecontroller

// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

import (
	"context"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar"
)

// AccessorMock implements Accessor
type AccessorMock struct {
	t minimock.Tester

	funcBegin          func(ctx context.Context, p1 insolar.PulseNumber) (done func(), err error)
	inspectFuncBegin   func(ctx context.Context, p1 insolar.PulseNumber)
	afterBeginCounter  uint64
	beforeBeginCounter uint64
	BeginMock          mAccessorMockBegin

	funcWaitOpened          func(ctx context.Context)
	inspectFuncWaitOpened   func(ctx context.Context)
	afterWaitOpenedCounter  uint64
	beforeWaitOpenedCounter uint64
	WaitOpenedMock          mAccessorMockWaitOpened
}

// NewAccessorMock returns a mock for Accessor
func NewAccessorMock(t minimock.Tester) *AccessorMock {
	m := &AccessorMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.BeginMock = mAccessorMockBegin{mock: m}
	m.BeginMock.callArgs = []*AccessorMockBeginParams{}

	m.WaitOpenedMock = mAccessorMockWaitOpened{mock: m}
	m.WaitOpenedMock.callArgs = []*AccessorMockWaitOpenedParams{}

	return m
}

type mAccessorMockBegin struct {
	mock               *AccessorMock
	defaultExpectation *AccessorMockBeginExpectation
	expectations       []*AccessorMockBeginExpectation

	callArgs []*AccessorMockBeginParams
	mutex    sync.RWMutex
}

// AccessorMockBeginExpectation specifies expectation struct of the Accessor.Begin
type AccessorMockBeginExpectation struct {
	mock    *AccessorMock
	params  *AccessorMockBeginParams
	results *AccessorMockBeginResults
	Counter uint64
}

// AccessorMockBeginParams contains parameters of the Accessor.Begin
type AccessorMockBeginParams struct {
	ctx context.Context
	p1  insolar.PulseNumber
}

// AccessorMockBeginResults contains results of the Accessor.Begin
type AccessorMockBeginResults struct {
	done func()
	err  error
}

// Expect sets up expected params for Accessor.Begin
func (mmBegin *mAccessorMockBegin) Expect(ctx context.Context, p1 insolar.PulseNumber) *mAccessorMockBegin {
	if mmBegin.mock.funcBegin != nil {
		mmBegin.mock.t.Fatalf("AccessorMock.Begin mock is already set by Set")
	}

	if mmBegin.defaultExpectation == nil {
		mmBegin.defaultExpectation = &AccessorMockBeginExpectation{}
	}

	mmBegin.defaultExpectation.params = &AccessorMockBeginParams{ctx, p1}
	for _, e := range mmBegin.expectations {
		if minimock.Equal(e.params, mmBegin.defaultExpectation.params) {
			mmBegin.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmBegin.defaultExpectation.params)
		}
	}

	return mmBegin
}

// Inspect accepts an inspector function that has same arguments as the Accessor.Begin
func (mmBegin *mAccessorMockBegin) Inspect(f func(ctx context.Context, p1 insolar.PulseNumber)) *mAccessorMockBegin {
	if mmBegin.mock.inspectFuncBegin != nil {
		mmBegin.mock.t.Fatalf("Inspect function is already set for AccessorMock.Begin")
	}

	mmBegin.mock.inspectFuncBegin = f

	return mmBegin
}

// Return sets up results that will be returned by Accessor.Begin
func (mmBegin *mAccessorMockBegin) Return(done func(), err error) *AccessorMock {
	if mmBegin.mock.funcBegin != nil {
		mmBegin.mock.t.Fatalf("AccessorMock.Begin mock is already set by Set")
	}

	if mmBegin.defaultExpectation == nil {
		mmBegin.defaultExpectation = &AccessorMockBeginExpectation{mock: mmBegin.mock}
	}
	mmBegin.defaultExpectation.results = &AccessorMockBeginResults{done, err}
	return mmBegin.mock
}

//Set uses given function f to mock the Accessor.Begin method
func (mmBegin *mAccessorMockBegin) Set(f func(ctx context.Context, p1 insolar.PulseNumber) (done func(), err error)) *AccessorMock {
	if mmBegin.defaultExpectation != nil {
		mmBegin.mock.t.Fatalf("Default expectation is already set for the Accessor.Begin method")
	}

	if len(mmBegin.expectations) > 0 {
		mmBegin.mock.t.Fatalf("Some expectations are already set for the Accessor.Begin method")
	}

	mmBegin.mock.funcBegin = f
	return mmBegin.mock
}

// When sets expectation for the Accessor.Begin which will trigger the result defined by the following
// Then helper
func (mmBegin *mAccessorMockBegin) When(ctx context.Context, p1 insolar.PulseNumber) *AccessorMockBeginExpectation {
	if mmBegin.mock.funcBegin != nil {
		mmBegin.mock.t.Fatalf("AccessorMock.Begin mock is already set by Set")
	}

	expectation := &AccessorMockBeginExpectation{
		mock:   mmBegin.mock,
		params: &AccessorMockBeginParams{ctx, p1},
	}
	mmBegin.expectations = append(mmBegin.expectations, expectation)
	return expectation
}

// Then sets up Accessor.Begin return parameters for the expectation previously defined by the When method
func (e *AccessorMockBeginExpectation) Then(done func(), err error) *AccessorMock {
	e.results = &AccessorMockBeginResults{done, err}
	return e.mock
}

// Begin implements Accessor
func (mmBegin *AccessorMock) Begin(ctx context.Context, p1 insolar.PulseNumber) (done func(), err error) {
	mm_atomic.AddUint64(&mmBegin.beforeBeginCounter, 1)
	defer mm_atomic.AddUint64(&mmBegin.afterBeginCounter, 1)

	if mmBegin.inspectFuncBegin != nil {
		mmBegin.inspectFuncBegin(ctx, p1)
	}

	params := &AccessorMockBeginParams{ctx, p1}

	// Record call args
	mmBegin.BeginMock.mutex.Lock()
	mmBegin.BeginMock.callArgs = append(mmBegin.BeginMock.callArgs, params)
	mmBegin.BeginMock.mutex.Unlock()

	for _, e := range mmBegin.BeginMock.expectations {
		if minimock.Equal(e.params, params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.done, e.results.err
		}
	}

	if mmBegin.BeginMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmBegin.BeginMock.defaultExpectation.Counter, 1)
		want := mmBegin.BeginMock.defaultExpectation.params
		got := AccessorMockBeginParams{ctx, p1}
		if want != nil && !minimock.Equal(*want, got) {
			mmBegin.t.Errorf("AccessorMock.Begin got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		results := mmBegin.BeginMock.defaultExpectation.results
		if results == nil {
			mmBegin.t.Fatal("No results are set for the AccessorMock.Begin")
		}
		return (*results).done, (*results).err
	}
	if mmBegin.funcBegin != nil {
		return mmBegin.funcBegin(ctx, p1)
	}
	mmBegin.t.Fatalf("Unexpected call to AccessorMock.Begin. %v %v", ctx, p1)
	return
}

// BeginAfterCounter returns a count of finished AccessorMock.Begin invocations
func (mmBegin *AccessorMock) BeginAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmBegin.afterBeginCounter)
}

// BeginBeforeCounter returns a count of AccessorMock.Begin invocations
func (mmBegin *AccessorMock) BeginBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmBegin.beforeBeginCounter)
}

// Calls returns a list of arguments used in each call to AccessorMock.Begin.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmBegin *mAccessorMockBegin) Calls() []*AccessorMockBeginParams {
	mmBegin.mutex.RLock()

	argCopy := make([]*AccessorMockBeginParams, len(mmBegin.callArgs))
	copy(argCopy, mmBegin.callArgs)

	mmBegin.mutex.RUnlock()

	return argCopy
}

// MinimockBeginDone returns true if the count of the Begin invocations corresponds
// the number of defined expectations
func (m *AccessorMock) MinimockBeginDone() bool {
	for _, e := range m.BeginMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.BeginMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterBeginCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcBegin != nil && mm_atomic.LoadUint64(&m.afterBeginCounter) < 1 {
		return false
	}
	return true
}

// MinimockBeginInspect logs each unmet expectation
func (m *AccessorMock) MinimockBeginInspect() {
	for _, e := range m.BeginMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to AccessorMock.Begin with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.BeginMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterBeginCounter) < 1 {
		if m.BeginMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to AccessorMock.Begin")
		} else {
			m.t.Errorf("Expected call to AccessorMock.Begin with params: %#v", *m.BeginMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcBegin != nil && mm_atomic.LoadUint64(&m.afterBeginCounter) < 1 {
		m.t.Error("Expected call to AccessorMock.Begin")
	}
}

type mAccessorMockWaitOpened struct {
	mock               *AccessorMock
	defaultExpectation *AccessorMockWaitOpenedExpectation
	expectations       []*AccessorMockWaitOpenedExpectation

	callArgs []*AccessorMockWaitOpenedParams
	mutex    sync.RWMutex
}

// AccessorMockWaitOpenedExpectation specifies expectation struct of the Accessor.WaitOpened
type AccessorMockWaitOpenedExpectation struct {
	mock   *AccessorMock
	params *AccessorMockWaitOpenedParams

	Counter uint64
}

// AccessorMockWaitOpenedParams contains parameters of the Accessor.WaitOpened
type AccessorMockWaitOpenedParams struct {
	ctx context.Context
}

// Expect sets up expected params for Accessor.WaitOpened
func (mmWaitOpened *mAccessorMockWaitOpened) Expect(ctx context.Context) *mAccessorMockWaitOpened {
	if mmWaitOpened.mock.funcWaitOpened != nil {
		mmWaitOpened.mock.t.Fatalf("AccessorMock.WaitOpened mock is already set by Set")
	}

	if mmWaitOpened.defaultExpectation == nil {
		mmWaitOpened.defaultExpectation = &AccessorMockWaitOpenedExpectation{}
	}

	mmWaitOpened.defaultExpectation.params = &AccessorMockWaitOpenedParams{ctx}
	for _, e := range mmWaitOpened.expectations {
		if minimock.Equal(e.params, mmWaitOpened.defaultExpectation.params) {
			mmWaitOpened.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmWaitOpened.defaultExpectation.params)
		}
	}

	return mmWaitOpened
}

// Inspect accepts an inspector function that has same arguments as the Accessor.WaitOpened
func (mmWaitOpened *mAccessorMockWaitOpened) Inspect(f func(ctx context.Context)) *mAccessorMockWaitOpened {
	if mmWaitOpened.mock.inspectFuncWaitOpened != nil {
		mmWaitOpened.mock.t.Fatalf("Inspect function is already set for AccessorMock.WaitOpened")
	}

	mmWaitOpened.mock.inspectFuncWaitOpened = f

	return mmWaitOpened
}

// Return sets up results that will be returned by Accessor.WaitOpened
func (mmWaitOpened *mAccessorMockWaitOpened) Return() *AccessorMock {
	if mmWaitOpened.mock.funcWaitOpened != nil {
		mmWaitOpened.mock.t.Fatalf("AccessorMock.WaitOpened mock is already set by Set")
	}

	if mmWaitOpened.defaultExpectation == nil {
		mmWaitOpened.defaultExpectation = &AccessorMockWaitOpenedExpectation{mock: mmWaitOpened.mock}
	}

	return mmWaitOpened.mock
}

//Set uses given function f to mock the Accessor.WaitOpened method
func (mmWaitOpened *mAccessorMockWaitOpened) Set(f func(ctx context.Context)) *AccessorMock {
	if mmWaitOpened.defaultExpectation != nil {
		mmWaitOpened.mock.t.Fatalf("Default expectation is already set for the Accessor.WaitOpened method")
	}

	if len(mmWaitOpened.expectations) > 0 {
		mmWaitOpened.mock.t.Fatalf("Some expectations are already set for the Accessor.WaitOpened method")
	}

	mmWaitOpened.mock.funcWaitOpened = f
	return mmWaitOpened.mock
}

// WaitOpened implements Accessor
func (mmWaitOpened *AccessorMock) WaitOpened(ctx context.Context) {
	mm_atomic.AddUint64(&mmWaitOpened.beforeWaitOpenedCounter, 1)
	defer mm_atomic.AddUint64(&mmWaitOpened.afterWaitOpenedCounter, 1)

	if mmWaitOpened.inspectFuncWaitOpened != nil {
		mmWaitOpened.inspectFuncWaitOpened(ctx)
	}

	params := &AccessorMockWaitOpenedParams{ctx}

	// Record call args
	mmWaitOpened.WaitOpenedMock.mutex.Lock()
	mmWaitOpened.WaitOpenedMock.callArgs = append(mmWaitOpened.WaitOpenedMock.callArgs, params)
	mmWaitOpened.WaitOpenedMock.mutex.Unlock()

	for _, e := range mmWaitOpened.WaitOpenedMock.expectations {
		if minimock.Equal(e.params, params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return
		}
	}

	if mmWaitOpened.WaitOpenedMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmWaitOpened.WaitOpenedMock.defaultExpectation.Counter, 1)
		want := mmWaitOpened.WaitOpenedMock.defaultExpectation.params
		got := AccessorMockWaitOpenedParams{ctx}
		if want != nil && !minimock.Equal(*want, got) {
			mmWaitOpened.t.Errorf("AccessorMock.WaitOpened got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		return

	}
	if mmWaitOpened.funcWaitOpened != nil {
		mmWaitOpened.funcWaitOpened(ctx)
		return
	}
	mmWaitOpened.t.Fatalf("Unexpected call to AccessorMock.WaitOpened. %v", ctx)

}

// WaitOpenedAfterCounter returns a count of finished AccessorMock.WaitOpened invocations
func (mmWaitOpened *AccessorMock) WaitOpenedAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmWaitOpened.afterWaitOpenedCounter)
}

// WaitOpenedBeforeCounter returns a count of AccessorMock.WaitOpened invocations
func (mmWaitOpened *AccessorMock) WaitOpenedBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmWaitOpened.beforeWaitOpenedCounter)
}

// Calls returns a list of arguments used in each call to AccessorMock.WaitOpened.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmWaitOpened *mAccessorMockWaitOpened) Calls() []*AccessorMockWaitOpenedParams {
	mmWaitOpened.mutex.RLock()

	argCopy := make([]*AccessorMockWaitOpenedParams, len(mmWaitOpened.callArgs))
	copy(argCopy, mmWaitOpened.callArgs)

	mmWaitOpened.mutex.RUnlock()

	return argCopy
}

// MinimockWaitOpenedDone returns true if the count of the WaitOpened invocations corresponds
// the number of defined expectations
func (m *AccessorMock) MinimockWaitOpenedDone() bool {
	for _, e := range m.WaitOpenedMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.WaitOpenedMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterWaitOpenedCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcWaitOpened != nil && mm_atomic.LoadUint64(&m.afterWaitOpenedCounter) < 1 {
		return false
	}
	return true
}

// MinimockWaitOpenedInspect logs each unmet expectation
func (m *AccessorMock) MinimockWaitOpenedInspect() {
	for _, e := range m.WaitOpenedMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to AccessorMock.WaitOpened with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.WaitOpenedMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterWaitOpenedCounter) < 1 {
		if m.WaitOpenedMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to AccessorMock.WaitOpened")
		} else {
			m.t.Errorf("Expected call to AccessorMock.WaitOpened with params: %#v", *m.WaitOpenedMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcWaitOpened != nil && mm_atomic.LoadUint64(&m.afterWaitOpenedCounter) < 1 {
		m.t.Error("Expected call to AccessorMock.WaitOpened")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *AccessorMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockBeginInspect()

		m.MinimockWaitOpenedInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *AccessorMock) MinimockWait(timeout mm_time.Duration) {
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

func (m *AccessorMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockBeginDone() &&
		m.MinimockWaitOpenedDone()
}
