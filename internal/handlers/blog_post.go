package handlers

import (
	"bytes"
	"context"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"golang.org/x/net/html"

	fwctx "statigo/framework/context"
	"statigo/framework/templates"
	"statigo/internal/services"
)

type TOCItem struct {
	Text  string
	ID    string
	Level int
}

type BlogPostData struct {
	Slug      string
	Cover     string
	Title     string
	Category  string
	Date      string
	DateISO   string
	ReadTime  string
	Excerpt   string
	Content   template.HTML
	Tags      []string
	Canonical string
	TOCItems  []TOCItem
	ViewCount int
}

type BlogPostHandler struct {
	renderer     *templates.Renderer
	bloggo       *services.BloggoService
	apiBase      string
	viewTracker  *services.ViewTracker
	viewsHandler *ViewsHandler
}

func NewBlogPostHandler(renderer *templates.Renderer, bloggo *services.BloggoService, apiBase string, viewTracker *services.ViewTracker, viewsHandler *ViewsHandler) *BlogPostHandler {
	return &BlogPostHandler{
		renderer:     renderer,
		bloggo:       bloggo,
		apiBase:      apiBase,
		viewTracker:  viewTracker,
		viewsHandler: viewsHandler,
	}
}

func (h *BlogPostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	const lang = "en"
	t := func(key string) string {
		return h.renderer.GetTranslation(lang, key)
	}

	// Extract slug from path
	path := r.URL.Path
	slug := strings.TrimPrefix(path, "/blogs/")
	slug = strings.TrimSuffix(slug, "/")

	// Fetch post from API
	post, err := h.bloggo.GetPost(r.Context(), slug)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		data := BaseData(lang, t)
		data["Title"] = t("pages.notfound.title")
		data["Content"] = map[string]string{
			"heading": t("pages.notfound.heading"),
			"message": t("pages.notfound.message"),
			"action":  t("pages.notfound.action"),
		}
		h.renderer.Render(w, "notfound.html", data)
		return
	}

	cover := ""
	if post.CoverImage != nil {
		cover = h.apiBase + *post.CoverImage
	}
	excerpt := ""
	if post.Description != nil {
		excerpt = *post.Description
	} else if post.Spot != nil {
		excerpt = *post.Spot
	}

	var tags []string
	for _, tag := range post.Tags {
		tags = append(tags, tag.Name)
	}

	content := markdownToHTML(post.Content)
	tocItems := extractTOCItems(string(content))

	// Fetch view count from cache
	viewCount := 0
	if h.viewsHandler != nil {
		views, err := h.viewsHandler.Cache.Get()
		if err == nil {
			viewCount = views[slug]
		}
	}

	blogPost := BlogPostData{
		Slug:      slug,
		Cover:     cover,
		Title:     post.Title,
		Category:  post.Category.Name,
		Date:      post.PublishedAt.Format("January 2, 2006"),
		DateISO:   post.PublishedAt.Format("2006-01-02"),
		ReadTime:  strconv.Itoa(post.ReadTime),
		Excerpt:   excerpt,
		Content:   content,
		Tags:      tags,
		Canonical: fwctx.GetCanonicalPath(r.Context()),
		TOCItems:  tocItems,
		ViewCount: viewCount,
	}

	data := BaseData(lang, t)
	data["Canonical"] = blogPost.Canonical
	data["Title"] = blogPost.Title + " | Furkan Baytekin"
	data["Meta"] = map[string]string{
		"description": blogPost.Excerpt,
	}
	data["BlogPost"] = blogPost
	data["JSONLD"] = mustMarshalJSON(struct {
		Context       string `json:"@context"`
		Type          string `json:"@type"`
		Headline      string `json:"headline"`
		Description   string `json:"description,omitempty"`
		Image         string `json:"image,omitempty"`
		DatePublished string `json:"datePublished"`
		URL           string `json:"url"`
		Author        struct {
			Type string `json:"@type"`
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"author"`
		Publisher struct {
			Type string `json:"@type"`
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"publisher"`
	}{
		Context:       "https://schema.org",
		Type:          "BlogPosting",
		Headline:      blogPost.Title,
		Description:   blogPost.Excerpt,
		Image:         blogPost.Cover,
		DatePublished: blogPost.DateISO,
		URL:           SiteBaseURL + blogPost.Canonical,
		Author: struct {
			Type string `json:"@type"`
			Name string `json:"name"`
			URL  string `json:"url"`
		}{Type: "Person", Name: SiteName, URL: SiteBaseURL},
		Publisher: struct {
			Type string `json:"@type"`
			Name string `json:"name"`
			URL  string `json:"url"`
		}{Type: "Person", Name: SiteName, URL: SiteBaseURL},
	})

	// Fetch related posts (same category)
	related, err := h.bloggo.ListPosts(r.Context(), services.ListPostsParams{
		Category: post.Category.Slug,
		Limit:    4,
	})
	if err == nil {
		var relatedPosts []map[string]string
		for _, p := range related.Data {
			if p.Slug == slug {
				continue
			}
			if len(relatedPosts) >= 3 {
				break
			}
			relCover := ""
			if p.CoverImage != nil {
				relCover = h.apiBase + *p.CoverImage
			}
			relatedPosts = append(relatedPosts, map[string]string{
				"Slug":     "/blogs/" + p.Slug,
				"Cover":    relCover,
				"Title":    p.Title,
				"Category": p.Category.Name,
				"Date":     p.PublishedAt.Format("January 2, 2006"),
			})
		}
		data["RelatedPosts"] = relatedPosts
	}

	h.renderer.Render(w, "blog-post.html", data)
}

