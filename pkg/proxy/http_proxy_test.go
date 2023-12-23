package proxy

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/Zentech-Development/conductor-proxy/domain"
)

func getValidRequest() domain.ProxyRequest {
	return domain.ProxyRequest{
		RequestID: "123",
		Method:    "GET",
		Endpoint:  "Get All Things",
		Resource:  getValidResource(),
		Service:   getValidService(),
	}
}

func getValidResource() domain.Resource {
	return domain.Resource{
		ID:           "15",
		Name:         "things",
		FriendlyName: "Things",
		ServiceID:    "72",
		Properties: []domain.Property{
			{
				Name:         "prop1",
				FriendlyName: "Prop 1",
				DataType:     domain.DataTypeString,
				Required:     false,
				DefaultValue: nil,
				HasDefault:   false,
			},
		},
		Endpoints: []domain.Endpoint{
			{
				Name:   "Get All Things",
				Path:   "/things",
				Method: "GET",
				Parameters: []domain.Parameter{
					{
						Name:         "param1",
						FriendlyName: "Param 1",
						DataType:     domain.DataTypeString,
						Required:     false,
						DefaultValue: "",
						HasDefault:   false,
						Type:         domain.ParameterTypeHeader,
					},
				},
			},
		},
	}
}

func getValidService() domain.Service {
	return domain.Service{
		ID:           "72",
		Name:         "my-app",
		FriendlyName: "My App",
		Host:         "localhost:8001",
		AdminGroups:  []string{"my-app-admins"},
		UserGroups:   []string{"my-app-users", "my-app-other-users"},
		Type:         domain.ServiceTypeHTTP,
	}
}

func TestNewHTTPProxy(t *testing.T) {
	p := NewHTTPProxy(domain.ProxyRequest{
		Method: "get",
	})

	if p.Request.Method != "GET" {
		t.Fatal("Should have converted method to uppercase")
	}
}

func TestGetEndpointDefinition(t *testing.T) {
	endpoints := getValidResource().Endpoints
	endpointName := endpoints[0].Name

	result, err := getEndpointDefinition(endpointName, endpoints)
	if err != nil {
		t.Fatal("Unexpected error")
	}

	if !reflect.DeepEqual(result, endpoints[0]) {
		t.Fatal("Got wrong endpoint")
	}
}

func TestGetEndpointDefinitionNotFound(t *testing.T) {
	endpoints := getValidResource().Endpoints
	endpointName := "not the right name"

	_, err := getEndpointDefinition(endpointName, endpoints)
	if err == nil {
		t.Fatal("Expected an error")
	}
}

func TestSetDefaultsForMissingRequiredParams(t *testing.T) {
	definitions := []domain.Parameter{
		{
			Name:         "param1",
			FriendlyName: "Param 1",
			DataType:     domain.DataTypeString,
			Required:     false,
			DefaultValue: "",
			HasDefault:   false,
			Type:         domain.ParameterTypeHeader,
		},
		{
			Name:         "param2",
			FriendlyName: "Param 2",
			DataType:     domain.DataTypeString,
			Required:     false,
			DefaultValue: "DEFAULT",
			HasDefault:   true,
			Type:         domain.ParameterTypeHeader,
		},
		{
			Name:         "param3",
			FriendlyName: "Param 3",
			DataType:     domain.DataTypeString,
			Required:     true,
			DefaultValue: "DEFAULT",
			HasDefault:   true,
			Type:         domain.ParameterTypeHeader,
		},
	}

	params := map[string]any{
		"param2": "val2",
		"param3": "val3",
	}

	params, err := setDefaultsForMissingRequiredParams(params, definitions)
	if err != nil {
		t.Fatal("Unexpected error")
	}

	if params["param2"] != "val2" || params["param3"] != "val3" {
		t.Fatal("Parameter changed unexpectedly")
	}

	if _, ok := params["param1"]; ok {
		t.Fatal("Parameter added unexpectedly")
	}

	definitions = []domain.Parameter{
		{
			Name:         "param1",
			FriendlyName: "Param 1",
			DataType:     domain.DataTypeString,
			Required:     true,
			DefaultValue: "",
			HasDefault:   false,
			Type:         domain.ParameterTypeHeader,
		},
	}

	params = map[string]any{
		"param2": "val2",
	}

	_, err = setDefaultsForMissingRequiredParams(params, definitions)
	if err == nil {
		t.Fatal("Expected error for missing required param1")
	}

	definitions = []domain.Parameter{
		{
			Name:         "param1",
			FriendlyName: "Param 1",
			DataType:     domain.DataTypeString,
			Required:     true,
			DefaultValue: "DEFAULT",
			HasDefault:   true,
			Type:         domain.ParameterTypeHeader,
		},
	}

	params = map[string]any{
		"param2": "val2",
	}

	params, err = setDefaultsForMissingRequiredParams(params, definitions)
	if err != nil {
		t.Fatal("Unexpected error")
	}

	if params["param2"] != "val2" {
		t.Fatal("Parameter changed unexpectedly")
	}

	if params["param1"] != "DEFAULT" {
		t.Fatal("Failed to set default value for missing required parameter")
	}
}

