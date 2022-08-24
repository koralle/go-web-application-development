package store

import (
	"context"

	"github.com/koralle/go-web-application-development/go_todo_app/entity"
)

// 全てのタスクを取得するメソッド
func (r *Repository) ListTasks(ctx context.Context, db Queryer) (entity.Tasks, error) {
	tasks := entity.Tasks{}

	sql := `SELECT id, title, _status, created, modified FROM task;`
	if err := db.SelectContext(ctx, &tasks, sql); err != nil {
		return nil, err
	}

	return tasks, nil
}

// タスクを保存するメソッド
func (r *Repository) AddTask(ctx context.Context, db Execer, t *entity.Task) error {

	t.Created = r.Clocker.Now()
	t.Modified = r.Clocker.Now()

	sql := `INSERT INTO task (title, _status, created, modified) VALUES (?, ?, ?, ?);`

	result, err := db.ExecContext(ctx, sql, t.Title, t.Status, t.Created, t.Modified)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	t.ID = entity.TaskID(id)

	return nil
}
