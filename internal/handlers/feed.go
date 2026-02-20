package handlers

import (
	"encoding/xml"
	"net/http"
	"time"

	"statigo/internal/services"
)

// RSS 2.0 XML structs

type rssFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Channel rssChannel `xml:"channel"`
}

type rssChannel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Items       []rssItem `xml:"item"`
}

type rssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	GUID        string `xml:"guid"`
}

type FeedHandler struct {
	bloggo  *services.BloggoService
	apiBase string
	siteURL string
}

func NewFeedHandler(bloggo *services.BloggoService, apiBase string, siteURL string) *FeedHandler {
	return &FeedHandler{
		bloggo:  bloggo,
		apiBase: apiBase,
		siteURL: siteURL,
	}
}

func (h *FeedHandler) RSS(w http.ResponseWriter, r *http.Request) {
	posts, err := h.fetchPosts(r)
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}

	items := make([]rssItem, 0, len(posts))
	for _, p := range posts {
		desc := ""
		if p.Description != nil {
			desc = *p.Description
		} else if p.Spot != nil {
			desc = *p.Spot
		}
		items = append(items, rssItem{
			Title:       p.Title,
			Link:        h.siteURL + "/blogs/" + p.Slug,
			Description: desc,
			PubDate:     p.PublishedAt.Format(time.RFC1123Z),
			GUID:        h.siteURL + "/blogs/" + p.Slug,
		})
	}

	feed := rssFeed{
		Version: "2.0",
		Channel: rssChannel{
			Title:       SiteName,
			Link:        h.siteURL,
			Description: SiteName + " Blog",
			Items:       items,
		},
	}

	w.Header().Set("Content-Type", "application/rss+xml; charset=utf-8")
	w.Write([]byte(xml.Header))
	xml.NewEncoder(w).Encode(feed)
}

func (h *FeedHandler) fetchPosts(r *http.Request) ([]services.PostSummary, error) {
	var all []services.PostSummary

	for page := 1; ; page++ {
		resp, err := h.bloggo.ListPosts(r.Context(), services.ListPostsParams{
			Page:  page,
			Limit: 100,
		})
		if err != nil || len(resp.Data) == 0 {
			break
		}
		all = append(all, resp.Data...)
		if len(resp.Data) < 100 {
			break
		}
	}

	return all, nil
}
