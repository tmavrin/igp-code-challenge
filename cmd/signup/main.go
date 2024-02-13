package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tmavrin/igp-code-challenge/internal/http/api"
)

func main() {
	ctx := context.Background()

	server, err := api.NewSignupServer(ctx)
	if err != nil {
		log.Fatalf("failed to initialize signup api: %s", err)
	}

	sCtx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()
	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen and serve failed: %s", err)
		}
	}()
	<-sCtx.Done()

	err = server.Close(ctx, 10*time.Second)
	if err != nil {
		log.Fatalf("failed to stop signup api: %s", err)
	}
}
