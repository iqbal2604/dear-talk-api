package redis

import (
	"context"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type TokenBlacklist struct {
	client *goredis.Client
}

func NewTokenBlacklist(client *goredis.Client) *TokenBlacklist {
	return &TokenBlacklist{client: client}
}

// Tambah token ke blacklist dengan expiry sama seperti token expire
func (b *TokenBlacklist) Add(ctx context.Context, token string, expiry time.Duration) error {
	key := fmt.Sprintf("blacklist:%s", token)
	return b.client.Set(ctx, key, "1", expiry).Err()
}

// Cek apakah token ada di blacklist
func (b *TokenBlacklist) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	key := fmt.Sprintf("blacklist:%s", token)
	result, err := b.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}
