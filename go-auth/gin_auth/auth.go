package gin_auth

import (
	"github.com/gin-gonic/gin"
	goauth "pkg.tanyudii.me/go-pkg/go-auth"
)

func Authenticate(
	authService goauth.Service,
	tokenService goauth.TokenService,
	args ...ConfigFunc,
) func(c *gin.Context) {
	svc := newService(authService, tokenService, args...)
	return func(c *gin.Context) {
		newCtx, err := svc.authenticate(c)
		if err != nil {
			c.Errors = append(c.Errors, c.Error(err))
			c.Abort()
		} else if newCtx != nil {
			c.Request = c.Request.WithContext(newCtx)
		}
		c.Next()
	}
}
