package testutils

// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

import (
	"context"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock"
	mm_insolar "github.com/insolar/insolar/insolar"
)

// LogicRunnerMock implements insolar.LogicRunner
type LogicRunnerMock struct {
	t minimock.Tester

	funcAddUnwantedResponse          func(ctx context.Context, msg mm_insolar.Payload) (err error)
	inspectFuncAddUnwantedResponse   func(ctx context.Context, msg mm_insolar.Payload)
	afterAddUnwantedResponseCounter  uint64
	beforeAddUnwantedResponseCounter uint64
	AddUnwantedResponseMock          mLogicRunnerMockAddUnwantedResponse

	funcLRI          func()
	inspectFuncLRI   func()
	afterLRICounter  uint64
	beforeLRICounter uint64
	LRIMock          mLogicRunnerMockLRI

	funcOnPulse          func(ctx context.Context, p1 mm_insolar.Pulse, p2 mm_insolar.Pulse) (err error)
	inspectFuncOnPulse   func(ctx context.Context, p1 mm_insolar.Pulse, p2 mm_insolar.Pulse)
	afterOnPulseCounter  uint64
	beforeOnPulseCounter uint64
	OnPulseMock          mLogicRunnerMockOnPulse
}

// NewLogicRunnerMock returns a mock for insolar.LogicRunner
func NewLogicRunnerMock(t minimock.Tester) *LogicRunnerMock {
	m := &LogicRunnerMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AddUnwantedResponseMock = mLogicRunnerMockAddUnwantedResponse{mock: m}
	m.AddUnwantedResponseMock.callArgs = []*LogicRunnerMockAddUnwantedResponseParams{}

	m.LRIMock = mLogicRunnerMockLRI{mock: m}

	m.OnPulseMock = mLogicRunnerMockOnPulse{mock: m}
	m.OnPulseMock.callArgs = []*LogicRunnerMockOnPulseParams{}

	return m
}

type mLogicRunnerMockAddUnwantedResponse struct {
	mock               *LogicRunnerMock
	defaultExpectation *LogicRunnerMockAddUnwantedResponseExpectation
	expectations       []*LogicRunnerMockAddUnwantedResponseExpectation

	callArgs []*LogicRunnerMockAddUnwantedResponseParams
	mutex    sync.RWMutex
}

// LogicRunnerMockAddUnwantedResponseExpectation specifies expectation struct of the LogicRunner.AddUnwantedResponse
type LogicRunnerMockAddUnwantedResponseExpectation struct {
	mock    *LogicRunnerMock
	params  *LogicRunnerMockAddUnwantedResponseParams
	results *LogicRunnerMockAddUnwantedResponseResults
	Counter uint64
}

// LogicRunnerMockAddUnwantedResponseParams contains parameters of the LogicRunner.AddUnwantedResponse
type LogicRunnerMockAddUnwantedResponseParams struct {
	ctx context.Context
	msg mm_insolar.Payload
}

// LogicRunnerMockAddUnwantedResponseResults contains results of the LogicRunner.AddUnwantedResponse
type LogicRunnerMockAddUnwantedResponseResults struct {
	err error
}

// Expect sets up expected params for LogicRunner.AddUnwantedResponse
func (mmAddUnwantedResponse *mLogicRunnerMockAddUnwantedResponse) Expect(ctx context.Context, msg mm_insolar.Payload) *mLogicRunnerMockAddUnwantedResponse {
	if mmAddUnwantedResponse.mock.funcAddUnwantedResponse != nil {
		mmAddUnwantedResponse.mock.t.Fatalf("LogicRunnerMock.AddUnwantedResponse mock is already set by Set")
	}

	if mmAddUnwantedResponse.defaultExpectation == nil {
		mmAddUnwantedResponse.defaultExpectation = &LogicRunnerMockAddUnwantedResponseExpectation{}
	}

	mmAddUnwantedResponse.defaultExpectation.params = &LogicRunnerMockAddUnwantedResponseParams{ctx, msg}
	for _, e := range mmAddUnwantedResponse.expectations {
		if minimock.Equal(e.params, mmAddUnwantedResponse.defaultExpectation.params) {
			mmAddUnwantedResponse.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmAddUnwantedResponse.defaultExpectation.params)
		}
	}

	return mmAddUnwantedResponse
}

