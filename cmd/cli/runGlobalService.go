package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/netip"
)

func runGlobalService(ctx context.Context, endCh chan<- string, app *application, args map[string]string) {
	slog.Info("Starting Global Controller mock service", "arguments", args)

	var listenIP netip.Addr
	strListenIP, ok := args["listen-ip"]
	if !ok {
		listenIP = netip.IPv4Unspecified()
	} else {
		var err error
		listenIP, err = netip.ParseAddr(strListenIP)
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to parse listen-ip argument (%s): %s", strListenIP, err.Error()))
			endCh <- "Global Controller mock service failed"
			return
		}
	}

	// web: TCP port 80
	app.wg.Add(1)
	go app.startHTTPServer(ctx, listenIP, 80)

	// websecure: TCP port 443
	app.wg.Add(1)
	go app.startHTTPSServer(ctx, listenIP, 443)

	// event-service: TCP port 9003
	// app.wg.Add(1)
	// go app.startTCPServer(ctx, listenIP, 9003)

	// gateway-api: TCP port 9009
	// app.wg.Add(1)
	// go app.startTCPServer(ctx, listenIP, 9009)

	// tunnel control messages: TCP port 9010
	app.wg.Add(1)
	go app.startWebSocketServer(ctx, listenIP, 9010)

	// tunnel-9011: TCP port 9011
	// app.wg.Add(1)
	// go app.startTCPServer(ctx, listenIP, 9011)

	// tunnel HTTP proxy: TCP port 9090
	app.wg.Add(1)
	go app.startHTTPServer(ctx, listenIP, 9090)

	// tunnel TCP proxy: TCP port 9091
	app.wg.Add(1)
	go app.startTCPServer(ctx, listenIP, 9091)

	// power-dns: UDP port 53
	app.wg.Add(1)
	go app.startUDPServer(ctx, listenIP, 53)
}
