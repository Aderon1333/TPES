package manager

import (
	"context"

	"github.com/Aderon1333/TPES/internal/api/rest/handlers"
	"github.com/Aderon1333/TPES/internal/models"
	"github.com/Aderon1333/TPES/pkg/utils/logfacade"
)

// интерфейс для работы с репозиторием
type TaskRepository interface {
	Create(ctx context.Context, task *models.Task, l *logfacade.LogFacade) error
	FindOne(ctx context.Context, id int, l *logfacade.LogFacade) (*models.Task, error)
}

type TaskManager struct {
	tri TaskRepository
}

func NewTaskManager(taskRepositoryInterface TaskRepository) handlers.TaskManagerService {
	return &TaskManager{
		tri: taskRepositoryInterface,
	}
}

func (tm *TaskManager) GetTaskFromDB(ctx context.Context, id int, l *logfacade.LogFacade) (*models.Task, error) {
	task, err := tm.tri.FindOne(ctx, id, l)
	if err != nil {
		//l.Error("failed to get task from a repository")
		return nil, err
	}

	return task, nil
}

func (tm *TaskManager) PutTaskInDB(ctx context.Context, task *models.Task, l *logfacade.LogFacade) error {
	err := tm.tri.Create(ctx, task, l)
	if err != nil {
		//l.Error("failed to write task to repository")
	}

	return err
}
