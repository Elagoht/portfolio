// Package router provides URL routing with canonical path management.
package router

import (
	"net/http"

	"github.com/go-chi/chi"
)

// RouteDefinition represents a route definition.
type RouteDefinition struct {
	Canonical string           // Canonical path, e.g., "/features"
	Path      string           // URL path: "/features"
	Handler   http.HandlerFunc // Handler function for this route
	Template  string           // Template name (e.g., "content.html")
	Title     string           // Translation key for page title (e.g., "main.title")
	Strategy  string           // Caching strategy: "static", "incremental", "dynamic", "immutable"
	Interval  string           // Revalidation interval for incremental strategy (e.g., "24h")
}

// Registry maintains the mapping between canonical paths and route definitions.
type Registry struct {
	routes       []RouteDefinition
	pathToRoute  map[string]*RouteDefinition // Maps actual paths to route definitions
	canonicalMap map[string]*RouteDefinition // Maps canonical paths to route definitions
}

// NewRegistry creates a new route registry.
func NewRegistry() *Registry {
	return &Registry{
		routes:       make([]RouteDefinition, 0),
		pathToRoute:  make(map[string]*RouteDefinition),
		canonicalMap: make(map[string]*RouteDefinition),
	}
}

// AddRoute registers a new route definition.
func (r *Registry) AddRoute(def RouteDefinition) error {
	// Store in registry
	r.routes = append(r.routes, def)
	routePtr := &r.routes[len(r.routes)-1]
	r.canonicalMap[def.Canonical] = routePtr

	// Map path to this definition
	r.pathToRoute[def.Path] = routePtr
	// Also map the path with trailing slash (unless it's the root path)
	if def.Path != "/" {
		r.pathToRoute[def.Path+"/"] = routePtr
	}

	return nil
}

// GetByPath returns the route definition for a given path.
func (r *Registry) GetByPath(path string) *RouteDefinition {
	return r.pathToRoute[path]
}

// GetByCanonical returns the route definition for a canonical path.
func (r *Registry) GetByCanonical(canonical string) *RouteDefinition {
	return r.canonicalMap[canonical]
}

// GetAll returns all registered routes.
func (r *Registry) GetAll() []RouteDefinition {
	return r.routes
}

// RegisterRoutes automatically registers all routes from the registry with a chi router.
// The canonicalMiddleware wraps each handler to inject canonical path context.
func (r *Registry) RegisterRoutes(router chi.Router, canonicalMiddleware func(http.Handler) http.Handler) {
	for _, route := range r.routes {
		// Wrap the handler once per route with canonical middleware
		wrappedHandler := canonicalMiddleware(route.Handler)

		// Convert to HandlerFunc
		handlerFunc := func(w http.ResponseWriter, req *http.Request) {
			wrappedHandler.ServeHTTP(w, req)
		}

		// Register the path
		router.Get(route.Path, handlerFunc)

		// Also register with trailing slash
		if route.Path != "/" {
			router.Get(route.Path+"/", handlerFunc)
		}
	}
}
