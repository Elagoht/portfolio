package handlers

import (
	"bytes"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/yuin/goldmark"

	fwctx "statigo/framework/context"
	"statigo/framework/templates"
	"statigo/internal/services"
)

type BlogPostData struct {
	Slug      string
	Cover     string
	Title     string
	Category  string
	Date      string
	ReadTime  string
	Excerpt   string
	Content   template.HTML
	Tags      []string
	Canonical string
}

type BlogPostHandler struct {
	renderer *templates.Renderer
	bloggo   *services.BloggoService
	apiBase  string
}

func NewBlogPostHandler(renderer *templates.Renderer, bloggo *services.BloggoService, apiBase string) *BlogPostHandler {
	return &BlogPostHandler{
		renderer: renderer,
		bloggo:   bloggo,
		apiBase:  apiBase,
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

	// Track view in background
	go func() {
		ua := r.Header.Get("User-Agent")
		_ = h.bloggo.TrackView(r.Context(), slug, ua)
	}()

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

	blogPost := BlogPostData{
		Slug:      slug,
		Cover:     cover,
		Title:     post.Title,
		Category:  post.Category.Name,
		Date:      post.PublishedAt.Format("January 2, 2006"),
		ReadTime:  strconv.Itoa(post.ReadTime),
		Excerpt:   excerpt,
		Content:   markdownToHTML(post.Content),
		Tags:      tags,
		Canonical: fwctx.GetCanonicalPath(r.Context()),
	}

	data := BaseData(lang, t)
	data["Canonical"] = blogPost.Canonical
	data["Title"] = blogPost.Title + " | Furkan Baytekin"
	data["Meta"] = map[string]string{
		"description": blogPost.Excerpt,
	}
	data["BlogPost"] = blogPost

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

func markdownToHTML(md string) template.HTML {
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(md), &buf); err != nil {
		return template.HTML(md)
	}
	return template.HTML(buf.String())
}
