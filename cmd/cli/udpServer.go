package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/netip"
	"strings"
	"time"
)

func (app *application) startUDPServer(ctx context.Context, ip netip.Addr, port uint16) {
	defer app.wg.Done()

	address := netip.AddrPortFrom(ip, port).String()

	slog.Info(fmt.Sprintf("Starting UDP server on port %s", address))

	ln, err := net.ListenPacket("udp", net.UDPAddrFromAddrPort(netip.AddrPortFrom(ip, port)).String())
	if err != nil {
		slog.Error(fmt.Sprintf("Error starting UDP server on %s - %s", address, err.Error()))
		return
	}
	defer ln.Close()

	go func() {
		<-ctx.Done()

		slog.Info(fmt.Sprintf("Shutting down UDP server on %s", address))

		if err := ln.Close(); err != nil {
			slog.Error(fmt.Sprintf("Error shutting down UDP server on %s - %s", address, err.Error()))
		}
	}()

	buffer := make([]byte, 1024)

	for {
		bytesRead, conn, err := ln.ReadFrom(buffer)
		if err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				slog.Info(fmt.Sprintf("UDP server on %s shut down", address))
				return
			}
			slog.Error(fmt.Sprintf("Could not read UDP packet on %s - %s", address, err.Error()))
			time.Sleep(5 * time.Second)
			continue
		}

		slog.Debug(fmt.Sprintf("Data received with UDP packet: %s", string(buffer[:bytesRead])))

		_, err = ln.WriteTo([]byte("PONG"), conn)
		if err != nil {
			slog.Error(fmt.Sprintf("Error writing UDP packet on %s - %s", address, err.Error()))
			continue
		}

		slog.Debug("Wrote UDP packet: PONG")
	}
}
