package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const loggerMessage = "http_request"

func RequestLogger(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		// Process request
		c.Next()

		// Skip logging static assets
		if strings.Contains(path, "/static/") {
			return
		}

		status := c.Writer.Status()
		end := time.Now()
		latency := end.Sub(start)

		fields := logrus.Fields{
			"status":      status,
			"method":      c.Request.Method,
			"path":        path,
			"remote_addr": c.ClientIP(),
			"duration":    latency,
		}

		if err := c.Errors.Last(); err != nil {
			fields["error"] = err.Error()
		}

		entry := logrus.WithFields(fields)

		switch {
		case status >= http.StatusBadRequest && status < http.StatusInternalServerError:
			entry.Warn(loggerMessage)
		case status >= http.StatusInternalServerError:
			entry.Error(loggerMessage)
		default:
			entry.Info(loggerMessage)
		}
	}
}