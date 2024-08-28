package mg

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Aderon1333/TPES/internal/models"
	"github.com/Aderon1333/TPES/internal/service/manager"
	"github.com/Aderon1333/TPES/pkg/utils/logfacade"
)

type UrlDAO struct {
	c *mongo.Collection
}

type mongorepo struct {
	urlDAO UrlDAO
}

func NewUrlDAO(ctx context.Context, client *mongo.Client) *UrlDAO {
	return &UrlDAO{
		c: client.Database("mongodb").Collection("tasks"),
	}
}

func NewRepository(urlDAO *UrlDAO) manager.TaskRepository {
	return &mongorepo{
		urlDAO: *urlDAO,
	}
}

func (r *mongorepo) Create(ctx context.Context, task *models.Task, logger *logfacade.LogFacade) error {
	_, err := r.urlDAO.c.InsertOne(ctx, task)
	if err != nil {
		return err
	}
	return nil
}

func (r *mongorepo) FindOne(ctx context.Context, id int, logger *logfacade.LogFacade) (*models.Task, error) {
	var task models.Task
	filter := bson.M{"_id": id}

	err := r.urlDAO.c.FindOne(ctx, filter).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		} else {
			return nil, err
		}
	}

	return &task, nil
}
