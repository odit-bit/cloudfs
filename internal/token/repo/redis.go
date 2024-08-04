package repo

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/odit-bit/cloudfs/internal/token"
	"github.com/redis/go-redis/v9"
)

// var _ service.TokenStore = (*RedisToken)(nil)

func redisClientBuilder(url string) *redis.Client {
	opt, err := redis.ParseURL(url)
	if err != nil {
		log.Fatal(err)
	}
	cli := redis.NewClient(opt)
	cli.Conn()
	return cli
}

type RedisToken struct {
	cli *redis.Client
}

func NewRedisToken(redisURL string) RedisToken {
	cli := redisClientBuilder(redisURL)
	return RedisToken{
		cli: cli,
	}
}

type redisShareToken struct {
	UserID     string
	Filename   string
	ValidUntil time.Time
}

// Get implements service.TokenStore.
func (r *RedisToken) Get(ctx context.Context, tokenString string) (*token.ShareToken, bool, error) {
	res := r.cli.Get(ctx, tokenString)
	err := res.Err()
	if err != nil {
		if err == redis.Nil {
			return &token.ShareToken{}, false, nil
		}
		return &token.ShareToken{}, false, err
	}

	b, err := res.Bytes()
	if err != nil {
		return nil, false, err
	}

	var rst redisShareToken
	if err := json.Unmarshal(b, &rst); err != nil {
		return nil, false, err
	}
	tkn := token.FromStore(tokenString, rst.UserID, rst.Filename, rst.ValidUntil)
	return tkn, true, err
}

// Put implements service.TokenStore.
func (r *RedisToken) Put(ctx context.Context, token *token.ShareToken) error {
	dur := time.Until(token.ValidUntil())
	rst := redisShareToken{
		UserID:     token.UserID(),
		Filename:   token.Filename(),
		ValidUntil: token.ValidUntil(),
	}

	b, err := json.Marshal(rst)
	if err != nil {
		return err
	}
	res := r.cli.Set(ctx, token.Key(), b, dur)
	return res.Err()
}
