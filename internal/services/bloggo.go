package services

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"time"

	"statigo/framework/client"
)

// BloggoTime handles the non-standard datetime format from the Bloggo API ("2025-10-12 13:15:38").
type BloggoTime struct {
	time.Time
}

func (bt *BloggoTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "null" || s == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		// Fallback to RFC3339
		t, err = time.Parse(time.RFC3339, s)
		if err != nil {
			return fmt.Errorf("cannot parse %q as BloggoTime: %w", s, err)
		}
	}
	bt.Time = t
	return nil
}

// BloggoService communicates with the Bloggo headless CMS API.
type BloggoService struct {
	client *client.Client
	logger *slog.Logger
}

// NewBloggoService creates a new Bloggo API service.
func NewBloggoService(c *client.Client, logger *slog.Logger) *BloggoService {
	return &BloggoService{
		client: c,
		logger: logger,
	}
}

// --- Response types ---

type PostsResponse struct {
	Data  []PostSummary `json:"data"`
	Page  int           `json:"page"`
	Take  int           `json:"take"`
	Total int           `json:"total"`
}

type PostSummary struct {
	Slug        string     `json:"slug"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	Spot        *string    `json:"spot"`
	CoverImage  *string    `json:"coverImage"`
	ReadCount   int        `json:"readCount"`
	ReadTime    int        `json:"readTime"`
	PublishedAt BloggoTime `json:"publishedAt"`
	Author      Author     `json:"author"`
	Category    Category   `json:"category"`
	Tags        []TagShort `json:"tags"`
}

type PostDetail struct {
	Slug        string     `json:"slug"`
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	Description *string    `json:"description"`
	Spot        *string    `json:"spot"`
	CoverImage  *string    `json:"coverImage"`
	ReadCount   int        `json:"readCount"`
	ReadTime    int        `json:"readTime"`
	PublishedAt BloggoTime `json:"publishedAt"`
	UpdatedAt   BloggoTime `json:"updatedAt"`
	Author      Author     `json:"author"`
	Category    Category   `json:"category"`
	Tags        []TagShort `json:"tags"`
}

type Author struct {
	ID     int     `json:"id"`
	Name   string  `json:"name"`
	Avatar *string `json:"avatar"`
}

type AuthorDetail struct {
	ID                 int        `json:"id"`
	Name               string     `json:"name"`
	Avatar             *string    `json:"avatar"`
	PublishedPostCount int        `json:"publishedPostCount"`
	MemberSince        BloggoTime `json:"memberSince"`
}

type AuthorsResponse struct {
	Authors []AuthorDetail `json:"authors"`
}

type Category struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

type CategoryDetail struct {
	Slug        string  `json:"slug"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Spot        *string `json:"spot"`
	PostCount   int     `json:"postCount"`
}

type CategoriesResponse struct {
	Categories []CategoryDetail `json:"categories"`
}

type TagShort struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

type TagDetail struct {
	Slug      string `json:"slug"`
	Name      string `json:"name"`
	PostCount int    `json:"postCount"`
}

type TagsResponse struct {
	Tags []TagDetail `json:"tags"`
}

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type TrackViewRequest struct {
	UserAgent string `json:"userAgent"`
}

// --- Query parameters ---

type ListPostsParams struct {
	Page     int
	Limit    int
	Category string
	Tag      string
	Author   string
	Search   string
}

// --- API methods ---

func (s *BloggoService) ListPosts(ctx context.Context, params ListPostsParams) (*PostsResponse, error) {
	query := url.Values{}

	if params.Page > 0 {
		query.Set("page", fmt.Sprintf("%d", params.Page))
	}
	if params.Limit > 0 {
		query.Set("limit", fmt.Sprintf("%d", params.Limit))
	}
	if params.Category != "" {
		query.Set("category", params.Category)
	}
	if params.Tag != "" {
		query.Set("tag", params.Tag)
	}
	if params.Author != "" {
		query.Set("author", params.Author)
	}
	if params.Search != "" {
		query.Set("search", params.Search)
	}

	path := "/api/posts?" + query.Encode()

	var resp PostsResponse
	if err := s.client.Get(ctx, path, &resp); err != nil {
		s.logger.Error("failed to list posts", "error", err)
		return nil, err
	}
	return &resp, nil
}

func (s *BloggoService) GetPost(ctx context.Context, slug string) (*PostDetail, error) {
	var resp PostDetail
	if err := s.client.Get(ctx, "/api/posts/"+slug, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *BloggoService) TrackView(ctx context.Context, slug string, userAgent string) error {
	body := TrackViewRequest{UserAgent: userAgent}
	return s.client.Post(ctx, "/api/posts/"+slug+"/view", body, nil)
}

func (s *BloggoService) GetViewCounts(ctx context.Context) (map[string]int, error) {
	var resp map[string]int
	if err := s.client.Get(ctx, "/api/posts/views", &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *BloggoService) ListCategories(ctx context.Context) ([]CategoryDetail, error) {
	var resp CategoriesResponse
	if err := s.client.Get(ctx, "/api/categories", &resp); err != nil {
		return nil, err
	}
	return resp.Categories, nil
}

func (s *BloggoService) GetCategory(ctx context.Context, slug string) (*CategoryDetail, error) {
	var resp CategoryDetail
	if err := s.client.Get(ctx, "/api/categories/"+slug, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *BloggoService) ListTags(ctx context.Context) ([]TagDetail, error) {
	var resp TagsResponse
	if err := s.client.Get(ctx, "/api/tags", &resp); err != nil {
		return nil, err
	}
	return resp.Tags, nil
}

func (s *BloggoService) GetTag(ctx context.Context, slug string) (*TagDetail, error) {
	var resp TagDetail
	if err := s.client.Get(ctx, "/api/tags/"+slug, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *BloggoService) ListAuthors(ctx context.Context) ([]AuthorDetail, error) {
	var resp AuthorsResponse
	if err := s.client.Get(ctx, "/api/authors", &resp); err != nil {
		return nil, err
	}
	return resp.Authors, nil
}

func (s *BloggoService) GetAuthor(ctx context.Context, id int) (*AuthorDetail, error) {
	var resp AuthorDetail
	if err := s.client.Get(ctx, fmt.Sprintf("/api/authors/%d", id), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *BloggoService) GetKeyValues(ctx context.Context, key, starting string) ([]KeyValue, error) {
	query := url.Values{}
	if key != "" {
		query.Set("key", key)
	}
	if starting != "" {
		query.Set("starting", starting)
	}

	path := "/api/key-values?" + query.Encode()

	var resp []KeyValue
	if err := s.client.Get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}
