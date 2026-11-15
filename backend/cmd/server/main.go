package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"text/tabwriter"
	"time"

	"github.com/will-steinleitner/will-steinleitner.de/internal/common/build"
	"github.com/will-steinleitner/will-steinleitner.de/internal/common/config"
	"github.com/will-steinleitner/will-steinleitner.de/internal/web/server"
)

func banner(cfg config.Config) {
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw)
	fmt.Fprintln(tw, "  Version:\t"+build.Version)
	fmt.Fprintln(tw, "  Git hash:\t"+build.GitHash)
	fmt.Fprintln(tw, "  Address:\thttp://"+net.JoinHostPort(cfg.Host, cfg.Port))
	fmt.Fprintln(tw)
	_ = tw.Flush()
}

func run(
	ctx context.Context,
	cfg config.Config,
) error {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	log.Println("Starting server at port", cfg.Port)

	rootMux := http.NewServeMux()
	server.AddRoute(rootMux, ctx, logger)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<h1>Seite 2 is running</h1>"))
	})

	httpServer := &http.Server{
		Addr:         net.JoinHostPort(cfg.Host, cfg.Port),
		Handler:      rootMux,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	done := make(chan os.Signal, 1)
	go func() {
		<-ctx.Done()
		log.Println(ctx.Err())

		log.Println("Shutting down server...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		httpServer.SetKeepAlivesEnabled(false)
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("Could not gracefully shutdown the server: %v\n", err)
		}

		close(done)
	}()

	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("ListenAndServe error: %v\n", err)
		return err
	}

	<-done
	log.Println("Server stopped")

	return nil
}

func main() {
	cfg := config.New()
	banner(cfg)

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	if err := run(ctx, cfg); err != nil {
		os.Exit(1)
	}
	log.Println("Main Context", ctx.Err())
}
