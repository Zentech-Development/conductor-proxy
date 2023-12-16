package bindings

import (
	"fmt"
	"net/http"

	"github.com/Zentech-Development/conductor-proxy/domain"
	"github.com/gin-gonic/gin"
)

type AppsGinBindings struct {
	Handlers domain.Handlers
}

func newAppsGinBinding(handlers domain.Handlers) *AppsGinBindings {
	return &AppsGinBindings{
		Handlers: handlers,
	}
}

func (b *AppsGinBindings) Post(c *gin.Context) {
	var appInput domain.AppInput

	if err := c.ShouldBindJSON(&appInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    fmt.Sprintf("[Request ID: %s]: Failed to parse request", c.GetString("requestId")),
			"data":       map[string]any{},
		})
		return
	}

	app, err := b.Handlers.Apps.Add(appInput, []string{"admin"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"message":    fmt.Sprintf("[Request ID: %s]: Unexpected error while adding app %s", c.GetString("requestId"), err.Error()),
			"data":       map[string]any{},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"statusCode": http.StatusCreated,
		"message":    fmt.Sprintf("[Request ID: %s]: Added app successfully", c.GetString("requestId")),
		"data": map[string]any{
			"app": app,
		},
	})
}
