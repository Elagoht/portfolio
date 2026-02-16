package handlers

import (
	"net/http"
	"strings"

	fwctx "statigo/framework/context"
	"statigo/framework/templates"
)

type BlogPostData struct {
	Slug      string
	Cover     string
	Title     string
	Category  string
	Date      string
	ReadTime  string
	Excerpt   string
	Content   string
	Tags      []string
	Canonical string
}

type BlogPostHandler struct {
	renderer *templates.Renderer
}

func NewBlogPostHandler(renderer *templates.Renderer) *BlogPostHandler {
	return &BlogPostHandler{
		renderer: renderer,
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

	// For now, return the same mock content for all slugs
	// TODO: Fetch actual blog post data based on slug
	blogPost := BlogPostData{
		Slug:      slug,
		Cover:     "https://furkanbaytekin.dev/blogs/software/using-makefile-in-go-projects-the-practical-way/opengraph-image",
		Title:     "Using Makefile in Go Projects (The Practical Way)",
		Category:  "Software",
		Date:      "December 21, 2025",
		ReadTime:  "5 min read",
		Excerpt:   "Learn how to use Makefile in Go projects with practical examples. Build, run, dev with live reload for your Go development workflow.",
		Tags:      []string{"go", "makefile", "devtools", "development"},
		Canonical: fwctx.GetCanonicalPath(r.Context()),
	}

	// Mock content - in production this would come from a database or CMS
	content := `<p>Go already gives you a solid toolchain, but once a project grows, repeating commands gets old fast. A simple <code>Makefile</code> gives you <strong>standardized commands</strong>, <strong>less typing</strong>, and <strong>zero dependencies</strong>.</p>

<p>This post shows a <strong>minimal, useful Makefile</strong> for Go projects:</p>

<ul>
<li><code>make build</code> &rarr; builds to <code>bin/projectname</code></li>
<li><code>make run</code> &rarr; runs the built binary</li>
<li><code>make dev</code> &rarr; live reload with <code>air</code></li>
<li><code>make clean</code> &rarr; removes build artifacts</li>
<li><code>make help</code> &rarr; self-documented commands</li>
</ul>

<hr>

<h2>Project Structure</h2>

<pre><code>.
├── cmd/
│   └── projectname/
│       └── main.go
├── bin/
├── air.toml
├── go.mod
├── Makefile</code></pre>

<h2>Why Use Makefile with Go?</h2>

<ul>
<li>One command everyone remembers</li>
<li>Same workflow across OSes</li>
<li>CI-friendly</li>
<li>Self-documenting</li>
<li>No bash scripts scattered around</li>
</ul>

<p>Go stays simple, Make just <strong>glues things together</strong>.</p>

<h2>The Makefile</h2>

<pre><code>APP_NAME := projectname
BIN_DIR := bin
BIN_PATH := $(BIN_DIR)/$(APP_NAME)

.PHONY: build run dev clean help

## build: Build the Go binary
build:
  @mkdir -p $(BIN_DIR)
  go build -o $(BIN_PATH) ./cmd/$(APP_NAME)

## run: Run the built binary
run: build
  @./$(BIN_PATH)

## dev: Run with live reload using air
dev:
  @air

## clean: Remove build artifacts
clean:
  @rm -rf $(BIN_DIR)

## help: Show available commands
help:
  @echo ""
  @echo "Available commands:"
  @grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
    | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-10s %s\n", $$1, $$2}'</code></pre>

<h2><code>make build</code></h2>

<pre><code>make build</code></pre>

<ul>
<li>Compiles your app</li>
<li>Output goes to <code>bin/projectname</code></li>
<li>No surprises, no magic paths</li>
</ul>

<p>This keeps your repo clean and binaries out of the way.</p>

<h2><code>make run</code></h2>

<pre><code>make run</code></pre>

<ul>
<li>Ensures the binary is built</li>
<li>Executes <code>./bin/projectname</code></li>
</ul>

<p>Simple dependency chain:</p>

<pre><code>run: build</code></pre>

<h2><code>make dev</code> (Live Reload)</h2>

<p>Uses <strong>Air</strong> for hot reload.</p>

<pre><code>make dev</code></pre>

<p>Make sure Air is installed:</p>

<pre><code>go install github.com/air-verse/air@latest</code></pre>

<p>And you have an <code>air.toml</code>:</p>

<pre><code>root = "."
tmp_dir = "tmp"

[build]
cmd = "go build -o bin/projectname ./cmd/projectname"
bin = "bin/projectname"</code></pre>

<h2><code>make clean</code></h2>

<pre><code>make clean</code></pre>

<p>Deletes: <code>bin/</code></p>

<p>Useful for:</p>

<ul>
<li>fresh builds</li>
<li>CI</li>
<li>sanity resets</li>
</ul>

<h2><code>make help</code></h2>

<pre><code>make help</code></pre>

<p>Outputs something like:</p>

<pre><code>Available commands:
  build      Build the Go binary
  run        Run the built binary
  dev        Run with live reload using air
  clean      Remove build artifacts
  help       Show available commands</code></pre>

<p>This works because of the <code>##</code> comments &mdash; small trick, big win.</p>

<hr>

<h2>Final Notes</h2>

<ul>
<li>Keep Makefiles <strong>boring</strong></li>
<li>Don't over-abstract</li>
<li>If a command needs explaining, it's probably wrong</li>
</ul>

<p>For Go projects, this setup hits the sweet spot: <strong>simple, readable, and production-safe</strong>.</p>

<p>Done.</p>`

	blogPost.Content = content

	data := BaseData(lang, t)
	data["Canonical"] = blogPost.Canonical
	data["Title"] = blogPost.Title + " | Furkan Baytekin"
	data["Meta"] = map[string]string{
		"description": blogPost.Excerpt,
	}
	data["BlogPost"] = blogPost

	// Mock navigation posts
	data["PrevPost"] = nil
	data["NextPost"] = map[string]string{
		"Slug":  "/blogs/why-backend-flows-must-be-restartable",
		"Title": "Why Backend Flows Must Be Restartable",
	}

	// Mock related posts
	data["RelatedPosts"] = []map[string]string{
		{
			"Slug":     "/blogs/bff-is-your-bff",
			"Cover":    "https://furkanbaytekin.dev/_next/image?url=http%3A%2Fdebian%3A2998%2Fuploads%2Fcovers%2Fbff-is-your-bff-why-backend-for-frontend-is-your-best-friend-forever%2B1750836318658.webp&w=1080&q=75",
			"Title":    "BFF is your BFF: Why Backend for Frontend is Your Best Friend Forever",
			"Category": "Software",
			"Date":     "June 25, 2025",
		},
		{
			"Slug":     "/blogs/understanding-eventual-consistency",
			"Cover":    "https://furkanbaytekin.dev/_next/image?url=http%3A%2Fdebian%3A2998%2Fuploads%2Fcovers%2Funderstanding-eventual-consistency%2B1750404639398.webp&w=1080&q=75",
			"Title":    "Understanding Eventual Consistency",
			"Category": "Software",
			"Date":     "June 20, 2025",
		},
		{
			"Slug":     "/blogs/idempotency-in-api-design",
			"Cover":    "https://furkanbaytekin.dev/_next/image?url=http%3A%2Fdebian%3A2998%2Fuploads%2Fcovers%2Fidempotency-in-api-design-why-it-matters-and-how-to-implement-it%2B1749384753379.webp&w=1080&q=75",
			"Title":    "Idempotency in API Design: Why it Matters and How to Implement It",
			"Category": "Software",
			"Date":     "June 8, 2025",
		},
	}

	h.renderer.Render(w, "blog-post.html", data)
}
