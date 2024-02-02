package main

import (
	"context"
	"strings"
)

func checkSiteSwitchManagement(ctx context.Context, endCh chan<- string, app *application, args map[string]string) {
	app.logger.Info().Msgf("Starting Site Controller switch management check with arguments (%v)", args)

	switchVendor := strings.ToLower(args["vendor"])
	switchIP := args["management-ip"]
	username := args["username"]
	password := args["password"]

	errors := 0

	// HTTP - TCP port 80
	errors += app.testHTTPConnection(ctx, switchIP, 80)

	// HTTPS - TCP port 443
	errors += app.testHTTPSConnection(ctx, switchIP, 443)

	// SSH - TCP port 22
	errors += app.testSSHConnection(ctx, switchIP, 22, username, password)

	if switchVendor == "junos" {
		// NETCONF/SSH - TCP port 830
		errors += app.testSSHConnection(ctx, switchIP, 830, username, password)
	}

	if errors > 0 {
		app.logger.Error().Msgf("The Site Controller switch management check detected %d problems", errors)
	}

	endCh <- "Site Controller switch management check completed"
}
