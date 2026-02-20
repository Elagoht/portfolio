package router

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"

	fwctx "statigo/framework/context"
	"statigo/framework/templates"
)

// RouteConfig represents a single route configuration from JSON.
type RouteConfig struct {
	Canonical string `json:"canonical"`
	Path      string `json:"path"`
	Template  string `json:"template"`
	Handler   string `json:"handler"`  // Handler name (e.g., "index", "content")
	Title     string `json:"title"`    // Translation key for page title
	Strategy  string `json:"strategy"` // Caching strategy: "static", "incremental", "dynamic", "immutable"
	Interval  string `json:"interval"` // Revalidation interval for incremental strategy (e.g., "24h")
}

// RoutesConfig represents the complete routes configuration file.
type RoutesConfig struct {
	Routes []RouteConfig `json:"routes"`
}

// LoadRoutesFromJSON loads route configurations from a JSON file.
func LoadRoutesFromJSON(
	configFS fs.FS,
	filePath string,
	registry *Registry,
	renderer *templates.Renderer,
	customHandlers map[string]http.HandlerFunc,
	logger *slog.Logger,
) error {
	// Read JSON file
	data, err := fs.ReadFile(configFS, filePath)
	if err != nil {
		return fmt.Errorf("failed to read routes file: %w", err)
	}

	// Parse JSON
	var config RoutesConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse routes JSON: %w", err)
	}

	logger.Info("Loading routes from JSON", "file", filePath, "routes", len(config.Routes))

	// Register each static route
	for _, routeConfig := range config.Routes {
		// Strategy-only entries (no handler) are stored in the registry for
		// wildcard pattern lookups but are not registered as chi routes.
		if routeConfig.Handler == "" {
			if err := registry.AddRoute(RouteDefinition{
				Path:     routeConfig.Path,
				Strategy: routeConfig.Strategy,
			}); err != nil {
				return fmt.Errorf("failed to add route %s: %w", routeConfig.Path, err)
			}
			logger.Debug("Registered strategy-only route",
				"path", routeConfig.Path,
				"strategy", routeConfig.Strategy)
			continue
		}

		var handler http.HandlerFunc

		// Determine which handler to use
		switch routeConfig.Handler {
		case "content":
			// Create inline content handler
			templateName := routeConfig.Template
			handler = func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				layoutData := fwctx.GetLayoutData(ctx)
				canonical := fwctx.GetCanonicalPath(ctx)

				renderer.Render(w, templateName, map[string]interface{}{
					"Lang":      "en",
					"Data":      map[string]interface{}{},
					"Layout":    layoutData,
					"Canonical": canonical,
					"WebAppURL": os.Getenv("WEBAPP_URL"),
				})
			}
		default:
			// Check if custom handler exists
			if customHandlers != nil {
				if customHandler, exists := customHandlers[routeConfig.Handler]; exists {
					handler = customHandler
				} else {
					// Fall back to content handler
					logger.Warn("Custom handler not found, using content handler",
						"handler", routeConfig.Handler,
						"canonical", routeConfig.Canonical)
					templateName := routeConfig.Template
					handler = func(w http.ResponseWriter, r *http.Request) {
						ctx := r.Context()
						layoutData := fwctx.GetLayoutData(ctx)
						canonical := fwctx.GetCanonicalPath(ctx)

						renderer.Render(w, templateName, map[string]interface{}{
							"Lang":      "en",
							"Data":      map[string]interface{}{},
							"Layout":    layoutData,
							"Canonical": canonical,
							"WebAppURL": os.Getenv("WEBAPP_URL"),
						})
					}
				}
			} else {
				// No custom handlers provided, use content handler
				templateName := routeConfig.Template
				handler = func(w http.ResponseWriter, r *http.Request) {
					ctx := r.Context()
					layoutData := fwctx.GetLayoutData(ctx)
					canonical := fwctx.GetCanonicalPath(ctx)

					renderer.Render(w, templateName, map[string]interface{}{
						"Lang":      "en",
						"Data":      map[string]interface{}{},
						"Layout":    layoutData,
						"Canonical": canonical,
						"WebAppURL": os.Getenv("WEBAPP_URL"),
					})
				}
			}
		}

		// Add route to registry
		if err := registry.AddRoute(RouteDefinition{
			Canonical: routeConfig.Canonical,
			Path:      routeConfig.Path,
			Handler:   handler,
			Template:  routeConfig.Template,
			Title:     routeConfig.Title,
			Strategy:  routeConfig.Strategy,
			Interval:  routeConfig.Interval,
		}); err != nil {
			return fmt.Errorf("failed to add route %s: %w", routeConfig.Canonical, err)
		}

		logger.Debug("Registered route",
			"canonical", routeConfig.Canonical,
			"handler", routeConfig.Handler,
			"template", routeConfig.Template)
	}

	logger.Info("Successfully loaded all routes", "routes", len(config.Routes))
	return nil
}
