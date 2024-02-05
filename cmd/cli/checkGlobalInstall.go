package main

import (
	"context"
	"fmt"
	"log/slog"
)

func checkGlobalInstall(ctx context.Context, endCh chan<- string, app *application, args map[string]string) {
	slog.Info("Starting Global Controller installation check", "arguments", args)

	msRepo := args["ms-repo"]
	msRepoSecure := args["ms-repo-secure"]
	msRegistry := args["ms-registry"]

	errors := 0

	// http://repo.metalsoft.io 80 tcp
	errors += app.testLink(ctx, msRepo)

	// https://repo.metalsoft.io 443 tcp
	errors += app.testLink(ctx, msRepoSecure)

	// https://registry.metalsoft.dev 443 tcp
	errors += app.testLink(ctx, msRegistry)

	// 1.1.1.1 Public ICMP
	errors += app.testICMPConnection(ctx, "1.1.1.1")

	// 1.1.1.1 Public port - TCP 80
	errors += app.testLink(ctx, "http://1.1.1.1/")

	// 1.1.1.1 Public port - TCP 443
	errors += app.testLink(ctx, "https://1.1.1.1/")

	// https://downloads.dell.com - TCP 443
	errors += app.testLink(ctx, "https://downloads.dell.com/")

	// http://downloads.linux.hpe.com - TCP 80
	errors += app.testLink(ctx, "http://downloads.linux.hpe.com/")

	// https://quay.io - TCP 443
	errors += app.testLink(ctx, "https://quay.io/")

	// https://gcr.io - TCP 443
	errors += app.testLink(ctx, "https://gcr.io/")

	// https://cloud.google.com - TCP 443
	errors += app.testLink(ctx, "https://cloud.google.com/")

	// https://helm.traefik.io - TCP 443
	errors += app.testLink(ctx, "https://helm.traefik.io/")

	// https://k8s.io - TCP  443
	errors += app.testLink(ctx, "https://k8s.io/")

	// smtp.office365.com - TCP 587
	errors += app.testTCPConnection(ctx, "smtp.office365.com", 587)

	// http://archive.ubuntu.com , http://security.ubuntu.com  80 tcp -> for base OS package updates

	if errors > 0 {
		slog.Error(fmt.Sprintf("Global Controller installation check detected %d problems", errors))
	} else {
		slog.Info("The Global Controller installation check detected no problems")
	}

	endCh <- "Global Controller installation check completed"
}
