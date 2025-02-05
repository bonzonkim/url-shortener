package storage

import (
	"context"
	"fmt"
	"hash/fnv"
	"net/url"
	"os"
	"strings"

	"github.com/go-redis/redis/v8"
)

type RedisStore struct {
	Client *redis.Client
}

var ctx = context.Background()

func NewRedisStore() *RedisStore {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}
	rdb := redis.NewClient(&redis.Options{
		Addr: redisHost + ":6379",
		DB:   0,
	})
	return &RedisStore{Client: rdb}
}

func NewTestRedisStore() *RedisStore {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}
	rdb := redis.NewClient(&redis.Options{
		Addr: redisHost + ":6379",
		DB:   1,
	})
	return &RedisStore{Client: rdb}
}

// SaveURL stores the originalURL and shortened version in Redis
func (s *RedisStore) SaveURL(originalURL string) (string, error) {
	shortURL, err := s.Client.Get(ctx, originalURL).Result()

	if err == redis.Nil {
		shortURL = s.generateShortURL(originalURL)
		err = s.Client.Set(ctx, originalURL, shortURL, 0).Err()
		if err != nil {
			return "", err
		}

		// Store the originalURL and the shortURL in Redis
		err = s.Client.Set(ctx, shortURL, originalURL, 0).Err()
		if err != nil {
			return "", err
		}

		// Increment the domain count in Redis
		domain, err := s.getDomain(originalURL)
		if err != nil {
			return "", err
		}

		err = s.Client.Incr(ctx, fmt.Sprintf("domain:%s", domain)).Err()
		if err != nil {
			return "", err
		}
	} else if err != nil {
		return "", err
	}
	return shortURL, nil
}

// GetOriginalURL retrieves the OriginalURL from redis using the shortURL
func (s *RedisStore) GetOriginalURL(shortURL string) (string, error) {
	originalURL, err := s.Client.Get(ctx, shortURL).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("URL not found")
	} else if err != nil {
		return "", err
	}

	return originalURL, nil
}

// GetDomainCount retrieves the counts of shortened URLs per domain from redis
func (s *RedisStore) GetDomainCount() (map[string]int, error) {
	keys, err := s.Client.Keys(ctx, "domain:*").Result()
	if err != nil {
		return nil, err
	}

	domainCounts := make(map[string]int)

	for _, key := range keys {
		count, err := s.Client.Get(ctx, key).Int()
		if err != nil {
			return nil, err
		}

		domain := strings.TrimPrefix(key, "domain:")
		domainCounts[domain] = count
	}

	return domainCounts, nil
}

// generateShortURL creates hexademical shortened URL string using hash function
func (s *RedisStore) generateShortURL(originalURL string) string {
	h := fnv.New32a()
	h.Write([]byte(originalURL))
	return fmt.Sprintf("%x", h.Sum32())
}

// getDomain extracts domain name from url
func (s *RedisStore) getDomain(originalURL string) (string, error) {
	parsedURL, err := url.Parse(originalURL)
	if err != nil {
		return "", err
	}

	return strings.TrimPrefix(parsedURL.Host, "www."), nil
}

func (s *RedisStore) FlushTestDB() {
	s.Client.FlushDB(ctx)
}
