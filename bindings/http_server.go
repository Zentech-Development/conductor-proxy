package bindings

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Zentech-Development/conductor-proxy/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type HTTPServerBindingConfig struct {
	Host      string
	SecretKey string
	GinMode   string
}

func NewHTTPServerBinding(h domain.Handlers, config HTTPServerBindingConfig) *gin.Engine {
	gin.SetMode(config.GinMode)

	r := gin.Default()
	r.SetTrustedProxies(nil)

	r.Use(func(c *gin.Context) {
		c.Set("requestId", uuid.NewString())
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": http.StatusOK,
			"message":    "Success",
			"data":       map[string]any{},
		})
	})

	apiRouter := r.Group("/api")
	{
		apiRouter.POST("/login", func(c *gin.Context) {
			var loginInput domain.LoginInput

			if err := c.ShouldBindJSON(&loginInput); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"statusCode": http.StatusBadRequest,
					"message":    fmt.Sprintf("[RequestID: %s]: Failed to parse request", c.GetString("requestId")),
					"data":       map[string]any{},
				})
				return
			}

			// TODO: FINISH
		})
	}

	r.POST("/proxy", func(c *gin.Context) {
		var requestInput domain.ProxyRequestInput

		if err := c.ShouldBindJSON(&requestInput); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"statusCode": http.StatusBadRequest,
				"message":    fmt.Sprintf("[RequestID: %s]: Failed to parse request", c.GetString("requestId")),
				"data":       map[string]any{},
			})
			return
		}

		resource, err := h.Resources.GetByID(requestInput.ResourceID, make([]string, 0))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"statusCode": http.StatusNotFound,
				"message":    fmt.Sprintf("[RequestID: %s]: Resource with ID %s not found", c.GetString("requestId"), requestInput.ResourceID),
				"data":       map[string]any{},
			})
			return
		}

		app, err := h.Apps.GetByID(requestInput.ResourceID, make([]string, 0))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"statusCode": http.StatusInternalServerError,
				"message":    fmt.Sprintf("[RequestID: %s]: Internal error occurred, possible orphaned resource", c.GetString("requestId")),
				"data":       map[string]any{},
			})
			return
		}

		request := domain.ProxyRequest{
			RequestID: c.GetString("requestId"),
			Resource:  resource,
			App:       app,
			Endpoint:  requestInput.Endpoint,
			Method:    requestInput.Method,
			Data:      requestInput.Data,
			Params:    requestInput.Params,
		}

		response, statusCode := h.Proxy.ProxyRequest(request, make([]string, 0))

		c.JSON(statusCode, gin.H{
			"statusCode": response.StatusCode,
			"message":    response.Message,
			"data":       response.Data,
		})
	})

	return r
}

func sendResult(w http.ResponseWriter, response domain.ProxyResponse, statusCode int) {
	responseSerialized, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(statusCode)

	w.Write(responseSerialized)
}
