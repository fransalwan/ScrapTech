package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"scrapflow-backend/pkg/response"
)

type idempotencyEntry struct {
	CreatedAt time.Time
}

var (
	idempotencyStore = make(map[string]idempotencyEntry)
	storeMutex       sync.Mutex
)

// IdempotencyMiddleware ensures mutating requests with an Idempotency-Key header are executed once
func IdempotencyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		idempotencyKey := c.GetHeader("Idempotency-Key")
		if idempotencyKey == "" {
			response.BadRequest(c, "Missing required 'Idempotency-Key' header for financial operations", nil)
			c.Abort()
			return
		}

		storeMutex.Lock()
		entry, exists := idempotencyStore[idempotencyKey]
		if exists && time.Since(entry.CreatedAt) < 5*time.Minute {
			storeMutex.Unlock()
			response.Error(c, 409, "Conflict: A request with this Idempotency-Key has already been processed or is in-flight", nil)
			c.Abort()
			return
		}

		// Reserve key
		idempotencyStore[idempotencyKey] = idempotencyEntry{CreatedAt: time.Now()}
		storeMutex.Unlock()

		c.Set("idempotency_key", idempotencyKey)
		c.Next()
	}
}
