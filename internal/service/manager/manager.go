package manager

import (
	"context"

	"github.com/Aderon1333/TPES/internal/core"
	"github.com/Aderon1333/TPES/internal/transport/rest/handlers"
	"github.com/Aderon1333/TPES/pkg/utils/logfacade"
)

// интерфейс для работы с репозиторием
type TaskRepositoryInterfaceDB interface {
	Create(ctx context.Context, task *core.Task, logger *logfacade.LogFacade) error
	FindOne(ctx context.Context, id int, logger *logfacade.LogFacade) (*core.Task, error)
}

type TaskManagerDB struct {
	tri TaskRepositoryInterfaceDB
}

func NewTaskManagerDB(taskRepositoryInterfaceDB TaskRepositoryInterfaceDB) handlers.TaskManagerInterfaceDB {
	return &TaskManagerDB{
		tri: taskRepositoryInterfaceDB,
	}
}

func (tm *TaskManagerDB) GetTaskFromDB(ctx context.Context, id int, logger *logfacade.LogFacade) (*core.Task, error) {
	// вызов логики из repository
	task, err := tm.tri.FindOne(ctx, id, logger)
	if err != nil {
		logger.Error("failed to get task from a repository")
		return nil, err
	}

	return task, nil
}

func (tm *TaskManagerDB) PutTaskInDB(ctx context.Context, task *core.Task, logger *logfacade.LogFacade) error {
	// вызов логики из repository
	err := tm.tri.Create(ctx, task, logger)
	if err != nil {
		logger.Error("failed to write task to repository")
	}

	return err
}
