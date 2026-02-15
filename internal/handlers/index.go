package handlers

import (
	"net/http"

	"statigo/framework/middleware"
	"statigo/framework/router"
	"statigo/framework/templates"
)

type IndexHandler struct {
	renderer *templates.Renderer
	registry *router.Registry
}

func NewIndexHandler(renderer *templates.Renderer, registry *router.Registry) *IndexHandler {
	return &IndexHandler{
		renderer: renderer,
		registry: registry,
	}
}

type TechGroup struct {
	Title        string
	Technologies []string
}

type BlogCategory struct {
	Name  string
	Count int
	Href  string
}

type Stat struct {
	Number string
	Label  string
}

type ExperienceItem struct {
	Title       string
	Company     string
	Date        string
	Description []string
}

type Education struct {
	University string
	Programme  string
	Date       string
}

func (h *IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lang := middleware.GetLanguage(r.Context())
	canonical := router.GetCanonicalPath(r.Context())
	t := func(key string) string {
		return h.renderer.GetTranslation(lang, key)
	}

	data := BaseData(lang)
	data["Canonical"] = canonical
	data["Title"] = SiteName
	data["Meta"] = map[string]string{
		"description": t("hero.subtitle"),
	}
	data["Titles"] = []string{"Lead Product Developer", "Fullstack Developer"}
	data["Stats"] = []Stat{
		{Number: "160+", Label: t("stats.blogPosts")},
		{Number: "90+", Label: t("stats.youtubeVideos")},
		{Number: "12", Label: t("stats.languages")},
		{Number: "1", Label: t("stats.udemyCourse")},
	}
	data["BlogCategories"] = []BlogCategory{
		{Name: t("categories.software"), Count: 142, Href: "/blogs?category=software"},
		{Name: t("categories.music"), Count: 8, Href: "/blogs?category=music"},
		{Name: t("categories.techtales"), Count: 5, Href: "/blogs?category=techtales"},
		{Name: t("categories.myLife"), Count: 3, Href: "/blogs?category=my-life"},
		{Name: t("categories.uxui"), Count: 2, Href: "/blogs?category=ux-ui"},
	}
	h.renderer.Render(w, "index.html", data)
}