// Inspect accepts an inspector function that has same arguments as the LogicRunner.AddUnwantedResponse
func (mmAddUnwantedResponse *mLogicRunnerMockAddUnwantedResponse) Inspect(f func(ctx context.Context, msg mm_insolar.Payload)) *mLogicRunnerMockAddUnwantedResponse {
	if mmAddUnwantedResponse.mock.inspectFuncAddUnwantedResponse != nil {
		mmAddUnwantedResponse.mock.t.Fatalf("Inspect function is already set for LogicRunnerMock.AddUnwantedResponse")
	}

	mmAddUnwantedResponse.mock.inspectFuncAddUnwantedResponse = f

	return mmAddUnwantedResponse
}

// Return sets up results that will be returned by LogicRunner.AddUnwantedResponse
func (mmAddUnwantedResponse *mLogicRunnerMockAddUnwantedResponse) Return(err error) *LogicRunnerMock {
	if mmAddUnwantedResponse.mock.funcAddUnwantedResponse != nil {
		mmAddUnwantedResponse.mock.t.Fatalf("LogicRunnerMock.AddUnwantedResponse mock is already set by Set")
	}

	if mmAddUnwantedResponse.defaultExpectation == nil {
		mmAddUnwantedResponse.defaultExpectation = &LogicRunnerMockAddUnwantedResponseExpectation{mock: mmAddUnwantedResponse.mock}
	}
	mmAddUnwantedResponse.defaultExpectation.results = &LogicRunnerMockAddUnwantedResponseResults{err}
	return mmAddUnwantedResponse.mock
}

//Set uses given function f to mock the LogicRunner.AddUnwantedResponse method
func (mmAddUnwantedResponse *mLogicRunnerMockAddUnwantedResponse) Set(f func(ctx context.Context, msg mm_insolar.Payload) (err error)) *LogicRunnerMock {
	if mmAddUnwantedResponse.defaultExpectation != nil {
		mmAddUnwantedResponse.mock.t.Fatalf("Default expectation is already set for the LogicRunner.AddUnwantedResponse method")
	}

	if len(mmAddUnwantedResponse.expectations) > 0 {
		mmAddUnwantedResponse.mock.t.Fatalf("Some expectations are already set for the LogicRunner.AddUnwantedResponse method")
	}

	mmAddUnwantedResponse.mock.funcAddUnwantedResponse = f
	return mmAddUnwantedResponse.mock
}

// When sets expectation for the LogicRunner.AddUnwantedResponse which will trigger the result defined by the following
// Then helper
func (mmAddUnwantedResponse *mLogicRunnerMockAddUnwantedResponse) When(ctx context.Context, msg mm_insolar.Payload) *LogicRunnerMockAddUnwantedResponseExpectation {
	if mmAddUnwantedResponse.mock.funcAddUnwantedResponse != nil {
		mmAddUnwantedResponse.mock.t.Fatalf("LogicRunnerMock.AddUnwantedResponse mock is already set by Set")
	}

	expectation := &LogicRunnerMockAddUnwantedResponseExpectation{
		mock:   mmAddUnwantedResponse.mock,
		params: &LogicRunnerMockAddUnwantedResponseParams{ctx, msg},
	}
	mmAddUnwantedResponse.expectations = append(mmAddUnwantedResponse.expectations, expectation)
	return expectation
}

// Then sets up LogicRunner.AddUnwantedResponse return parameters for the expectation previously defined by the When method
func (e *LogicRunnerMockAddUnwantedResponseExpectation) Then(err error) *LogicRunnerMock {
	e.results = &LogicRunnerMockAddUnwantedResponseResults{err}
	return e.mock
}

