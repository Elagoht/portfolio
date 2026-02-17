package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"statigo/internal/services"
)

// ViewsCache holds cached view counts with expiration.
type ViewsCache struct {
	mu         sync.RWMutex
	data       map[string]int
	expiresAt  time.Time
	ttl        time.Duration
	bloggo     *services.BloggoService
}

// NewViewsCache creates a new views cache.
func NewViewsCache(bloggo *services.BloggoService, ttl time.Duration) *ViewsCache {
	return &ViewsCache{
		data:   make(map[string]int),
		ttl:    ttl,
		bloggo: bloggo,
	}
}

// Get retrieves view counts, refreshing from API if cache is expired.
func (vc *ViewsCache) Get() (map[string]int, error) {
	vc.mu.RLock()
	needsRefresh := time.Now().After(vc.expiresAt)
	currentData := vc.data
	vc.mu.RUnlock()

	if !needsRefresh && len(currentData) > 0 {
		return currentData, nil
	}

	// Cache is empty or expired, fetch fresh data
	vc.mu.Lock()
	defer vc.mu.Unlock()

	// Double-check after acquiring write lock
	if !time.Now().After(vc.expiresAt) && len(vc.data) > 0 {
		return vc.data, nil
	}

	// Fetch from API
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	freshData, err := vc.bloggo.GetViewCounts(ctx)
	if err != nil {
		// Return stale data if available
		if len(vc.data) > 0 {
			return vc.data, nil
		}
		return nil, err
	}

	vc.data = freshData
	vc.expiresAt = time.Now().Add(vc.ttl)

	return vc.data, nil
}

// Invalidate clears the cached data.
func (vc *ViewsCache) Invalidate() {
	vc.mu.Lock()
	defer vc.mu.Unlock()
	vc.data = make(map[string]int)
	vc.expiresAt = time.Time{}
}

// ViewsHandler handles the /api/posts/views endpoint.
type ViewsHandler struct {
	Cache *ViewsCache
}

// NewViewsHandler creates a new views handler.
func NewViewsHandler(bloggo *services.BloggoService, ttl time.Duration) *ViewsHandler {
	return &ViewsHandler{
		Cache: NewViewsCache(bloggo, ttl),
	}
}

// GetSlug returns view count for a specific slug.
func (h *ViewsHandler) GetSlug(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Extract slug from path /api/posts/views/{slug}
	slug := r.URL.Path
	slug = strings.TrimPrefix(slug, "/api/posts/views/")
	slug = strings.TrimSuffix(slug, "/")

	if slug == "" || slug == "views" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "slug is required",
		})
		return
	}

	views, err := h.Cache.Get()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "failed to fetch view counts",
		})
		return
	}

	count, exists := views[slug]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]any{
			"slug":  slug,
			"views": 0,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"slug":  slug,
		"views": count,
	})
}

// InvalidateCache triggers a cache refresh (useful for webhooks).
func (h *ViewsHandler) InvalidateCache() {
	h.Cache.Invalidate()
}
