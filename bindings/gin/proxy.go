package bindings

import (
	"fmt"
	"net/http"

	"github.com/Zentech-Development/conductor-proxy/domain"
	"github.com/gin-gonic/gin"
)

type ProxyGinBindings struct {
	Handlers domain.Handlers
}

func newProxyGinBinding(handlers domain.Handlers) *ProxyGinBindings {
	return &ProxyGinBindings{
		Handlers: handlers,
	}
}

func (b *ProxyGinBindings) Post(c *gin.Context) {
	var requestInput domain.ProxyRequestInput

	if err := c.ShouldBindJSON(&requestInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    fmt.Sprintf("[Request ID: %s]: Failed to parse request", c.GetString("requestId")),
			"data":       map[string]any{},
		})
		return
	}

	userGroups, _ := c.Get("userGroups")

	resource, err := b.Handlers.Resources.GetByID(requestInput.ResourceID, userGroups.([]string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"statusCode": http.StatusNotFound,
			"message":    fmt.Sprintf("[Request ID: %s]: Resource with ID %s not found", c.GetString("requestId"), requestInput.ResourceID),
			"data":       map[string]any{},
		})
		return
	}

	service, err := b.Handlers.Services.GetByID(resource.ServiceID, userGroups.([]string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusInternalServerError,
			"message":    fmt.Sprintf("[Request ID: %s]: Internal error occurred, possible orphaned resource", c.GetString("requestId")),
			"data":       map[string]any{},
		})
		return
	}

	request := domain.ProxyRequest{
		RequestID: c.GetString("requestId"),
		Resource:  resource,
		Service:   service,
		Endpoint:  requestInput.Endpoint,
		Method:    requestInput.Method,
		Data:      requestInput.Data,
		Params:    requestInput.Params,
	}

	response, statusCode := b.Handlers.Proxy.ProxyRequest(request, userGroups.([]string))

	c.JSON(statusCode, gin.H{
		"statusCode": response.StatusCode,
		"message":    response.Message,
		"data":       response.Data,
	})
}
