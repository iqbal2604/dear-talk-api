package middleware

import (
	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"
	"github.com/ulule/limiter/v3"
	ginlimiter "github.com/ulule/limiter/v3/drivers/middleware/gin"
	redisstore "github.com/ulule/limiter/v3/drivers/store/redis"
)

func RateLimiterMiddleware(redisClient *goredis.Client) gin.HandlerFunc {
	// 60 request per menit per IP
	rate, err := limiter.NewRateFromFormatted("60-M")
	if err != nil {
		panic(err)
	}

	// Pakai Redis sebagai store
	store, err := redisstore.NewStoreWithOptions(
		redisClient,
		limiter.StoreOptions{
			Prefix: "rate_limit",
		},
	)
	if err != nil {
		panic(err)
	}

	instance := limiter.New(store, rate)
	return ginlimiter.NewMiddleware(instance)
}

func StrictRateLimiterMiddleware(redisClient *goredis.Client) gin.HandlerFunc {
	// 10 request per menit untuk endpoint sensitif (login, register)
	rate, err := limiter.NewRateFromFormatted("10-M")
	if err != nil {
		panic(err)
	}

	store, err := redisstore.NewStoreWithOptions(
		redisClient,
		limiter.StoreOptions{
			Prefix: "strict_rate_limit",
		},
	)
	if err != nil {
		panic(err)
	}

	instance := limiter.New(store, rate)
	return ginlimiter.NewMiddleware(instance)
}
