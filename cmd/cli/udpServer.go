package main

import (
	"context"
	"fmt"
	"net"
	"net/netip"
	"strings"
	"time"
)

func (app *application) startUDPServer(ctx context.Context, ip netip.Addr, port uint16) {
	defer app.wg.Done()

	address := fmt.Sprintf(":%d", port)

	app.logger.Info().Str("address", address).Msg("Starting UDP server")

	ln, err := net.ListenPacket("udp", net.UDPAddrFromAddrPort(netip.AddrPortFrom(ip, port)).String())
	if err != nil {
		app.logger.Error().Err(err).Str("address", address).Msg("Error starting UDP server")
		return
	}
	defer ln.Close()

	go func() {
		<-ctx.Done()

		app.logger.Info().Str("address", address).Msg("Shutting down UDP server")

		if err := ln.Close(); err != nil {
			app.logger.Error().Err(err).Str("address", address).Msg("Error shutting down UDP server")
		}
	}()

	buffer := make([]byte, 1024)

	for {
		bytesRead, conn, err := ln.ReadFrom(buffer)
		if err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				app.logger.Info().Str("address", address).Msg("UDP server shut down")
				return
			}
			app.logger.Error().Err(err).Str("address", address).Msg("Could not read UDP packet")
			time.Sleep(5 * time.Second)
			continue
		}

		app.logger.Info().Msgf("Data received with UDP packet: %s", string(buffer[:bytesRead]))

		_, err = ln.WriteTo([]byte("PONG"), conn)
		if err != nil {
			app.logger.Error().Err(err).Str("address", address).Msg("Error writing UDP packet")
			continue
		}

		app.logger.Info().Msg("Wrote UDP packet: PONG")
	}
}
