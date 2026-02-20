package handlers

import (
	"encoding/xml"
	"net/http"

	"statigo/internal/services"
)

type sitemapURLSet struct {
	XMLName xml.Name     `xml:"urlset"`
	Xmlns   string       `xml:"xmlns,attr"`
	URLs    []sitemapURL `xml:"url"`
}

type sitemapURL struct {
	Loc        string `xml:"loc"`
	ChangeFreq string `xml:"changefreq"`
	Priority   string `xml:"priority"`
}

type SitemapHandler struct {
	bloggo  *services.BloggoService
	siteURL string
}

func NewSitemapHandler(bloggo *services.BloggoService, siteURL string) *SitemapHandler {
	return &SitemapHandler{bloggo: bloggo, siteURL: siteURL}
}

func (h *SitemapHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	urls := []sitemapURL{
		{Loc: h.siteURL + "/", ChangeFreq: "weekly", Priority: "1.0"},
		{Loc: h.siteURL + "/about", ChangeFreq: "monthly", Priority: "0.8"},
		{Loc: h.siteURL + "/blogs", ChangeFreq: "daily", Priority: "0.9"},
	}

	page := 1
	for {
		resp, err := h.bloggo.ListPosts(r.Context(), services.ListPostsParams{
			Page:  page,
			Limit: 100,
		})
		if err != nil || len(resp.Data) == 0 {
			break
		}
		for _, p := range resp.Data {
			urls = append(urls, sitemapURL{
				Loc:        h.siteURL + "/blogs/" + p.Slug,
				ChangeFreq: "monthly",
				Priority:   "0.7",
			})
		}
		if len(resp.Data) < 100 {
			break
		}
		page++
	}

	urlset := sitemapURLSet{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  urls,
	}

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Write([]byte(xml.Header))
	xml.NewEncoder(w).Encode(urlset)
}
