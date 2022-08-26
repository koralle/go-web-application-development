package service

import (
	"context"

	"github.com/koralle/go-web-application-development/go_todo_app/entity"
	"github.com/koralle/go-web-application-development/go_todo_app/store"
)

type TaskAddr interface {
	AddTask(ctx context.Context, db store.Execer, t *entity.Task) error
}

type TaskLister interface {
	ListTasks(ctx context.Context, db store.Queryer) (entity.Tasks, error)
}
