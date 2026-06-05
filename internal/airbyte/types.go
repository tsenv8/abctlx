package airbyte

type AbctlxResponse struct {
	msg      string
	body     []byte
	data     any
	endpoint *string
}

// sources
type CreateSourceRequest struct {
	Name          string                         `json:"name"`
	WorkspaceId   string                         `json:"workspaceId"`
	Configuration PostgresConfigurationParameter `json:"configuration"`
}

type CreateSourceParams struct {
	name            string
	hostName        string
	dbName          string
	username        string
	password        string
	replicationSlot string
	publicationName string
	schemas         []string
	port            int
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
	Method          string `json:"method"`
	Plugin          string `json:"plugin"`
	ReplicationSlot string `json:"replication_slot"`
	Publication     string `json:"publication"`
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

// workspaces
type ListWorkspacesResponse struct {
	Next     string
	Previous string
	Data     []WorkspaceData
}

type WorkspaceData struct {
	WorkspaceId   string
	Name          string
	DataResidency string
}
