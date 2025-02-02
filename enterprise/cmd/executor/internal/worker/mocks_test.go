// Code generated by go-mockgen 1.3.7; DO NOT EDIT.
//
// This file was generated by running `sg generate` (or `go-mockgen`) at the root of
// this repository. To add additional mocks to this or another package, add a new entry
// to the mockgen.yaml file in the root of this repository.

package worker

import (
	"context"
	"io"
	"sync"

	command "github.com/sourcegraph/sourcegraph/enterprise/cmd/executor/internal/command"
	workspace "github.com/sourcegraph/sourcegraph/enterprise/cmd/executor/internal/worker/workspace"
	types "github.com/sourcegraph/sourcegraph/enterprise/internal/executor/types"
	executor "github.com/sourcegraph/sourcegraph/internal/executor"
)

// MockExecutionLogEntryStore is a mock implementation of the
// ExecutionLogEntryStore interface (from the package
// github.com/sourcegraph/sourcegraph/enterprise/cmd/executor/internal/command)
// used for unit testing.
type MockExecutionLogEntryStore struct {
	// AddExecutionLogEntryFunc is an instance of a mock function object
	// controlling the behavior of the method AddExecutionLogEntry.
	AddExecutionLogEntryFunc *ExecutionLogEntryStoreAddExecutionLogEntryFunc
	// UpdateExecutionLogEntryFunc is an instance of a mock function object
	// controlling the behavior of the method UpdateExecutionLogEntry.
	UpdateExecutionLogEntryFunc *ExecutionLogEntryStoreUpdateExecutionLogEntryFunc
}

// NewMockExecutionLogEntryStore creates a new mock of the
// ExecutionLogEntryStore interface. All methods return zero values for all
// results, unless overwritten.
func NewMockExecutionLogEntryStore() *MockExecutionLogEntryStore {
	return &MockExecutionLogEntryStore{
		AddExecutionLogEntryFunc: &ExecutionLogEntryStoreAddExecutionLogEntryFunc{
			defaultHook: func(context.Context, types.Job, executor.ExecutionLogEntry) (r0 int, r1 error) {
				return
			},
		},
		UpdateExecutionLogEntryFunc: &ExecutionLogEntryStoreUpdateExecutionLogEntryFunc{
			defaultHook: func(context.Context, types.Job, int, executor.ExecutionLogEntry) (r0 error) {
				return
			},
		},
	}
}

// NewStrictMockExecutionLogEntryStore creates a new mock of the
// ExecutionLogEntryStore interface. All methods panic on invocation, unless
// overwritten.
func NewStrictMockExecutionLogEntryStore() *MockExecutionLogEntryStore {
	return &MockExecutionLogEntryStore{
		AddExecutionLogEntryFunc: &ExecutionLogEntryStoreAddExecutionLogEntryFunc{
			defaultHook: func(context.Context, types.Job, executor.ExecutionLogEntry) (int, error) {
				panic("unexpected invocation of MockExecutionLogEntryStore.AddExecutionLogEntry")
			},
		},
		UpdateExecutionLogEntryFunc: &ExecutionLogEntryStoreUpdateExecutionLogEntryFunc{
			defaultHook: func(context.Context, types.Job, int, executor.ExecutionLogEntry) error {
				panic("unexpected invocation of MockExecutionLogEntryStore.UpdateExecutionLogEntry")
			},
		},
	}
}

// NewMockExecutionLogEntryStoreFrom creates a new mock of the
// MockExecutionLogEntryStore interface. All methods delegate to the given
// implementation, unless overwritten.
func NewMockExecutionLogEntryStoreFrom(i command.ExecutionLogEntryStore) *MockExecutionLogEntryStore {
	return &MockExecutionLogEntryStore{
		AddExecutionLogEntryFunc: &ExecutionLogEntryStoreAddExecutionLogEntryFunc{
			defaultHook: i.AddExecutionLogEntry,
		},
		UpdateExecutionLogEntryFunc: &ExecutionLogEntryStoreUpdateExecutionLogEntryFunc{
			defaultHook: i.UpdateExecutionLogEntry,
		},
	}
}

