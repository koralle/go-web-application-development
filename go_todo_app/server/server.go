package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/koralle/go-web-application-development/go_todo_app/config"
	"golang.org/x/sync/errgroup"
)

func Run(ctx context.Context) error {

	cfg, err := config.New()
	if err != nil {
		return err
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen port %d: %v", cfg.Port, err)
	}

	url := fmt.Sprintf("http://%s", l.Addr().String())
	log.Printf("start with: %v", url)

	// HTTPサーバーの作成
	s := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
		}),
	}

	eg, ctx := errgroup.WithContext(ctx)

	// fan-out
	// 別のゴルーチンで作成したHTTPサーバーの起動
	eg.Go(func() error {
		// HTTPサーバーの起動
		if err := s.Serve(l); err != nil && err != http.ErrServerClosed {
			log.Printf("failed to close: %+v", err)
			return err
		}

		return nil
	})

	// fan-in
	// チャネルからの終了通知の待機
	<-ctx.Done()

	if err := s.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown: %v", err)
	}

	// ゴルーチンの終了を待つ
	return eg.Wait()
}
