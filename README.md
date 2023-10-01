# api-mapping-customization-guide
Custom API Mapping Guide: Configure plugins, request/response mappings, and examples for tailored API interactions

This guide explains how to configure the API Gateway system for your project. The API Gateway acts as a middleware that handles incoming requests and routes them to the appropriate target service while performing data transformations.

## 1. Plugin Configurations

Plugin configurations define the plugins used by the API Gateway for request/response processing. Each plugin has a name, path to the plugin library file, and an instance name.

Example:
```json
"pluginConfigs": [
    {
        "name": "GenerateTransactionIDPlugin",
        "path": "plugins/generate_transaction_id_plugin.so",
        "instanceName": "GenerateTransactionIDPluginInstance"
    },
    {
        "name": "SimSwapPlugin",
        "path": "plugins/sim_swap_plugin.so",
        "instanceName": "SimSwapPluginInstance"
    }
]
```

- `name`: The name of the plugin.
- `path`: The path to the plugin library file.
- `instanceName`: The name of the plugin- instance.

---

## 2. API Mappings

API mappings define how incoming requests are mapped to target services, along with request and response data transformations.

Example:

```json
"apiMappings": [
    {
        "name": "ExampleRequestMapping",
        "source": {
            "url": "/example/request",
            "method": "POST"
        },
        "target": {
            "url": "http://localhost:8083/example/target",
            "method": "POST",
            "headers": {
                "Authorization": "Bearer YOUR_ACCESS_TOKEN"
            }
        },
        "requestMapping": {
            "queryParam": {},
            "requestBody": {
                "phoneNumber": "src:req_body|newPhoneNumber",
                "maxAge": "src:req_body|age",
                "countryCode": "src:req_body|location",
                "categoryCode": "src:req_body|productCodes.1"
            }
        },
        "responseMapping": {
            "200": {
                "sim_swap": "src:func|SimSwapPlugin.Execute(src:req_body|age,src:res_body|data.score)",
                "transaction_id": "src:func|GenerateTransactionIDPlugin.Execute()",
                "score": "src:res_body|data.score",
                "age": "src:req_body|age"
            },
            "400": {
                "status": "src:static|error",
                "message": "src:res_body|errorMessage",
                "data": "src:res_body|errorData"
            }
        }
    },
    // ... More API Mappings ...
]

```

- `name`: A unique name for the API mapping.
- `source`: Defines the source of incoming requests, including URL and HTTP method.
- `target`: Specifies the target service's URL, HTTP method, and headers for forwarding requests.
- `requestMapping`: Maps request data from the source to the target service.
- `responseMapping`: Maps response data from the target service back to the source.

---

## 3. Request Mapping Examples

### Request Mapping with Data Transformations

#### Mapping a Nested Field in the Request Body to a Query Parameter:

- Data Request Received:
```json
{
  "xyz": {
    "nested": {
      "field": "nested_value"
    }
  }
}
```

- Mapping Configuration:
```json
"param_nested_value": "src:req_body|xyz.nested.field"
```

- Data Transformed for Target:
```json
/example/target?param_nested_value=nested_value
```

#### Mapping a Specific Element of an Array in the Request Body to a Query Parameter:

- Data Request Received:
```json
{
  "xyz": {
    "array": ["value1", "value2"]
  }
}
```

- Mapping Configuration:
```json
"param_array_value": "src:req_body|xyz.array.1"
```

- Data Transformed for Target:
```json
/example/target?param_array_value=value2
```

#### Mapping an Entire Nested Object in the Request Body:

- Data Request Received:
```json
{
  "xyz": {
    "nested": {
      "field1": "value1",
      "field2": "value2"
    }
  }
}
```

- Mapping Configuration:
```json
"param_nested_key": "src:req_body|xyz.nested"
```

- Data Transformed for Target:
```json
{
  "param_nested_key": {
    "field1": "value1",
    "field2": "value2"
  }
}
```

#### Mapping a JSON Array in the Request Body to a Field:

- Data Request Received:
```json
{
  "xyz": {
    "array": [1, 2, 3]
  }
}
```

- Mapping Configuration:
```json
"param_array_values": "src:req_body|xyz.array"
```

- Data Transformed for Target:
```json
{
  "param_array_values": [1, 2, 3]
}
```


#### Mapping a JSON Object in the Request Body to a Field:

