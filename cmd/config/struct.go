package config

// Configuration represents the entire JSON configuration.
type Configuration struct {
	APIMappings   []APIEndpoint  `json:"apiMappings"`
	PluginConfigs []PluginConfig `json:"pluginConfigs"`
}

// APIEndpoint represents an API mapping configuration.
type APIEndpoint struct {
	Name                string          `json:"name"`
	ResponseMappingType string          `json:"responseMappingType"`
	Source              APITarget       `json:"source"`
	Target              APITarget       `json:"target"`
	RequestMapping      RequestMapping  `json:"requestMapping"`
	ResponseMapping     ResponseMapping `json:"responseMapping"`
}

// APITarget represents the target API configuration.
type APITarget struct {
	URL     string                 `json:"url"`
	Method  string                 `json:"method"`
	Headers map[string]interface{} `json:"headers"`
}

// RequestMapping defines how to map request data.
type RequestMapping struct {
	QueryParam  map[string]interface{} `json:"queryParam,omitempty"`
	RequestBody map[string]interface{} `json:"requestBody"`
}

// ResponseMapping defines how to map response data.
type ResponseMapping struct {
	ByHTTPStatusCode ByHTTPStatusCode `json:"byHTTPStatusCode"`
	ByBodyResponse   ByBodyResponse   `json:"byBodyResponse"`
}

// ByBodyResponse defines custom response mappings.
type ByBodyResponse struct {
	Default Default                   `json:"default"`
	Custom  map[string][]BodyResponse `json:"custom"`
}

// ByHTTPStatusCode defines response mappings by HTTP status code.
type ByHTTPStatusCode struct {
	Default Default                `json:"default"`
	Custom  map[string]interface{} `json:"custom"`
}

// Default defines the default response.
type Default struct {
	Response Response `json:"response"`
}

// BodyResponse defines the structure for response mapping.
type BodyResponse struct {
	Values   []string `json:"values"`
	Response Response `json:"response"`
}

// Response defines the structure for response mapping.
type Response struct {
	HTTPStatusCode int                    `json:"http_status_code"`
	JSONBody       map[string]interface{} `json:"json_body"`
}

// PluginConfig defines the configuration for custom plugins.
type PluginConfig struct {
	Name         string `json:"name"`
	Path         string `json:"path"`
	InstanceName string `json:"instanceName"`
}
