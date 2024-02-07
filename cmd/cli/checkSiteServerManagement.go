package main

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
)

func checkSiteServerManagement(ctx context.Context, endCh chan<- string, app *application, args map[string]string) {
	slog.Info("Starting Site Controller server management check", "arguments", args)

	serverVendor := strings.ToLower(args["vendor"])
	bmcIP := args["bmc-ip"]
	username := args["username"]
	password := args["password"]
	vncPort, err := strconv.Atoi(args["vnc-port"])
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to parse vnc-port argument (%s): %s", args["vnc-port"], err.Error()))
		endCh <- "Site Controller server management check failed"
		return
	}
	vncPassword := args["vnc-password"]

	errors := 0

	// Redfish - TCP port 443
	errors += app.testHTTPSConnection(ctx, bmcIP, 443)
	result, err := app.testRedfishAPI(ctx, bmcIP, 443, username, password, "/redfish/v1")
	if err != nil {
		errors++
	} else {
		data := result.(map[string]interface{})
		slog.Debug(fmt.Sprintf("Received Redfish response from %s:%d\n  RedfishVersion: %s\n  Vendor: %s",
			bmcIP,
			443,
			safeConvert(data, "RedfishVersion"),
			safeConvert(data, "Vendor")))
	}

	// SSH - TCP port 22
	errors += app.testSSHConnection(ctx, bmcIP, 22, username, password)

	// IPMI - UDP port 623
	errors += app.testIPMIConnection(ctx, bmcIP, 623, username, password)

	if serverVendor == "dell" && vncPassword != "" {
		// Dell iDRAC VNC - TCP port 5901
		errors += app.testVNCConnection(ctx, bmcIP, vncPort, vncPassword)
	}

	// TODO: Add test for virtual media mounting

	if errors > 0 {
		slog.Error(fmt.Sprintf("Site Controller server management check detected %d problems", errors))
	} else {
		slog.Info("Site Controller server management check detected no problems")
	}

	endCh <- "Site Controller server management check completed"
}
