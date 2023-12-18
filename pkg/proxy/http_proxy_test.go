package proxy

import (
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

	params, err = setDefaultsForMissingRequiredParams(params, definitions)
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

	headers, err = getHeaders(params, definitions)
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

	headers, err = getHeaders(params, definitions)
	if err == nil {
		t.Fatal("Expected error for wrong value type")
	}
}
