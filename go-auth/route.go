package go_auth

import "context"

type RouteService interface {
	GetRouteConfig(ctx context.Context, fm string) (RouteConfig, error)
	CheckRoutePermission(ctx context.Context, fm string) error
}

type RouteConfig interface {
	GetPermissions() []string
	GetScopes() []string
	GetUserTypes() []string
}
