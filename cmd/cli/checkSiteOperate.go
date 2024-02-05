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

	// Metalsoft Controller TCP ports 80/443
	errors += app.testHTTPConnection(ctx, globalControllerHostname, 80)
	errors += app.testHTTPSConnection(ctx, globalControllerHostname, 443)

	// Metalsoft Controller → TCP ports 9003,9009,9090,9091,9011,9010
	errors += app.testTCPConnection(ctx, globalControllerHostname, 9003)
	errors += app.testTCPConnection(ctx, globalControllerHostname, 9009)
	errors += app.testTCPConnection(ctx, globalControllerHostname, 9010)
	errors += app.testTCPConnection(ctx, globalControllerHostname, 9011)
	errors += app.testTCPConnection(ctx, globalControllerHostname, 9090)
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
