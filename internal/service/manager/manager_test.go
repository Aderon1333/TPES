package manager_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Aderon1333/TPES/internal/models"
	"github.com/Aderon1333/TPES/internal/service/manager"
	"github.com/Aderon1333/TPES/mocks"
	"github.com/Aderon1333/TPES/pkg/utils/logfacade"
)

func TestGetTaskFromDB_Success(t *testing.T) {
	mockRepo := new(mocks.TaskRepository)
	taskManager := manager.NewTaskManager(mockRepo)

	// что давать как логгер?
	// Тестовый случай: успешное получение задачи
	mockRepo.On("FindOne", mock.Anything, 1, mock.Anything).Return(&models.Task{ID: 1, Status: "new", Item: "test"}, nil)
	_, err := taskManager.GetTaskFromDB(context.Background(), 1, &logfacade.LogFacade{})
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetTaskFromDB_Error(t *testing.T) {
	mockRepo := new(mocks.TaskRepository)
	taskManager := manager.NewTaskManager(mockRepo)

	// что давать как логгер?
	// Тестовый случай: ошибка при получении задачи
	mockRepo.On("FindOne", mock.Anything, 1, mock.Anything).Return(nil, errors.New("error"))
	_, err := taskManager.GetTaskFromDB(context.Background(), 1, nil)
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestPutTaskInDB_Success(t *testing.T) {

}

func TestPutTaskInDB_Error(t *testing.T) {

}
