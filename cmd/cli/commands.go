package main

import (
	"context"
)

type argumentDetails struct {
	description  string
	required     bool
	defaultValue string
}

type commandDetails struct {
	description string
	arguments   map[string]argumentDetails
	handler     func(context.Context, chan<- string, *application, map[string]string)
}

var commands = map[string]commandDetails{
	"global-install": {
		description: "Checks prerequisites for installing global controller.",
		arguments: map[string]argumentDetails{
			"ms-repo": {
				description:  "URL of the repository with MetalSoft packages.",
				required:     false,
				defaultValue: "http://repo.metalsoft.io",
			},
			"ms-repo-secure": {
				description:  "Secure URL of the repository with MetalSoft packages.",
				required:     false,
				defaultValue: "https://repo.metalsoft.io",
			},
			"ms-registry": {
				description:  "URL of the MetalSoft registry.",
				required:     false,
				defaultValue: "https://registry.metalsoft.dev",
			},
		},
		handler: checkGlobalInstall,
	},
	"global-operate": {
		description: "Checks prerequisites for operating global controller.",
		arguments:   map[string]argumentDetails{},
		handler:     checkGlobalOperate,
	},
	"global-service": {
		description: "Runs global controller service.",
		arguments: map[string]argumentDetails{
			"listen-ip": {
				description:  "IP address to listen on.",
				required:     false,
				defaultValue: "0.0.0.0",
			},
		},
		handler: runGlobalService,
	},
	"site-install": {
		description: "Checks prerequisites for installing site controller.",
		arguments:   map[string]argumentDetails{},
		handler:     checkSiteInstall,
	},
	"site-operate": {
		description: "Checks prerequisites for operating site controller.",
		arguments: map[string]argumentDetails{
			"global-controller-hostname": {
				description: "IP address or hostname of the global controller.",
				required:    true,
			},
			"nfs-server": {
				description: "NFS server for use by the site controller.",
				required:    false,
			},
		},
		handler: checkSiteOperate,
	},
	"site-manage-switch": {
		description: "Checks site controller access to manage switch.",
		arguments: map[string]argumentDetails{
			"nos": {
				description: "The switch NOS - one of (OS10, SONiC, JunOS, Cisco).",
				required:    true,
			},
			"management-ip": {
				description: "IP address of the switch management port.",
				required:    true,
			},
			"username": {
				description: "Username of the switch management admin user.",
				required:    true,
			},
			"password": {
				description: "Password of the switch management admin user.",
				required:    true,
			},
		},
		handler: checkSiteSwitchManagement,
	},
	"site-manage-server": {
		description: "Checks site controller access to manage server.",
		arguments: map[string]argumentDetails{
			"vendor": {
				description: "The server vendor - one of (Dell, HP, Lenovo).",
				required:    true,
			},
			"bmc-ip": {
				description: "IP address of the server BMC interface.",
				required:    true,
			},
			"username": {
				description: "Username of the server BMC admin user.",
				required:    true,
			},
			"password": {
				description: "Password of the server BMC admin user.",
				required:    true,
			},
			"iso-link": {
				description: "Link to an ISO to test mounting virtual media.",
				required:    false,
			},
		},
		handler: checkSiteServerManagement,
	},
}
