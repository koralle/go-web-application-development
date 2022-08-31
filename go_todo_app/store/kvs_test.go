package store

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/koralle/go-web-application-development/go_todo_app/entity"
	"github.com/koralle/go-web-application-development/go_todo_app/testutil"
)

func TestKVS_Save(t *testing.T) {
	t.Parallel()

	cli := testutil.OpenRedisForTest(t)
	sut := &KVS{Cli: cli}

	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		key := "TestKVS_load_ok"
		uid := entity.UserID(1234)

		ctx := context.Background()
		cli.Set(ctx, key, int64(uid), 30*time.Minute)

		t.Cleanup(func() {
			cli.Del(ctx, key)
		})

		got, err := sut.Load(ctx, key)

		if err != nil {
			t.Fatalf("want no error, but got %v", err)
		}

		if got != uid {
			t.Errorf("want %d, but got %d", uid, got)
		}
	})

	t.Run("notfound", func(t *testing.T) {
		t.Parallel()

		key := "TestKVS_Save_notfound"
		ctx := context.Background()

		got, err := sut.Load(ctx, key)

		if err == nil || !errors.Is(err, ErrNotFound) {
			t.Errorf("want %v, but got %v(value = %d)", ErrNotFound, err, got)
		}
	})
}
