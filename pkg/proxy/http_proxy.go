package proxy

import (
	"bytes"
	"encoding/json"
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

	pRequest.URL, err = getURL(p.Request.Params, endpoint.Parameters, p.Request.Service.Type, p.Request.Service.Host, endpoint.Path)
	if err != nil {
		return makeErrorResponse(makeMessage(err.Error(), p.Request.RequestID), http.StatusBadRequest)
	}

	pRequest.Data, err = getBody(p.Request.Params, endpoint.Parameters, p.Request.Data)
	if err != nil {
		return makeErrorResponse(makeMessage(err.Error(), p.Request.RequestID), http.StatusBadRequest)
	}

	pResponse := makeProxyRequest(&pRequest)

	if pResponse.StatusCode != http.StatusOK {
		return pResponse, pResponse.StatusCode
	}

	return pResponse, http.StatusOK
}

func getEndpointDefinition(requestedEndpointName string, resourceEndpoints []domain.Endpoint) (domain.Endpoint, error) {
	for _, endpoint := range resourceEndpoints {
		if endpoint.Name == requestedEndpointName {
			return endpoint, nil
		}
	}

	return domain.Endpoint{}, errors.New("endpoint not found in resource endpoints")
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
			return map[string]any{}, fmt.Errorf("missing required parameter: %s", definition.Name)
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
			return map[string]string{}, fmt.Errorf("unsupported data type for header param: %s", definition.Name)
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

		return map[string]string{}, fmt.Errorf("failed to parse value for %s as %s", definition.Name, definition.DataType)
	}

	return headers, nil
}

func getURL(
	params map[string]any,
	definitions []domain.Parameter,
	scheme string,
	host string,
	path string,
) (string, error) {
	url := fmt.Sprintf("%s://%s%s", scheme, host, path)

	supportedDataTypes := []string{domain.DataTypeString, domain.DataTypeInt, domain.DataTypeBool}

	for _, definition := range definitions {
		val, present := params[definition.Name]

		valToSet := ""

		if !present || definition.Type != domain.ParameterTypePath {
			continue
		}

		if !slices.Contains(supportedDataTypes, definition.DataType) {
			return "", fmt.Errorf("unsupported data type for path param: %s", definition.Name)
		}

		if stringVal, ok := val.(string); definition.DataType == domain.DataTypeString && ok {
			valToSet = stringVal
		} else if intVal, ok := val.(int); definition.DataType == domain.DataTypeInt && ok {
			valToSet = strconv.Itoa(intVal)
		} else if boolVal, ok := val.(bool); definition.DataType == domain.DataTypeBool && ok {
			valToSet = strconv.FormatBool(boolVal)
		} else {
			return "", fmt.Errorf("parameter value was bad data type: %s", definition.Name)
		}

		url = strings.ReplaceAll(url, fmt.Sprintf(":%s", definition.Name), valToSet)
	}

	return url, nil
}

func getBody(params map[string]any, definitions []domain.Parameter, originalBody any) (any, error) {
	body := originalBody

	for _, definition := range definitions {
		val, present := params[definition.Name]

		if !present || (definition.Type != domain.ParameterTypeBody && definition.Type != domain.ParameterTypeBodyFlat) {
			continue
		}

		if definition.Type == domain.ParameterTypeBodyFlat {
			return val, nil
		}

		if _, ok := originalBody.(map[string]any); !ok {
			return nil, errors.New("supplied a non-object as body with body parameter")
		}

		body.(map[string]any)[definition.Name] = val
	}

	return body, nil
}

func makeProxyRequest(proxyRequest *PRequest) *domain.ProxyResponse {
	serialized, err := json.Marshal(proxyRequest.Data)
	if err != nil {
		return &domain.ProxyResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Something very weird happened",
			Data:       map[string]any{},
		}
	}

	req, err := http.NewRequest(proxyRequest.Method, proxyRequest.URL, bytes.NewReader(serialized))
	if err != nil {
		return &domain.ProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Something very weird happened",
			Data:       map[string]any{},
		}
	}

	for key, val := range proxyRequest.Headers {
		req.Header.Add(key, val)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &domain.ProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Something very weird happened",
			Data:       map[string]any{},
		}
	}

	pResponse := &domain.ProxyResponse{
		StatusCode: http.StatusOK,
		Message:    "Conductor Proxy request success",
		Data:       make([]byte, 0),
	}

	_, err = resp.Body.Read(pResponse.Data.([]byte))
	if err != nil {
		return &domain.ProxyResponse{
			StatusCode: http.StatusOK,
			Message:    "Response did not parse as bytes",
			Data:       map[string]any{},
		}
	}

	return pResponse
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
