package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/netip"
)

func runSiteService(ctx context.Context, endCh chan<- string, app *application, args map[string]string) {
	slog.Info("Starting Site Controller mock service", "arguments", args)

	var listenIP netip.Addr
	strListenIP, ok := args["listen-ip"]
	if !ok {
		listenIP = netip.IPv4Unspecified()
	} else {
		var err error
		listenIP, err = netip.ParseAddr(strListenIP)
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to parse listen-ip argument (%s): %s", strListenIP, err.Error()))
			endCh <- "Site Controller mock service failed"
			return
		}
	}

	// DNS: UDP port 53
	app.wg.Add(1)
	go app.startDHCPServer(ctx, listenIP, 53)
}
