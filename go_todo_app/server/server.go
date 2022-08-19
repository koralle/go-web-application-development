package server

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

type Server struct {
	srv *http.Server
	l   net.Listener
}

func NewServer(l net.Listener, mux http.Handler) *Server {
	return &Server{
		srv: &http.Server{Handler: mux},
		l:   l,
	}
}

func (s *Server) Run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	eg, ctx := errgroup.WithContext(ctx)

	// fan-out
	// 別のゴルーチンで作成したHTTPサーバーの起動
	eg.Go(func() error {
		// HTTPサーバーの起動
		if err := s.srv.Serve(s.l); err != nil && err != http.ErrServerClosed {
			log.Printf("failed to close: %+v", err)
			return err
		}

		return nil
	})

	// fan-in
	// チャネルからの終了通知の待機
	<-ctx.Done()

	if err := s.srv.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown: %v", err)
	}

	// ゴルーチンの終了を待つ
	return eg.Wait()
}
