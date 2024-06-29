package gzipmiddleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GzipMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Content-Encoding") == "gzip" {
			body, err := io.ReadAll(c.Request.Body)
			if err != nil {
				c.String(http.StatusInternalServerError, "Error reading request body")
				c.Abort()
				return
			}

			gzipReader, err := gzip.NewReader(bytes.NewReader(body))
			if err != nil {
				c.String(http.StatusInternalServerError, "Error creating gzip reader")
				c.Abort()
				return
			}
			defer gzipReader.Close()

			uncompressed, err := io.ReadAll(gzipReader)
			if err != nil {
				c.String(http.StatusInternalServerError, "Error reading uncompressed data")
				c.Abort()
				return
			}

			c.Request.Body = io.NopCloser(bytes.NewBuffer(uncompressed))
			c.Request.ContentLength = int64(len(uncompressed))
			c.Request.Header.Del("Content-Encoding")
		}

		c.Next()
	}
}
