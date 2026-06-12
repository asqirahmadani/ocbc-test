package router

import (
	"boilerplate/config"
	"boilerplate/middleware"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg config.Config) *gin.Engine {
	r := gin.New()

	r.Use(
		middleware.CORSMiddleware(),
		gin.Recovery(),
		middleware.LoggerMiddleware(),
		middleware.ErrorMiddleware(),
		middleware.TimeoutMiddleware(time.Duration(cfg.App.RequestTimeout)*time.Second),
	)

	root := r.Group("/api/v1")
	{
		root.GET("")
	}

	return r
}