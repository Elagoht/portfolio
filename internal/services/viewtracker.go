package services

import (
	"net"
	"net/http"
	"sync"
	"time"
)

// ViewTracker tracks blog post views with IP-based rate limiting.
type ViewTracker struct {
	mu           sync.RWMutex
	views        map[string]time.Time // key: "ip:slug", value: last view time
	cleanupTimer *time.Timer
	logger       Logger
}

// Logger interface for view tracker logging.
type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

// NewViewTracker creates a new view tracker.
func NewViewTracker(logger Logger) *ViewTracker {
	vt := &ViewTracker{
		views:  make(map[string]time.Time),
		logger: logger,
	}

	// Start periodic cleanup of old entries
	vt.startCleanup()

	return vt
}

const viewCooldown = 10 * time.Minute
const cleanupInterval = 5 * time.Minute

// ShouldTrackView returns true if the view should be tracked (not within cooldown period).
func (vt *ViewTracker) ShouldTrackView(r *http.Request, slug string) bool {
	// Get client IP
	ip := getClientIP(r)
	if ip == "" {
		vt.logger.Debug("no client IP found, skipping view tracking")
		return false
	}

	key := ip + ":" + slug

	vt.mu.RLock()
	lastView, exists := vt.views[key]
	vt.mu.RUnlock()

	if exists {
		if time.Since(lastView) < viewCooldown {
			vt.logger.Debug("view within cooldown period, skipping",
				"ip", ip,
				"slug", slug,
				"remaining", viewCooldown-time.Since(lastView))
			return false
		}
	}

	// Mark this view
	vt.mu.Lock()
	vt.views[key] = time.Now()
	vt.mu.Unlock()

	return true
}

// getClientIP extracts the client IP from the request, checking headers for proxies.
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (for proxies/load balancers)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one (original client)
		if idx := len(xff); idx > 0 {
			if ip, _, err := net.SplitHostPort(xff); err == nil {
				return ip
			}
			return xff
		}
	}

	// Check X-Real-IP header (common with nginx)
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		if ip, _, err := net.SplitHostPort(xri); err == nil {
			return ip
		}
		return xri
	}

	// Fall back to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

// startCleanup begins the periodic cleanup of old entries.
func (vt *ViewTracker) startCleanup() {
	vt.cleanupTimer = time.AfterFunc(cleanupInterval, func() {
		vt.cleanup()
		vt.startCleanup() // Reschedule
	})
}

// cleanup removes entries older than the cooldown period.
func (vt *ViewTracker) cleanup() {
	vt.mu.Lock()
	defer vt.mu.Unlock()

	cutoff := time.Now().Add(-viewCooldown)
	removed := 0

	for key, lastView := range vt.views {
		if lastView.Before(cutoff) {
			delete(vt.views, key)
			removed++
		}
	}

	if removed > 0 {
		vt.logger.Debug("cleaned up old view entries", "removed", removed, "remaining", len(vt.views))
	}
}

// Stop stops the cleanup timer.
func (vt *ViewTracker) Stop() {
	if vt.cleanupTimer != nil {
		vt.cleanupTimer.Stop()
	}
}