// AddUnwantedResponse implements insolar.LogicRunner
func (mmAddUnwantedResponse *LogicRunnerMock) AddUnwantedResponse(ctx context.Context, msg mm_insolar.Payload) (err error) {
	mm_atomic.AddUint64(&mmAddUnwantedResponse.beforeAddUnwantedResponseCounter, 1)
	defer mm_atomic.AddUint64(&mmAddUnwantedResponse.afterAddUnwantedResponseCounter, 1)

	if mmAddUnwantedResponse.inspectFuncAddUnwantedResponse != nil {
		mmAddUnwantedResponse.inspectFuncAddUnwantedResponse(ctx, msg)
	}

	params := &LogicRunnerMockAddUnwantedResponseParams{ctx, msg}

	// Record call args
	mmAddUnwantedResponse.AddUnwantedResponseMock.mutex.Lock()
	mmAddUnwantedResponse.AddUnwantedResponseMock.callArgs = append(mmAddUnwantedResponse.AddUnwantedResponseMock.callArgs, params)
	mmAddUnwantedResponse.AddUnwantedResponseMock.mutex.Unlock()

	for _, e := range mmAddUnwantedResponse.AddUnwantedResponseMock.expectations {
		if minimock.Equal(e.params, params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.err
		}
	}

	if mmAddUnwantedResponse.AddUnwantedResponseMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmAddUnwantedResponse.AddUnwantedResponseMock.defaultExpectation.Counter, 1)
		want := mmAddUnwantedResponse.AddUnwantedResponseMock.defaultExpectation.params
		got := LogicRunnerMockAddUnwantedResponseParams{ctx, msg}
		if want != nil && !minimock.Equal(*want, got) {
			mmAddUnwantedResponse.t.Errorf("LogicRunnerMock.AddUnwantedResponse got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		results := mmAddUnwantedResponse.AddUnwantedResponseMock.defaultExpectation.results
		if results == nil {
			mmAddUnwantedResponse.t.Fatal("No results are set for the LogicRunnerMock.AddUnwantedResponse")
		}
		return (*results).err
	}
	if mmAddUnwantedResponse.funcAddUnwantedResponse != nil {
		return mmAddUnwantedResponse.funcAddUnwantedResponse(ctx, msg)
	}
	mmAddUnwantedResponse.t.Fatalf("Unexpected call to LogicRunnerMock.AddUnwantedResponse. %v %v", ctx, msg)
	return
}

// AddUnwantedResponseAfterCounter returns a count of finished LogicRunnerMock.AddUnwantedResponse invocations
func (mmAddUnwantedResponse *LogicRunnerMock) AddUnwantedResponseAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmAddUnwantedResponse.afterAddUnwantedResponseCounter)
}

// AddUnwantedResponseBeforeCounter returns a count of LogicRunnerMock.AddUnwantedResponse invocations
func (mmAddUnwantedResponse *LogicRunnerMock) AddUnwantedResponseBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmAddUnwantedResponse.beforeAddUnwantedResponseCounter)
}

// Calls returns a list of arguments used in each call to LogicRunnerMock.AddUnwantedResponse.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmAddUnwantedResponse *mLogicRunnerMockAddUnwantedResponse) Calls() []*LogicRunnerMockAddUnwantedResponseParams {
	mmAddUnwantedResponse.mutex.RLock()

	argCopy := make([]*LogicRunnerMockAddUnwantedResponseParams, len(mmAddUnwantedResponse.callArgs))
	copy(argCopy, mmAddUnwantedResponse.callArgs)

	mmAddUnwantedResponse.mutex.RUnlock()

	return argCopy
}

// MinimockAddUnwantedResponseDone returns true if the count of the AddUnwantedResponse invocations corresponds
// the number of defined expectations
func (m *LogicRunnerMock) MinimockAddUnwantedResponseDone() bool {
	for _, e := range m.AddUnwantedResponseMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.AddUnwantedResponseMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterAddUnwantedResponseCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcAddUnwantedResponse != nil && mm_atomic.LoadUint64(&m.afterAddUnwantedResponseCounter) < 1 {
		return false
	}
	return true
}

