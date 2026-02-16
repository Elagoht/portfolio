package handlers

import (
	"encoding/json"
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

// JSON Feed 1.1 structs

type jsonFeed struct {
	Version     string         `json:"version"`
	Title       string         `json:"title"`
	HomePageURL string         `json:"home_page_url"`
	FeedURL     string         `json:"feed_url"`
	Items       []jsonFeedItem `json:"items"`
}

type jsonFeedItem struct {
	ID            string `json:"id"`
	URL           string `json:"url"`
	Title         string `json:"title"`
	DatePublished string `json:"date_published"`
	Summary       string `json:"summary,omitempty"`
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

func (h *FeedHandler) JSON(w http.ResponseWriter, r *http.Request) {
	posts, err := h.fetchPosts(r)
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}

	items := make([]jsonFeedItem, 0, len(posts))
	for _, p := range posts {
		summary := ""
		if p.Description != nil {
			summary = *p.Description
		} else if p.Spot != nil {
			summary = *p.Spot
		}
		items = append(items, jsonFeedItem{
			ID:            h.siteURL + "/blogs/" + p.Slug,
			URL:           h.siteURL + "/blogs/" + p.Slug,
			Title:         p.Title,
			DatePublished: p.PublishedAt.Format(time.RFC3339),
			Summary:       summary,
		})
	}

	feed := jsonFeed{
		Version:     "https://jsonfeed.org/version/1.1",
		Title:       SiteName,
		HomePageURL: h.siteURL,
		FeedURL:     h.siteURL + "/feed.json",
		Items:       items,
	}

	w.Header().Set("Content-Type", "application/feed+json; charset=utf-8")
	json.NewEncoder(w).Encode(feed)
}

func (h *FeedHandler) fetchPosts(r *http.Request) ([]services.PostSummary, error) {
	resp, err := h.bloggo.ListPosts(r.Context(), services.ListPostsParams{
		Page:  1,
		Limit: 20,
	})
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}
