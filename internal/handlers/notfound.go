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

	t := func(key string) string {
		return h.renderer.GetTranslation(lang, key)
	}

	data := BaseData(lang, t)
	data["Title"] = t("pages.notfound.title")
	data["Content"] = map[string]string{
		"heading": t("pages.notfound.heading"),
		"message": t("pages.notfound.message"),
		"action":  t("pages.notfound.action"),
	}

	h.renderer.Render(w, "notfound.html", data)
}
