package airbyte

type AbctlxResponse struct {
	Msg      string
	Body     []byte
	Data     any
	Endpoint *string
	Status   int
}

// connections
type ConnectionData struct {
	ConnectionId                     string                      `json:"connectionId"`
	Name                             string                      `json:"name"`
	SourceId                         string                      `json:"sourceId"`
	DestinationId                    string                      `json:"destinationId"`
	WorkspaceId                      string                      `json:"workspaceId"`
	Status                           string                      `json:"status"`
	Schedule                         ConnectionScheduleParameter `json:"schedule"`
	DataResidency                    string                      `json:"dataResidency"`
	NonBreakingSchemaUpdatesBehavior string                      `json:"nonBreakingSchemaUpdatesBehavior"`
	NamespaceDefinition              string                      `json:"namespaceDefinition"`
	NamespaceFormat                  string                      `json:"namespaceFormat"`
	Prefix                           string                      `json:"prefix"`
	Configurations                   StreamConfigurations        `json:"configurations"`
	CreatedAt                        int                         `json:"createdAt"`
}
type UpdateConnectionRequest struct {
	Name                             string                      `json:"name,omitempty"`
	Configurations                   StreamConfigurations        `json:"configurations"`
	Schedule                         ConnectionScheduleParameter `json:"schedule"`
	DataResidency                    string                      `json:"dataResidency,omitempty"`
	NamespaceDefinition              string                      `json:"namespaceDefinition,omitempty"`
	NamespaceFormat                  string                      `json:"namespaceFormat,omitempty"`
	Prefix                           string                      `json:"prefix,omitempty"`
	NonBreakingSchemaUpdatesBehavior string                      `json:"nonBreakingSchemaUpdatesBehavior,omitempty"`
	Status                           string                      `json:"status,omitempty"`
}

type CreateConnectionRequest struct {
	Name                             string                      `json:"name"`
	SourceId                         string                      `json:"sourceId"`
	Configurations                   StreamConfigurations        `json:"configurations"`
	DestinationId                    string                      `json:"destinationId"`
	Schedule                         ConnectionScheduleParameter `json:"schedule"`
	Residency                        string                      `json:"residency"`
	NamespaceDefinition              string                      `json:"namespaceDefinition"`
	Prefix                           string                      `json:"prefix,omitempty"`
	NonBreakingSchemaUpdatesBehavior string                      `json:"nonBreakingSchemaUpdatesBehavior,omitempty"`
	Status                           string                      `json:"status"`
}

type ListConnectionResponse struct {
	Previous string           `json:"previous"`
	Next     string           `json:"next"`
	Data     []ConnectionData `json:"data"`
}

type StreamConfigurations struct {
	streams []StreamConfiguration
}

type StreamConfiguration struct {
	Name           string   `json:"name"`
	SyncMode       string   `json:"syncMode"`
	CursorField    []string `json:"cursorField"`
	SelectedFields []string `json:"selectedFields"`
}
type ConnectionScheduleParameter struct {
	Type    string `json:"scheduleType"`
	CronExp string `json:"cronExpression"`
}

// destinations

type CreateDestinationFlags struct {
	Name       string
	ConfigType string
}

type UpdateDestinationFlags struct {
	DestName     *string
	Name         *string
	ConfigType   *string
	Host         *string
	Port         *int
	Database     *string
	Username     *string
	Password     *string
	TunnelMethod *string
}

type CreateDestinationRequest struct {
	Name          string                            `json:"name"`
	WorkspaceId   string                            `json:"workspaceId"`
	Configuration DestinationConfigurationParameter `json:"configuration"`
}

type DestinationConfigurationParameter struct {
	Host            string                `json:"host,omitempty"`
	Port            int                   `json:"port,omitempty"`
	Database        string                `json:"database,omitempty"`
	Username        string                `json:"username,omitempty"`
	Password        string                `json:"password,omitempty"`
	TunnelMethod    TunnelMethodParameter `json:"tunnel_method,omitempty"`
	DestinationType string                `json:"destinationType,omitempty"`
}

type UpdateDestinationRequest struct {
	Name          string                            `json:"name,omitempty"`
	Configuration DestinationConfigurationParameter `json:"configuration"`
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
	Next     string       `json:"next"`
	Previous string       `json:"previous"`
	Data     []SourceData `json:"data"`
}

type SourceData struct {
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
