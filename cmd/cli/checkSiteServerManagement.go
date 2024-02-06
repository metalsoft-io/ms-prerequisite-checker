package main

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
)

func checkSiteServerManagement(ctx context.Context, endCh chan<- string, app *application, args map[string]string) {
	slog.Info("Starting Site Controller server management check", "arguments", args)

	serverVendor := strings.ToLower(args["vendor"])
	bmcIP := args["bmc-ip"]
	username := args["username"]
	password := args["password"]

	errors := 0

	// Redfish - TCP port, 443
	errors += app.testHTTPSConnection(ctx, bmcIP, 443)
	response, err := app.testRedfishAPI(ctx, bmcIP, 443, username, password, "/redfish/v1")
	if err != nil {
		errors++
	} else {
		slog.Debug(fmt.Sprintf("Read %d bytes from Redfish %s:%d - %s", len(response), bmcIP, 443, string(response)))
	}

	// SSH - TCP port 22
	errors += app.testSSHConnection(ctx, bmcIP, 22, username, password)

	if serverVendor == "dell" {
		// Dell iDRAC VNC - TCP port 5901
		// TODO: Open the VNC port before performing the test
		errors += app.testTCPConnection(ctx, bmcIP, 5901)
	}

	// IPMI - UDP port 623
	errors += app.testUDPConnection(ctx, bmcIP, 623)

	if errors > 0 {
		slog.Error(fmt.Sprintf("Site Controller server management check detected %d problems", errors))
	} else {
		slog.Info("Site Controller server management check detected no problems")
	}

	endCh <- "Site Controller server management check completed"
}
