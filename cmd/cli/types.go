package main

type AgentRegistrationRequest struct {
	AgentId      string       `json:"agent_id"`
	AgentType    string       `json:"agent_type"`
	AgentVersion string       `json:"agent_version"`
	DatacenterId string       `json:"datacenter_id"`
	SharedSecret string       `json:"shared_secret"`
	Capabilities Capabilities `json:"capabilities"`
}

type Capabilities struct {
	HttpProxyEnabled          bool `json:"http_proxy_enabled"`
	InBandHttpProxyEnabled    bool `json:"inband_http_proxy_enabled"`
	FileTransferEnabled       bool `json:"file_transfer_enabled"`
	InBandFileTransferEnabled bool `json:"inband_file_transfer_enabled"`
	SwitchSubscriptionEnabled bool `json:"switch_subscription_enabled"`
	CommandExecutionEnabled   bool `json:"command_execution_enabled"`
	VncEnabled                bool `json:"vnc_enabled"`
	SpiceEnabled              bool `json:"spice_enabled"`
}
