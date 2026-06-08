package airbyte

type AbctlxResponse struct {
	Msg      string
	Body     []byte
	Data     any
	Endpoint *string
	Status   int
}

// sources
type CreateSourceRequest struct {
	Name          string                         `json:"name"`
	WorkspaceId   string                         `json:"workspaceId"`
	Configuration PostgresConfigurationParameter `json:"configuration"`
}

type CreateSourceResponse struct {
	SourceId      string `json:"sourceId"`
	Name          string `json:"name"`
	SourceType    string `json:"sourceType"`
	WorkspaceId   string `json:"workspaceId"`
	Configuration any    `json:"configuration"`
}

type UpdateSourceRequest struct {
	SourceName    string              `json:"name,omitempty"`
	WorkspaceId   string              `json:"workspaceId,omitempty"`
	Configuration *UpdateSourceFields `json:"configuration,omitempty"`
}

type UpdateSourceFields struct {
	SourceType        string                         `json:"sourceType"`
	HostName          string                         `json:"host,omitempty"`
	Port              int                            `json:"port,omitempty"`
	DBName            string                         `json:"database,omitempty"`
	Username          string                         `json:"username,omitempty"`
	Password          string                         `json:"password,omitempty"`
	Schemas           []string                       `json:"schemas,omitempty"`
	SSLMode           *SSLModeParameter              `json:"ssl_mode"`
	ReplicationMethod *CDCReplicationMethodParameter `json:"replication_method"`
	TunnelMethod      *TunnelMethodParameter         `json:"tunnel_method"`
}

type UpdateSourceResponse struct {
	SourceId    string `json:"sourceId"`
	Name        string `json:"name"`
	SourceType  string `json:"sourceType"`
	WorkspaceId string `json:"workspaceId"`
}

type CreateSourceParams struct {
	Name            string
	HostName        string
	DBName          string
	Username        string
	Password        string
	ReplicationSlot string
	PublicationName string
	Schemas         []string
	Port            int
}

type ListSourcesResponse struct {
	Next     string            `json:"next"`
	Previous string            `json:"previous"`
	Data     []ListSourcesData `json:"data"`
}

type ListSourcesData struct {
	SourceId    string `json:"sourceId"`
	Name        string `json:"name"`
	SourceType  string `json:"sourceType"`
	WorkspaceId string `json:"workspaceId"`
}

type PostgresConfigurationParameter struct {
	SourceType        string                        `json:"sourceType"`
	Host              string                        `json:"host"`
	Port              int                           `json:"port"`
	Database          string                        `json:"database"`
	Schemas           []string                      `json:"schemas"`
	Username          string                        `json:"username"`
	Password          string                        `json:"password"`
	SSlMode           *SSLModeParameter             `json:"ssl_mode"`
	ReplicationMethod CDCReplicationMethodParameter `json:"replication_method"`
	TunnelMethod      TunnelMethodParameter         `json:"tunnel_method"`
}

type CDCReplicationMethodParameter struct {
	Method          string `json:"method,omitempty"`
	Plugin          string `json:"plugin,omitempty"`
	ReplicationSlot string `json:"replication_slot,omitempty"`
	Publication     string `json:"publication,omitempty"`
}

type SSLModeParameter struct {
	Mode string `json:"mode"`
}

type TunnelMethodParameter struct {
	TunnelMethod string `json:"tunnel_method"`
}

// applications
type GenerateAccessTokenRequest struct {
	ClientId  string `json:"client_id"`
	ClientKey string `json:"client_secret"`
}

type GenerateAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type HealthCheckResponse struct {
	Status bool
}

// workspaces
type ListWorkspacesResponse struct {
	Next     string          `json:"next"`
	Previous string          `json:"previous"`
	Data     []WorkspaceData `json:"data"`
}

type WorkspaceData struct {
	WorkspaceId   string
	Name          string
	DataResidency string
}
