package main

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"
)

func (app *application) startTCPServer(ctx context.Context, port int) {
	defer app.wg.Done()

	address := fmt.Sprintf(":%d", port)

	app.logger.Info().Str("address", address).Msg("Starting TCP server")

	ln, err := net.Listen("tcp", address)
	if err != nil {
		app.logger.Error().Err(err).Str("address", address).Msg("Error starting TCP server")
		return
	}
	defer ln.Close()

	go func() {
		<-ctx.Done()

		app.logger.Info().Str("address", address).Msg("Shutting down TCP server")

		if err := ln.Close(); err != nil {
			app.logger.Error().Err(err).Str("address", address).Msg("Error shutting down TCP server")
		}
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				app.logger.Info().Str("address", address).Msg("TCP server shut down")
				return
			}
			app.logger.Error().Err(err).Str("address", address).Msg("Could not accept TCP connection")
			time.Sleep(5 * time.Second)
			continue
		}
		go app.tcpConnectionHandler(ctx, conn)
	}
}

func (app *application) tcpConnectionHandler(ctx context.Context, socket net.Conn) {
	defer socket.Close()
	app.logger.Info().Msgf("Processing TCP connection from %s", socket.RemoteAddr())

	socket.SetReadDeadline(time.Now().Add(5 * time.Second))

	data := make([]byte, 1024)
	bytesRead, err := socket.Read(data)
	if err != nil {
		app.logger.Error().Err(err).Msg("Error reading from TCP connection")
		return
	}

	app.logger.Info().Msgf("Data received from TCP connection: %s", string(data[:bytesRead]))

	bytesWritten, err := socket.Write([]byte("PONG"))
	if err != nil {
		app.logger.Error().Err(err).Msg("Error writing to TCP connection")
		return
	}

	app.logger.Info().Msgf("Data written to TCP connection: %d bytes", bytesWritten)
}
