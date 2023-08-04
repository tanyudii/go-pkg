package gin_auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	goauth "pkg.tanyudii.me/go-pkg/go-auth"
	goerr "pkg.tanyudii.me/go-pkg/go-err"
	gotex "pkg.tanyudii.me/go-pkg/go-tex"
	"strings"
)

var (
	ErrUnauthenticated = goerr.NewUnauthenticatedErrorWithName("unauthenticated", "UNAUTHENTICATED")
)

type Service interface {
	authenticate(c *gin.Context) (newCtx context.Context, err error)
	authenticateBearer(c *gin.Context) (context.Context, error)
}

type service struct {
	authService  goauth.Service
	tokenService goauth.TokenService
	cfg          *Config
}

func newService(
	authService goauth.Service,
	tokenService goauth.TokenService,
	args ...ConfigFunc,
) Service {
	return &service{
		authService:  authService,
		tokenService: tokenService,
		cfg:          generate(args...),
	}
}

func (s *service) authenticate(c *gin.Context) (context.Context, error) {
	fullMethod := fmt.Sprintf("[%s] %s", c.Request.Method, c.Request.RequestURI)

	//skip when route is public routes
	ok, err := s.authService.IsPublicRoute(fullMethod)
	if err != nil {
		return nil, err
	} else if ok {
		return c, nil
	}

	newCtx, err := s.authenticateBearer(c)
	if err != nil {
		// skip when is graphqlMode
		// GraphQL will validate on resolver
		if errors.Is(err, ErrUnauthenticated) && s.cfg.graphqlMode {
			return c, nil
		}
		return nil, err
	}

	return s.authService.Authenticate(newCtx, fullMethod)
}

func (s *service) authenticateBearer(c *gin.Context) (context.Context, error) {
	token := c.GetHeader("Authorization")
	if token == "" {
		return nil, ErrUnauthenticated
	}
	splitToken := strings.Split(token, "Bearer ")
	if len(splitToken) != 2 {
		return nil, ErrUnauthenticated
	}
	respToken, err := s.tokenService.TokenInfo(context.Background(), splitToken[1])
	if err != nil {
		return nil, err
	}
	gtx, ok := gotex.FromContext(c.Request.Context())
	if !ok {
		gtx = &gotex.Gotex{}
	}
	gtx.Scopes = respToken.Scope
	gtx.Authorization = token
	if ti := respToken.TokenInfo; ti != nil {
		gtx.UserID = ti.UserID
		gtx.UserName = ti.UserName
		gtx.UserEmail = ti.UserEmail
		gtx.UserType = ti.UserType
		gtx.CompanyID = ti.CompanyID
		gtx.CompanyName = ti.CompanyName
		gtx.Permissions = strings.Join(ti.Permissions, gotex.PermissionSeparator)
	}
	if ci := respToken.ClientInfo; ci != nil {
		gtx.ClientID = ci.ClientID
		gtx.ClientName = ci.ClientName
	}
	return gotex.NewContext(c.Request.Context(), gtx), nil

}
