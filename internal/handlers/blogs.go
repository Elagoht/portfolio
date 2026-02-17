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
	for i := 1; i <= totalPages; i++ {
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

	h.renderer.Render(w, "blogs.html", data)
}
