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

	// Metalsoft Controller TCP port 9003 - event-service - unused
	// errors += app.testTCPConnection(ctx, globalControllerHostname, 9003)

	// Metalsoft Controller TCP port 9009 - gateway-api - unused
	// errors += app.testTCPConnection(ctx, globalControllerHostname, 9009)

	// Metalsoft Controller WebSocket port 9010 - tunnel control messages
	errors += app.testWebSocketConnection(ctx, globalControllerHostname, 9010, "/tunnel-ctrl")

	// Metalsoft Controller TCP port 9011 - unused
	// errors += app.testTCPConnection(ctx, globalControllerHostname, 9011)

	// Metalsoft Controller HTTP port 9090 - HTTP Proxy
	errors += app.testHTTPConnection(ctx, globalControllerHostname, 9090)

	// Metalsoft Controller TCP port 9091 - TCP Proxy
	errors += app.testTCPConnection(ctx, globalControllerHostname, 9091)

	// Metalsoft Controller - UDP ports 53
	errors += app.testUDPConnection(ctx, globalControllerHostname, 53)

	if nfs := args["nfs-server"]; nfs != "" {
		// NFS server - TCP/UDP ports 111 and 2049
		errors += app.testTCPConnection(ctx, nfs, 111)
		errors += app.testUDPConnection(ctx, nfs, 111)
		errors += app.testTCPConnection(ctx, nfs, 2049)
		errors += app.testUDPConnection(ctx, nfs, 2049)
	}

	if errors > 0 {
		slog.Error(fmt.Sprintf("Site Operation test detected %d problems", errors))
	} else {
		slog.Info("Site Operation test detected no problems")
	}

	endCh <- "Site Controller operation check completed"
}