func TestGetHeaders(t *testing.T) {
	definitions := []domain.Parameter{
		{
			Name:         "param1",
			FriendlyName: "Param 1",
			DataType:     domain.DataTypeString,
			Required:     false,
			DefaultValue: "",
			HasDefault:   false,
			Type:         domain.ParameterTypeHeader,
		},
		{
			Name:         "param2",
			FriendlyName: "Param 2",
			DataType:     domain.DataTypeString,
			Required:     false,
			DefaultValue: "DEFAULT",
			HasDefault:   true,
			Type:         domain.ParameterTypePath,
		},
		{
			Name:         "param3",
			FriendlyName: "Param 3",
			DataType:     domain.DataTypeInt,
			Required:     false,
			DefaultValue: 0,
			HasDefault:   true,
			Type:         domain.ParameterTypeHeader,
		},
		{
			Name:         "param4",
			FriendlyName: "Param 4",
			DataType:     domain.DataTypeBool,
			Required:     false,
			DefaultValue: false,
			HasDefault:   false,
			Type:         domain.ParameterTypeHeader,
		},
	}

	params := map[string]any{
		"param1": "val1",
		"param2": "val2",
		"param3": 123,
		"param4": true,
	}

	headers, err := getHeaders(params, definitions)
	if err != nil {
		t.Fatal("Unexpected error")
	}

	if headers["param1"] != "val1" || headers["param3"] != "123" || headers["param4"] != "true" {
		t.Fatal("Failed to set a header value")
	}

	if _, ok := headers["param2"]; ok {
		t.Fatal("Set a non-header param in the header")
	}

	definitions = []domain.Parameter{
		{
			Name:         "param1",
			FriendlyName: "Param 1",
			DataType:     domain.DataTypeArray,
			Required:     false,
			DefaultValue: "",
			HasDefault:   false,
			Type:         domain.ParameterTypeHeader,
		},
	}

	params = map[string]any{
		"param1": []string{"val1"},
	}

	_, err = getHeaders(params, definitions)
	if err == nil {
		t.Fatal("Expected error for unsupported type")
	}

	definitions = []domain.Parameter{
		{
			Name:         "param1",
			FriendlyName: "Param 1",
			DataType:     domain.DataTypeString,
			Required:     false,
			DefaultValue: "",
			HasDefault:   false,
			Type:         domain.ParameterTypeHeader,
		},
	}

	params = map[string]any{
		"param1": 123,
	}

	_, err = getHeaders(params, definitions)
	if err == nil {
		t.Fatal("Expected error for wrong value type")
	}
}

