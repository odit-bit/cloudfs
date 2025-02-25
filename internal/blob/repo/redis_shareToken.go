package repo

// var _ storage.TokenStorer = (*tokenStore)(nil)

// type tokenStore struct {
// 	redisCli *redis.Client
// }

// func NewRedisObjectToken(cli *redis.Client) *tokenStore {
// 	return &tokenStore{
// 		redisCli: cli,
// 	}
// }

// // Delete implements storage.TokenStorer.
// func (t *tokenStore) Delete(ctx context.Context, tokenKey string) error {
// 	return t.redisCli.Del(ctx, tokenKey).Err()
// }

// // Get implements storage.TokenStorer.
// func (t *tokenStore) Get(ctx context.Context, tokenKey string) (*storage.Token, bool, error) {
// 	tkn := &storage.Token{}
// 	if err := t.redisCli.HGetAll(ctx, tokenKey).Scan(tkn); err != nil {
// 		return nil, false, err
// 	}
// 	return tkn, true, nil
// }

// // Put implements storage.TokenStorer.
// func (t *tokenStore) Put(ctx context.Context, token *storage.Token) storage.OpErr {

// 	pipe := t.redisCli.TxPipeline()
// 	pipe.HSet(ctx, token.Key, token)

// 	dur := time.Until(token.ValidUntil())
// 	pipe.SetNX(ctx, token.Bucket, token.Key, dur)

// 	if _, err := pipe.Exec(ctx); err != nil {
// 		return storage.NewException(err)
// 	}
// 	return nil
// }

// func (t *tokenStore) GetByBucket(ctx context.Context, bucket string) (*storage.Token, bool, error) {
// 	key, err := t.redisCli.Get(ctx, bucket).Result()
// 	if err != nil {
// 		return nil, false, nil
// 	}

// 	tkn := &storage.Token{}
// 	if err := t.redisCli.HGetAll(ctx, key).Scan(tkn); err != nil {
// 		return nil, false, err
// 	}
// 	return tkn, true, nil

// }

// //////////////////
