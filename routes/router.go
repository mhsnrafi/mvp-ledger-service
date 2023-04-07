package routes

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"ledger-service/docs"
	"ledger-service/middlewares"
	"ledger-service/models"
	"ledger-service/services"
	"net/http"
)

func New() *gin.Engine {
	r := gin.New()
	initRoute(r)

	r.Use(gin.LoggerWithWriter(middlewares.LogWriter()))
	r.Use(gin.CustomRecovery(middlewares.AppRecovery()))
	r.Use(middlewares.CORSMiddleware())
	r.Use(middlewares.ResponseTimeMiddleware())
	r.Use(middlewares.Pagination())

	v1 := r.Group("/v1")
	{
		AuthRoute(v1)
		Legder(v1)

	}

	docs.SwaggerInfo.BasePath = v1.BasePath()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	// Add Prometheus endpoint for metrics collection
	r.GET("/metrics", func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	})

	// Add pprof endpoint for profiling
	pprof.Register(r)

	return r
}

func initRoute(r *gin.Engine) {
	_ = r.SetTrustedProxies(nil)
	r.RedirectTrailingSlash = false
	r.HandleMethodNotAllowed = true

	r.NoRoute(func(c *gin.Context) {
		models.SendErrorResponse(c, http.StatusNotFound, c.Request.RequestURI+" not found")
	})

	r.NoMethod(func(c *gin.Context) {
		models.SendErrorResponse(c, http.StatusMethodNotAllowed, c.Request.Method+" is not allowed here")
	})
}

func InitGin() {
	gin.DisableConsoleColor()
	gin.SetMode(services.Config.Mode)
}
