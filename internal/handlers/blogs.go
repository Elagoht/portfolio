package handlers

import (
	"net/http"
	"strconv"

	"statigo/framework/middleware"
	"statigo/framework/router"
	"statigo/framework/templates"
)

type BlogPost struct {
	Slug     string
	Cover    string
	Title    string
	Category string
	Date     string
	Excerpt  string
}

type PageNumber struct {
	Number      int
	Href        string
	IsCurrent   bool
	IsEllipsis bool
}

type BlogsHandler struct {
	renderer *templates.Renderer
	registry *router.Registry
}

func NewBlogsHandler(renderer *templates.Renderer, registry *router.Registry) *BlogsHandler {
	return &BlogsHandler{
		renderer: renderer,
		registry: registry,
	}
}

func (h *BlogsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lang := middleware.GetLanguage(r.Context())
	canonical := router.GetCanonicalPath(r.Context())
	t := func(key string) string {
		return h.renderer.GetTranslation(lang, key)
	}

	data := BaseData(lang)
	data["Canonical"] = canonical
	data["Title"] = t("pages.blogs.title")
	data["Meta"] = map[string]string{
		"description": t("sections.blogsDescription"),
	}

	// TODO: Replace with real data source when available
	data["BlogCategories"] = []BlogCategory{
		{Name: t("categories.software"), Count: 142, Href: "/blogs?category=software"},
		{Name: t("categories.music"), Count: 8, Href: "/blogs?category=music"},
		{Name: t("categories.techtales"), Count: 5, Href: "/blogs?category=techtales"},
		{Name: t("categories.myLife"), Count: 3, Href: "/blogs?category=my-life"},
		{Name: t("categories.uxui"), Count: 2, Href: "/blogs?category=ux-ui"},
	}

	// TODO: Replace with real blog posts from your data source
	data["Blogs"] = []BlogPost{
		{
			Slug:     "/blogs/how-to-deploy-a-nextjs-app-with-cloudflare",
			Cover:    "https://furkanbaytekin.dev/_next/image?url=http%3A%2Fdebian%3A2998%2Fuploads%2Fcovers%2Fhow-to-deploy-a-nextjs-app-with-cloudflare-a-step-by-step-guide%2B1759231398074.webp&w=1080&q=75",
			Title:    "How to Deploy a Next.js App with Cloudflare: A Step-by-Step Guide",
			Category: t("categories.software"),
			Date:     "September 30, 2025",
			Excerpt:  "Deploy your Next.js app with Cloudflare for free SSL & global access.",
		},
		{
			Slug:     "/blogs/why-backend-flows-must-be-restartable",
			Cover:    "https://furkanbaytekin.dev/_next/image?url=http%3A%2Fdebian%3A2998%2Fuploads%2Fcovers%2Fwhy-backend-flows-must-be-restartable%2B1750922555473.webp&w=1080&q=75",
			Title:    "Why Backend Flows Must Be Restartable",
			Category: t("categories.software"),
			Date:     "June 26, 2025",
			Excerpt:  "Build resilient backend flows that users can restart or resume at any time",
		},
		{
			Slug:     "/blogs/bff-is-your-bff",
			Cover:    "https://furkanbaytekin.dev/_next/image?url=http%3A%2Fdebian%3A2998%2Fuploads%2Fcovers%2Fbff-is-your-bff-why-backend-for-frontend-is-your-best-friend-forever%2B1750836318658.webp&w=1080&q=75",
			Title:    "BFF is your BFF: Why Backend for Frontend is Your Best Friend Forever",
			Category: t("categories.software"),
			Date:     "June 25, 2025",
			Excerpt:  "Master the BFF pattern to build faster, more maintainable applications",
		},
		{
			Slug:     "/blogs/rarely-used-obscure-html-tags",
			Cover:    "https://furkanbaytekin.dev/_next/image?url=http%3A%2Fdebian%3A2998%2Fuploads%2Fcovers%2Frarely-used-obscure-html-tags%2B1750747535832.webp&w=1080&q=75",
			Title:    "Rarely Used Obscure HTML Tags",
			Category: t("categories.software"),
			Date:     "June 24, 2025",
			Excerpt:  "Discover 10 powerful HTML tags to enhance your web contents.",
		},
		{
			Slug:     "/blogs/exploring-lesser-used-semantic-html-elements",
			Cover:    "https://furkanbaytekin.dev/_next/image?url=http%3A%2Fdebian%3A2998%2Fuploads%2Fcovers%2Fexploring-lesser-used-by-jrs-but-powerful-semantic-html-elements%2B1750662679611.webp&w=1080&q=75",
			Title:    "Exploring Lesser-Used (By JRs) but Powerful Semantic HTML Elements",
			Category: t("categories.software"),
			Date:     "June 23, 2025",
			Excerpt:  "Master semantic HTML elements to build better, more accessible websites",
		},
		{
			Slug:     "/blogs/understanding-eventual-consistency",
			Cover:    "https://furkanbaytekin.dev/_next/image?url=http%3A%2Fdebian%3A2998%2Fuploads%2Fcovers%2Funderstanding-eventual-consistency%2B1750404639398.webp&w=1080&q=75",
			Title:    "Understanding Eventual Consistency",
			Category: t("categories.software"),
			Date:     "June 20, 2025",
			Excerpt:  "Understanding eventual consistency and its role in distributed systems",
		},
		{
			Slug:     "/blogs/replay-attacks",
			Cover:    "https://furkanbaytekin.dev/_next/image?url=http%3A%2Fdebian%3A2998%2Fuploads%2Fcovers%2Funderstanding-replay-attacks-what-they-are-and-how-to-protect-yourself%2B1750317284667.webp&w=1080&q=75",
			Title:    "Replay Attacks: What They Are and How to Protect Yourself",
			Category: t("categories.software"),
			Date:     "June 19, 2025",
			Excerpt:  "Defend against replay attacks via cybersecurity strategies, best practices",
		},
		{
			Slug:     "/blogs/accepting-files-on-the-backend",
			Cover:    "https://furkanbaytekin.dev/_next/image?url=http%3A%2Fdebian%3A2998%2Fuploads%2Fcovers%2Faccepting-files-on-the-backend-mime-extensions-path-traversal-and-quotas-explained%2B1750230423863.webp&w=1080&q=75",
			Title:    "Accepting Files on the Backend: MIME, Extensions, Path Traversal, and Quotas Explained",
			Category: t("categories.software"),
			Date:     "June 18, 2025",
			Excerpt:  "Secure file upload guide: Validation, quotas and security best practices",
		},
		{
			Slug:     "/blogs/stale-while-revalidate-swr",
			Cover:    "https://furkanbaytekin.dev/_next/image?url=http%3A%2Fdebian%3A2998%2Fuploads%2Fcovers%2Fstale-while-revalidate-swr-in-frontend-development%2B1750147336299.webp&w=1080&q=75",
			Title:    "Stale-While-Revalidate (SWR) in Frontend Development",
			Category: t("categories.software"),
			Date:     "June 17, 2025",
			Excerpt:  "Optimize frontend performance with Stale-While-Revalidate caching strategy",
		},
		{
			Slug:     "/blogs/why-standardized-hardware-wins",
			Cover:    "https://furkanbaytekin.dev/_next/image?url=http%3A%2Fdebian%3A2998%2Fuploads%2Fcovers%2Fwhy-standardized-hardware-wins-better-software-support-for-steam-deck-iphones-and-beyond%2B1750055170900.webp&w=1080&q=75",
			Title:    "Why Standardized Hardware Wins: Better Software Support for Steam Deck, iPhones, and Beyond",
			Category: t("categories.techtales"),
			Date:     "June 16, 2025",
			Excerpt:  "Standardized hardware: The key to better experiences",
		},
		{
			Slug:     "/blogs/idempotency-in-api-design",
			Cover:    "https://furkanbaytekin.dev/_next/image?url=http%3A%2Fdebian%3A2998%2Fuploads%2Fcovers%2Fidempotency-in-api-design-why-it-matters-and-how-to-implement-it%2B1749384753379.webp&w=1080&q=75",
			Title:    "Idempotency in API Design: Why it Matters and How to Implement It",
			Category: t("categories.software"),
			Date:     "June 8, 2025",
			Excerpt:  "Building reliable APIs with idempotency",
		},
		{
			Slug:     "/blogs/the-myth-of-code-speaks-for-itself",
			Cover:    "https://furkanbaytekin.dev/_next/image?url=http%3A%2Fdebian%3A2998%2Fuploads%2Fcovers%2Fthe-myth-of-code-speaks-for-itself%2B1749384690504.webp&w=1080&q=75",
			Title:    "The Myth of \"Code Speaks for Itself\"",
			Category: t("categories.software"),
			Date:     "June 8, 2025",
			Excerpt:  "Why code alone isn't enough: Tips for better software documentation",
		},
	}

	// TODO: Replace with real pagination logic from your data source
	currentPage := 1
	totalPages := 10 // Example: 142 posts / 12 per page ≈ 12 pages
	pagePrefix := "/" + lang + "/blogs"

	// Generate page numbers for pagination
	var pageNumbers []PageNumber
	for i := 1; i <= totalPages; i++ {
		isCurrent := i == currentPage
		href := ""
		if !isCurrent {
			href = pagePrefix + "?page=" + strconv.Itoa(i)
		}
		pageNumbers = append(pageNumbers, PageNumber{
			Number:    i,
			Href:      href,
			IsCurrent: isCurrent,
		})
	}

	data["HasPrev"] = currentPage > 1
	data["HasNext"] = currentPage < totalPages
	if currentPage > 1 {
		data["PrevPage"] = pagePrefix + "?page=" + strconv.Itoa(currentPage-1)
	} else {
		data["PrevPage"] = ""
	}
	if currentPage < totalPages {
		data["NextPage"] = pagePrefix + "?page=" + strconv.Itoa(currentPage+1)
	} else {
		data["NextPage"] = ""
	}
	data["PageNumbers"] = pageNumbers

	h.renderer.Render(w, "blogs.html", data)
}
