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

	"github.com/coder/websocket"
	"metalsoft.io/prerequisite-check/certs"
)

func (app *application) startWebSocketServer(ctx context.Context, ip netip.Addr, port uint16) {
	defer app.wg.Done()

	address := netip.AddrPortFrom(ip, port).String()

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

	http.HandleFunc("/tunnel-ctrl", app.wsRequestHandler)
	http.HandleFunc("/", app.wsRequestHandlerDefault)

	srv := &http.Server{
		Addr:         address,
		Handler:      http.DefaultServeMux,
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		<-ctx.Done()
		ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		slog.Info(fmt.Sprintf("Shutting down WebSocket server on %s", address))

		if err := srv.Shutdown(ctxShutdown); err != nil {
			slog.Error("Error shutting down WebSocket server on %s - %s", address, err.Error())
		}
	}()

	slog.Info(fmt.Sprintf("Starting WebSocket server on %s", address))

	err = srv.ListenAndServeTLS("", "")
	if !errors.Is(err, http.ErrServerClosed) {
		slog.Error(fmt.Sprintf("Error starting WebSocket server on %s - %s", address, err.Error()))
		return
	}

	slog.Info(fmt.Sprintf("WebSocket server on %s shut down", address))
}

func (app *application) wsRequestHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug(fmt.Sprintf("WebSocket request received from %s: %s %s%s", r.RemoteAddr, r.Method, r.Host, r.URL.Path))

	ws, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: true})
	if err != nil {
		slog.Error(fmt.Sprintf("Error accepting WebSocket connection: %s", err.Error()))
		return
	}
	defer ws.Close(websocket.StatusNormalClosure, "OK")

	slog.Debug(fmt.Sprintf("Connected WebSocket from %s", r.RemoteAddr))
	slog.Debug(fmt.Sprintf("WebSocket request: %s %+v", r.URL, r.Header))

	messageType, message, err := ws.Read(r.Context())
	if err != nil {
		slog.Error(fmt.Sprintf("Error reading WebSocket message: %s", err.Error()))
		return
	}
	slog.Debug(fmt.Sprintf("Received WebSocket message of type %v: %+v", messageType, message))

	err = ws.Write(r.Context(), websocket.MessageText, []byte("OK"))
	if err != nil {
		slog.Error(fmt.Sprintf("Error writing WebSocket message: %s", err.Error()))
		return
	}

	slog.Debug("Sent WebSocket message")
}

func (app *application) wsRequestHandlerDefault(w http.ResponseWriter, r *http.Request) {
	slog.Debug(fmt.Sprintf("HTTPS request received from %s: %s %s%s", r.RemoteAddr, r.Method, r.Host, r.URL.Path))

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
