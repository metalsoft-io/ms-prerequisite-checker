package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/netip"
	"time"
)

func (app *application) startHTTPServer(ctx context.Context, ip netip.Addr, port uint16) {
	defer app.wg.Done()

	address := netip.AddrPortFrom(ip, port).String()

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      http.HandlerFunc(app.httpRequestHandler),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		<-ctx.Done()
		ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		slog.Info(fmt.Sprintf("Shutting down HTTP server on %s", address))

		if err := srv.Shutdown(ctxShutdown); err != nil {
			slog.Error("Error shutting down HTTP server on %s - %s", address, err.Error())
		}
	}()

	slog.Info(fmt.Sprintf("Starting HTTP server on %s", address))

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		slog.Error(fmt.Sprintf("Error starting HTTP server on %s - %s", address, err.Error()))
	}

	slog.Info(fmt.Sprintf("HTTP server on %s shut down", address))
}

func (app *application) httpRequestHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug(fmt.Sprintf("HTTP request received from %s: %s %s%s", r.RemoteAddr, r.Method, r.Host, r.URL.Path))

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
