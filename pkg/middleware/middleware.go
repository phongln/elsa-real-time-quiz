package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		logrus.Infof("Handled request [%s] %s in %v", c.Request.Method, c.Request.URL.Path, duration)
	}
}

func MetricsHandler() http.Handler {
	return promhttp.Handler()
}
