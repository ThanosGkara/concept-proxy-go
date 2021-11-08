package intercept

import (
	"net/http"

	"github.com/patrickmn/go-cache"
)

// Needed to pass also the cache_key and the pointer of cache
type customWriter struct {
	http.ResponseWriter
	pcache    *cache.Cache
	cache_key string
}

func NewCustomWriter(w http.ResponseWriter, pcache *cache.Cache, cache_key string) *customWriter {
	return &customWriter{w, pcache, cache_key}
}

func (c *customWriter) Header() http.Header {
	return c.ResponseWriter.Header()
}

func (c *customWriter) Write(data []byte) (int, error) {

	// First write the data to the cache and then return the responce to the client
	c.pcache.Set(c.cache_key, string(data), cache.DefaultExpiration)

	return c.ResponseWriter.Write(data)
}

func (c *customWriter) WriteHeader(i int) {
	c.ResponseWriter.WriteHeader(i)
}
