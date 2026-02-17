package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"statigo/framework/cache"
)

// WebhookPayload matches the payload structure sent by Bloggo CMS.
type WebhookPayload struct {
	Event     string                 `json:"event"`
	Entity    string                 `json:"entity"`
	ID        *int64                 `json:"id"`
	Slug      *string                `json:"slug"`
	OldSlug   *string                `json:"oldSlug,omitempty"`
	Action    string                 `json:"action"`
	Timestamp string                 `json:"timestamp"`
	Data      map[string]any `json:"data"`
}

// WebhookHandler handles incoming webhooks from Bloggo CMS for cache invalidation.
type WebhookHandler struct {
	cacheManager *cache.Manager
	viewsHandler *ViewsHandler
	logger       *slog.Logger
}

// NewWebhookHandler creates a new webhook handler.
func NewWebhookHandler(cacheManager *cache.Manager, viewsHandler *ViewsHandler, logger *slog.Logger) *WebhookHandler {
	return &WebhookHandler{
		cacheManager: cacheManager,
		viewsHandler: viewsHandler,
		logger:       logger,
	}
}

// Handle processes incoming webhook events from Bloggo CMS.
func (h *WebhookHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var payload WebhookPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.logger.Warn("webhook: invalid payload", slog.String("error", err.Error()))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid payload"})
		return
	}

	h.logger.Info("webhook received",
		slog.String("event", payload.Event),
		slog.String("entity", payload.Entity),
		slog.String("action", payload.Action),
	)

	invalidated := 0

	switch payload.Entity {
	case "post":
		// Post changes affect blog listing and home page (recent posts)
		if h.cacheManager != nil {
			invalidated = h.cacheManager.MarkStale("static", true)
		}
		// Also invalidate views cache since post list may change
		if h.viewsHandler != nil {
			h.viewsHandler.InvalidateCache()
		}

	case "category", "tag":
		// Category/tag changes affect blog listing filters
		if h.cacheManager != nil {
			invalidated = h.cacheManager.MarkStale("static", true)
		}

	case "author":
		// Author data appears on cached blog post detail pages
		if h.cacheManager != nil {
			invalidated = h.cacheManager.MarkStale("static", true)
		}

	case "keyvalue":
		// Site-wide config changes, invalidate everything
		if h.cacheManager != nil {
			invalidated = h.cacheManager.MarkAllStale(true)
		}

	case "cms":
		// Manual sync - full invalidation
		if h.cacheManager != nil {
			invalidated = h.cacheManager.MarkAllStale(true)
		}
		if h.viewsHandler != nil {
			h.viewsHandler.InvalidateCache()
		}

	default:
		h.logger.Warn("webhook: unknown entity", slog.String("entity", payload.Entity))
	}

	h.logger.Info("webhook processed",
		slog.String("event", payload.Event),
		slog.Int("invalidated", invalidated),
	)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"success":     true,
		"invalidated": invalidated,
	})
}
