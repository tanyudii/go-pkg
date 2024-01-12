package go_auth

import "context"

type RouteService interface {
	GetRouteConfig(ctx context.Context, fullMethod string) (RouteConfig, error)
}

type RouteConfig interface {
	GetPermissions() []string
	GetScopes() []string
	GetUserTypes() []string
}
