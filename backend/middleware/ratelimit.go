package middleware

import (
	"fmt"
	"math"
	"net"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	visitors = make(map[string]*visitor)
	mu       sync.Mutex
	r        = rate.Every(200 * time.Millisecond) // 5 req/detik
	burst    = 10
)

func getLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	v, ok := visitors[ip]
	if !ok {
		lim := rate.NewLimiter(r, burst)
		visitors[ip] = &visitor{limiter: lim, lastSeen: time.Now()}
		return lim
	}

	v.lastSeen = time.Now()
	return v.limiter
}

func cleanupVisitors() {
	mu.Lock()
	defer mu.Unlock()

	for ip, v := range visitors {
		if time.Since(v.lastSeen) > 3*time.Minute {
			delete(visitors, ip)
		}
	}
}

func init() {
	go func() {
		for {
			time.Sleep(1 * time.Minute)
			cleanupVisitors()
		}
	}()
}

func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := clientIP(c)
		limiter := getLimiter(ip)

		if !limiter.Allow() {
			reservation := limiter.Reserve()
			delay := reservation.Delay()
			reservation.Cancel()
			retryAfter := int(math.Ceil(delay.Seconds()))
			if retryAfter <= 0 {
				retryAfter = 1
			}

			c.Header("Retry-After", fmt.Sprintf("%d", retryAfter))
			c.AbortWithStatusJSON(429, gin.H{"error": "RATE_LIMITED"})
			return
		}
		c.Next()
	}
}

func clientIP(c *gin.Context) string {
	ip := c.ClientIP()
	if parsed := net.ParseIP(ip); parsed != nil {
		return parsed.String()
	}
	return "0.0.0.0"
}
