package pg

import (
	"context"

	"gorm.io/gorm"

	"github.com/Aderon1333/TPES/internal/models"
	"github.com/Aderon1333/TPES/internal/service/manager"
	"github.com/Aderon1333/TPES/pkg/utils/logfacade"
)

type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) manager.TaskRepository {
	return &repo{
		db: db,
	}
}

func (r *repo) Create(ctx context.Context, task *models.Task, logger *logfacade.LogFacade) error {
	result := r.db.Create(&task)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *repo) FindOne(ctx context.Context, id int, logger *logfacade.LogFacade) (*models.Task, error) {
	var task models.Task

	result := r.db.Where("id = ?", id).First(&task)
	if result.Error != nil {
		return nil, result.Error
	}

	return &task, nil
}
