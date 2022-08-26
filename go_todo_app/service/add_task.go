package service

import (
	"context"
	"fmt"

	"github.com/koralle/go-web-application-development/go_todo_app/entity"
	"github.com/koralle/go-web-application-development/go_todo_app/store"
)

type AddTask struct {
	DB   store.Execer
	Repo TaskAddr
}

func (a *AddTask) AddTask(ctx context.Context, title string) (*entity.Task, error) {
	t := &entity.Task{
		Title:  title,
		Status: entity.TaskStatusTodo,
	}

	if err := a.Repo.AddTask(ctx, a.DB, t); err != nil {
		return nil, fmt.Errorf("failed to register: %v", err)
	}

	return t, nil
}
