package handlers

import (
	"net/http"

	"github.com/Zentech-Development/conductor-proxy/domain"
	"github.com/Zentech-Development/conductor-proxy/proxy"
)

type ProxyHandler struct {
	Adapters *domain.Adapters
}

func NewProxyHandler(adapters *domain.Adapters) AppHandler {
	return AppHandler{
		Adapters: adapters,
	}
}

func (h ProxyHandler) ProxyRequest(request domain.ProxyRequest, userGroups []string) (*domain.ProxyResponse, int) {
	if !checkForGroupMatch(userGroups, request.App.UserGroups) && !checkForGroupMatch(userGroups, request.App.AdminGroups) {
		return &domain.ProxyResponse{
			StatusCode: http.StatusForbidden,
			Message:    "Account not authorized to access this resource",
			Data:       map[string]any{},
		}, http.StatusForbidden
	}

	if request.App.Type == domain.AppTypeHTTP || request.App.Type == domain.AppTypeHTTPS {
		proxy := proxy.NewHTTPProxy(request)
		response, statusCode := proxy.GetResponse()
		return response, statusCode
	}

	return &domain.ProxyResponse{
		StatusCode: http.StatusNotFound,
		Message:    "No proxy found for that protocol",
		Data:       map[string]any{},
	}, http.StatusNotFound
}
