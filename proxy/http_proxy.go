package proxy

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/Zentech-Development/conductor-proxy/domain"
)

type HTTPProxy struct {
	Request domain.ProxyRequest
}

type PRequest struct {
	Method  string
	Data    any
	URL     string
	Headers map[string]string
}

func NewHTTPProxy(request domain.ProxyRequest) *HTTPProxy {
	request.Method = strings.ToUpper(request.Method)

	return &HTTPProxy{
		Request: request,
	}
}

func (p *HTTPProxy) GetResponse() (*domain.ProxyResponse, int) {
	pRequest := PRequest{
		Method:  p.Request.Method,
		Data:    map[string]any{},
		URL:     "",
		Headers: map[string]string{},
	}

	endpoint, err := getEndpointDefinition(p.Request.Endpoint, p.Request.Resource.Endpoints)
	if err != nil {
		return makeErrorResponse(makeMessage(fmt.Sprintf("Endpoint with name %s not found in resource", p.Request.Endpoint), p.Request.RequestID), http.StatusNotFound)
	}

	p.Request.Params, err = setDefaultsForMissingRequiredParams(p.Request.Params, endpoint.Parameters)
	if err != nil {
		return makeErrorResponse(makeMessage(err.Error(), p.Request.RequestID), http.StatusBadRequest)
	}

	pRequest.Headers, err = getHeaders(p.Request.Params, endpoint.Parameters)
	if err != nil {
		return makeErrorResponse(makeMessage(err.Error(), p.Request.RequestID), http.StatusBadRequest)
	}

	// TODO: get URL

	// TODO: get body from request (default) or from body or bodyflat params

	// TODO: make request and process response

	return &domain.ProxyResponse{}, http.StatusOK
}

func getEndpointDefinition(requestedEndpointName string, resourceEndpoints []domain.Endpoint) (domain.Endpoint, error) {
	for _, endpoint := range resourceEndpoints {
		if endpoint.Name == requestedEndpointName {
			return endpoint, nil
		}
	}

	return domain.Endpoint{}, errors.New("Endpoint not found in resource endpoints")
}

func setDefaultsForMissingRequiredParams(params map[string]any, definitions []domain.Parameter) (map[string]any, error) {
	for _, definition := range definitions {
		if !definition.Required {
			continue
		}

		_, present := params[definition.Name]

		if present {
			continue
		}

		if !definition.HasDefault {
			return map[string]any{}, errors.New(fmt.Sprintf("Missing required parameter: %s", definition.Name))
		}

		params[definition.Name] = definition.DefaultValue
	}

	return params, nil
}

func getHeaders(params map[string]any, definitions []domain.Parameter) (map[string]string, error) {
	supportedDataTypes := []string{domain.DataTypeString, domain.DataTypeInt, domain.DataTypeBool}

	headers := map[string]string{}

	for _, definition := range definitions {
		val, present := params[definition.Name]

		if !present || definition.Type != domain.ParameterTypeHeader {
			continue
		}

		if !slices.Contains(supportedDataTypes, definition.DataType) {
			return map[string]string{}, errors.New(fmt.Sprintf("Unsupported data type for header param: %s", definition.Name))
		}

		if definition.DataType == domain.DataTypeString {
			if stringVal, ok := val.(string); ok {
				headers[definition.Name] = stringVal
				continue
			}
		}

		if definition.DataType == domain.DataTypeInt {
			if intVal, ok := val.(int); ok {
				headers[definition.Name] = strconv.Itoa(intVal)
				continue
			}
		}

		if definition.DataType == domain.DataTypeBool {
			if boolVal, ok := val.(bool); ok {
				headers[definition.Name] = strconv.FormatBool(boolVal)
				continue
			}
		}

		return map[string]string{}, errors.New(fmt.Sprintf("Failed to parse value for %s as %s", definition.Name, definition.DataType))
	}

	return headers, nil
}

func makeMessage(message string, requestId string) string {
	return fmt.Sprintf("[Request ID: %s] %s", requestId, message)
}

func makeErrorResponse(message string, statusCode int) (*domain.ProxyResponse, int) {
	return &domain.ProxyResponse{
		StatusCode: statusCode,
		Message:    message,
		Data:       map[string]any{},
	}, statusCode
}
