package bindings

import (
	"fmt"
	"net/http"

	"github.com/Zentech-Development/conductor-proxy/domain"
	"github.com/gin-gonic/gin"
)

type AccountsGinBindings struct {
	Handlers domain.Handlers
}

func newAccountsGinBinding(handlers domain.Handlers) *AccountsGinBindings {
	return &AccountsGinBindings{
		Handlers: handlers,
	}
}

func (b *AccountsGinBindings) Post(c *gin.Context) {
	var accountInput domain.AccountInput

	if err := c.ShouldBindJSON(&accountInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    fmt.Sprintf("[Request ID: %s]: Failed to parse request", c.GetString("requestId")),
			"data":       map[string]any{},
		})
		return
	}

	account, err := b.Handlers.Accounts.Add(accountInput, []string{"admin"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"message":    fmt.Sprintf("[Request ID: %s]: Unexpected error while adding account %s", c.GetString("requestId"), err.Error()),
			"data":       map[string]any{},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"statusCode": http.StatusCreated,
		"message":    fmt.Sprintf("[Request ID: %s]: Added account successfully", c.GetString("requestId")),
		"data": map[string]any{
			"account": account,
		},
	})
}

type updateAccountGroupsInput struct {
	GroupsToAdd    []string `json:"groupsToAdd"`
	GroupsToRemove []string `json:"groupsToRemove"`
}

func (b *AccountsGinBindings) UpdateGroups(c *gin.Context) {
	var input updateAccountGroupsInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    fmt.Sprintf("[Request ID: %s]: Failed to parse request", c.GetString("requestId")),
			"data":       map[string]any{},
		})
		return
	}

	accountID := c.Param("id")

	err := b.Handlers.Accounts.UpdateGroups(accountID, input.GroupsToAdd, input.GroupsToRemove, []string{"admin"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"message":    fmt.Sprintf("[Request ID: %s]: Unexpected error while updating account groups %s", c.GetString("requestId"), err.Error()),
			"data":       map[string]any{},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"statusCode": http.StatusCreated,
		"message":    fmt.Sprintf("[Request ID: %s]: Updated account groups successfully", c.GetString("requestId")),
		"data":       map[string]any{},
	})
}
