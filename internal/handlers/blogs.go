package handlers

import (
	"math"
	"net/http"
	"strconv"

	fwctx "statigo/framework/context"
	"statigo/framework/templates"
	"statigo/framework/utils"
	"statigo/internal/services"
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
	Number     int
	Href       string
	IsCurrent  bool
	IsEllipsis bool
}

type BlogsHandler struct {
	renderer     *templates.Renderer
	bloggo       *services.BloggoService
	apiBase      string
	postsPerPage int
}

func NewBlogsHandler(renderer *templates.Renderer, bloggo *services.BloggoService, apiBase string) *BlogsHandler {
	return &BlogsHandler{
		renderer:     renderer,
		bloggo:       bloggo,
		apiBase:      apiBase,
		postsPerPage: utils.GetEnvInt("BLOGS_PAGE_SIZE", 12),
	}
}

func (h *BlogsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	const lang = "en"
	canonical := fwctx.GetCanonicalPath(r.Context())
	t := func(key string) string {
		return h.renderer.GetTranslation(lang, key)
	}

	data := BaseData(lang, t)
	data["Canonical"] = canonical
	data["Title"] = t("pages.blogs.title")
	data["Meta"] = map[string]string{
		"description": t("sections.blogsDescription"),
	}

	// Parse query parameters
	query := r.URL.Query()
	currentPage := 1
	if p := query.Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			currentPage = parsed
		}
	}
	category := query.Get("category")
	tag := query.Get("tag")
	search := query.Get("search")

	// Pass current filter values to template
	data["CurrentSearch"] = search
	data["CurrentCategory"] = category
	data["CurrentTag"] = tag
	data["HasActiveFilters"] = category != "" || tag != "" || search != ""

	// Helper to build filter URLs preserving other params
	buildFilterURL := func(cat, tg, srch string) string {
		url := "/blogs"
		sep := "?"
		if cat != "" {
			url += sep + "category=" + cat
			sep = "&"
		}
		if tg != "" {
			url += sep + "tag=" + tg
			sep = "&"
		}
		if srch != "" {
			url += sep + "search=" + srch
		}
		return url
	}

	// Clear-all URL (remove all filters)
	data["ClearFiltersHref"] = "/blogs"

	// Fetch categories — clicking active deselects
	categories, err := h.bloggo.ListCategories(r.Context())
	if err == nil {
		var blogCategories []BlogCategory
		for _, cat := range categories {
			active := category == cat.Slug
			href := buildFilterURL(cat.Slug, tag, search)
			if active {
				href = buildFilterURL("", tag, search)
			}
			blogCategories = append(blogCategories, BlogCategory{
				Slug:   cat.Slug,
				Name:   cat.Name,
				Count:  cat.PostCount,
				Href:   href,
				Active: active,
			})
		}
		data["BlogCategories"] = blogCategories
	}

	// Fetch tags — clicking active deselects
	tags, err := h.bloggo.ListTags(r.Context())
	if err == nil {
		var blogTags []BlogTag
		for _, tg := range tags {
			active := tag == tg.Slug
			href := buildFilterURL(category, tg.Slug, search)
			if active {
				href = buildFilterURL(category, "", search)
			}
			blogTags = append(blogTags, BlogTag{
				Slug:   tg.Slug,
				Name:   tg.Name,
				Count:  tg.PostCount,
				Href:   href,
				Active: active,
			})
		}
		data["BlogTags"] = blogTags
	}

	// Fetch posts
	postsResp, err := h.bloggo.ListPosts(r.Context(), services.ListPostsParams{
		Page:     currentPage,
		Limit:    h.postsPerPage,
		Category: category,
		Tag:      tag,
		Search:   search,
	})

	if err != nil {
		h.renderer.Render(w, "blogs.html", data)
		return
	}

	var blogs []BlogPost
	for _, p := range postsResp.Data {
		cover := ""
		if p.CoverImage != nil {
			cover = h.apiBase + *p.CoverImage
		}
		excerpt := ""
		if p.Description != nil {
			excerpt = *p.Description
		} else if p.Spot != nil {
			excerpt = *p.Spot
		}
		blogs = append(blogs, BlogPost{
			Slug:     "/blogs/" + p.Slug,
			Cover:    cover,
			Title:    p.Title,
			Category: p.Category.Name,
			Date:     p.PublishedAt.Format("January 2, 2006"),
			Excerpt:  excerpt,
		})
	}
	data["Blogs"] = blogs

	// Pagination
	totalPages := int(math.Ceil(float64(postsResp.Total) / float64(h.postsPerPage)))
	if totalPages < 1 {
		totalPages = 1
	}

	pagePrefix := "/blogs"
	// Preserve filter params in pagination links
	filterQuery := ""
	if category != "" {
		filterQuery += "&category=" + category
	}
	if tag != "" {
		filterQuery += "&tag=" + tag
	}
	if search != "" {
		filterQuery += "&search=" + search
	}

	var pageNumbers []PageNumber
	addPage := func(i int) {
		isCurrent := i == currentPage
		href := ""
		if !isCurrent {
			href = pagePrefix + "?page=" + strconv.Itoa(i) + filterQuery
		}
		pageNumbers = append(pageNumbers, PageNumber{
			Number:    i,
			Href:      href,
			IsCurrent: isCurrent,
		})
	}
	addEllipsis := func() {
		pageNumbers = append(pageNumbers, PageNumber{IsEllipsis: true})
	}

	const delta = 2
	if totalPages <= 2*delta+5 {
		for i := 1; i <= totalPages; i++ {
			addPage(i)
		}
	} else {
		left := currentPage - delta
		right := currentPage + delta

		addPage(1)
		if left > 2 {
			addEllipsis()
		}
		for i := max(2, left); i <= min(totalPages-1, right); i++ {
			addPage(i)
		}
		if right < totalPages-1 {
			addEllipsis()
		}
		addPage(totalPages)
	}

	data["HasPrev"] = currentPage > 1
	data["HasNext"] = currentPage < totalPages
	if currentPage > 1 {
		data["PrevPage"] = pagePrefix + "?page=" + strconv.Itoa(currentPage-1) + filterQuery
	} else {
		data["PrevPage"] = ""
	}
	if currentPage < totalPages {
		data["NextPage"] = pagePrefix + "?page=" + strconv.Itoa(currentPage+1) + filterQuery
	} else {
		data["NextPage"] = ""
	}
	data["PageNumbers"] = pageNumbers
	data["JSONLD"] = mustMarshalJSON(struct {
		Context     string `json:"@context"`
		Type        string `json:"@type"`
		Name        string `json:"name"`
		Description string `json:"description"`
		URL         string `json:"url"`
		Author      struct {
			Type string `json:"@type"`
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"author"`
	}{
		Context:     "https://schema.org",
		Type:        "Blog",
		Name:        t("pages.blogs.title"),
		Description: t("sections.blogsDescription"),
		URL:         SiteBaseURL + "/blogs",
		Author: struct {
			Type string `json:"@type"`
			Name string `json:"name"`
			URL  string `json:"url"`
		}{
			Type: "Person",
			Name: SiteName,
			URL:  SiteBaseURL,
		},
	})

	h.renderer.Render(w, "blogs.html", data)
}
