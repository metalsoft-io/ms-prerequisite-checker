package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

func (app *application) startHTTPServer(ctx context.Context, port int) {
	defer app.wg.Done()

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      http.HandlerFunc(app.httpRequestHandler),
		ErrorLog:     log.New(app.logger, "", 0),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		<-ctx.Done()
		ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		app.logger.Info().Str("address", srv.Addr).Msg("Shutting down HTTP server")

		if err := srv.Shutdown(ctxShutdown); err != nil {
			app.logger.Error().Err(err).Str("address", srv.Addr).Msg("Error shutting down HTTP server")
		}
	}()

	app.logger.Info().Str("address", srv.Addr).Msg("Starting HTTP server")

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		app.logger.Error().Err(err).Str("address", srv.Addr).Msg("Error starting HTTP server")
	}

	app.logger.Info().Str("address", srv.Addr).Msg("HTTP server shut down")
}

func (app *application) httpRequestHandler(w http.ResponseWriter, r *http.Request) {
	app.logger.Info().Str("host", r.Host).Str("remote", r.RemoteAddr).Str("method", r.Method).Str("path", r.URL.Path).Msg("HTTP request received")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
