package proxy

import "github.com/Zentech-Development/conductor-proxy/domain"

type HTTPProxy struct {
	Request domain.ProxyRequest
}

func NewHTTPProxy(request domain.ProxyRequest) *HTTPProxy {
	return &HTTPProxy{
		Request: request,
	}
}

func (p *HTTPProxy) GetResponse() (*domain.ProxyResponse, int) {
	return &domain.ProxyResponse{}, 200
}
