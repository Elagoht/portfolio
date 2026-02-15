package handlers

import (
	"net/http"

	"statigo/framework/middleware"
	"statigo/framework/router"
	"statigo/framework/templates"
)

type AboutHandler struct {
	renderer *templates.Renderer
}

func NewAboutHandler(renderer *templates.Renderer) *AboutHandler {
	return &AboutHandler{
		renderer: renderer,
	}
}

func (h *AboutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lang := middleware.GetLanguage(r.Context())
	canonical := router.GetCanonicalPath(r.Context())
	t := func(key string) string {
		return h.renderer.GetTranslation(lang, key)
	}

	data := BaseData(lang, t)
	data["Canonical"] = canonical
	data["Title"] = t("pages.about.title")
	data["Meta"] = map[string]string{
		"description": t("about.text"),
	}
	data["AboutText"] = t("about.text")
	data["Expertise"] = []string{
		t("expertise.systemDesign"),
		t("expertise.reactFrontend"),
		t("expertise.modularMonolith"),
		t("expertise.restApi"),
		t("expertise.authSecurity"),
		t("expertise.performance"),
		t("expertise.dbMigration"),
		t("expertise.queues"),
		t("expertise.cicd"),
		t("expertise.linux"),
	}
	data["Languages"] = []string{
		"Go", "TypeScript", "JavaScript", "C#", "Python",
		"Bash", "SQL", "GraphQL", "HTML", "CSS", "GDScript", "AWK",
		"MDX", "Roff", "Pug", "Lua",
	}
	data["TechStack"] = []TechGroup{
		{Title: t("stack.frontend"), Technologies: []string{"Next.js", "React.js", "Solid.js", "Zustand", "Redux", "Tailwind CSS", "SASS"}},
		{Title: t("stack.backend"), Technologies: []string{"Go (Chi, Concurrency, Middleware)", "Nest.js (REST APIs)", "Node.js (Express.js)", "Django", "Flask", "Strapi (Headless CMS)"}},
		{Title: t("stack.infra"), Technologies: []string{"Docker", "GitHub Actions", "GitLab CI", "Systemd", "Caddy", "Nginx", "Cloudflare", "Monitoring & Statistics", "Backups"}},
		{Title: t("stack.databases"), Technologies: []string{"SQLite", "PostgreSQL", "Redis", "MYSQL", "MSSQL", "MongoDB"}},
		{Title: t("stack.testing"), Technologies: []string{"Go Test", "Jest", "Selenium"}},
		{Title: t("stack.apiConcepts"), Technologies: []string{"HTTP", "OAuth2", "JWT (RSA, JWKS)", "Input Validation", "API Versioning", "REST API & CRUD", "RBAC & API Security", "Rate Limiting", "Swagger / Postman"}},
	}
	data["Experience"] = []ExperienceItem{
		{
			Title:   "Lead Product Developer",
			Company: "Abonesepeti",
			Date:    "2025 - Present",
			Description: []string{
				"Designed and built 7 Nest.js <b>microservices</b> with <b>RESTful API endpoints</b>, then consolidated into 1 <b>modular monolith</b> to reduce infrastructure overhead.",
				"Built <b>versioned RESTful APIs</b> with <b>JWT (RSA/JWKS)</b> authentication and <b>Swagger/OpenAPI</b> documentation across all platform services.",
				"Migrated from per-tenant Vercel deploys to a single Host-header-based <b>multi-tenant API server</b> with <b>~15MB</b> memory usage with Go.",
				"Re-architected the entire platform to run within a single <b>ARM-based Linux server with 2GB RAM</b>. Reduced hosting costs by 90%.",
				"Built custom <b>Go web framework</b> with built-in <b>routing, middleware, API layer</b>, due to CVEs. Reduced RAM usage from <b>300MB to 15MB, bandwidth by ~95%, response times by 60%</b>.",
				"Reduced <b>CI build times from ~3 minutes to ~3 seconds</b> using Go's native compilation and implemented near-instant service restarts via systemd.",
				"Designed automated <b>image resizing pipeline</b> converting uploads to WebP to <b>reduce bandwidth usage by ~75-80%</b> and improve load performance.",
				"Implemented <b>Blue-Green</b> and <b>Rolling Update/Rollback</b> deployment strategies with <b>automated updates</b> for zero-downtime releases.",
				"Managed <b>ARM-based Linux</b> production infrastructure, <b>reverse proxy</b> configuration (Caddy), and full <b>CI/CD automation</b> end-to-end.",
				"Developed a shared <b>React-based UI component library</b> used across multiple internal projects, <b>improving development speed (~30%)</b>, creating a consistent UX.",
			},
		},
		{
			Title:   "Full Stack Developer",
			Company: "Abonesepeti",
			Date:    "2024 - 2025",
			Description: []string{
				"Rewrote the website using <b>GraphQL API</b> and <b>Next.js</b> with server-side data fetching.",
				"Developed <b>automation tools</b> including <b>PDF and WebP generation</b>.",
				"Built lead-page infrastructure with <b>webhook API integrations</b> and real-time updates using Vercel.",
				"Adapted the web application into İş Bankası İşCep <b>mini-app via API integration</b> (+8,000 active users).",
			},
		},
		{
			Title:   "Frontend Developer",
			Company: "Abonesepeti",
			Date:    "2023 - 2024",
			Description: []string{
				"<b>Rebuilt web application</b> using <b>Next.js</b> and <b>Tailwind CSS</b>.",
				"Increased <b>uptime from ~90% to ~99.99%</b>.",
				"Reduced release time from <b>30 to 7 minutes</b> through <b>CI/CD improvements</b>.",
				"Improved organic traffic through <b>technical SEO optimization.</b>",
			},
		},
	}
	data["Education"] = []Education{
		{University: "Anadolu University", Programme: "Web Design and Development (A.D.)", Date: "2021 - 2023"},
		{University: "Gazi University", Programme: "Turkish Language Education (B.A.)", Date: "2019 - 2023"},
	}

	h.renderer.Render(w, "about.html", data)
}
