package main

import (
	"context"
	"fmt"
	"log/slog"
)

func checkSiteInstall(ctx context.Context, endCh chan<- string, app *application, args map[string]string) {
	slog.Info("Starting Site Controller installation check", "arguments", args)

	errors := 0

	errors += app.testLink(ctx, "https://registry.metalsoft.dev")
	errors += app.testLink(ctx, "http://repo.metalsoft.io/")
	errors += app.testLink(ctx, "https://repo.metalsoft.io/")

	if errors > 0 {
		slog.Error(fmt.Sprintf("Site Controller installation check detected %d problems", errors))
	} else {
		slog.Info("Site Controller installation check detected no problems")
	}

	endCh <- "Site Controller installation check completed"
}
