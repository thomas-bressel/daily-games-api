package rss

import (
	"context"
	"crypto/md5"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/tbressel/daily-games-api/pkg"
)

// Parser handles fetching and parsing of RSS/Atom feeds.
type Parser struct {
	httpClient *http.Client
	maxItems   int
}

// New creates a new RSS Parser with a configured HTTP client and item limit.
//
// timeoutSeconds is the maximum duration for a single RSS feed HTTP request.
// maxItems is the maximum number of articles extracted per feed.
// userAgent mimics a real browser to bypass Cloudflare bot protection on some feeds.
const userAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"

// New creates a new RSS Parser with a configured HTTP client and item limit.
//
// timeoutSeconds is the maximum duration for a single RSS feed HTTP request.
// maxItems is the maximum number of articles extracted per feed.
func New(timeoutSeconds, maxItems int) *Parser {
	return &Parser{
		httpClient: &http.Client{
			Timeout: time.Duration(timeoutSeconds) * time.Second,
			Transport: &userAgentTransport{
				wrapped:   http.DefaultTransport,
				userAgent: userAgent,
			},
		},
		maxItems: maxItems,
	}
}

// userAgentTransport is an http.RoundTripper that injects a custom User-Agent header
// on every outgoing request, overriding the default Go HTTP client agent.
type userAgentTransport struct {
	wrapped   http.RoundTripper
	userAgent string
}

// RoundTrip executes the HTTP request after injecting browser-like headers
// to reduce the chance of being blocked by Cloudflare or similar protections.
func (t *userAgentTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
	req.Header.Set("User-Agent", t.userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "fr-FR,fr;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Referer", "https://www.google.com/")
	req.Header.Set("Cache-Control", "no-cache")
	return t.wrapped.RoundTrip(req)
}

// ParseFeed fetches and parses a single RSS feed, returning a slice of Articles.
// Items are limited to p.maxItems. Network or parse errors are returned directly.
func (p *Parser) ParseFeed(ctx context.Context, feed pkg.Feed) ([]pkg.Article, error) {
	fp := gofeed.NewParser()
	fp.Client = p.httpClient

	parsed, err := fp.ParseURLWithContext(feed.URL, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to parse feed %s: %w", feed.ID, err)
	}

	limit := min(p.maxItems, len(parsed.Items))

	articles := make([]pkg.Article, 0, limit)
	for _, item := range parsed.Items[:limit] {
		article := p.transformItem(ctx, item, feed)
		articles = append(articles, article)
	}

	return articles, nil
}

// ParseFeeds concurrently fetches and parses multiple RSS feeds.
// Results from all feeds are merged into a single article slice.
// Individual feed errors are logged but do not stop the other feeds from being fetched.
func (p *Parser) ParseFeeds(ctx context.Context, feeds []pkg.Feed) []pkg.Article {
	type result struct {
		articles []pkg.Article
		err      error
		feedID   string
	}

	ch := make(chan result, len(feeds))

	for _, feed := range feeds {
		go func(f pkg.Feed) {
			articles, err := p.ParseFeed(ctx, f)
			ch <- result{articles: articles, err: err, feedID: f.ID}
		}(feed)
	}

	var all []pkg.Article
	for range feeds {
		r := <-ch
		if r.err != nil {
			slog.Warn("[RSS] Feed fetch failed", "feedID", r.feedID, "err", r.err)
			continue
		}
		all = append(all, r.articles...)
	}

	return all
}

// transformItem converts a raw gofeed.Item into a pkg.Article.
// It extracts the best available image, cleans the description, and generates tags.
func (p *Parser) transformItem(ctx context.Context, item *gofeed.Item, feed pkg.Feed) pkg.Article {
	pubDate := time.Now()
	if item.PublishedParsed != nil {
		pubDate = *item.PublishedParsed
	} else if item.UpdatedParsed != nil {
		pubDate = *item.UpdatedParsed
	}

	creator := extractCreator(item)
	rawDescription := extractDescription(item)
	imageURL := extractImage(item)

	if enricher, ok := enrichers[feed.ID]; ok {
		enrichedImage, enrichedDesc := enricher(item)
		if enrichedImage != "" {
			imageURL = enrichedImage
		}
		if enrichedDesc != "" {
			rawDescription = enrichedDesc
		}
	}

	description := pkg.CleanDescription(rawDescription, 200)
	tags := pkg.ExtractTags(item.Title, description)

	return pkg.Article{
		ID:          generateID(feed.ID, item.Link),
		Title:       strings.TrimSpace(item.Title),
		Link:        item.Link,
		PubDate:     pubDate,
		Creator:     creator,
		Description: description,
		ImageURL:    imageURL,
		Source:      feed.ID,
		SourceName:  feed.Name,
		Category:    feed.Category,
		Tags:        tags,
	}
}

// extractCreator returns the best available author name from a feed item.
// It checks the Authors slice, then the Dublin Core creator extension.
func extractCreator(item *gofeed.Item) string {
	if len(item.Authors) > 0 && item.Authors[0].Name != "" {
		return item.Authors[0].Name
	}

	// Dublin Core fallback
	if dc := item.Extensions["dc"]; dc != nil {
		if creators, ok := dc["creator"]; ok && len(creators) > 0 {
			return creators[0].Value
		}
	}

	return ""
}

// extractDescription returns the best available description text from a feed item.
// Priority: Description > ContentSnippet > Content (first 500 chars).
func extractDescription(item *gofeed.Item) string {
	if item.Description != "" {
		return item.Description
	}
	return item.Content
}

// extractImage returns the best available image URL from a feed item.
// Priority: media:thumbnail (YouTube) > enclosure image > feed image.
func extractImage(item *gofeed.Item) string {
	// YouTube / media namespace thumbnail
	if media := item.Extensions["media"]; media != nil {
		if thumbnails, ok := media["thumbnail"]; ok && len(thumbnails) > 0 {
			if url := thumbnails[0].Attrs["url"]; url != "" {
				return url
			}
		}
		// media:group > media:thumbnail
		if groups, ok := media["group"]; ok && len(groups) > 0 {
			for _, ext := range groups[0].Children["thumbnail"] {
				if url := ext.Attrs["url"]; url != "" {
					return url
				}
			}
		}
	}

	// RSS enclosure (podcast / image attachment)
	if item.Image != nil && item.Image.URL != "" {
		return item.Image.URL
	}

	for _, enc := range item.Enclosures {
		if strings.HasPrefix(enc.Type, "image/") && enc.URL != "" {
			return enc.URL
		}
	}

	return ""
}

// generateID creates a deterministic MD5-based ID from the feed ID and article URL.
// This ensures the same article always gets the same ID across fetches.
func generateID(feedID, link string) string {
	h := md5.New()
	h.Write([]byte(feedID + link))
	return fmt.Sprintf("%x", h.Sum(nil))
}
