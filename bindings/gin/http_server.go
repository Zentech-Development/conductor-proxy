package bindings

import (
	"encoding/json"
	"net/http"

	"github.com/Zentech-Development/conductor-proxy/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type HTTPServerBindingConfig struct {
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

	accountsBindings := newAccountsGinBinding(h)
	appsBindings := newAppsGinBinding(h)
	groupsBindings := newGroupsGinBinding(h)
	resourcesBindings := newResourcesGinBinding(h)
	proxyBindings := newProxyGinBinding(h)

	r.POST("/proxy", proxyBindings.Post)

	apiRouter := r.Group("/api")
	{
		apiRouter.POST("/login", accountsBindings.Login)

		accountsRouter := apiRouter.Group("/accounts")
		{
			accountsRouter.POST("/", accountsBindings.Post)
			accountsRouter.PUT("/:id", accountsBindings.UpdateGroups)
			accountsRouter.DELETE("/:id", func(c *gin.Context) {
				c.AbortWithStatus(http.StatusNotImplemented)
			})
		}

		groupsRouter := apiRouter.Group("/groups")
		{
			groupsRouter.POST("/", groupsBindings.Post)
			groupsRouter.DELETE("/:id", func(c *gin.Context) {
				c.AbortWithStatus(http.StatusNotImplemented)
			})
		}

		appsRouter := apiRouter.Group("/apps")
		{
			appsRouter.POST("/", appsBindings.Post)
		}

		resourceRouter := apiRouter.Group("/resources")
		{
			resourceRouter.POST("/", resourcesBindings.Post)
		}
	}

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
