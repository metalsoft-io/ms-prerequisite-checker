package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"metalsoft.io/prerequisite-check/certs"
)

func (app *application) startHTTPSServer(ctx context.Context, port int) {
	defer app.wg.Done()

	var err error

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		MinVersion:       tls.VersionTLS13,
		Certificates:     make([]tls.Certificate, 1),
	}
	certPEMBlock, err := certs.GetCert("cert.pem")
	if err != nil {
		app.logger.Error().Err(err).Msg("Error loading TLS certificate")
		return
	}
	keyPEMBlock, err := certs.GetCert("key.pem")
	if err != nil {
		app.logger.Error().Err(err).Msg("Error loading TLS key")
		return
	}
	tlsConfig.Certificates[0], err = tls.X509KeyPair(certPEMBlock, keyPEMBlock)
	if err != nil {
		app.logger.Error().Err(err).Msg("Error loading TLS certificate and key")
		return
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      http.HandlerFunc(app.httpsRequestHandler),
		TLSConfig:    tlsConfig,
		ErrorLog:     log.New(app.logger, "", 0),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		<-ctx.Done()
		ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		app.logger.Info().Str("address", srv.Addr).Msg("Shutting down HTTPS server")

		if err := srv.Shutdown(ctxShutdown); err != nil {
			app.logger.Error().Err(err).Str("address", srv.Addr).Msg("Error shutting down HTTPS server")
		}
	}()

	app.logger.Info().Str("address", srv.Addr).Msg("Starting HTTPS server")

	err = srv.ListenAndServeTLS("", "")
	if !errors.Is(err, http.ErrServerClosed) {
		app.logger.Error().Err(err).Str("address", srv.Addr).Msg("Error starting HTTPS server")
	}

	app.logger.Info().Str("address", srv.Addr).Msg("HTTPS server shut down")
}

func (app *application) httpsRequestHandler(w http.ResponseWriter, r *http.Request) {
	app.logger.Info().Str("host", r.Host).Str("remote", r.RemoteAddr).Str("method", r.Method).Str("path", r.URL.Path).Msg("HTTPS request received")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
