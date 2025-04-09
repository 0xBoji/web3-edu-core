package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

// Cache represents the Redis cache service
type Cache struct {
	client *redis.Client
}

// NewCache creates a new Redis cache service
func NewCache() *Cache {
	return &Cache{
		client: Client,
	}
}

// SetJSON sets a JSON value in the cache
func (c *Cache) SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, expiration).Err()
}

// GetJSON gets a JSON value from the cache
func (c *Cache) GetJSON(ctx context.Context, key string, dest interface{}) error {
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

// Delete deletes a key from the cache
func (c *Cache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

// SetUserSession sets a user session in the cache
func (c *Cache) SetUserSession(ctx context.Context, sessionID string, userData interface{}, expiration time.Duration) error {
	key := "user:session:" + sessionID
	return c.SetJSON(ctx, key, userData, expiration)
}

// GetUserSession gets a user session from the cache
func (c *Cache) GetUserSession(ctx context.Context, sessionID string, dest interface{}) error {
	key := "user:session:" + sessionID
	return c.GetJSON(ctx, key, dest)
}

// DeleteUserSession deletes a user session from the cache
func (c *Cache) DeleteUserSession(ctx context.Context, sessionID string) error {
	key := "user:session:" + sessionID
	return c.Delete(ctx, key)
}

// SetCourse sets a course in the cache
func (c *Cache) SetCourse(ctx context.Context, courseID string, courseData interface{}, expiration time.Duration) error {
	key := "course:" + courseID
	return c.SetJSON(ctx, key, courseData, expiration)
}

// GetCourse gets a course from the cache
func (c *Cache) GetCourse(ctx context.Context, courseID string, dest interface{}) error {
	key := "course:" + courseID
	return c.GetJSON(ctx, key, dest)
}

// DeleteCourse deletes a course from the cache
func (c *Cache) DeleteCourse(ctx context.Context, courseID string) error {
	key := "course:" + courseID
	return c.Delete(ctx, key)
}

// AddToSortedSet adds an item to a sorted set
func (c *Cache) AddToSortedSet(ctx context.Context, key string, score float64, member string) error {
	return c.client.ZAdd(ctx, key, redis.Z{
		Score:  score,
		Member: member,
	}).Err()
}

// GetSortedSet gets items from a sorted set
func (c *Cache) GetSortedSet(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return c.client.ZRevRange(ctx, key, start, stop).Result()
}

// SetCategories sets the categories list in the cache
func (c *Cache) SetCategories(ctx context.Context, categories interface{}, expiration time.Duration) error {
	return c.SetJSON(ctx, "categories:all", categories, expiration)
}

// GetCategories gets the categories list from the cache
func (c *Cache) GetCategories(ctx context.Context, dest interface{}) error {
	return c.GetJSON(ctx, "categories:all", dest)
}

// SetUserProgress sets a user's progress for a course in the cache
func (c *Cache) SetUserProgress(ctx context.Context, userID, courseID string, progressData interface{}, expiration time.Duration) error {
	key := "user:" + userID + ":progress:" + courseID
	return c.SetJSON(ctx, key, progressData, expiration)
}

// GetUserProgress gets a user's progress for a course from the cache
func (c *Cache) GetUserProgress(ctx context.Context, userID, courseID string, dest interface{}) error {
	key := "user:" + userID + ":progress:" + courseID
	return c.GetJSON(ctx, key, dest)
}

// DeleteUserProgress deletes a user's progress for a course from the cache
func (c *Cache) DeleteUserProgress(ctx context.Context, userID, courseID string) error {
	key := "user:" + userID + ":progress:" + courseID
	return c.Delete(ctx, key)
}

// IncrementRateLimit increments the rate limit counter for an IP
func (c *Cache) IncrementRateLimit(ctx context.Context, ip string, expiration time.Duration) (int64, error) {
	key := "rate:limit:" + ip
	pipe := c.client.Pipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, expiration)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}
	return incr.Val(), nil
}

// GetRateLimit gets the rate limit counter for an IP
func (c *Cache) GetRateLimit(ctx context.Context, ip string) (int64, error) {
	key := "rate:limit:" + ip
	val, err := c.client.Get(ctx, key).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	return val, err
}
