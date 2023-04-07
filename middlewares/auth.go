package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"ledger-service/models"
	"ledger-service/services"
	"net/http"
	"time"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Bearer-Token")
		tokenModel, err := services.VerifyToken(token, models.TokenTypeAccess)
		if err != nil {
			models.SendErrorResponse(c, http.StatusUnauthorized, err.Error())
			return
		}

		c.Set("userIdHex", tokenModel.ID)
		c.Set("userId", tokenModel.ID)

		c.Next()
	}
}

func ResponseTimeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		responseTime := prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "response_time_seconds",
				Help:    "Time taken to serve a request",
				Buckets: []float64{0.1, 0.5, 1.0, 2.5, 5.0, 10.0},
			},
			[]string{"method", "endpoint"},
		)

		start := time.Now()
		c.Next()
		elapsed := time.Since(start).Seconds()
		responseTime.WithLabelValues(c.Request.Method, c.FullPath()).Observe(elapsed)
	}
}
