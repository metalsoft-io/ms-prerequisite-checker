package main

import (
	"context"
	"strings"
)

func checkSiteServerManagement(ctx context.Context, endCh chan<- string, app *application, args map[string]string) {
	app.logger.Info().Msgf("Starting Site Controller server management check with arguments (%v)", args)

	serverVendor := strings.ToLower(args["vendor"])
	bmcIP := args["bmc-ip"]
	username := args["username"]
	password := args["password"]

	errors := 0

	// Redfish - TCP port, 443
	errors += app.testHTTPSConnection(ctx, bmcIP, 443)

	// SSH - TCP port 22
	errors += app.testSSHConnection(ctx, bmcIP, 22, username, password)

	if serverVendor == "dell" {
		// Dell iDRAC VNC - TCP port 5900
		errors += app.testHTTPConnection(ctx, bmcIP, 5900)
	}

	// IPMI - UDP port 623
	errors += app.testUDPConnection(ctx, bmcIP, 623)

	if errors > 0 {
		app.logger.Error().Msgf("The Site Controller server management check detected %d problems", errors)
	}

	endCh <- "Site Controller server management check completed"
}
