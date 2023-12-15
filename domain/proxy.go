package domain

type ProxyResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Data       any    `json:"data"`
}

type ProxyRequest struct {
	RequestID string         `json:"requestId"`
	Resource  Resource       `json:"resource"`
	App       App            `json:"app"`
	Endpoint  string         `json:"endpoint"`
	Method    string         `json:"method"`
	Params    map[string]any `json:"params"`
	Data      map[string]any `json:"data"`
}

type ProxyRequestInput struct {
	ResourceID string         `json:"resourceId"`
	Endpoint   string         `json:"endpoint"`
	Method     string         `json:"method"`
	Params     map[string]any `json:"params"`
	Data       map[string]any `json:"data"`
}

type ProxyHandlers interface {
	HTTPRequest(request ProxyRequest) (ProxyResponse, int)
}
