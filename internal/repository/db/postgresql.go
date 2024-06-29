package db

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgconn"
	"github.com/sirupsen/logrus"

	"github.com/Aderon1333/TPES/internal/core"
	"github.com/Aderon1333/TPES/internal/service/manager"
	"github.com/Aderon1333/TPES/pkg/repository/postgresql"
	"github.com/Aderon1333/TPES/pkg/utils/logfacade"
)

type repo struct {
	client postgresql.Client
}

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}

func (r *repo) Create(ctx context.Context, task *core.Task, logger *logfacade.LogFacade) error {
	q := `
		INSERT INTO tasks 
		    (id, status, item) 
		VALUES 
		       ($1, $2, $3) 
		RETURNING id
	`
	logrus.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))
	if err := r.client.QueryRow(ctx, q, task.ID, task.Status, task.Item).Scan(&task.ID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			logrus.Error(newErr)
			return newErr
		}
		return err
	}

	return nil
}

func (r *repo) FindOne(ctx context.Context, id int, logger *logfacade.LogFacade) (*core.Task, error) {
	q := `
		SELECT id, status, item FROM public.tasks WHERE id = $1
	`
	logrus.Trace(fmt.Sprintf("SQL Query: %s", q))

	var task core.Task
	err := r.client.QueryRow(ctx, q, id).Scan(&task.ID, &task.Status, &task.Item)
	if err != nil {
		return &core.Task{}, err
	}

	return &task, nil
}

func NewRepository(client postgresql.Client) manager.TaskRepositoryInterfaceDB {
	return &repo{
		client: client,
	}
}
