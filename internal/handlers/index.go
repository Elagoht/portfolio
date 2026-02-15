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

type Link struct {
	Title string
	Href  string
}

type TechGroup struct {
	Title        string
	Technologies []string
}

type Project struct {
	Title string
	Repo  string
	Stack []string
}

type AboutCard struct {
	Title       string
	Description string
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

func (h *IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lang := middleware.GetLanguage(r.Context())
	canonical := router.GetCanonicalPath(r.Context())
	t := func(key string) string {
		return h.renderer.GetTranslation(lang, key)
	}

	data := map[string]any{
		"Lang":      lang,
		"Canonical": canonical,
		"Title":     "Furkan Baytekin",
		"Meta": map[string]string{
			"description": t("hero.subtitle"),
		},
		"Name":  "Furkan Baytekin",
		"Email": "furkan@baytekin.dev",
		"Links": []Link{
			{Title: "GitHub", Href: "https://github.com/Elagoht"},
			{Title: "LinkedIn", Href: "https://linkedin.com/in/furkan-baytekin"},
			{Title: "YouTube", Href: "https://youtube.com/@furkanbytekin"},
			{Title: "X", Href: "https://x.com/furkanbytekin"},
			{Title: "Telegram", Href: "https://t.me/furkanbytekin"},
			{Title: "Reddit", Href: "https://reddit.com/u/furkanbytekin"},
			{Title: "Spotify", Href: "https://open.spotify.com/user/furkanbytekin"},
			{Title: "Udemy", Href: "https://www.udemy.com/user/furkan-baytekin/"},
			{Title: "Itch.io", Href: "https://elagoht.itch.io"},
		},
		"Stats": []Stat{
			{Number: "160+", Label: t("stats.blogPosts")},
			{Number: "90+", Label: t("stats.youtubeVideos")},
			{Number: "12", Label: t("stats.languages")},
			{Number: "1", Label: t("stats.udemyCourse")},
		},
		"Expertise": []string{
			t("expertise.backend"),
			t("expertise.frontend"),
			t("expertise.devops"),
			t("expertise.database"),
			t("expertise.api"),
			t("expertise.testing"),
		},
		"Languages": []string{
			"Go", "TypeScript", "JavaScript", "C#", "Python",
			"Bash", "SQL", "HTML", "CSS", "GDScript", "AWK",
		},
		"TechStack": []TechGroup{
			{
				Title:        t("stack.backend"),
				Technologies: []string{"Go", "Chi", ".NET", "Node.js", "Express"},
			},
			{
				Title:        t("stack.frontend"),
				Technologies: []string{"React", "Next.js", "Astro", "TailwindCSS"},
			},
			{
				Title:        t("stack.devops"),
				Technologies: []string{"Docker", "Nginx", "Linux", "GitHub Actions", "CI/CD"},
			},
			{
				Title:        t("stack.database"),
				Technologies: []string{"PostgreSQL", "SQLite", "Redis", "MongoDB"},
			},
			{
				Title:        t("stack.tools"),
				Technologies: []string{"Git", "Neovim", "Tmux", "Make", "Air"},
			},
			{
				Title:        t("stack.other"),
				Technologies: []string{"REST", "WebSocket", "gRPC", "OAuth2", "JWT"},
			},
		},
		"Projects": []Project{
			{
				Title: "StatiGo",
				Repo:  "https://github.com/Elagoht/StatiGo",
				Stack: []string{"Go", "Chi", "HTML Templates"},
			},
			{
				Title: "Passenger",
				Repo:  "https://github.com/Elagoht/Passenger",
				Stack: []string{"C#", ".NET", "CLI"},
			},
			{
				Title: "SelfMark",
				Repo:  "https://github.com/Elagoht/SelfMark",
				Stack: []string{"TypeScript", "React", "Vite"},
			},
			{
				Title: "Inventa",
				Repo:  "https://github.com/Elagoht/Inventa",
				Stack: []string{"Python", "Flask", "SQLite"},
			},
		},
		"BlogCategories": []BlogCategory{
			{Name: t("categories.software"), Count: 142, Href: "/blogs?category=software"},
			{Name: t("categories.music"), Count: 8, Href: "/blogs?category=music"},
			{Name: t("categories.techtales"), Count: 5, Href: "/blogs?category=techtales"},
			{Name: t("categories.myLife"), Count: 3, Href: "/blogs?category=my-life"},
			{Name: t("categories.uxui"), Count: 2, Href: "/blogs?category=ux-ui"},
		},
		"About": []AboutCard{
			{Title: t("about.bass.title"), Description: t("about.bass.desc")},
			{Title: t("about.vinyl.title"), Description: t("about.vinyl.desc")},
			{Title: t("about.books.title"), Description: t("about.books.desc")},
			{Title: t("about.education.title"), Description: t("about.education.desc")},
		},
	}

	h.renderer.Render(w, "index.html", data)
}
