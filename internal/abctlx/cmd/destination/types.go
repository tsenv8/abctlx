package destination

type CreateDestinationFlags struct {
	Name       string
	ConfigType string
}

type UpdateDestinationFlags struct {
	DestName     string
	Name         string
	ConfigType   string
	Host         string
	Port         string
	Database     string
	Username     string
	Password     string
	TunnelMethod string
}

type CreateDestinationRequest struct {
	Name          string                            `json:"name"`
	WorkspaceId   string                            `json:"workspaceId"`
	Configuration DestinationConfigurationParameter `json:"configuration"`
}

type DestinationConfigurationParameter struct {
	Host            string                `json:"host,omitempty"`
	Port            string                `json:"port,omitempty"`
	Database        string                `json:"database,omitempty"`
	Protocol        string                `json:"protocol,omitempty"`
	Username        string                `json:"username,omitempty"`
	Password        string                `json:"password,omitempty"`
	TunnelMethod    TunnelMethodParameter `json:"tunnel_method"`
	DestinationType string                `json:"destinationType,omitempty"`
}

type UpdateDestinationRequest struct {
	Name          string                             `json:"name,omitempty"`
	Configuration *DestinationConfigurationParameter `json:"configuration"`
}

type ListDestinationResponse struct {
	Previous string
	Next     string
	Data     []DestinationData
}

type DestinationData struct {
	DestinationId   string `json:"destinationId"`
	Name            string `json:"name"`
	DestinationType string `json:"destinationType"`
	WorkspaceId     string `json:"workspaceId"`
	Configuration   any    `json:"configuration,omitempty"`
}

type TunnelMethodParameter struct {
	TunnelMethod string `json:"tunnel_method"`
}
