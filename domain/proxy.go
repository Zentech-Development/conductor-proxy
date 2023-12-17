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
	ResourceID string         `json:"resourceId" binding:"required"`
	Endpoint   string         `json:"endpoint" binding:"required"`
	Method     string         `json:"method" binding:"required"`
	Params     map[string]any `json:"params" binding:"required"`
	Data       map[string]any `json:"data" binding:"required"`
}

type ProxyHandlers interface {
	ProxyRequest(request ProxyRequest, userGroups []string) (*ProxyResponse, int)
}
