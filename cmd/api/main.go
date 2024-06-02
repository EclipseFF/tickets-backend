package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	dsn := "postgres://postgres:asd123@localhost:5432/tap2go?sslmode=disable"
	//dsn := "postgres://default:y30QtohsfVwg@ep-green-wildflower-a282vvrx.eu-central-1.aws.neon.tech:5432/verceldb?sslmode=require"
	port := ":4000"
	app, err := NewApp(&dsn, &port)
	if err != nil {
		panic(err)
	}

	go func() {
		if err := app.server.Start(*app.config.port); err != nil && !errors.Is(err, http.ErrServerClosed) {
			app.server.Logger.Fatal("shutting down the server")
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := app.server.Shutdown(ctx); err != nil {
		app.server.Logger.Fatal(err)
	}
}
