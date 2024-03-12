// Package router provides functions for managing the app's router via Gin.
package router

import (
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xtt28/shortener/api/handlers"
)

// InitRouter creates a Gin router with the default middleware, then registering
// all custom templates, middleware and routes.
func InitRouter() (r *gin.Engine) {
	r = gin.Default()
	r.LoadHTMLGlob("web/templates/*.go.html")
	r.Static("/static", "web/static")

	r.GET("/v/:id", handlers.Redirect)
	r.GET("/", handlers.CreateView)

	api := r.Group("/api")
	{
		api.POST("/create", handlers.Create)
	}
	return
}

// InitAndStartRouter creates a Gin router with the default middleware, then
// registering all custom templates, middleware and routes. It continues by
// running the application on the port specified in the environment variable.
func InitAndStartRouter() {
	r := InitRouter()

	portStr := ":" + os.Getenv("PORT")
	if res, err := strconv.ParseBool(os.Getenv("TLS_ENABLED")); res && err == nil {
		r.RunTLS(portStr, os.Getenv("TLS_CERT"), os.Getenv("TLS_KEY"))
	} else {
		r.Run(portStr)
	}
}
