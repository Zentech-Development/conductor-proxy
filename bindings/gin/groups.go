package bindings

import (
	"fmt"
	"net/http"

	"github.com/Zentech-Development/conductor-proxy/domain"
	"github.com/gin-gonic/gin"
)

type GroupsGinBindings struct {
	Handlers domain.Handlers
}

func newGroupsGinBinding(handlers domain.Handlers) *GroupsGinBindings {
	return &GroupsGinBindings{
		Handlers: handlers,
	}
}

func (b *GroupsGinBindings) Post(c *gin.Context) {
	var groupInput domain.GroupInput

	if err := c.ShouldBindJSON(&groupInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    fmt.Sprintf("[Request ID: %s]: Failed to parse request", c.GetString("requestId")),
			"data":       map[string]any{},
		})
		return
	}

	userGroups, _ := c.Get("userGroups")

	group, err := b.Handlers.Groups.Add(groupInput, userGroups.([]string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"message":    fmt.Sprintf("[Request ID: %s]: Unexpected error while adding group: %s", c.GetString("requestId"), err.Error()),
			"data":       map[string]any{},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"statusCode": http.StatusCreated,
		"message":    fmt.Sprintf("[Request ID: %s]: Added group successfully", c.GetString("requestId")),
		"data": map[string]any{
			"group": group,
		},
	})
}
