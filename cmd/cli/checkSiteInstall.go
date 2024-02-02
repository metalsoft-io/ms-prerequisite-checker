package main

import (
	"context"
)

func checkSiteInstall(ctx context.Context, endCh chan<- string, app *application, args map[string]string) {
	app.logger.Info().Msgf("Starting Site Controller installation check with arguments (%v)", args)

	errors := 0

	errors += app.testLink(ctx, "https://registry.metalsoft.dev")
	errors += app.testLink(ctx, "http://repo.metalsoft.io/")
	errors += app.testLink(ctx, "https://repo.metalsoft.io/")

	if errors > 0 {
		app.logger.Error().Msgf("The Site Controller installation check detected %d problems", errors)
	}

	endCh <- "Site Controller installation check completed"
}