- Data Request Received:
```json
{
  "xyz": {
    "field1": "value1",
    "field2": "value2"
  }
}
```

- Mapping Configuration:
```json
"param_object": "src:req_body|xyz"
```

- Data Transformed for Target:
```json
{
  "param_object": {
    "field1": "value1",
    "field2": "value2"
  }
}
```

#### Mapping a Nested Array Element in the Request Body to a Field:

- Data Request Received:
```json
{
  "xyz": {
    "nested": [
      {
        "field": "value1"
      },
      {
        "field": "value2"
      }
    ]
  }
}
```

- Mapping Configuration:
```json
"param_nested_array_element": "src:req_body|xyz.nested.1.field"
```

- Data Transformed for Target:
```json
{
  "param_nested_array_element": "value2"
}
```

## 4. Response Mapping Examples
### Response Mapping with Data Transformations

#### Mapping a Nested Field in the Response Body to a Field:

- Data Response Received:
```json
{
  "xyz": {
    "nested": {
      "field": "nested_value"
    }
  }
}
```

- Mapping Configuration:
```json
"param_nested_value": "src:res_body|xyz.nested.field"
```

- Data Transformed for Response:
```json
{
  "param_nested_value": "nested_value"
}
```

#### Mapping a Specific Element of an Array in the Response Body to a Field:

- Data Response Received:
```json
{
  "xyz": {
    "array": ["value1", "value2"]
  }
}
```

- Mapping Configuration:
```json
"param_array_value": "src:res_body|xyz.array.1"
```

- Data Transformed for Response:
```json
{
  "param_array_value": "value2"
}
```

#### Mapping an Entire Nested Object in the Response Body:

- Data Response Received:
```json
{
  "xyz": {
    "nested": {
      "field1": "value1",
      "field2": "value2"
    }
  }
}
```

- Mapping Configuration:
```json
"param_nested_key": "src:res_body|xyz.nested"
```

- Data Transformed for Response:
```json
{
  "param_nested_key": {
    "field1": "value1",
    "field2": "value2"
  }
}
```

#### Mapping a JSON Array in the Response Body to a Field:

- Data Response Received:
```json
{
  "xyz": {
    "array": [1, 2, 3]
  }
}
```

- Mapping Configuration:
```json
"param_array_values": "src:res_body|xyz.array"
```

- Data Transformed for Response:
```json
{
  "param_array_values": [1, 2, 3]
}
```

#### Mapping a JSON Object in the Response Body to a Field:

- Data Response Received:
```json
{
  "xyz": {
    "field1": "value1",
    "field2": "value2"
  }
}
```

- Mapping Configuration:
```json
"param_object": "src:res_body|xyz"
```

- Data Transformed for Response:
```json
{
  "param_object": {
    "field1": "value1",
    "field2": "value2"
  }
}
```

#### Mapping a Nested Array Element in the Response Body to a Field:

- Data Response Received:
```json
{
  "xyz": {
    "nested": [
      {
        "field": "value1"
      },
      {
        "field": "value2"
      }
    ]
  }
}
```

- Mapping Configuration:
```json
"param_nested_array_element": "src:res_body|xyz.nested.1.field"
```
- Data Transformed for Response:
```json
{
  "param_nested_array_element": "value2"
}
```

### HTTP Status Code Mapping
#### 200 - Successful Response

For a successful response with an HTTP status code of 200, you can define response mappings as follows:
```json
"responseMapping": {
  "200": {
    "field1": "src:res_body|data.field1",
    "field2": "src:res_body|data.field2"
  }
}
```

In this example:

- `field1` in the response will be mapped to the value of `data.field1` in the target API's response body.
- `field2` in the response will be mapped to the value of `data.field2` in the target API's response body.

#### 400 - Client Error Response

For a client error response with an HTTP status code of 400, you can define response mappings as follows:
```json
"responseMapping": {
  "400": {
    "error": "src:static|Bad Request",
    "details": "src:res_body|error.details"
  }
}
```

In this example:

- `error` in the response will be set to the static string "Bad Request."
- `details` in the response will be mapped to the value of `error.details` in the target API's response body.

#### 401 - Unauthorized Response
For an unauthorized response with an HTTP status code of 401, you can define response mappings as follows:
```json
"responseMapping": {
  "401": {
    "error": "src:static|Unauthorized",
    "message": "src:res_body|error.message"
  }
}

```

