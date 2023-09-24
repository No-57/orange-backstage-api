package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"orange-backstage-api/infra/util/convert"
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func Logger(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// Ignore health-check to avoid polluting API logs
		if path == "/api/v1/health" {
			c.Next()
			return
		}

		start := time.Now()
		rawQuery := c.Request.URL.RawQuery
		requestID := requestid.Get(c)

		logger := zerolog.Ctx(ctx).With().
			Str("req_id", requestID).
			Str("path", path).
			Str("method", c.Request.Method).
			Logger()
		c.Request = c.Request.WithContext(logger.WithContext(ctx))

		var buf bytes.Buffer
		tee := io.TeeReader(c.Request.Body, &buf)
		c.Request.Body = io.NopCloser(tee)

		c.Next()

		end := time.Now()
		status := c.Writer.Status()

		logLevel := zerolog.InfoLevel
		switch {
		case status >= http.StatusInternalServerError:
			logLevel = zerolog.ErrorLevel
		case status >= http.StatusBadRequest:
			logLevel = zerolog.WarnLevel
		}

		l := logger.WithLevel(logLevel).
			Str("ip", c.ClientIP()).
			Str("user_agent", c.Request.UserAgent()).
			Str("query", rawQuery).
			Int("status", status).
			Time("callTime", end).
			Dur("duration", end.Sub(start)).
			Str("errors", c.Errors.ByType(gin.ErrorTypePrivate).String())

		if buf.Len() > 0 {
			data := buf.Bytes()

			data = filterSensitiveAPI(path, data)

			var jsonBuf bytes.Buffer
			if err := json.Compact(&jsonBuf, data); err == nil {
				l = l.RawJSON("body", jsonBuf.Bytes())
			}
		}

		l.Send()
	}
}

var sensitiveAPIs = map[string]struct{}{
	"/api/v1/login": {},
}

func filterSensitiveAPI(path string, data []byte) []byte {
	_, ok := sensitiveAPIs[path]
	if ok {
		return convert.StrToBytes(`{"data":"removed due to sensitive data"}`)
	}

	return data
}
