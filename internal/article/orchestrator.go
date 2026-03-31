package article

import (
	"context"
	"log/slog"
	"sort"

	"github.com/tbressel/daily-games-api/internal/cache"
	"github.com/tbressel/daily-games-api/internal/feed"
	"github.com/tbressel/daily-games-api/internal/metrics"
	"github.com/tbressel/daily-games-api/internal/rss"
	"github.com/tbressel/daily-games-api/pkg"
)

// Orchestrator coordinates feed selection, caching, RSS parsing,
// filtering, sorting, and pagination of articles.
type Orchestrator struct {
	parser *rss.Parser
	cache  *cache.Client
}

// New creates a new Orchestrator with the given RSS parser and Redis cache client.
func New(parser *rss.Parser, cache *cache.Client) *Orchestrator {
	return &Orchestrator{
		parser: parser,
		cache:  cache,
	}
}

// GetArticles is the main pipeline entry point.
// It resolves which feeds to fetch based on filters, checks the Redis cache,
// fetches from RSS if needed, then filters, sorts, and paginates the results.
func (o *Orchestrator) GetArticles(ctx context.Context, filters pkg.ArticleFilters) (pkg.ArticlesData, error) {
	// Step 1  -- resolve the feeds to fetch based on source/category filters
	feeds := o.resolveFeeds(filters)
	feeds = excludeFeeds(feeds, filters.ExcludeSources, filters.ExcludeCategories)
	if len(feeds) == 0 {
		return emptyData(filters), nil
	}

	// Step 2  -- check Redis cache (skip if refresh is forced)
	var articles []pkg.Article
	if !filters.Refresh {
		cached, err := o.cache.GetArticles(ctx, filters.Source, filters.Category, filters.Lang)
		if err != nil {
			slog.Warn("[Orchestrator] Cache read error", "err", err)
		}
		if cached != nil {
			slog.Info("[Orchestrator] Cache hit", "source", filters.Source, "category", filters.Category, "lang", filters.Lang)
			metrics.CacheHits.WithLabelValues(filters.Category, filters.Lang).Inc()
			articles = cached
		}
	}

	// Step 3  -- fetch from RSS feeds if cache missed or refresh forced
	if articles == nil {
		slog.Info("[Orchestrator] Cache miss  -- fetching RSS", "source", filters.Source, "category", filters.Category, "lang", filters.Lang)
		metrics.CacheMisses.WithLabelValues(filters.Category, filters.Lang).Inc()
		articles = o.parser.ParseFeeds(ctx, feeds)

		// Store fresh results in cache
		if err := o.cache.SetArticles(ctx, filters.Source, filters.Category, filters.Lang, articles); err != nil {
			slog.Warn("[Orchestrator] Cache write error", "err", err)
		}
	}

	// Step 4  -- sort by publication date descending (newest first)
	sortByDateDesc(articles)

	// Step 4b -- in "all feeds" mode (no source/category filter), cap to 5 articles per source
	// Applies whether or not a lang filter is active
	if filters.Source == "" && filters.Category == "" {
		articles = nPerSource(articles, 5)
	}

	// Step 5  -- paginate
	total := len(articles)
	paginated := paginate(articles, filters.Offset, filters.Limit)

	return pkg.ArticlesData{
		Articles: paginated,
		Metadata: pkg.ArticleMetadata{
			Offset:  filters.Offset,
			Limit:   filters.Limit,
			Total:   total,
			HasMore: filters.Offset+filters.Limit < total,
		},
	}, nil
}

// resolveFeeds returns the list of feeds to fetch based on the active filters.
// If a source ID is provided, only that feed is returned.
// If a category is provided, feeds are filtered by category (and lang if set).
// If only lang is provided, feeds are filtered by lang.
// Otherwise all active feeds are returned (filtered by lang if set).
func (o *Orchestrator) resolveFeeds(filters pkg.ArticleFilters) []pkg.Feed {
	if filters.Source != "" {
		f, ok := feed.GetByID(filters.Source)
		if !ok {
			return nil
		}
		return []pkg.Feed{f}
	}

	if filters.Category != "" && filters.Lang != "" {
		return feed.GetByCategoryAndLang(filters.Category, filters.Lang)
	}

	if filters.Category != "" {
		return feed.GetByCategory(filters.Category)
	}

	if filters.Lang != "" {
		return feed.GetByLang(filters.Lang)
	}

	return feed.GetActive()
}

// excludeFeeds removes feeds whose ID is in excludeSources or whose category is in excludeCategories.
func excludeFeeds(feeds []pkg.Feed, excludeSources, excludeCategories []string) []pkg.Feed {
	if len(excludeSources) == 0 && len(excludeCategories) == 0 {
		return feeds
	}
	srcSet := make(map[string]bool, len(excludeSources))
	for _, s := range excludeSources {
		srcSet[s] = true
	}
	catSet := make(map[string]bool, len(excludeCategories))
	for _, c := range excludeCategories {
		catSet[c] = true
	}
	result := feeds[:0:0]
	for _, f := range feeds {
		if !srcSet[f.ID] && !catSet[f.Category] {
			result = append(result, f)
		}
	}
	return result
}

// nPerSource returns at most n articles per unique source.
// Articles must already be sorted by date descending before calling this.
func nPerSource(articles []pkg.Article, n int) []pkg.Article {
	counts := make(map[string]int)
	result := make([]pkg.Article, 0, len(articles))
	for _, a := range articles {
		if counts[a.Source] < n {
			counts[a.Source]++
			result = append(result, a)
		}
	}
	return result
}

// sortByDateDesc sorts articles in place, newest publication date first.
func sortByDateDesc(articles []pkg.Article) {
	sort.Slice(articles, func(i, j int) bool {
		return articles[i].PubDate.After(articles[j].PubDate)
	})
}

// paginate returns a slice of articles starting at offset with a maximum of limit items.
// If offset exceeds the slice length, an empty slice is returned.
func paginate(articles []pkg.Article, offset, limit int) []pkg.Article {
	if offset >= len(articles) {
		return []pkg.Article{}
	}
	end := offset + limit
	if end > len(articles) {
		end = len(articles)
	}
	return articles[offset:end]
}

// emptyData returns an ArticlesData with zero articles and correct pagination metadata.
func emptyData(filters pkg.ArticleFilters) pkg.ArticlesData {
	return pkg.ArticlesData{
		Articles: []pkg.Article{},
		Metadata: pkg.ArticleMetadata{
			Offset:  filters.Offset,
			Limit:   filters.Limit,
			Total:   0,
			HasMore: false,
		},
	}
}
