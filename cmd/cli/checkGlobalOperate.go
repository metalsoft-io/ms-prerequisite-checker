package main

import (
	"context"
)

func checkGlobalOperate(ctx context.Context, endCh chan<- string, app *application, args map[string]string) {
	app.logger.Info().Msgf("Starting Global Controller operation check with arguments (%v)", args)

	errors := 0

	// Metalsoft image registry → TCP port 443 registry.metalsoft.dev
	errors += app.testHTTPSConnection(ctx, "registry.metalsoft.dev", 443)

	// Metalsoft assets repo → TCP ports 80,443 repo.metalsoft.io
	errors += app.testHTTPConnection(ctx, "repo.metalsoft.io", 80)
	errors += app.testHTTPSConnection(ctx, "repo.metalsoft.io", 443)

	if errors > 0 {
		app.logger.Error().Msgf("The Global Controller operation check detected %d problems", errors)
	}

	endCh <- "Global Controller operation check completed"
}
