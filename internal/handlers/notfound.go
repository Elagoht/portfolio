package handlers

import (
	"net/http"

	"statigo/framework/middleware"
	"statigo/framework/templates"
)

type NotFoundHandler struct {
	renderer *templates.Renderer
}

func NewNotFoundHandler(renderer *templates.Renderer) *NotFoundHandler {
	return &NotFoundHandler{
		renderer: renderer,
	}
}

func (h *NotFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lang := middleware.GetLanguage(r.Context())
	if lang == "" {
		lang = "en"
	}

	w.WriteHeader(http.StatusNotFound)

	data := BaseData(lang)
	data["Title"] = h.renderer.GetTranslation(lang, "pages.notfound.title")
	data["Content"] = map[string]string{
		"heading": h.renderer.GetTranslation(lang, "pages.notfound.heading"),
		"message": h.renderer.GetTranslation(lang, "pages.notfound.message"),
		"action":  h.renderer.GetTranslation(lang, "pages.notfound.action"),
	}

	h.renderer.Render(w, "notfound.html", data)
}
