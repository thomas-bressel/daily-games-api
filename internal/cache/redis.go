package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tbressel/daily-games-api/config"
	"github.com/tbressel/daily-games-api/pkg"
)

// Client wraps the Redis client and cache configuration.
type Client struct {
	rdb *redis.Client
	ttl time.Duration
}

// New creates and returns a new Redis cache Client.
// It pings Redis to verify the connection is alive.
func New(cfg *config.Config) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	return &Client{
		rdb: rdb,
		ttl: time.Duration(cfg.CacheTTLMinutes) * time.Minute,
	}, nil
}

// buildKey returns a namespaced Redis key for articles cache entries.
// source and category are optional — empty strings are included as-is.
func buildKey(source, category string) string {
	return fmt.Sprintf("daily-games:articles:%s:%s", source, category)
}

// GetArticles retrieves a cached article list from Redis.
// Returns nil, nil if the key does not exist (cache miss).
func (c *Client) GetArticles(ctx context.Context, source, category string) ([]pkg.Article, error) {
	key := buildKey(source, category)

	val, err := c.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		// Cache miss — not an error
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("redis GET failed: %w", err)
	}

	var articles []pkg.Article
	if err := json.Unmarshal([]byte(val), &articles); err != nil {
		return nil, fmt.Errorf("cache decode failed: %w", err)
	}

	return articles, nil
}

// SetArticles stores an article list in Redis with the configured TTL.
func (c *Client) SetArticles(ctx context.Context, source, category string, articles []pkg.Article) error {
	key := buildKey(source, category)

	data, err := json.Marshal(articles)
	if err != nil {
		return fmt.Errorf("cache encode failed: %w", err)
	}

	if err := c.rdb.Set(ctx, key, data, c.ttl).Err(); err != nil {
		return fmt.Errorf("redis SET failed: %w", err)
	}

	return nil
}

// DeleteArticles removes a specific cache entry from Redis.
// Used to force a refresh for a given source/category combination.
func (c *Client) DeleteArticles(ctx context.Context, source, category string) error {
	key := buildKey(source, category)
	return c.rdb.Del(ctx, key).Err()
}

// FlushAll removes all daily-games cache entries from Redis.
// Used when a full refresh is requested.
func (c *Client) FlushAll(ctx context.Context) error {
	pattern := "daily-games:articles:*"

	keys, err := c.rdb.Keys(ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("redis KEYS failed: %w", err)
	}

	if len(keys) == 0 {
		return nil
	}

	return c.rdb.Del(ctx, keys...).Err()
}

// IncrTrack atomically increments the counter for a given article event (share or bookmark).
// Key format: daily-games:track:<articleID>:<event>
// The counter has no TTL — it persists indefinitely.
func (c *Client) IncrTrack(ctx context.Context, articleID, event string) (int64, error) {
	key := fmt.Sprintf("daily-games:track:%s:%s", articleID, event)
	count, err := c.rdb.Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("redis INCR failed: %w", err)
	}
	return count, nil
}

// GetTrack returns the current counter value for a given article event.
// Returns 0 if the key does not exist.
func (c *Client) GetTrack(ctx context.Context, articleID, event string) (int64, error) {
	key := fmt.Sprintf("daily-games:track:%s:%s", articleID, event)
	count, err := c.rdb.Get(ctx, key).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("redis GET failed: %w", err)
	}
	return count, nil
}

// Close gracefully closes the Redis connection.
func (c *Client) Close() error {
	return c.rdb.Close()
}
