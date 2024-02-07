package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/netip"
	"time"

	"metalsoft.io/prerequisite-check/certs"
)

func (app *application) startHTTPSServer(ctx context.Context, ip netip.Addr, port uint16) {
	defer app.wg.Done()

	address := netip.AddrPortFrom(ip, port).String()

	var err error

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		MinVersion:       tls.VersionTLS13,
		Certificates:     make([]tls.Certificate, 1),
	}
	certPEMBlock, err := certs.GetCert("cert.pem")
	if err != nil {
		slog.Error("Error loading TLS certificate", err)
		return
	}
	keyPEMBlock, err := certs.GetCert("key.pem")
	if err != nil {
		slog.Error("Error loading TLS key", err)
		return
	}
	tlsConfig.Certificates[0], err = tls.X509KeyPair(certPEMBlock, keyPEMBlock)
	if err != nil {
		slog.Error("Error loading TLS certificate and key", err)
		return
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      http.HandlerFunc(app.httpsRequestHandler),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		<-ctx.Done()
		ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		slog.Info(fmt.Sprintf("Shutting down HTTPS server on %s", address))

		if err := srv.Shutdown(ctxShutdown); err != nil {
			slog.Error(fmt.Sprintf("Error shutting down HTTPS server on %s - %s", address, err.Error()))
		}
	}()

	slog.Info(fmt.Sprintf("Starting HTTPS server on %s", address))

	err = srv.ListenAndServeTLS("", "")
	if !errors.Is(err, http.ErrServerClosed) {
		slog.Error(fmt.Sprintf("Error starting HTTPS server on %s - %s", address, err.Error()))
		return
	}

	slog.Info(fmt.Sprintf("HTTPS server on %s shut down", address))
}

func (app *application) httpsRequestHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug(fmt.Sprintf("HTTPS request received from %s: %s %s%s", r.RemoteAddr, r.Method, r.Host, r.URL.Path))

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