// MinimockAddUnwantedResponseInspect logs each unmet expectation
func (m *LogicRunnerMock) MinimockAddUnwantedResponseInspect() {
	for _, e := range m.AddUnwantedResponseMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to LogicRunnerMock.AddUnwantedResponse with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.AddUnwantedResponseMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterAddUnwantedResponseCounter) < 1 {
		if m.AddUnwantedResponseMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to LogicRunnerMock.AddUnwantedResponse")
		} else {
			m.t.Errorf("Expected call to LogicRunnerMock.AddUnwantedResponse with params: %#v", *m.AddUnwantedResponseMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcAddUnwantedResponse != nil && mm_atomic.LoadUint64(&m.afterAddUnwantedResponseCounter) < 1 {
		m.t.Error("Expected call to LogicRunnerMock.AddUnwantedResponse")
	}
}

type mLogicRunnerMockLRI struct {
	mock               *LogicRunnerMock
	defaultExpectation *LogicRunnerMockLRIExpectation
	expectations       []*LogicRunnerMockLRIExpectation
}

// LogicRunnerMockLRIExpectation specifies expectation struct of the LogicRunner.LRI
type LogicRunnerMockLRIExpectation struct {
	mock *LogicRunnerMock

	Counter uint64
}

// Expect sets up expected params for LogicRunner.LRI
func (mmLRI *mLogicRunnerMockLRI) Expect() *mLogicRunnerMockLRI {
	if mmLRI.mock.funcLRI != nil {
		mmLRI.mock.t.Fatalf("LogicRunnerMock.LRI mock is already set by Set")
	}

	if mmLRI.defaultExpectation == nil {
		mmLRI.defaultExpectation = &LogicRunnerMockLRIExpectation{}
	}

	return mmLRI
}

// Inspect accepts an inspector function that has same arguments as the LogicRunner.LRI
func (mmLRI *mLogicRunnerMockLRI) Inspect(f func()) *mLogicRunnerMockLRI {
	if mmLRI.mock.inspectFuncLRI != nil {
		mmLRI.mock.t.Fatalf("Inspect function is already set for LogicRunnerMock.LRI")
	}

	mmLRI.mock.inspectFuncLRI = f

	return mmLRI
}

// Return sets up results that will be returned by LogicRunner.LRI
func (mmLRI *mLogicRunnerMockLRI) Return() *LogicRunnerMock {
	if mmLRI.mock.funcLRI != nil {
		mmLRI.mock.t.Fatalf("LogicRunnerMock.LRI mock is already set by Set")
	}

	if mmLRI.defaultExpectation == nil {
		mmLRI.defaultExpectation = &LogicRunnerMockLRIExpectation{mock: mmLRI.mock}
	}

	return mmLRI.mock
}

//Set uses given function f to mock the LogicRunner.LRI method
func (mmLRI *mLogicRunnerMockLRI) Set(f func()) *LogicRunnerMock {
	if mmLRI.defaultExpectation != nil {
		mmLRI.mock.t.Fatalf("Default expectation is already set for the LogicRunner.LRI method")
	}

	if len(mmLRI.expectations) > 0 {
		mmLRI.mock.t.Fatalf("Some expectations are already set for the LogicRunner.LRI method")
	}

	mmLRI.mock.funcLRI = f
	return mmLRI.mock
}

// LRI implements insolar.LogicRunner
func (mmLRI *LogicRunnerMock) LRI() {
	mm_atomic.AddUint64(&mmLRI.beforeLRICounter, 1)
	defer mm_atomic.AddUint64(&mmLRI.afterLRICounter, 1)

	if mmLRI.inspectFuncLRI != nil {
		mmLRI.inspectFuncLRI()
	}

	if mmLRI.LRIMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmLRI.LRIMock.defaultExpectation.Counter, 1)

		return

	}
	if mmLRI.funcLRI != nil {
		mmLRI.funcLRI()
		return
	}
	mmLRI.t.Fatalf("Unexpected call to LogicRunnerMock.LRI.")

}

// LRIAfterCounter returns a count of finished LogicRunnerMock.LRI invocations
func (mmLRI *LogicRunnerMock) LRIAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmLRI.afterLRICounter)
}

// LRIBeforeCounter returns a count of LogicRunnerMock.LRI invocations
func (mmLRI *LogicRunnerMock) LRIBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmLRI.beforeLRICounter)
}

// MinimockLRIDone returns true if the count of the LRI invocations corresponds
// the number of defined expectations
func (m *LogicRunnerMock) MinimockLRIDone() bool {
	for _, e := range m.LRIMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.LRIMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterLRICounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcLRI != nil && mm_atomic.LoadUint64(&m.afterLRICounter) < 1 {
		return false
	}
	return true
}

// MinimockLRIInspect logs each unmet expectation
func (m *LogicRunnerMock) MinimockLRIInspect() {
	for _, e := range m.LRIMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Error("Expected call to LogicRunnerMock.LRI")
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.LRIMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterLRICounter) < 1 {
		m.t.Error("Expected call to LogicRunnerMock.LRI")
	}
	// if func was set then invocations count should be greater than zero
	if m.funcLRI != nil && mm_atomic.LoadUint64(&m.afterLRICounter) < 1 {
		m.t.Error("Expected call to LogicRunnerMock.LRI")
	}
}

type mLogicRunnerMockOnPulse struct {
	mock               *LogicRunnerMock
	defaultExpectation *LogicRunnerMockOnPulseExpectation
	expectations       []*LogicRunnerMockOnPulseExpectation

	callArgs []*LogicRunnerMockOnPulseParams
	mutex    sync.RWMutex
}

// LogicRunnerMockOnPulseExpectation specifies expectation struct of the LogicRunner.OnPulse
type LogicRunnerMockOnPulseExpectation struct {
	mock    *LogicRunnerMock
	params  *LogicRunnerMockOnPulseParams
	results *LogicRunnerMockOnPulseResults
	Counter uint64
}

