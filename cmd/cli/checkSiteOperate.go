package main

import (
	"context"
	"fmt"
	"log/slog"
)

func checkSiteOperate(ctx context.Context, endCh chan<- string, app *application, args map[string]string) {
	slog.Info("Starting Site Controller operation check", "arguments", args)

	globalControllerHostname := args["global-controller-hostname"]

	errors := 0

	// Metalsoft Controller HTTP port 80
	errors += app.testHTTPConnection(ctx, globalControllerHostname, 80)

	// Metalsoft Controller HTTPS port 443
	errors += app.testHTTPSConnection(ctx, globalControllerHostname, 443)

	// Metalsoft Controller WebSocketSecure port 9010 - tunnel control messages
	errors += app.testWebSocketConnection(ctx, globalControllerHostname, 443, "/tunnel-ctrl", true)

	// // Metalsoft Controller HTTPS port 9010 - HTTP Proxy
	// errors += app.testHTTPSConnection(ctx, globalControllerHostname, 9010)

	// Metalsoft Controller TCP port 9091 - TCP Proxy
	errors += app.testTCPConnection(ctx, globalControllerHostname, 9091)

	// Metalsoft Controller - UDP ports 53
	errors += app.testUDPConnection(ctx, globalControllerHostname, 53, "dns")

	if nfs := args["nfs-server"]; nfs != "" {
		// NFS server - TCP/UDP ports 111 and 2049
		errors += app.testTCPConnection(ctx, nfs, 111)
		errors += app.testUDPConnection(ctx, nfs, 111, "nfs")
		errors += app.testTCPConnection(ctx, nfs, 2049)
		errors += app.testUDPConnection(ctx, nfs, 2049, "nfs")
	}

	if errors > 0 {
		slog.Error(fmt.Sprintf("Site Operation test detected %d problems", errors))
	} else {
		slog.Info("Site Operation test detected no problems")
	}

	endCh <- "Site Controller operation check completed"
}
