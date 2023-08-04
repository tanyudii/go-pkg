package grpc_auth

import (
	"context"
	"google.golang.org/grpc"
	goauth "pkg.tanyudii.me/go-pkg/go-auth"
	goerr "pkg.tanyudii.me/go-pkg/go-err"
	gotex "pkg.tanyudii.me/go-pkg/go-tex"
	"strings"
)

var (
	ErrUnauthenticated = goerr.NewUnauthenticatedErrorWithName("unauthenticated", "UNAUTHENTICATED")
)

type service struct {
	authService  goauth.Service
	tokenService goauth.TokenService
	cfg          *Config
}

func newService(
	authService goauth.Service,
	tokenService goauth.TokenService,
	args ...ConfigFunc,
) *service {
	return &service{
		authService:  authService,
		tokenService: tokenService,
		cfg:          generate(args...),
	}
}

func (s *service) authenticate(ctx context.Context, info *grpc.UnaryServerInfo) (context.Context, error) {
	//skip when route is public routes
	ok, err := s.authService.IsPublicRoute(info.FullMethod)
	if err != nil {
		return nil, err
	} else if ok {
		return ctx, nil
	}

	if newCtx, ok := s.authorizedInternalCall(ctx); ok {
		return newCtx, nil
	}

	newCtx, err := s.authenticateBearer(ctx)
	if err != nil {
		return nil, err
	}

	return s.authService.Authenticate(newCtx, info.FullMethod)
}

func (s *service) authenticateBearer(ctx context.Context) (context.Context, error) {
	md := gotex.FromIncoming(ctx)
	token := md.Get(strings.ToLower(gotex.RequestHeaderKeyAuthorization))
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
	md.Set(strings.ToLower(gotex.RequestHeaderKeyScopes), respToken.Scope)
	md.Set(strings.ToLower(gotex.RequestHeaderKeyAuthorization), token)
	if ti := respToken.TokenInfo; ti != nil {
		md.Set(strings.ToLower(gotex.RequestHeaderKeyUserID), ti.UserID)
		md.Set(strings.ToLower(gotex.RequestHeaderKeyUserName), ti.UserName)
		md.Set(strings.ToLower(gotex.RequestHeaderKeyUserEmail), ti.UserEmail)
		md.Set(strings.ToLower(gotex.RequestHeaderKeyUserType), ti.UserType)
		md.Set(strings.ToLower(gotex.RequestHeaderKeyCompanyID), ti.CompanyID)
		md.Set(strings.ToLower(gotex.RequestHeaderKeyCompanyName), ti.CompanyName)
		md.Set(strings.ToLower(gotex.RequestHeaderKeyPermissions), strings.Join(ti.Permissions, gotex.PermissionSeparator))
	}
	if ci := respToken.ClientInfo; ci != nil {
		md.Set(strings.ToLower(gotex.RequestHeaderKeyClientID), ci.ClientID)
		md.Set(strings.ToLower(gotex.RequestHeaderKeyClientName), ci.ClientName)
	}
	return md.ToIncoming(gotex.NewContext(ctx, gotex.NewGotex(md))), nil
}

func (s *service) authorizedInternalCall(ctx context.Context) (context.Context, bool) {
	md := gotex.FromIncoming(ctx)
	gtx := gotex.NewGotex(md)
	return md.ToIncoming(gotex.NewContext(ctx, gtx)),
		// is internal call when internal call password is not empty and
		// internal call password is equal with internal call password from request
		s.cfg.InternalCallPassword != "" && gtx.InternalCallPassword == s.cfg.InternalCallPassword
}
