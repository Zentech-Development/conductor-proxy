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

	userGroups, _ := c.Get("userGroups")

	account, err := b.Handlers.Accounts.Add(accountInput, userGroups.([]string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"message":    fmt.Sprintf("[Request ID: %s]: Unexpected error while adding account: %s", c.GetString("requestId"), err.Error()),
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

	accountUsername := c.Param("id")

	userGroups, _ := c.Get("userGroups")

	err := b.Handlers.Accounts.UpdateGroups(accountUsername, input.GroupsToAdd, input.GroupsToRemove, userGroups.([]string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"message":    fmt.Sprintf("[Request ID: %s]: Unexpected error while updating account groups: %s", c.GetString("requestId"), err.Error()),
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

func (b *AccountsGinBindings) Login(c *gin.Context) {
	var input domain.LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    fmt.Sprintf("[Request ID: %s]: Failed to parse request", c.GetString("requestId")),
			"data":       map[string]any{},
		})
		return
	}

	account, err := b.Handlers.Accounts.Login(input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"statusCode": http.StatusUnauthorized,
			"message":    fmt.Sprintf("[Request ID: %s]: Failed to login", c.GetString("requestId")),
			"data":       map[string]any{},
		})
		return
	}

	token, err := getAccessToken(input.Username, account.Groups, account.TokenExpiration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"message":    fmt.Sprintf("[Request ID: %s]: Unexpected error while logging in: %s", c.GetString("requestId"), err.Error()),
			"data":       map[string]any{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"message":    fmt.Sprintf("[Request ID: %s]: Login successful", c.GetString("requestId")),
		"data": map[string]any{
			"token": token,
		},
	})
}
