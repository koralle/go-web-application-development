package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"golang.org/x/sync/errgroup"
)

func Run(ctx context.Context, l net.Listener) error {

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
