package handlers

import (
	"net/http"
	"sort"

	"github.com/bonzonkim/url-shortener/models"
	"github.com/bonzonkim/url-shortener/storage"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	store *storage.RedisStore
}

func NewHandlers(store *storage.RedisStore) *Handlers {
	return &Handlers{store: store}
}


func (h *Handlers) ShortenURL(c *gin.Context) {
	var req models.ShortenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shortURL, err := h.store.SaveURL(req.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.ShortenResponse{ShortURL: shortURL})
}

func (h *Handlers) RedirectURL(c *gin.Context) {
	shortURL := c.Param("shortURL")

	originalURL, err := h.store.GetOriginalURL(shortURL)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	c.Redirect(http.StatusMovedPermanently, originalURL)
}

func (h *Handlers) GetTopDomains(c *gin.Context) {
	domainCounts, err := h.store.GetDomainCount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erorr": err.Error()})
		return
	}

	// sort the domains by count
	type kv struct {
		Key string
		Value int
	}

	var sortedDomains []kv

	for k, v := range domainCounts {
		sortedDomains = append(sortedDomains, kv{k, v})
	}

	sort.Slice(sortedDomains, func(i, j int) bool {
		return sortedDomains[i].Value > sortedDomains[j].Value
	})

	topDomains := make(map[string]int)

	for i, domain := range sortedDomains {
		if i >= 3 {
			break
		}
		topDomains[domain.Key] = domain.Value
	}

	c.JSON(http.StatusOK, gin.H{"top domains": topDomains})
}
