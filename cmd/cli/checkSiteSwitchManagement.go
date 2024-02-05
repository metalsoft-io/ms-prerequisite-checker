package main

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
)

func checkSiteSwitchManagement(ctx context.Context, endCh chan<- string, app *application, args map[string]string) {
	slog.Info("Starting Site Controller switch management check", "arguments", args)

	switchNos := strings.ToLower(args["nos"])
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

	if switchNos == "junos" {
		// NETCONF/SSH - TCP port 830
		errors += app.testSSHConnection(ctx, switchIP, 830, username, password)
	}

	if errors > 0 {
		slog.Error(fmt.Sprintf("Site Controller switch management check detected %d problems", errors))
	} else {
		slog.Info("Site Controller switch management check detected no problems")
	}

	endCh <- "Site Controller switch management check completed"
}
