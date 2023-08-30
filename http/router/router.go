package router

import (
	"github.com/labstack/echo/v4"
)

// Route is the information for every URI.
type Route struct {
	// Name is the name of this Route.
	Name string

	// Method is the string for the HTTP method. ex) GET, POST etc..
	Method string

	// Pattern is the pattern of the URI.
	Pattern string

	// HandlerFunc is the handler function of this route.
	HandlerFunc echo.HandlerFunc

	// MiddlewareFunc is route-level middleware.
	MiddlewareFunc []echo.MiddlewareFunc
}

type Router struct {
	Routes []*Route
}

func (r *Router) FindRouteByName(name string) *Route {
	for _, route := range r.Routes {
		if route.Name == name {
			return route
		}
	}
	return nil
}

// Register routes with Echo.
func Register(e *echo.Echo, router *Router) {
	for _, route := range router.Routes {
		e.Add(route.Method, route.Pattern, route.HandlerFunc, route.MiddlewareFunc...)
	}
}
