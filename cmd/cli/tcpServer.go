package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net"
	"net/netip"
	"strings"
	"time"
)

func (app *application) startTCPServer(ctx context.Context, ip netip.Addr, port uint16, encrypted bool) {
	defer app.wg.Done()

	address := netip.AddrPortFrom(ip, port).String()

	slog.Info(fmt.Sprintf("Starting TCP server on port %s", address))

	var err error
	var ln net.Listener
	if encrypted {
		var cert tls.Certificate
		cert, err = tls.X509KeyPair(getServerCertificate(), getServerPrivateKey())
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to load server certificate: %s", err.Error()))
			return
		}

		ln, err = tls.Listen("tcp", address, &tls.Config{Certificates: []tls.Certificate{cert}})
	} else {
		ln, err = net.Listen("tcp", address)
	}
	if err != nil {
		slog.Error(fmt.Sprintf("Error starting TCP server on %s - %s", address, err.Error()))
		return
	}
	defer ln.Close()

	go func() {
		<-ctx.Done()

		slog.Info(fmt.Sprintf("Shutting down TCP server on %s", address))

		if err := ln.Close(); err != nil {
			slog.Error(fmt.Sprintf("Error shutting down TCP server on %s - %s", address, err.Error()))
		}
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				slog.Info(fmt.Sprintf("TCP server on %s shut down", address))
				return
			}
			slog.Error(fmt.Sprintf("Could not accept TCP connection on %s - %s", address, err.Error()))
			time.Sleep(5 * time.Second)
			continue
		}
		go app.tcpConnectionHandler(ctx, conn)
	}
}

func (app *application) tcpConnectionHandler(ctx context.Context, socket net.Conn) {
	defer socket.Close()
	slog.Debug(fmt.Sprintf("Processing TCP connection from %s", socket.RemoteAddr()))

	socket.SetReadDeadline(time.Now().Add(5 * time.Second))

	data := make([]byte, 1024)
	bytesRead, err := socket.Read(data)
	if err != nil {
		slog.Error(fmt.Sprintf("Error reading from TCP connection - %s", err.Error()))
		return
	}

	slog.Debug(fmt.Sprintf("Data received from TCP connection: %s", string(data[:bytesRead])))

	bytesWritten, err := socket.Write([]byte("PONG"))
	if err != nil {
		slog.Error(fmt.Sprintf("Error writing to TCP connection - %s", err.Error()))
		return
	}

	slog.Debug(fmt.Sprintf("Data written to TCP connection: %d bytes", bytesWritten))
}