// ViewTrackingMiddleware tracks blog post views before the cache layer,
// so view counts are incremented even when serving cached responses.
func (h *BlogPostHandler) ViewTrackingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.HasPrefix(path, "/blogs/") {
			slug := strings.TrimSuffix(strings.TrimPrefix(path, "/blogs/"), "/")
			if slug != "" {
				ua := r.Header.Get("User-Agent")
				go func() {
					if h.viewTracker.ShouldTrackView(r, slug) {
						ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
						defer cancel()
						_ = h.bloggo.TrackView(ctx, slug, ua)
					}
				}()
			}
		}
		next.ServeHTTP(w, r)
	})
}

func markdownToHTML(md string) template.HTML {
	mdParser := goldmark.New(
		goldmark.WithExtensions(
			extension.Table,
			extension.Strikethrough,
			extension.Linkify,
			highlighting.NewHighlighting(
				highlighting.WithStyle("dracula"),
				highlighting.WithCSSWriter(htmlEscapeWriter{}),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)

	var buf bytes.Buffer
	if err := mdParser.Convert([]byte(md), &buf); err != nil {
		return template.HTML(md)
	}
	return template.HTML(buf.String())
}

func extractTOCItems(htmlContent string) []TOCItem {
	var items []TOCItem
	tokenizer := html.NewTokenizer(strings.NewReader(htmlContent))

	for {
		tt := tokenizer.Next()
		if tt == html.ErrorToken {
			break
		}
		if tt != html.StartTagToken {
			continue
		}

		tn, _ := tokenizer.TagName()
		tagName := string(tn)

		var level int
		switch tagName {
		case "h2":
			level = 2
		case "h3":
			level = 3
		case "h4":
			level = 4
		default:
			continue
		}

		// Extract id attribute
		var id string
		for {
			key, val, more := tokenizer.TagAttr()
			if string(key) == "id" {
				id = string(val)
			}
			if !more {
				break
			}
		}
		if id == "" {
			continue
		}

		// Collect all text content until closing tag
		var text strings.Builder
		depth := 1
		for depth > 0 {
			next := tokenizer.Next()
			switch next {
			case html.TextToken:
				text.Write(tokenizer.Text())
			case html.StartTagToken:
				depth++
			case html.EndTagToken:
				depth--
			case html.ErrorToken:
				depth = 0
			}
		}

		items = append(items, TOCItem{
			Text:  strings.TrimSpace(text.String()),
			ID:    id,
			Level: level,
		})
	}

	return items
}

// htmlEscapeWriter wraps a bytes.Buffer to escape HTML output for CSS
type htmlEscapeWriter struct{}

func (w htmlEscapeWriter) Write(p []byte) (int, error) {
	return 0, nil // We don't need CSS output since we'll use our own stylesheet
}