In this example:

- `error` in the response will be set to the static string "Unauthorized."
- `message` in the response will be mapped to the value of `error.message` in the target API's response body.

#### 403 - Forbidden Response

For a forbidden response with an HTTP status code of 403, you can define response mappings as follows:
```json
"responseMapping": {
  "403": {
    "error": "src:static|Forbidden",
    "message": "src:res_body|error.message"
  }
}

```

In this example:

- `error` in the response will be set to the static string "Forbidden."
- `message` in the response will be mapped to the value of `error.message` in the target API's response body.

#### 404 - Not Found Response

For a "Not Found" response with an HTTP status code of 404, you can define response mappings as follows:
```json
"responseMapping": {
  "404": {
    "error": "src:static|Resource Not Found",
    "message": "src:static|The requested resource was not found."
  }
}
```

In this example:

- `error` in the response will be set to the static string "Resource Not Found."
- `message` in the response will be set to the static string "The requested resource was not found."

#### 500 - Server Error Response

For a server error response with an HTTP status code of 500, you can define response mappings as follows:
```json
"responseMapping": {
  "500": {
    "error": "src:static|Internal Server Error",
    "details": "src:res_body|error.details"
  }
}

```

In this example:

- `error` in the response will be set to the static string "Internal Server Error."
- `details` in the response will be mapped to the value of `error.details` in the target API's response body.

#### 503 - Service Unavailable Response

For a service unavailable response with an HTTP status code of 503, you can define response mappings as follows:
```json
"responseMapping": {
  "503": {
    "error": "src:static|Service Unavailable",
    "message": "src:res_body|error.message"
  }
}

```

In this example:

- `error` in the response will be set to the static string "Service Unavailable."
- `message` in the response will be mapped to the value of `error.message` in the target API's response body.

#### 504 - Timeout Response

For a timeout response with an HTTP status code of 504, you can define response mappings as follows:
```json
"responseMapping": {
  "504": {
    "error": "src:static|Request Timeout",
    "message": "src:static|The request timed out."
  }
}

```

In this example:

- `error` in the response will be set to the static string "Request Timeout."
- `message` in the response will be set to the static string "The request timed out."

These additional examples cover various HTTP status codes and demonstrate how you can structure the response data for different scenarios. You can customize response mappings to match your API Gateway's specific requirements.
## 5. Source Types (`src`) in Mapping

The following source types (`src`) can be used in request and response mapping configurations:

### Request Mapping

1. **Static Value (`src:static|value`)**: Use a static value as-is in the mapping. For example: `"phoneNumber": "src:static|12345"`.

2. **Query Parameter (`src:query|parameter_name`)**: Extract a value from the incoming query parameters of the request. For example: `"userId": "src:query|user_id"`.

3. **Request Body Field (`src:req_body|field_path`)**: Access a specific field in the request body JSON by providing the field's path. For example: `"name": "src:req_body|user.name"`.

4. **Request Header (`src:req_header|header_name`)**: Get the value of a specific request header. For example: `"apiKey": "src:req_header|X-API-Key"`.

5. **Function Call (`src:func|function_name(arguments)`)**: Invoke a custom function with specified arguments to generate the mapped value. For example: `"age": "src:func|calculateAge(src:req_body|dob)"`.


### Response Mapping

1. **Static Value (`src:static|value`)**: Use a static value as-is in the mapping. For example: `"status": "src:static|success"`.

2. **Response Body Field (`src:res_body|field_path`)**: Access a specific field in the response body JSON by providing the field's path. For example: `"userName": "src:res_body|user.name"`.

3. **Response Header (`src:res_header|header_name`)**: Get the value of a specific response header. For example: `"contentType": "src:res_header|Content-Type"`.

4. **Function Call (`src:func|function_name(arguments)`)**: Invoke a custom function with specified arguments to generate the mapped value. For example: `"totalScore": "src:func|calculateTotalScore(src:res_body|scores)"`.


### Examples

#### Request Mapping Examples

```json
"requestMapping": {
  "phoneNumber": "src:static|12345",
  "userId": "src:query|user_id",
  "name": "src:req_body|user.name",
  "apiKey": "src:req_header|X-API-Key",
  "age": "src:func|calculateAge(src:req_body|dob)"
}
```

#### Response Mapping Examples

