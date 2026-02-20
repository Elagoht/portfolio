package handlers

import (
	"net/http"

	fwctx "statigo/framework/context"
	"statigo/framework/templates"
)

type IndexHandler struct {
	renderer *templates.Renderer
}

func NewIndexHandler(renderer *templates.Renderer) *IndexHandler {
	return &IndexHandler{
		renderer: renderer,
	}
}

type TechGroup struct {
	Title        string
	Technologies []string
}

type BlogCategory struct {
	Slug   string
	Name   string
	Count  int
	Href   string
	Active bool
}

type BlogTag struct {
	Slug   string
	Name   string
	Count  int
	Href   string
	Active bool
}

type Stat struct {
	Number string
	Label  string
	Href   string
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
	const lang = "en"
	canonical := fwctx.GetCanonicalPath(r.Context())
	t := func(key string) string {
		return h.renderer.GetTranslation(lang, key)
	}

	data := BaseData(lang, t)
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
		{Number: "160+", Label: t("stats.blogPosts"), Href: "/blogs"},
		{Number: "90+", Label: t("stats.youtubeVideos"), Href: "https://www.youtube.com/@furkanbytekin"},
		{Number: "16", Label: t("stats.languages"), Href: "/about"},
		{Number: "1", Label: t("stats.udemyCourse"), Href: "https://www.udemy.com/user/furkan-baytekin/"},
	}
	data["AboutHeading"] = t("index.aboutHeading")
	data["AboutText"] = t("index.aboutText")
	data["Hobbies"] = []Hobby{
		{Icon: "ti ti-player-record", Title: t("index.hobbyVinylTitle"), Description: t("index.hobbyVinylDesc")},
		{Icon: "ti ti-guitar-pick", Title: t("index.hobbyBassTitle"), Description: t("index.hobbyBassDesc")},
		{Icon: "ti ti-book", Title: t("index.hobbyBooksTitle"), Description: t("index.hobbyBooksDesc")},
		{Icon: "ti ti-device-gamepad-2", Title: t("index.hobbyGameJamsTitle"), Description: t("index.hobbyGameJamsDesc")},
	}
	data["JSONLD"] = mustMarshalJSON(struct {
		Context     string   `json:"@context"`
		Type        string   `json:"@type"`
		Name        string   `json:"name"`
		URL         string   `json:"url"`
		Email       string   `json:"email"`
		Description string   `json:"description"`
		JobTitle    string   `json:"jobTitle"`
		SameAs      []string `json:"sameAs"`
	}{
		Context:     "https://schema.org",
		Type:        "Person",
		Name:        SiteName,
		URL:         SiteBaseURL,
		Email:       SiteEmail,
		Description: t("hero.subtitle"),
		JobTitle:    "Lead Product Developer",
		SameAs: []string{
			"https://github.com/Elagoht",
			"https://linkedin.com/in/furkan-baytekin",
			"https://youtube.com/@furkanbytekin",
			"https://x.com/furkanbytekin",
		},
	})
	h.renderer.Render(w, "index.html", data)
}
