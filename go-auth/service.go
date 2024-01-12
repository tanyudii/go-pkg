package go_auth

import (
	"context"
	goerr "pkg.tanyudii.me/go-pkg/go-err"
	gotex "pkg.tanyudii.me/go-pkg/go-tex"
)

type Service interface {
	IsPublicRoute(fullMethod string) (bool, error)
	Authenticate(ctx context.Context, fullMethod string) (context.Context, error)
}

type service struct {
	cfg *Config
}

func NewService(args ...ConfigFunc) Service {
	return &service{
		cfg: generate(args...),
	}
}

func (s *service) IsPublicRoute(fullMethod string) (bool, error) {
	return s.cfg.mapPublicRoutes[fullMethod], nil
}

func (s *service) Authenticate(ctx context.Context, fullMethod string) (context.Context, error) {
	session, err := gotex.FromContextWithErr(ctx)
	if err != nil {
		return nil, err
	}

	routeUserTypes := s.cfg.mapUserTypeRoutes[fullMethod]
	routePermissions := s.cfg.mapPermissionRoutes[fullMethod]
	routeScopes := s.cfg.mapScopeRoutes[fullMethod]

	routeConfig, err := s.getRouteConfig(ctx, fullMethod)
	if err != nil {
		return nil, err
	} else if routeConfig != nil {
		routeUserTypes = append(routeUserTypes, routeConfig.GetUserTypes()...)
		routePermissions = append(routePermissions, routeConfig.GetPermissions()...)
		routeScopes = append(routeScopes, routeConfig.GetScopes()...)
	}
	if err = s.checkRouterPermission(ctx, fullMethod); err != nil {
		return nil, err
	}

	//if user authorized with type, will be skip other middleware
	ok, err := s.authorizedUserType(session, routeUserTypes)
	if err != nil {
		return nil, goerr.NewUnauthorizedError(err.Error())
	} else if ok {
		return ctx, nil
	}

	if err = s.authorizedPermission(session, routePermissions); err != nil {
		return nil, goerr.NewUnauthorizedError(err.Error())
	}

	if err = s.authorizedScope(session, routeScopes); err != nil {
		return nil, goerr.NewUnauthorizedError(err.Error())
	}

	return ctx, nil
}

func (s *service) authorizedUserType(session *gotex.Gotex, userTypes []string) (bool, error) {
	ok, err := session.HasUserTypeByMapCode(s.cfg.mapUserTypeTrusted)
	if err != nil {
		return false, err
	} else if ok {
		return ok, nil
	}
	return session.HasUserType(userTypes)
}

func (s *service) authorizedPermission(session *gotex.Gotex, permissions []string) error {
	_, err := session.HasPermission(permissions)
	return err
}

func (s *service) authorizedScope(session *gotex.Gotex, scopes []string) error {
	_, err := session.HasScope(scopes)
	return err
}

func (s *service) getRouteConfig(ctx context.Context, fm string) (RouteConfig, error) {
	if s.cfg.routeService == nil {
		return nil, nil
	}
	return s.cfg.routeService.GetRouteConfig(ctx, fm)
}

func (s *service) checkRouterPermission(ctx context.Context, fm string) error {
	if s.cfg.routeService == nil {
		return nil
	}
	return s.cfg.routeService.CheckRoutePermission(ctx, fm)
}