```json
"responseMapping": {
  "status": "src:static|success",
  "userName": "src:res_body|user.name",
  "contentType": "src:res_header|Content-Type",
  "totalScore": "src:func|calculateTotalScore(src:res_body|scores)"
}
```

These source types provide flexibility in mapping values from various sources, including static values, query parameters, request and response body fields, headers, and custom function calls. Customize your mappings according to your specific API integration needs.
## 6. Creating and Using Custom Plugins

Custom plugins can extend the functionality of your API mapping configuration by allowing you to define custom functions that can be used in request and response mappings. Follow these steps to create and use custom plugins:

### Step 1: Plugin Creation

1. **Create a New Plugin File**: Begin by creating a new Go source file for your plugin. You can name it something like `my_custom_plugin.go`.

2. **Define Your Plugin**: In your Go source file, define your custom plugin as a struct with any necessary fields and methods. Ensure it follows the correct structure:

```go
package main

// MyCustomPlugin is a custom plugin that provides additional functionality.
type MyCustomPlugin struct{}

// Execute is the method that performs the custom logic.
func (p MyCustomPlugin) Execute(args ...interface{}) interface{} {
    // Implement your custom logic here.
    // Ensure the method signature matches and handles variable arguments.
}

// MyCustomPluginInstance is a variable that stores the plugin instance.
var MyCustomPluginInstance MyCustomPlugin
```

	- `MyCustomPlugin`: Your custom plugin struct.
	- `Execute(args ...interface{}) interface{}`: The method where your custom logic is implemented. It must accept variable arguments and return an interface{}.
	- `MyCustomPluginInstance`: Variable that stores the plugin instance.

3. **Build Your Plugin**: To compile your plugin into a shared object (`.so`) file, use the following command:
```bash
go build -buildmode=plugin -o my_custom_plugin.so my_custom_plugin.go
```

Replace `my_custom_plugin.go` and `my_custom_plugin.so` with your plugin's source file and desired output filename.


### Step 2: Plugin Configuration

1. **Add Plugin Configuration**: In your API mapping configuration file, specify the path to your custom plugin shared object file and create an instance of your plugin. For example:

```json
"pluginConfigs": [
{
"name": "MyCustomPlugin",
"path": "plugins/my_custom_plugin.so",
"instanceName": "MyCustomPluginInstance"
}
]

```

- `"name"`: A unique name for your plugin.
- `"path"`: The path to your compiled `.so` file.
- `"instanceName"`: A name to reference your plugin instance in mappings.

### Step 3: Using Custom Functions

1. **Invoke Custom Functions**: You can now use your custom functions in request and response mappings. For example:
```json
"requestMapping": {
"customValue": "src:func|MyCustomPluginInstance.Execute(arg1, arg2)"
}
```

### Example Plugin
#### `GenerateTransactionIDPlugin`
```go
package main

import (
	"github.com/google/uuid"
)

// GenerateTransactionIDPlugin is a plugin for generating transaction IDs.
type GenerateTransactionIDPlugin struct{}

// Execute generates a transaction ID.
func (p GenerateTransactionIDPlugin) Execute(args ...interface{}) interface{} {
	id := uuid.New()
	return id.String()
}

// GenerateTransactionIDPluginInstance is a variable that stores an instance of the plugin.
var GenerateTransactionIDPluginInstance GenerateTransactionIDPlugin

```

#### `SimSwapPlugin`
```go
package main

// SimSwapPlugin is a plugin for checking SIM swap.
type SimSwapPlugin struct{}

// Execute checks SIM swap with the arguments maxAge and score.
func (p SimSwapPlugin) Execute(args ...interface{}) interface{} {
	if len(args) != 2 {
		return "leng kurang"
	}

	maxAge, ok1 := args[0].(float64)
	score, ok2 := args[1].(float64)

	if !ok1 || !ok2 {
		return "Error"
	}

	if maxAge < 24 && score == 1 {
		return true
	} else if maxAge >= 24 && maxAge <= 48 && score == 2 {
		return true
	} else if maxAge > 48 && maxAge <= 72 && score == 3 {
		return true
	} else if maxAge > 72 && score == 4 {
		return true
	} else {
		return false
	}
}

// SimSwapPluginInstance is a variable that stores the plugin instance.
var SimSwapPluginInstance SimSwapPlugin 
```