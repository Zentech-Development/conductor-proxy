package bindings

import (
	"fmt"
	"net/http"

	"github.com/Zentech-Development/conductor-proxy/domain"
	"github.com/gin-gonic/gin"
)

type ServicesGinBindings struct {
	Handlers domain.Handlers
}

func newServicesGinBinding(handlers domain.Handlers) *ServicesGinBindings {
	return &ServicesGinBindings{
		Handlers: handlers,
	}
}

func (b *ServicesGinBindings) Post(c *gin.Context) {
	var serviceInput domain.ServiceInput

	if err := c.ShouldBindJSON(&serviceInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    fmt.Sprintf("[Request ID: %s]: Failed to parse request", c.GetString("requestId")),
			"data":       map[string]any{},
		})
		return
	}

	service, err := b.Handlers.Services.Add(serviceInput, []string{"admin"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"message":    fmt.Sprintf("[Request ID: %s]: Unexpected error while adding service %s", c.GetString("requestId"), err.Error()),
			"data":       map[string]any{},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"statusCode": http.StatusCreated,
		"message":    fmt.Sprintf("[Request ID: %s]: Added service successfully", c.GetString("requestId")),
		"data": map[string]any{
			"service": service,
		},
	})
}