func TestGetURL(t *testing.T) {
	definitions := []domain.Parameter{
		{
			Name:         "param1",
			FriendlyName: "Param 1",
			DataType:     domain.DataTypeString,
			Required:     false,
			DefaultValue: "",
			HasDefault:   false,
			Type:         domain.ParameterTypePath,
		},
	}

	params := map[string]any{
		"param1": "val1",
	}

	scheme := "http"

	host := "localhost:8091"

	path := "/api/v1/:param1"

	url, err := getURL(params, definitions, scheme, host, path)
	if err != nil {
		t.Fatal("Unexpected error occurred")
	}

	if url != "http://localhost:8091/api/v1/val1" {
		t.Fatal("URL generated incorrectly")
	}

	noURLParamsDefinitions := []domain.Parameter{
		{
			Name:         "param1",
			FriendlyName: "Param 1",
			DataType:     domain.DataTypeString,
			Required:     false,
			DefaultValue: "",
			HasDefault:   false,
			Type:         domain.ParameterTypePath,
		},
		{
			Name:         "param2",
			FriendlyName: "Param 2",
			DataType:     domain.DataTypeString,
			Required:     false,
			DefaultValue: "",
			HasDefault:   false,
			Type:         domain.ParameterTypeHeader,
		},
	}

	noURLParamsParams := map[string]any{
		"param2": "val2",
	}

	url, err = getURL(noURLParamsParams, noURLParamsDefinitions, scheme, host, "/api/v2/:val1")
	if err != nil {
		t.Fatal("Unexpected error occurred")
	}

	if url != "http://localhost:8091/api/v2/:val1" {
		t.Fatal("URL generated incorrectly")
	}

	multiParamPath := "/api/v1/:param1/:param1"

	url, err = getURL(params, definitions, scheme, host, multiParamPath)
	if err != nil {
		t.Fatal("Unexpected error occurred")
	}

	if url != "http://localhost:8091/api/v1/val1/val1" {
		t.Fatal("URL generated incorrectly")
	}

	badDataTypeParams := map[string]any{
		"param1": true,
	}

	if _, err = getURL(badDataTypeParams, definitions, scheme, host, path); err == nil {
		t.Fatal("Expected bad data type error")
	}
}

func TestGetBody(t *testing.T) {
	definitions := []domain.Parameter{
		{
			Name:         "param1",
			FriendlyName: "Param 1",
			DataType:     domain.DataTypeString,
			Required:     false,
			DefaultValue: map[string]any{},
			HasDefault:   false,
			Type:         domain.ParameterTypeBody,
		},
		{
			Name:         "param2",
			FriendlyName: "Param 2",
			DataType:     domain.DataTypeObject,
			Required:     false,
			DefaultValue: map[string]any{},
			HasDefault:   false,
			Type:         domain.ParameterTypeBodyFlat,
		},
		{
			Name:         "param3",
			FriendlyName: "Param 3",
			DataType:     domain.DataTypeString,
			Required:     false,
			DefaultValue: "",
			HasDefault:   false,
			Type:         domain.ParameterTypeHeader,
		},
	}

	params := map[string]any{
		"param1": "val1",
		"param3": "val3",
	}

	originalBody := map[string]any{
		"test":   "asdf",
		"param1": "old value",
	}

	body, err := getBody(params, definitions, originalBody)
	if err != nil {
		t.Fatal("Unexpcted error occurred")
	}

	if bodyDict, ok := body.(map[string]any); !ok || bodyDict["test"] != "asdf" || bodyDict["param1"] != "val1" {
		t.Fatal("Body generated incorrectly")
	}

	if _, err = getBody(params, definitions, []string{"invalid"}); err == nil {
		t.Fatal("Expected non-object body with body parameter error")
	}

	flatParams := map[string]any{
		"param2": map[string]any{
			"test": "1234",
		},
	}

	flatBody, err := getBody(flatParams, definitions, originalBody)
	if err != nil {
		t.Fatal("Unexpcted error occurred")
	}

	if flatBody, ok := flatBody.(map[string]any); !ok || flatBody["test"] != "1234" {
		t.Fatal("Body generated incorrectly")
	}
}

func TestMakeProxyRequest(t *testing.T) {
	request := PRequest{
		Method:  "GET",
		Data:    map[string]any{},
		URL:     "http://google.com",
		Headers: map[string]string{},
	}

	response := makeProxyRequest(&request)

	if response.StatusCode != http.StatusOK || response.Message != "Conductor Proxy request success" {
		t.Fatal("Failed to make simple request")
	}
}
