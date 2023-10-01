package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"plugin"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

// APITarget represents the target API configuration.
type APITarget struct {
	URL     string                 `json:"url"`
	Method  string                 `json:"method"`
	Headers map[string]interface{} `json:"headers"`
}

// APIEndpoint represents an API mapping configuration.
type APIEndpoint struct {
	Name            string                 `json:"name"`
	Source          APITarget              `json:"source"`
	Target          APITarget              `json:"target"`
	RequestMapping  map[string]interface{} `json:"requestMapping"`
	ResponseMapping map[string]interface{} `json:"responseMapping"`
}

// Configuration represents the entire JSON configuration.
type Configuration struct {
	APIMappings   []APIEndpoint  `json:"apiMappings"`
	PluginConfigs []PluginConfig `json:"pluginConfigs"`
}

type PluginConfig struct {
	Name         string `json:"name"`
	Path         string `json:"path"`
	InstanceName string `json:"instanceName"`
}

// PluginInterface is the interface for custom plugin functions.
type PluginInterface interface {
	Execute(args ...interface{}) interface{}
}

var requestBodyJSON gjson.Result
var responseBodyJSON gjson.Result
var config Configuration

// HandleAPIRequest handles incoming requests based on the API mapping configuration.
func HandleAPIRequest(w http.ResponseWriter, r *http.Request, endpoint APIEndpoint) {
	fmt.Println("Received request for API:", endpoint.Name)

	// Translate query parameters
	for key, value := range endpoint.RequestMapping["queryParam"].(map[string]interface{}) {
		if val, ok := value.(string); ok {
			parts := strings.SplitN(val, "|", 2)
			if len(parts) == 2 {
				srcType := parts[0]
				srcValue := parts[1]

				switch srcType {
				case "src:query":
					r.URL.Query().Set(key, r.URL.Query().Get(srcValue))
				case "src:static":
					r.URL.Query().Set(key, srcValue)
				}
			}
		}
	}

	rBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	requestBodyJSON = gjson.Parse(string(rBody))

	// Convert request body to JSON
	reqBody := mapData(endpoint.RequestMapping["requestBody"], r, r.Header)
	requestBody, err := json.Marshal(reqBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	code, err := performTargetRequest(endpoint.Target, requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Mapping the response based on HTTP status code.
	var responseMappingKey string

	switch code {
	case 200:
		responseMappingKey = "200"
	case 400:
		responseMappingKey = "400"
	case 401:
		responseMappingKey = "401"
	case 403:
		responseMappingKey = "403"
	case 404:
		responseMappingKey = "404"
	case 409:
		responseMappingKey = "409"
	case 500:
		responseMappingKey = "500"
	case 503:
		responseMappingKey = "503"
	case 504:
		responseMappingKey = "504"
	default:
		// Handle any other status codes if needed.
		http.Error(w, "Unsupported status code", http.StatusNotImplemented)
		return
	}

	// Get the response mapping for the matched status code.
	responseMapping := endpoint.ResponseMapping[responseMappingKey]

	// Initialize an empty response map.
	response := mapData(responseMapping, r, r.Header)

	// Convert response to JSON and send it.
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonResponse)
}

// performTargetRequest performs the HTTP request to the target API.
func performTargetRequest(target APITarget, reqBodyJSON []byte) (int, error) {
	// Create an HTTP client.
	client := &http.Client{}

	// Prepare the request based on the target config.uration.
	req, err := http.NewRequest(target.Method, target.URL, strings.NewReader(string(reqBodyJSON)))
	if err != nil {
		return 0, err
	}

	// Copy headers from the target configuration to the request.
	for key, value := range target.Headers {
		switch val := value.(type) {
		case string:
			req.Header.Set(key, val)
		}
	}

	fmt.Println("Request:")
	dump, err := httputil.DumpRequest(req, true)
	if err != nil {
		fmt.Println("Error dumping request:", err)
	}

	fmt.Println(string(dump))

	// Perform the HTTP request.
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	code := resp.StatusCode

	fmt.Println("RESPONSE CODE:", code)

	// Read the response body.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// Unmarshal the response JSON into a gjson.Result.
	responseBodyJSON = gjson.Parse(string(body))

	return code, nil
}

func mapData(dataMapping interface{}, r *http.Request, header http.Header) interface{} {
	switch mapping := dataMapping.(type) {
	case string:
		// Handle string data mapping
		result := handleStringDataMapping(mapping, r, header)
		return result
	case map[string]interface{}:
		result := make(map[string]interface{})
		for key, value := range mapping {
			switch val := value.(type) {
			case string:
				// Handle string data mapping
				result[key] = handleStringDataMapping(val, r, header)
			case map[string]interface{}:
				// Handle nested fields recursively
				nestedData := mapData(val, r, header)
				result[key] = nestedData
			case []interface{}:
				// Handle arrays recursively
				var arrayData []interface{}
				for _, arrayVal := range val {
					if arrayValMap, isArrayMap := arrayVal.(map[string]interface{}); isArrayMap {
						nestedData := mapData(arrayValMap, r, header)
						arrayData = append(arrayData, nestedData)
					} else {
						arrayData = append(arrayData, arrayVal)
					}
				}
				result[key] = arrayData
			default:
				result[key] = val
			}
		}
		return result
	default:
		// Handle unsupported type
		return nil
	}
}

func handleStringDataMapping(val string, r *http.Request, header http.Header) interface{} {
	parts := strings.SplitN(val, "|", 2)
	if len(parts) != 2 {
		fmt.Printf("Invalid format for data mapping: %s\n", val)
		return nil
	}

	srcType := parts[0]
	srcValue := parts[1]

	switch srcType {
	case "src:static":
		if srcValue == "true" {
			return true
		} else if srcValue == "false" {
			return false
		} else {
			intValue, err := strconv.Atoi(srcValue)
			if err == nil {
				return intValue
			}
			return srcValue
		}
	case "src:req_body":
		// Map from the request body
		if requestBodyJSON.Raw == "" {
			requestBody, err := io.ReadAll(r.Body)
			if err != nil {
				return nil
			}
			defer r.Body.Close()

			// Parse JSON data using gjson
			requestBodyJSON = gjson.ParseBytes(requestBody)
		}
		return getKeyValueReq(srcValue)
	case "src:res_body":
		// Map from the response body
		return getValueKeyRes(srcValue)
	case "src:query":
		// Map from query parameters
		if r.URL != nil {
			return r.URL.Query().Get(srcValue)
		}
	case "src:req_header":
		// Map from request headers
		if r.Header != nil {
			return r.Header.Get(srcValue)
		}
	case "src:res_header":
		// Map from response headers
		if header != nil {
			return header.Get(srcValue)
		} else {
			// Handle the case when response headers are nil
			return ""
		}
	case "src:func":
		// Split the function name and arguments using parentheses "(" and ")"
		partsX := strings.Split(srcValue, "(")
		if len(partsX) != 2 {
			fmt.Printf("Invalid function call format: %s\n", srcValue)
			return nil
		}
		pluginName := strings.TrimSpace(partsX[0])

		// Remove the ".Execute()" part from the function name if present
		if strings.HasSuffix(pluginName, ".Execute") {
			pluginName = pluginName[:len(pluginName)-len(".Execute")]
		}

		// If pluginArgsStr is not empty, split the arguments by ","
		pluginArgsStr := strings.TrimSuffix(partsX[1], ")")
		var pluginArgs []string
		if pluginArgsStr != "" {
			pluginArgs = strings.Split(pluginArgsStr, ",")
		}

		// Prepare arguments for the function call
		var args []interface{}
		for _, arg := range pluginArgs {
			trimmedArg := strings.TrimSpace(arg)
			if strings.HasPrefix(trimmedArg, "src:") {
				argValue := mapData(trimmedArg, r, header)
				args = append(args, argValue)
			} else {
				// Handle as a constant argument
				args = append(args, trimmedArg)
			}
		}
		return callCustomFunction(pluginName, args...)
	default:
		// Handle unsupported srcType
		return ""
	}

	// Handle unknown srcType
	return nil
}

func getKeyValueReq(key string) interface{} {
	result := requestBodyJSON.Get(key)
	if !result.Exists() {
		return nil
	}
	return result.Value()
}

func getValueKeyRes(key string) interface{} {
	result := responseBodyJSON.Get(key)
	if result.Exists() {
		return result.Value()
	}
	return nil
}

func main() {
	// Open and read the JSON configuration file.
	configFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println("Failed to open config.json:", err)
		return
	}
	defer configFile.Close()

	// Read the contents of the file.
	configData, err := io.ReadAll(configFile)
	if err != nil {
		fmt.Println("Failed to read config.json:", err)
		return
	}

	// Parse the JSON configuration.
	err = json.Unmarshal(configData, &config)
	if err != nil {
		fmt.Println("Failed to parse config.json:", err)
		return
	}

	// Start the HTTP server on port 8082.
	fmt.Println("Listening on :8082...")
	http.ListenAndServe(":8082", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Determine which API endpoint to use based on the request path or other criteria.
		var matchedEndpoint APIEndpoint

		// Iterate through all API mappings to find the appropriate endpoint.
		for _, endpoint := range config.APIMappings {
			if r.URL.Path == endpoint.Source.URL && r.Method == endpoint.Source.Method {
				matchedEndpoint = endpoint
				break
			}
		}

		// If no matching endpoint is found, return a 404 error.
		if matchedEndpoint.Name == "" {
			http.Error(w, "No matching API endpoint found", http.StatusNotFound)
			return
		}

		//body, err := ioutil.ReadAll(r.Body)
		//if err != nil {
		//	http.Error(w, "Error reading request body", http.StatusBadRequest)
		//	return
		//}
		//
		//fmt.Println("DEBUG1", string(body))

		// Handle the API request using the matched endpoint.
		HandleAPIRequest(w, r, matchedEndpoint)
	}))
}

func callCustomFunction(pluginName string, args ...interface{}) interface{} {
	var matchedPlugin PluginConfig
	for _, p := range config.PluginConfigs {
		if pluginName == p.Name {
			matchedPlugin = p
			break
		}
	}

	// If no matching plugin name is found, return nil.
	if matchedPlugin.Name == "" {
		return nil
	}

	p, err := plugin.Open(matchedPlugin.Path)
	if err != nil {
		fmt.Printf("Error loading plugin %s: %v\n", pluginName, err)
		return nil
	}

	// Lookup instance name
	fnSymbol, err := p.Lookup(matchedPlugin.InstanceName)
	if err != nil {
		fmt.Printf("Error looking up function %s: %v\n", pluginName, err)
		return nil
	}

	// Check if the function implements the common interface
	pluginInstance, ok := fnSymbol.(PluginInterface)
	if !ok {
		fmt.Println("Invalid function type from module symbol")
		return nil
	}

	// Call the Execute method and get the result
	result := pluginInstance.Execute(args...)
	return result
}
