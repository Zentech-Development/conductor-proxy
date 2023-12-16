package bindings

import (
	"fmt"
	"net/http"

	"github.com/Zentech-Development/conductor-proxy/domain"
	"github.com/gin-gonic/gin"
)

type ResourcesGinBindings struct {
	Handlers domain.Handlers
}

func newResourcesGinBinding(handlers domain.Handlers) *ResourcesGinBindings {
	return &ResourcesGinBindings{
		Handlers: handlers,
	}
}

func (b *ResourcesGinBindings) Post(c *gin.Context) {
	var resourceInput domain.ResourceInput

	if err := c.ShouldBindJSON(&resourceInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    fmt.Sprintf("[Request ID: %s]: Failed to parse request", c.GetString("requestId")),
			"data":       map[string]any{},
		})
		return
	}

	resource, err := b.Handlers.Resources.Add(resourceInput, []string{"admin"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"message":    fmt.Sprintf("[Request ID: %s]: Unexpected error while adding resource %s", c.GetString("requestId"), err.Error()),
			"data":       map[string]any{},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"statusCode": http.StatusCreated,
		"message":    fmt.Sprintf("[Request ID: %s]: Added resource successfully", c.GetString("requestId")),
		"data": map[string]any{
			"resource": resource,
		},
	})
}