// ExecutionLogEntryStoreAddExecutionLogEntryFunc describes the behavior
// when the AddExecutionLogEntry method of the parent
// MockExecutionLogEntryStore instance is invoked.
type ExecutionLogEntryStoreAddExecutionLogEntryFunc struct {
	defaultHook func(context.Context, types.Job, executor.ExecutionLogEntry) (int, error)
	hooks       []func(context.Context, types.Job, executor.ExecutionLogEntry) (int, error)
	history     []ExecutionLogEntryStoreAddExecutionLogEntryFuncCall
	mutex       sync.Mutex
}

// AddExecutionLogEntry delegates to the next hook function in the queue and
// stores the parameter and result values of this invocation.
func (m *MockExecutionLogEntryStore) AddExecutionLogEntry(v0 context.Context, v1 types.Job, v2 executor.ExecutionLogEntry) (int, error) {
	r0, r1 := m.AddExecutionLogEntryFunc.nextHook()(v0, v1, v2)
	m.AddExecutionLogEntryFunc.appendCall(ExecutionLogEntryStoreAddExecutionLogEntryFuncCall{v0, v1, v2, r0, r1})
	return r0, r1
}

// SetDefaultHook sets function that is called when the AddExecutionLogEntry
// method of the parent MockExecutionLogEntryStore instance is invoked and
// the hook queue is empty.
func (f *ExecutionLogEntryStoreAddExecutionLogEntryFunc) SetDefaultHook(hook func(context.Context, types.Job, executor.ExecutionLogEntry) (int, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// AddExecutionLogEntry method of the parent MockExecutionLogEntryStore
// instance invokes the hook at the front of the queue and discards it.
// After the queue is empty, the default hook function is invoked for any
// future action.
func (f *ExecutionLogEntryStoreAddExecutionLogEntryFunc) PushHook(hook func(context.Context, types.Job, executor.ExecutionLogEntry) (int, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *ExecutionLogEntryStoreAddExecutionLogEntryFunc) SetDefaultReturn(r0 int, r1 error) {
	f.SetDefaultHook(func(context.Context, types.Job, executor.ExecutionLogEntry) (int, error) {
		return r0, r1
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *ExecutionLogEntryStoreAddExecutionLogEntryFunc) PushReturn(r0 int, r1 error) {
	f.PushHook(func(context.Context, types.Job, executor.ExecutionLogEntry) (int, error) {
		return r0, r1
	})
}

func (f *ExecutionLogEntryStoreAddExecutionLogEntryFunc) nextHook() func(context.Context, types.Job, executor.ExecutionLogEntry) (int, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *ExecutionLogEntryStoreAddExecutionLogEntryFunc) appendCall(r0 ExecutionLogEntryStoreAddExecutionLogEntryFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of
// ExecutionLogEntryStoreAddExecutionLogEntryFuncCall objects describing the
// invocations of this function.
func (f *ExecutionLogEntryStoreAddExecutionLogEntryFunc) History() []ExecutionLogEntryStoreAddExecutionLogEntryFuncCall {
	f.mutex.Lock()
	history := make([]ExecutionLogEntryStoreAddExecutionLogEntryFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// ExecutionLogEntryStoreAddExecutionLogEntryFuncCall is an object that
// describes an invocation of method AddExecutionLogEntry on an instance of
// MockExecutionLogEntryStore.
type ExecutionLogEntryStoreAddExecutionLogEntryFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 types.Job
	// Arg2 is the value of the 3rd argument passed to this method
	// invocation.
	Arg2 executor.ExecutionLogEntry
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 int
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c ExecutionLogEntryStoreAddExecutionLogEntryFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1, c.Arg2}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c ExecutionLogEntryStoreAddExecutionLogEntryFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1}
}

// ExecutionLogEntryStoreUpdateExecutionLogEntryFunc describes the behavior
// when the UpdateExecutionLogEntry method of the parent
// MockExecutionLogEntryStore instance is invoked.
type ExecutionLogEntryStoreUpdateExecutionLogEntryFunc struct {
	defaultHook func(context.Context, types.Job, int, executor.ExecutionLogEntry) error
	hooks       []func(context.Context, types.Job, int, executor.ExecutionLogEntry) error
	history     []ExecutionLogEntryStoreUpdateExecutionLogEntryFuncCall
	mutex       sync.Mutex
}

// UpdateExecutionLogEntry delegates to the next hook function in the queue
// and stores the parameter and result values of this invocation.
func (m *MockExecutionLogEntryStore) UpdateExecutionLogEntry(v0 context.Context, v1 types.Job, v2 int, v3 executor.ExecutionLogEntry) error {
	r0 := m.UpdateExecutionLogEntryFunc.nextHook()(v0, v1, v2, v3)
	m.UpdateExecutionLogEntryFunc.appendCall(ExecutionLogEntryStoreUpdateExecutionLogEntryFuncCall{v0, v1, v2, v3, r0})
	return r0
}

// SetDefaultHook sets function that is called when the
// UpdateExecutionLogEntry method of the parent MockExecutionLogEntryStore
// instance is invoked and the hook queue is empty.
func (f *ExecutionLogEntryStoreUpdateExecutionLogEntryFunc) SetDefaultHook(hook func(context.Context, types.Job, int, executor.ExecutionLogEntry) error) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// UpdateExecutionLogEntry method of the parent MockExecutionLogEntryStore
// instance invokes the hook at the front of the queue and discards it.
// After the queue is empty, the default hook function is invoked for any
// future action.
func (f *ExecutionLogEntryStoreUpdateExecutionLogEntryFunc) PushHook(hook func(context.Context, types.Job, int, executor.ExecutionLogEntry) error) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *ExecutionLogEntryStoreUpdateExecutionLogEntryFunc) SetDefaultReturn(r0 error) {
	f.SetDefaultHook(func(context.Context, types.Job, int, executor.ExecutionLogEntry) error {
		return r0
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *ExecutionLogEntryStoreUpdateExecutionLogEntryFunc) PushReturn(r0 error) {
	f.PushHook(func(context.Context, types.Job, int, executor.ExecutionLogEntry) error {
		return r0
	})
}

func (f *ExecutionLogEntryStoreUpdateExecutionLogEntryFunc) nextHook() func(context.Context, types.Job, int, executor.ExecutionLogEntry) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *ExecutionLogEntryStoreUpdateExecutionLogEntryFunc) appendCall(r0 ExecutionLogEntryStoreUpdateExecutionLogEntryFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of
// ExecutionLogEntryStoreUpdateExecutionLogEntryFuncCall objects describing
// the invocations of this function.
func (f *ExecutionLogEntryStoreUpdateExecutionLogEntryFunc) History() []ExecutionLogEntryStoreUpdateExecutionLogEntryFuncCall {
	f.mutex.Lock()
	history := make([]ExecutionLogEntryStoreUpdateExecutionLogEntryFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// ExecutionLogEntryStoreUpdateExecutionLogEntryFuncCall is an object that
// describes an invocation of method UpdateExecutionLogEntry on an instance
// of MockExecutionLogEntryStore.
type ExecutionLogEntryStoreUpdateExecutionLogEntryFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 types.Job
	// Arg2 is the value of the 3rd argument passed to this method
	// invocation.
	Arg2 int
	// Arg3 is the value of the 4th argument passed to this method
	// invocation.
	Arg3 executor.ExecutionLogEntry
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c ExecutionLogEntryStoreUpdateExecutionLogEntryFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1, c.Arg2, c.Arg3}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c ExecutionLogEntryStoreUpdateExecutionLogEntryFuncCall) Results() []interface{} {
	return []interface{}{c.Result0}
}

// MockRunner is a mock implementation of the Runner interface (from the
// package
// github.com/sourcegraph/sourcegraph/enterprise/cmd/executor/internal/command)
// used for unit testing.
type MockRunner struct {
	// RunFunc is an instance of a mock function object controlling the
	// behavior of the method Run.
	RunFunc *RunnerRunFunc
	// SetupFunc is an instance of a mock function object controlling the
	// behavior of the method Setup.
	SetupFunc *RunnerSetupFunc
	// TeardownFunc is an instance of a mock function object controlling the
	// behavior of the method Teardown.
	TeardownFunc *RunnerTeardownFunc
}

// NewMockRunner creates a new mock of the Runner interface. All methods
// return zero values for all results, unless overwritten.
func NewMockRunner() *MockRunner {
	return &MockRunner{
		RunFunc: &RunnerRunFunc{
			defaultHook: func(context.Context, command.CommandSpec) (r0 error) {
				return
			},
		},
		SetupFunc: &RunnerSetupFunc{
			defaultHook: func(context.Context) (r0 error) {
				return
			},
		},
		TeardownFunc: &RunnerTeardownFunc{
			defaultHook: func(context.Context) (r0 error) {
				return
			},
		},
	}
}

// NewStrictMockRunner creates a new mock of the Runner interface. All
// methods panic on invocation, unless overwritten.
func NewStrictMockRunner() *MockRunner {
	return &MockRunner{
		RunFunc: &RunnerRunFunc{
			defaultHook: func(context.Context, command.CommandSpec) error {
				panic("unexpected invocation of MockRunner.Run")
			},
		},
		SetupFunc: &RunnerSetupFunc{
			defaultHook: func(context.Context) error {
				panic("unexpected invocation of MockRunner.Setup")
			},
		},
		TeardownFunc: &RunnerTeardownFunc{
			defaultHook: func(context.Context) error {
				panic("unexpected invocation of MockRunner.Teardown")
			},
		},
	}
}

// NewMockRunnerFrom creates a new mock of the MockRunner interface. All
// methods delegate to the given implementation, unless overwritten.
func NewMockRunnerFrom(i command.Runner) *MockRunner {
	return &MockRunner{
		RunFunc: &RunnerRunFunc{
			defaultHook: i.Run,
		},
		SetupFunc: &RunnerSetupFunc{
			defaultHook: i.Setup,
		},
		TeardownFunc: &RunnerTeardownFunc{
			defaultHook: i.Teardown,
		},
	}
}

// RunnerRunFunc describes the behavior when the Run method of the parent
// MockRunner instance is invoked.
type RunnerRunFunc struct {
	defaultHook func(context.Context, command.CommandSpec) error
	hooks       []func(context.Context, command.CommandSpec) error
	history     []RunnerRunFuncCall
	mutex       sync.Mutex
}

// Run delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockRunner) Run(v0 context.Context, v1 command.CommandSpec) error {
	r0 := m.RunFunc.nextHook()(v0, v1)
	m.RunFunc.appendCall(RunnerRunFuncCall{v0, v1, r0})
	return r0
}

// SetDefaultHook sets function that is called when the Run method of the
// parent MockRunner instance is invoked and the hook queue is empty.
func (f *RunnerRunFunc) SetDefaultHook(hook func(context.Context, command.CommandSpec) error) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// Run method of the parent MockRunner instance invokes the hook at the
// front of the queue and discards it. After the queue is empty, the default
// hook function is invoked for any future action.
func (f *RunnerRunFunc) PushHook(hook func(context.Context, command.CommandSpec) error) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *RunnerRunFunc) SetDefaultReturn(r0 error) {
	f.SetDefaultHook(func(context.Context, command.CommandSpec) error {
		return r0
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *RunnerRunFunc) PushReturn(r0 error) {
	f.PushHook(func(context.Context, command.CommandSpec) error {
		return r0
	})
}

func (f *RunnerRunFunc) nextHook() func(context.Context, command.CommandSpec) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *RunnerRunFunc) appendCall(r0 RunnerRunFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of RunnerRunFuncCall objects describing the
// invocations of this function.
func (f *RunnerRunFunc) History() []RunnerRunFuncCall {
	f.mutex.Lock()
	history := make([]RunnerRunFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// RunnerRunFuncCall is an object that describes an invocation of method Run
// on an instance of MockRunner.
type RunnerRunFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 command.CommandSpec
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c RunnerRunFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c RunnerRunFuncCall) Results() []interface{} {
	return []interface{}{c.Result0}
}

// RunnerSetupFunc describes the behavior when the Setup method of the
// parent MockRunner instance is invoked.
type RunnerSetupFunc struct {
	defaultHook func(context.Context) error
	hooks       []func(context.Context) error
	history     []RunnerSetupFuncCall
	mutex       sync.Mutex
}

// Setup delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockRunner) Setup(v0 context.Context) error {
	r0 := m.SetupFunc.nextHook()(v0)
	m.SetupFunc.appendCall(RunnerSetupFuncCall{v0, r0})
	return r0
}

// SetDefaultHook sets function that is called when the Setup method of the
// parent MockRunner instance is invoked and the hook queue is empty.
func (f *RunnerSetupFunc) SetDefaultHook(hook func(context.Context) error) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// Setup method of the parent MockRunner instance invokes the hook at the
// front of the queue and discards it. After the queue is empty, the default
// hook function is invoked for any future action.
func (f *RunnerSetupFunc) PushHook(hook func(context.Context) error) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *RunnerSetupFunc) SetDefaultReturn(r0 error) {
	f.SetDefaultHook(func(context.Context) error {
		return r0
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *RunnerSetupFunc) PushReturn(r0 error) {
	f.PushHook(func(context.Context) error {
		return r0
	})
}

func (f *RunnerSetupFunc) nextHook() func(context.Context) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *RunnerSetupFunc) appendCall(r0 RunnerSetupFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of RunnerSetupFuncCall objects describing the
// invocations of this function.
func (f *RunnerSetupFunc) History() []RunnerSetupFuncCall {
	f.mutex.Lock()
	history := make([]RunnerSetupFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// RunnerSetupFuncCall is an object that describes an invocation of method
// Setup on an instance of MockRunner.
type RunnerSetupFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c RunnerSetupFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c RunnerSetupFuncCall) Results() []interface{} {
	return []interface{}{c.Result0}
}

// RunnerTeardownFunc describes the behavior when the Teardown method of the
// parent MockRunner instance is invoked.
type RunnerTeardownFunc struct {
	defaultHook func(context.Context) error
	hooks       []func(context.Context) error
	history     []RunnerTeardownFuncCall
	mutex       sync.Mutex
}

// Teardown delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockRunner) Teardown(v0 context.Context) error {
	r0 := m.TeardownFunc.nextHook()(v0)
	m.TeardownFunc.appendCall(RunnerTeardownFuncCall{v0, r0})
	return r0
}

// SetDefaultHook sets function that is called when the Teardown method of
// the parent MockRunner instance is invoked and the hook queue is empty.
func (f *RunnerTeardownFunc) SetDefaultHook(hook func(context.Context) error) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// Teardown method of the parent MockRunner instance invokes the hook at the
// front of the queue and discards it. After the queue is empty, the default
// hook function is invoked for any future action.
func (f *RunnerTeardownFunc) PushHook(hook func(context.Context) error) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *RunnerTeardownFunc) SetDefaultReturn(r0 error) {
	f.SetDefaultHook(func(context.Context) error {
		return r0
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *RunnerTeardownFunc) PushReturn(r0 error) {
	f.PushHook(func(context.Context) error {
		return r0
	})
}

func (f *RunnerTeardownFunc) nextHook() func(context.Context) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *RunnerTeardownFunc) appendCall(r0 RunnerTeardownFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of RunnerTeardownFuncCall objects describing
// the invocations of this function.
func (f *RunnerTeardownFunc) History() []RunnerTeardownFuncCall {
	f.mutex.Lock()
	history := make([]RunnerTeardownFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// RunnerTeardownFuncCall is an object that describes an invocation of
// method Teardown on an instance of MockRunner.
type RunnerTeardownFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c RunnerTeardownFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c RunnerTeardownFuncCall) Results() []interface{} {
	return []interface{}{c.Result0}
}

// MockFilesStore is a mock implementation of the FilesStore interface (from
// the package
// github.com/sourcegraph/sourcegraph/enterprise/cmd/executor/internal/worker/workspace)
// used for unit testing.
type MockFilesStore struct {
	// ExistsFunc is an instance of a mock function object controlling the
	// behavior of the method Exists.
	ExistsFunc *FilesStoreExistsFunc
	// GetFunc is an instance of a mock function object controlling the
	// behavior of the method Get.
	GetFunc *FilesStoreGetFunc
}

// NewMockFilesStore creates a new mock of the FilesStore interface. All
// methods return zero values for all results, unless overwritten.
func NewMockFilesStore() *MockFilesStore {
	return &MockFilesStore{
		ExistsFunc: &FilesStoreExistsFunc{
			defaultHook: func(context.Context, types.Job, string, string) (r0 bool, r1 error) {
				return
			},
		},
		GetFunc: &FilesStoreGetFunc{
			defaultHook: func(context.Context, types.Job, string, string) (r0 io.ReadCloser, r1 error) {
				return
			},
		},
	}
}

// NewStrictMockFilesStore creates a new mock of the FilesStore interface.
// All methods panic on invocation, unless overwritten.
func NewStrictMockFilesStore() *MockFilesStore {
	return &MockFilesStore{
		ExistsFunc: &FilesStoreExistsFunc{
			defaultHook: func(context.Context, types.Job, string, string) (bool, error) {
				panic("unexpected invocation of MockFilesStore.Exists")
			},
		},
		GetFunc: &FilesStoreGetFunc{
			defaultHook: func(context.Context, types.Job, string, string) (io.ReadCloser, error) {
				panic("unexpected invocation of MockFilesStore.Get")
			},
		},
	}
}

// NewMockFilesStoreFrom creates a new mock of the MockFilesStore interface.
// All methods delegate to the given implementation, unless overwritten.
func NewMockFilesStoreFrom(i workspace.FilesStore) *MockFilesStore {
	return &MockFilesStore{
		ExistsFunc: &FilesStoreExistsFunc{
			defaultHook: i.Exists,
		},
		GetFunc: &FilesStoreGetFunc{
			defaultHook: i.Get,
		},
	}
}

// FilesStoreExistsFunc describes the behavior when the Exists method of the
// parent MockFilesStore instance is invoked.
type FilesStoreExistsFunc struct {
	defaultHook func(context.Context, types.Job, string, string) (bool, error)
	hooks       []func(context.Context, types.Job, string, string) (bool, error)
	history     []FilesStoreExistsFuncCall
	mutex       sync.Mutex
}

// Exists delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockFilesStore) Exists(v0 context.Context, v1 types.Job, v2 string, v3 string) (bool, error) {
	r0, r1 := m.ExistsFunc.nextHook()(v0, v1, v2, v3)
	m.ExistsFunc.appendCall(FilesStoreExistsFuncCall{v0, v1, v2, v3, r0, r1})
	return r0, r1
}

// SetDefaultHook sets function that is called when the Exists method of the
// parent MockFilesStore instance is invoked and the hook queue is empty.
func (f *FilesStoreExistsFunc) SetDefaultHook(hook func(context.Context, types.Job, string, string) (bool, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// Exists method of the parent MockFilesStore instance invokes the hook at
// the front of the queue and discards it. After the queue is empty, the
// default hook function is invoked for any future action.
func (f *FilesStoreExistsFunc) PushHook(hook func(context.Context, types.Job, string, string) (bool, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *FilesStoreExistsFunc) SetDefaultReturn(r0 bool, r1 error) {
	f.SetDefaultHook(func(context.Context, types.Job, string, string) (bool, error) {
		return r0, r1
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *FilesStoreExistsFunc) PushReturn(r0 bool, r1 error) {
	f.PushHook(func(context.Context, types.Job, string, string) (bool, error) {
		return r0, r1
	})
}

func (f *FilesStoreExistsFunc) nextHook() func(context.Context, types.Job, string, string) (bool, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *FilesStoreExistsFunc) appendCall(r0 FilesStoreExistsFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of FilesStoreExistsFuncCall objects describing
// the invocations of this function.
func (f *FilesStoreExistsFunc) History() []FilesStoreExistsFuncCall {
	f.mutex.Lock()
	history := make([]FilesStoreExistsFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// FilesStoreExistsFuncCall is an object that describes an invocation of
// method Exists on an instance of MockFilesStore.
type FilesStoreExistsFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 types.Job
	// Arg2 is the value of the 3rd argument passed to this method
	// invocation.
	Arg2 string
	// Arg3 is the value of the 4th argument passed to this method
	// invocation.
	Arg3 string
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 bool
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c FilesStoreExistsFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1, c.Arg2, c.Arg3}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c FilesStoreExistsFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1}
}

// FilesStoreGetFunc describes the behavior when the Get method of the
// parent MockFilesStore instance is invoked.
type FilesStoreGetFunc struct {
	defaultHook func(context.Context, types.Job, string, string) (io.ReadCloser, error)
	hooks       []func(context.Context, types.Job, string, string) (io.ReadCloser, error)
	history     []FilesStoreGetFuncCall
	mutex       sync.Mutex
}

// Get delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockFilesStore) Get(v0 context.Context, v1 types.Job, v2 string, v3 string) (io.ReadCloser, error) {
	r0, r1 := m.GetFunc.nextHook()(v0, v1, v2, v3)
	m.GetFunc.appendCall(FilesStoreGetFuncCall{v0, v1, v2, v3, r0, r1})
	return r0, r1
}

// SetDefaultHook sets function that is called when the Get method of the
// parent MockFilesStore instance is invoked and the hook queue is empty.
func (f *FilesStoreGetFunc) SetDefaultHook(hook func(context.Context, types.Job, string, string) (io.ReadCloser, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// Get method of the parent MockFilesStore instance invokes the hook at the
// front of the queue and discards it. After the queue is empty, the default
// hook function is invoked for any future action.
func (f *FilesStoreGetFunc) PushHook(hook func(context.Context, types.Job, string, string) (io.ReadCloser, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *FilesStoreGetFunc) SetDefaultReturn(r0 io.ReadCloser, r1 error) {
	f.SetDefaultHook(func(context.Context, types.Job, string, string) (io.ReadCloser, error) {
		return r0, r1
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *FilesStoreGetFunc) PushReturn(r0 io.ReadCloser, r1 error) {
	f.PushHook(func(context.Context, types.Job, string, string) (io.ReadCloser, error) {
		return r0, r1
	})
}

func (f *FilesStoreGetFunc) nextHook() func(context.Context, types.Job, string, string) (io.ReadCloser, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *FilesStoreGetFunc) appendCall(r0 FilesStoreGetFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of FilesStoreGetFuncCall objects describing
// the invocations of this function.
func (f *FilesStoreGetFunc) History() []FilesStoreGetFuncCall {
	f.mutex.Lock()
	history := make([]FilesStoreGetFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// FilesStoreGetFuncCall is an object that describes an invocation of method
// Get on an instance of MockFilesStore.
type FilesStoreGetFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 types.Job
	// Arg2 is the value of the 3rd argument passed to this method
	// invocation.
	Arg2 string
	// Arg3 is the value of the 4th argument passed to this method
	// invocation.
	Arg3 string
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 io.ReadCloser
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c FilesStoreGetFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1, c.Arg2, c.Arg3}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c FilesStoreGetFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1}
}