// LogicRunnerMockOnPulseParams contains parameters of the LogicRunner.OnPulse
type LogicRunnerMockOnPulseParams struct {
	ctx context.Context
	p1  mm_insolar.Pulse
	p2  mm_insolar.Pulse
}

// LogicRunnerMockOnPulseResults contains results of the LogicRunner.OnPulse
type LogicRunnerMockOnPulseResults struct {
	err error
}

// Expect sets up expected params for LogicRunner.OnPulse
func (mmOnPulse *mLogicRunnerMockOnPulse) Expect(ctx context.Context, p1 mm_insolar.Pulse, p2 mm_insolar.Pulse) *mLogicRunnerMockOnPulse {
	if mmOnPulse.mock.funcOnPulse != nil {
		mmOnPulse.mock.t.Fatalf("LogicRunnerMock.OnPulse mock is already set by Set")
	}

	if mmOnPulse.defaultExpectation == nil {
		mmOnPulse.defaultExpectation = &LogicRunnerMockOnPulseExpectation{}
	}

	mmOnPulse.defaultExpectation.params = &LogicRunnerMockOnPulseParams{ctx, p1, p2}
	for _, e := range mmOnPulse.expectations {
		if minimock.Equal(e.params, mmOnPulse.defaultExpectation.params) {
			mmOnPulse.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmOnPulse.defaultExpectation.params)
		}
	}

	return mmOnPulse
}

// Inspect accepts an inspector function that has same arguments as the LogicRunner.OnPulse
func (mmOnPulse *mLogicRunnerMockOnPulse) Inspect(f func(ctx context.Context, p1 mm_insolar.Pulse, p2 mm_insolar.Pulse)) *mLogicRunnerMockOnPulse {
	if mmOnPulse.mock.inspectFuncOnPulse != nil {
		mmOnPulse.mock.t.Fatalf("Inspect function is already set for LogicRunnerMock.OnPulse")
	}

	mmOnPulse.mock.inspectFuncOnPulse = f

	return mmOnPulse
}

// Return sets up results that will be returned by LogicRunner.OnPulse
func (mmOnPulse *mLogicRunnerMockOnPulse) Return(err error) *LogicRunnerMock {
	if mmOnPulse.mock.funcOnPulse != nil {
		mmOnPulse.mock.t.Fatalf("LogicRunnerMock.OnPulse mock is already set by Set")
	}

	if mmOnPulse.defaultExpectation == nil {
		mmOnPulse.defaultExpectation = &LogicRunnerMockOnPulseExpectation{mock: mmOnPulse.mock}
	}
	mmOnPulse.defaultExpectation.results = &LogicRunnerMockOnPulseResults{err}
	return mmOnPulse.mock
}

//Set uses given function f to mock the LogicRunner.OnPulse method
func (mmOnPulse *mLogicRunnerMockOnPulse) Set(f func(ctx context.Context, p1 mm_insolar.Pulse, p2 mm_insolar.Pulse) (err error)) *LogicRunnerMock {
	if mmOnPulse.defaultExpectation != nil {
		mmOnPulse.mock.t.Fatalf("Default expectation is already set for the LogicRunner.OnPulse method")
	}

	if len(mmOnPulse.expectations) > 0 {
		mmOnPulse.mock.t.Fatalf("Some expectations are already set for the LogicRunner.OnPulse method")
	}

	mmOnPulse.mock.funcOnPulse = f
	return mmOnPulse.mock
}

// When sets expectation for the LogicRunner.OnPulse which will trigger the result defined by the following
// Then helper
func (mmOnPulse *mLogicRunnerMockOnPulse) When(ctx context.Context, p1 mm_insolar.Pulse, p2 mm_insolar.Pulse) *LogicRunnerMockOnPulseExpectation {
	if mmOnPulse.mock.funcOnPulse != nil {
		mmOnPulse.mock.t.Fatalf("LogicRunnerMock.OnPulse mock is already set by Set")
	}

	expectation := &LogicRunnerMockOnPulseExpectation{
		mock:   mmOnPulse.mock,
		params: &LogicRunnerMockOnPulseParams{ctx, p1, p2},
	}
	mmOnPulse.expectations = append(mmOnPulse.expectations, expectation)
	return expectation
}

// Then sets up LogicRunner.OnPulse return parameters for the expectation previously defined by the When method
func (e *LogicRunnerMockOnPulseExpectation) Then(err error) *LogicRunnerMock {
	e.results = &LogicRunnerMockOnPulseResults{err}
	return e.mock
}

