package main

import (
	"cmp"
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	// postgres driver.
	_ "github.com/lib/pq"

	"github.com/otakakot/sample-go-server-db-test/internal/gateway"
	"github.com/otakakot/sample-go-server-db-test/internal/handler"
)

func main() {
	dsn := cmp.Or(os.Getenv("DATABASE_URL"), "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			slog.Error("failed to close database: %v", err)
		}
	}()

	gw := gateway.New(db)

	hdl := handler.New(gw)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /users", hdl.CreateUser)
	mux.HandleFunc("GET /users/{id}", hdl.ReadUser)
	mux.HandleFunc("PUT /users/{id}", hdl.UpdateUser)
	mux.HandleFunc("DELETE /users/{id}", hdl.DeleteUser)

	srv := &http.Server{
		Addr:              ":8888",
		Handler:           mux,
		ReadHeaderTimeout: 30 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	defer stop()

	go func() {
		slog.Info("server is running on " + srv.Addr)

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	<-ctx.Done()

	slog.Info("server is shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server shutdown failed: %v", err)
	}

	slog.Info("server is stopped")
}
