package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/netip"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/server4"
)

func (app *application) startDHCPServer(ctx context.Context, ip netip.Addr, port uint16) {
	defer app.wg.Done()

	address := &net.UDPAddr{
		IP:   net.ParseIP(ip.String()),
		Port: int(port),
	}

	slog.Info(fmt.Sprintf("Starting DHCP server on %s", address))

	srv, err := server4.NewServer("", address, dhcpHandler)
	if err != nil {
		slog.Error(fmt.Sprintf("Error starting DHCP server on %s - %s", address, err.Error()))
		return
	}

	go func() {
		<-ctx.Done()

		slog.Info(fmt.Sprintf("Shutting down DHCP server on %s", address))

		if err := srv.Close(); err != nil {
			slog.Error(fmt.Sprintf("Error shutting down DHCP server on %s - %s", address, err.Error()))
		}
	}()

	err = srv.Serve()
	if !errors.Is(err, net.ErrClosed) {
		slog.Error(fmt.Sprintf("Error starting DHCP server on %s - %s", address, err.Error()))
		return
	}

	slog.Info(fmt.Sprintf("DHCP server on %s shut down", address))
}

func dhcpHandler(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) {
	slog.Info(fmt.Sprintf("Received DHCP packet from %s - %s", peer.String(), m.Summary()))
}
