package main

import (
	"context"
	"fmt"
	"log/slog"
)

func checkGlobalOperate(ctx context.Context, endCh chan<- string, app *application, args map[string]string) {
	slog.Info("Starting Global Controller operation check", "arguments", args)

	errors := 0

	// Metalsoft image registry → TCP port 443 registry.metalsoft.dev
	errors += app.testHTTPSConnection(ctx, "registry.metalsoft.dev", 443)

	// Metalsoft assets repo → TCP ports 80,443 repo.metalsoft.io
	errors += app.testHTTPConnection(ctx, "repo.metalsoft.io", 80)
	errors += app.testHTTPSConnection(ctx, "repo.metalsoft.io", 443)

	if errors > 0 {
		slog.Error(fmt.Sprintf("Global Controller operation check detected %d problems", errors))
	} else {
		slog.Info("Global Controller operation check detected no problems")
	}

	endCh <- "Global Controller operation check completed"
}