// OnPulse implements insolar.LogicRunner
func (mmOnPulse *LogicRunnerMock) OnPulse(ctx context.Context, p1 mm_insolar.Pulse, p2 mm_insolar.Pulse) (err error) {
	mm_atomic.AddUint64(&mmOnPulse.beforeOnPulseCounter, 1)
	defer mm_atomic.AddUint64(&mmOnPulse.afterOnPulseCounter, 1)

	if mmOnPulse.inspectFuncOnPulse != nil {
		mmOnPulse.inspectFuncOnPulse(ctx, p1, p2)
	}

	params := &LogicRunnerMockOnPulseParams{ctx, p1, p2}

	// Record call args
	mmOnPulse.OnPulseMock.mutex.Lock()
	mmOnPulse.OnPulseMock.callArgs = append(mmOnPulse.OnPulseMock.callArgs, params)
	mmOnPulse.OnPulseMock.mutex.Unlock()

	for _, e := range mmOnPulse.OnPulseMock.expectations {
		if minimock.Equal(e.params, params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.err
		}
	}

	if mmOnPulse.OnPulseMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmOnPulse.OnPulseMock.defaultExpectation.Counter, 1)
		want := mmOnPulse.OnPulseMock.defaultExpectation.params
		got := LogicRunnerMockOnPulseParams{ctx, p1, p2}
		if want != nil && !minimock.Equal(*want, got) {
			mmOnPulse.t.Errorf("LogicRunnerMock.OnPulse got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		results := mmOnPulse.OnPulseMock.defaultExpectation.results
		if results == nil {
			mmOnPulse.t.Fatal("No results are set for the LogicRunnerMock.OnPulse")
		}
		return (*results).err
	}
	if mmOnPulse.funcOnPulse != nil {
		return mmOnPulse.funcOnPulse(ctx, p1, p2)
	}
	mmOnPulse.t.Fatalf("Unexpected call to LogicRunnerMock.OnPulse. %v %v %v", ctx, p1, p2)
	return
}

// OnPulseAfterCounter returns a count of finished LogicRunnerMock.OnPulse invocations
func (mmOnPulse *LogicRunnerMock) OnPulseAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmOnPulse.afterOnPulseCounter)
}

// OnPulseBeforeCounter returns a count of LogicRunnerMock.OnPulse invocations
func (mmOnPulse *LogicRunnerMock) OnPulseBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmOnPulse.beforeOnPulseCounter)
}

// Calls returns a list of arguments used in each call to LogicRunnerMock.OnPulse.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmOnPulse *mLogicRunnerMockOnPulse) Calls() []*LogicRunnerMockOnPulseParams {
	mmOnPulse.mutex.RLock()

	argCopy := make([]*LogicRunnerMockOnPulseParams, len(mmOnPulse.callArgs))
	copy(argCopy, mmOnPulse.callArgs)

	mmOnPulse.mutex.RUnlock()

	return argCopy
}

// MinimockOnPulseDone returns true if the count of the OnPulse invocations corresponds
// the number of defined expectations
func (m *LogicRunnerMock) MinimockOnPulseDone() bool {
	for _, e := range m.OnPulseMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.OnPulseMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterOnPulseCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcOnPulse != nil && mm_atomic.LoadUint64(&m.afterOnPulseCounter) < 1 {
		return false
	}
	return true
}

// MinimockOnPulseInspect logs each unmet expectation
func (m *LogicRunnerMock) MinimockOnPulseInspect() {
	for _, e := range m.OnPulseMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to LogicRunnerMock.OnPulse with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.OnPulseMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterOnPulseCounter) < 1 {
		if m.OnPulseMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to LogicRunnerMock.OnPulse")
		} else {
			m.t.Errorf("Expected call to LogicRunnerMock.OnPulse with params: %#v", *m.OnPulseMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcOnPulse != nil && mm_atomic.LoadUint64(&m.afterOnPulseCounter) < 1 {
		m.t.Error("Expected call to LogicRunnerMock.OnPulse")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *LogicRunnerMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockAddUnwantedResponseInspect()

		m.MinimockLRIInspect()

		m.MinimockOnPulseInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *LogicRunnerMock) MinimockWait(timeout mm_time.Duration) {
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

func (m *LogicRunnerMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockAddUnwantedResponseDone() &&
		m.MinimockLRIDone() &&
		m.MinimockOnPulseDone()
}
