package router

import (
	"html/template"
)

// SEOHelpers provides template functions for SEO optimization.
type SEOHelpers struct {
	registry  *Registry
	deployURL string // Base URL, e.g., "https://example.com"
}

// NewSEOHelpers creates a new SEO helpers instance.
func NewSEOHelpers(registry *Registry, deployURL string) *SEOHelpers {
	return &SEOHelpers{
		registry:  registry,
		deployURL: deployURL,
	}
}

// GetCanonicalURL returns the full canonical URL for the current page.
func (sh *SEOHelpers) GetCanonicalURL(canonical string, _ string) string {
	// Try to look up the route in the registry
	if route := sh.registry.GetByCanonical(canonical); route != nil {
		return sh.deployURL + route.Path
	}

	// If not found in registry, just prepend the deploy URL to the canonical path
	if canonical != "" && canonical[0] == '/' {
		return sh.deployURL + canonical
	}

	// Fallback to base deploy URL
	return sh.deployURL
}

// GetAlternateLinks returns empty string (no alternate links needed for single language).
func (sh *SEOHelpers) GetAlternateLinks(_ string) template.HTML {
	return ""
}

// GetAlternateURLs returns nil (no alternate URLs needed for single language).
func (sh *SEOHelpers) GetAlternateURLs(_ string) map[string]string {
	return nil
}

// GetLocalePath returns the URL path for a canonical path.
func (sh *SEOHelpers) GetLocalePath(canonical string, _ string) string {
	if route := sh.registry.GetByCanonical(canonical); route != nil {
		return route.Path
	}
	return "/"
}

// SEOFunctions holds SEO-related template functions.
// This struct is used to pass SEO functions to the template renderer.
type SEOFunctions struct {
	CanonicalURL   func(canonical, lang string) string
	AlternateLinks func(canonical string) template.HTML
	AlternateURLs  func(canonical string) map[string]string
	LocalePath     func(canonical, lang string) string
}

// ToTemplateFunctions converts SEOHelpers to a SEOFunctions struct.
func (sh *SEOHelpers) ToTemplateFunctions() *SEOFunctions {
	return &SEOFunctions{
		CanonicalURL:   sh.GetCanonicalURL,
		AlternateLinks: sh.GetAlternateLinks,
		AlternateURLs:  sh.GetAlternateURLs,
		LocalePath:     sh.GetLocalePath,
	}
}
