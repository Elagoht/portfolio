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

type Hobby struct {
	Icon        string
	Title       string
	Description string
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
	data["HeroLinks"] = []Link{
		{Title: "GitHub", Href: "https://github.com/Elagoht"},
		{Title: "LinkedIn", Href: "https://linkedin.com/in/furkan-baytekin"},
		{Title: "YouTube", Href: "https://youtube.com/@furkanbytekin"},
		{Title: "X", Href: "https://x.com/furkanbytekin"},
	}
	data["Stats"] = []Stat{
		{Number: "160+", Label: t("stats.blogPosts")},
		{Number: "90+", Label: t("stats.youtubeVideos")},
		{Number: "12", Label: t("stats.languages")},
		{Number: "1", Label: t("stats.udemyCourse")},
	}
	data["AboutHeading"] = t("index.aboutHeading")
	data["AboutText"] = t("index.aboutText")
	data["Hobbies"] = []Hobby{
		{Icon: "ti ti-player-record", Title: t("index.hobbyVinylTitle"), Description: t("index.hobbyVinylDesc")},
		{Icon: "ti ti-guitar-pick", Title: t("index.hobbyBassTitle"), Description: t("index.hobbyBassDesc")},
		{Icon: "ti ti-book", Title: t("index.hobbyBooksTitle"), Description: t("index.hobbyBooksDesc")},
	}
	h.renderer.Render(w, "index.html", data)
}
