package server_test

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"

	"github.com/koralle/go-web-application-development/go_todo_app/server"
	"golang.org/x/sync/errgroup"
)

func TestRun(t *testing.T) {
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to listen port %v", err)
	}

	// コンテキストと、キャンセル関数の取得
	ctx, cancel := context.WithCancel(context.Background())

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return server.Run(ctx, l)
	})

	in := "message"

	url := fmt.Sprintf("http://%s/%s", l.Addr().String(), in)
	t.Logf("try request to %q", url)
	rsp, err := http.Get(url)

	if err != nil {
		t.Errorf("failed to get: %+v", err)
	}

	defer rsp.Body.Close()

	got, err := io.ReadAll(rsp.Body)

	if err != nil {
		t.Fatalf("failed to read body: %+v", err)
	}

	want := fmt.Sprintf("Hello, %s!", in)

	if string(got) != want {
		t.Errorf("want %q, but got %q", want, got)
	}

	cancel()

	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}
}