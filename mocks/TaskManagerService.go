// Code generated by mockery v2.44.1. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/Aderon1333/TPES/internal/models"

	logfacade "github.com/Aderon1333/TPES/pkg/utils/logfacade"

	mock "github.com/stretchr/testify/mock"
)

// TaskManagerService is an autogenerated mock type for the TaskManagerService type
type TaskManagerService struct {
	mock.Mock
}

// GetTaskFromDB provides a mock function with given fields: ctx, id, logger
func (_m *TaskManagerService) GetTaskFromDB(ctx context.Context, id int, logger *logfacade.LogFacade) (*models.Task, error) {
	ret := _m.Called(ctx, id, logger)

	if len(ret) == 0 {
		panic("no return value specified for GetTaskFromDB")
	}

	var r0 *models.Task
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int, *logfacade.LogFacade) (*models.Task, error)); ok {
		return rf(ctx, id, logger)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, *logfacade.LogFacade) *models.Task); ok {
		r0 = rf(ctx, id, logger)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Task)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, *logfacade.LogFacade) error); ok {
		r1 = rf(ctx, id, logger)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutTaskInDB provides a mock function with given fields: ctx, task, logger
func (_m *TaskManagerService) PutTaskInDB(ctx context.Context, task *models.Task, logger *logfacade.LogFacade) error {
	ret := _m.Called(ctx, task, logger)

	if len(ret) == 0 {
		panic("no return value specified for PutTaskInDB")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Task, *logfacade.LogFacade) error); ok {
		r0 = rf(ctx, task, logger)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewTaskManagerService creates a new instance of TaskManagerService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTaskManagerService(t interface {
	mock.TestingT
	Cleanup(func())
}) *TaskManagerService {
	mock := &TaskManagerService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
