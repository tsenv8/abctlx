package airbyte

// API
const HEADER_CONTENT_TYPE = "application/json"
const REQUEST_FAIL = "API Request Failed"
const REQUEST_SUCCESS = "API Request Successful"
const JSON_MARSHAL_FAIL = "JSON Marshal Fail"
const JSON_UNMARSHAL_FAIL = "JSON Unmarshal Fail"

// API ENDPOINTS
const GENERATE_ACCESS_TOKEN_ENDPOINT = "/v1/applications/token"
const HEALTH_CHECK_ENDPOINT = "/v1/health"
const LIST_WORKSPACES_ENDPOINT = "/v1/workspaces"
const SOURCES_ENDPOINT = "/v1/sources"
const DESTINATION_ENDPOINT = "/v1/destinations/"
const CONNECTION_ENDPOINT = "/v1/connections"
