package go_auth

type MapUserTypeTrusted map[string]bool
type MapPublicRoutes map[string]bool
type MapUserTypeRoutes map[string][]string
type MapPermissionRoutes map[string][]string
type MapScopeRoutes map[string][]string

type Config struct {
	mapUserTypeTrusted  MapUserTypeTrusted
	mapPublicRoutes     MapPublicRoutes
	mapUserTypeRoutes   MapUserTypeRoutes
	mapPermissionRoutes MapPermissionRoutes
	mapScopeRoutes      MapScopeRoutes
	routeService        RouteService
}

type ConfigFunc func(c *Config)

func PublicRoutes(r MapPublicRoutes) ConfigFunc {
	return func(c *Config) {
		c.mapPublicRoutes = r
	}
}

func UserTypeTrusted(r MapUserTypeTrusted) ConfigFunc {
	return func(c *Config) {
		c.mapUserTypeTrusted = r
	}
}

func UserTypeRoutes(r MapUserTypeRoutes) ConfigFunc {
	return func(c *Config) {
		c.mapUserTypeRoutes = r
	}
}

func PermissionRoutes(r MapPermissionRoutes) ConfigFunc {
	return func(c *Config) {
		c.mapPermissionRoutes = r
	}
}

func ScopeRoutes(r MapScopeRoutes) ConfigFunc {
	return func(c *Config) {
		c.mapScopeRoutes = r
	}
}

func SetRouteService(r RouteService) ConfigFunc {
	return func(c *Config) {
		c.routeService = r
	}
}

func generate(args ...ConfigFunc) *Config {
	c := &Config{}
	for i := range args {
		args[i](c)
	}
	if c.mapPublicRoutes == nil {
		c.mapPublicRoutes = make(MapPublicRoutes)
	}

	// register route grpc health check
	c.mapPublicRoutes["/grpc.health.v1.Health/Check"] = true
	c.mapPublicRoutes["/grpc.health.v1.Health/Watch"] = true

	return c
}
