package main

import (
	"context"
)

type argumentDetails struct {
	key          string
	description  string
	required     bool
	defaultValue string
}

type argumentsList []argumentDetails

type commandDetails struct {
	key         string
	description string
	arguments   argumentsList
	handler     func(context.Context, chan<- string, *application, map[string]string)
}

type commandsList []commandDetails

var commands = commandsList{
	{
		key:         "global-install",
		description: "Checks prerequisites for installing global controller.",
		arguments: argumentsList{
			{
				key:          "ms-repo",
				description:  "URL of the repository with MetalSoft packages.",
				required:     false,
				defaultValue: "http://repo.metalsoft.io",
			},
			{
				key:          "ms-repo-secure",
				description:  "Secure URL of the repository with MetalSoft packages.",
				required:     false,
				defaultValue: "https://repo.metalsoft.io",
			},
			{
				key:          "ms-registry",
				description:  "URL of the MetalSoft registry.",
				required:     false,
				defaultValue: "https://registry.metalsoft.dev",
			},
		},
		handler: checkGlobalInstall,
	},
	{
		key:         "global-operate",
		description: "Checks prerequisites for operating global controller.",
		arguments:   argumentsList{},
		handler:     checkGlobalOperate,
	},
	{
		key:         "global-service",
		description: "Runs global controller emulation service.",
		arguments: argumentsList{
			{
				key:          "listen-ip",
				description:  "IP address to listen on.",
				required:     false,
				defaultValue: "0.0.0.0",
			},
		},
		handler: runGlobalService,
	},
	{
		key:         "site-install",
		description: "Checks prerequisites for installing site controller.",
		arguments:   argumentsList{},
		handler:     checkSiteInstall,
	},
	{
		key:         "site-operate",
		description: "Checks prerequisites for operating site controller.",
		arguments: argumentsList{
			{
				key:         "global-controller-hostname",
				description: "IP address or hostname of the global controller.",
				required:    true,
			},
			{
				key:         "nfs-server",
				description: "NFS server for use by the site controller.",
				required:    false,
			},
		},
		handler: checkSiteOperate,
	},
	{
		key:         "site-manage-switch",
		description: "Checks site controller access to manage switch.",
		arguments: argumentsList{
			{
				key:         "nos",
				description: "The switch NOS - one of (OS10, SONiC, JunOS, Cisco).",
				required:    true,
			},
			{
				key:         "management-ip",
				description: "IP address of the switch management port.",
				required:    true,
			},
			{
				key:         "username",
				description: "Username of the switch management admin user.",
				required:    true,
			},
			{
				key:         "password",
				description: "Password of the switch management admin user.",
				required:    true,
			},
		},
		handler: checkSiteSwitchManagement,
	},
	{
		key:         "site-manage-server",
		description: "Checks site controller access to manage server.",
		arguments: argumentsList{
			{
				key:         "vendor",
				description: "The server vendor - one of (Dell, HP, Lenovo).",
				required:    true,
			},
			{
				key:         "bmc-ip",
				description: "IP address of the server BMC interface.",
				required:    true,
			},
			{
				key:         "username",
				description: "Username of the server BMC admin user.",
				required:    true,
			},
			{
				key:         "password",
				description: "Password of the server BMC admin user.",
				required:    true,
			},
			{
				key:          "vnc-port",
				description:  "VNC service port.",
				required:     false,
				defaultValue: "5901",
			},
			{
				key:         "vnc-password",
				description: "VNC password.",
				required:    false,
			},
			{
				key:         "iso-link",
				description: "Link to an ISO to test mounting virtual media.",
				required:    false,
			},
		},
		handler: checkSiteServerManagement,
	},
	{
		key:         "site-service",
		description: "Runs global controller emulation service.",
		arguments: argumentsList{
			{
				key:          "listen-ip",
				description:  "IP address to listen on.",
				required:     false,
				defaultValue: "0.0.0.0",
			},
		},
		handler: runSiteService,
	},
}
