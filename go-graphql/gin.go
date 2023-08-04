package go_graphql

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	gotex "pkg.tanyudii.me/go-pkg/go-tex"
	"time"
)

func GinCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func GinAcceptLanguage() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		gtx, ok := gotex.FromContext(ctx)
		if !ok {
			gtx = &gotex.Gotex{}
		}
		if gtx.AcceptLanguage == "" {
			gtx.AcceptLanguage = c.GetHeader("Accept-Language")
		}
		c.Request = c.Request.WithContext(gotex.NewContext(ctx, gtx))
		c.Next()
	}
}

func GinRequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		gtx, ok := gotex.FromContext(ctx)
		if !ok {
			gtx = &gotex.Gotex{}
		}
		if gtx.RequestID == "" {
			requestID := c.GetHeader(gotex.RequestHeaderKeyRequestID)
			if requestID == "" {
				requestID = fmt.Sprintf("%s-%d", uuid.NewString(), time.Now().Unix())
			}
			gtx.RequestID = requestID
		}
		c.Request = c.Request.WithContext(gotex.NewContext(ctx, gtx))
		c.Next()
	}
}
