// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	context "context"

	logger "github.com/smartcontractkit/chainlink/core/logger"
	mock "github.com/stretchr/testify/mock"

	pg "github.com/smartcontractkit/chainlink/core/services/pg"

	pipeline "github.com/smartcontractkit/chainlink/core/services/pipeline"

	uuid "github.com/satori/go.uuid"
)

// Runner is an autogenerated mock type for the Runner type
type Runner struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *Runner) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ExecuteAndInsertFinishedRun provides a mock function with given fields: ctx, spec, vars, l, saveSuccessfulTaskRuns
func (_m *Runner) ExecuteAndInsertFinishedRun(ctx context.Context, spec pipeline.Spec, vars pipeline.Vars, l logger.Logger, saveSuccessfulTaskRuns bool) (int64, pipeline.FinalResult, error) {
	ret := _m.Called(ctx, spec, vars, l, saveSuccessfulTaskRuns)

	var r0 int64
	var r1 pipeline.FinalResult
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, pipeline.Spec, pipeline.Vars, logger.Logger, bool) (int64, pipeline.FinalResult, error)); ok {
		return rf(ctx, spec, vars, l, saveSuccessfulTaskRuns)
	}
	if rf, ok := ret.Get(0).(func(context.Context, pipeline.Spec, pipeline.Vars, logger.Logger, bool) int64); ok {
		r0 = rf(ctx, spec, vars, l, saveSuccessfulTaskRuns)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, pipeline.Spec, pipeline.Vars, logger.Logger, bool) pipeline.FinalResult); ok {
		r1 = rf(ctx, spec, vars, l, saveSuccessfulTaskRuns)
	} else {
		r1 = ret.Get(1).(pipeline.FinalResult)
	}

	if rf, ok := ret.Get(2).(func(context.Context, pipeline.Spec, pipeline.Vars, logger.Logger, bool) error); ok {
		r2 = rf(ctx, spec, vars, l, saveSuccessfulTaskRuns)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// ExecuteRun provides a mock function with given fields: ctx, spec, vars, l
func (_m *Runner) ExecuteRun(ctx context.Context, spec pipeline.Spec, vars pipeline.Vars, l logger.Logger) (pipeline.Run, pipeline.TaskRunResults, error) {
	ret := _m.Called(ctx, spec, vars, l)

	var r0 pipeline.Run
	var r1 pipeline.TaskRunResults
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, pipeline.Spec, pipeline.Vars, logger.Logger) (pipeline.Run, pipeline.TaskRunResults, error)); ok {
		return rf(ctx, spec, vars, l)
	}
	if rf, ok := ret.Get(0).(func(context.Context, pipeline.Spec, pipeline.Vars, logger.Logger) pipeline.Run); ok {
		r0 = rf(ctx, spec, vars, l)
	} else {
		r0 = ret.Get(0).(pipeline.Run)
	}

	if rf, ok := ret.Get(1).(func(context.Context, pipeline.Spec, pipeline.Vars, logger.Logger) pipeline.TaskRunResults); ok {
		r1 = rf(ctx, spec, vars, l)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(pipeline.TaskRunResults)
		}
	}

	if rf, ok := ret.Get(2).(func(context.Context, pipeline.Spec, pipeline.Vars, logger.Logger) error); ok {
		r2 = rf(ctx, spec, vars, l)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// HealthReport provides a mock function with given fields:
func (_m *Runner) HealthReport() map[string]error {
	ret := _m.Called()

	var r0 map[string]error
	if rf, ok := ret.Get(0).(func() map[string]error); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]error)
		}
	}

	return r0
}

// Healthy provides a mock function with given fields:
func (_m *Runner) Healthy() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// InsertFinishedRun provides a mock function with given fields: run, saveSuccessfulTaskRuns, qopts
func (_m *Runner) InsertFinishedRun(run *pipeline.Run, saveSuccessfulTaskRuns bool, qopts ...pg.QOpt) error {
	_va := make([]interface{}, len(qopts))
	for _i := range qopts {
		_va[_i] = qopts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, run, saveSuccessfulTaskRuns)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(*pipeline.Run, bool, ...pg.QOpt) error); ok {
		r0 = rf(run, saveSuccessfulTaskRuns, qopts...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// InsertFinishedRuns provides a mock function with given fields: runs, saveSuccessfulTaskRuns, qopts
func (_m *Runner) InsertFinishedRuns(runs []*pipeline.Run, saveSuccessfulTaskRuns bool, qopts ...pg.QOpt) error {
	_va := make([]interface{}, len(qopts))
	for _i := range qopts {
		_va[_i] = qopts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, runs, saveSuccessfulTaskRuns)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func([]*pipeline.Run, bool, ...pg.QOpt) error); ok {
		r0 = rf(runs, saveSuccessfulTaskRuns, qopts...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Name provides a mock function with given fields:
func (_m *Runner) Name() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// OnRunFinished provides a mock function with given fields: _a0
func (_m *Runner) OnRunFinished(_a0 func(*pipeline.Run)) {
	_m.Called(_a0)
}

// Ready provides a mock function with given fields:
func (_m *Runner) Ready() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ResumeRun provides a mock function with given fields: taskID, value, err
func (_m *Runner) ResumeRun(taskID uuid.UUID, value interface{}, err error) error {
	ret := _m.Called(taskID, value, err)

	var r0 error
	if rf, ok := ret.Get(0).(func(uuid.UUID, interface{}, error) error); ok {
		r0 = rf(taskID, value, err)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Run provides a mock function with given fields: ctx, run, l, saveSuccessfulTaskRuns, fn
func (_m *Runner) Run(ctx context.Context, run *pipeline.Run, l logger.Logger, saveSuccessfulTaskRuns bool, fn func(pg.Queryer) error) (bool, error) {
	ret := _m.Called(ctx, run, l, saveSuccessfulTaskRuns, fn)

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *pipeline.Run, logger.Logger, bool, func(pg.Queryer) error) (bool, error)); ok {
		return rf(ctx, run, l, saveSuccessfulTaskRuns, fn)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *pipeline.Run, logger.Logger, bool, func(pg.Queryer) error) bool); ok {
		r0 = rf(ctx, run, l, saveSuccessfulTaskRuns, fn)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *pipeline.Run, logger.Logger, bool, func(pg.Queryer) error) error); ok {
		r1 = rf(ctx, run, l, saveSuccessfulTaskRuns, fn)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Start provides a mock function with given fields: _a0
func (_m *Runner) Start(_a0 context.Context) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewRunner interface {
	mock.TestingT
	Cleanup(func())
}

// NewRunner creates a new instance of Runner. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRunner(t mockConstructorTestingTNewRunner) *Runner {
	mock := &Runner{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
