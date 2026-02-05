package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/shpdwx/claims/common"
)

type Cache interface {
	Rdb() *redis.Client
	LimitTokens(limit, uid int64, str string) error
}

type defCache struct {
	ctx context.Context
	rdb *redis.Client
}

func NewCache(ctx context.Context) Cache {
	return &defCache{
		ctx: ctx,
		rdb: common.InitRedis(ctx),
	}
}

func (c *defCache) Rdb() *redis.Client {
	return c.rdb
}

func (c *defCache) LimitTokens(limit, uid int64, str string) error {
	if limit < 1 {
		err := errors.New("Token数量未设置")
		return err
	}

	rdb := c.rdb
	key := fmt.Sprintf("user_tokens:%d", uid)

	num, _ := rdb.SCard(c.ctx, key).Result()
	if num >= limit {
		err := errors.New("超出限制的Token数量")
		return err
	}

	rdb.SAdd(c.ctx, key, str)
	return nil
}
