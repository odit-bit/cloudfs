package sessionredis

// func RedisClientBuilder(url string) *redis.Client {
// 	opt, err := redis.ParseURL(url)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	cli := redis.NewClient(opt)
// 	cli.Conn()
// 	return cli
// }

// var _ scs.Store = (*sessionStore)(nil)

// type sessionStore struct {
// 	cli *redis.Client
// }

// func NewSessionRedis(url string) *sessionStore {
// 	cli := RedisClientBuilder(url)
// 	return &sessionStore{
// 		cli: cli,
// 	}
// }

// // Commit implements scs.Store.
// func (s *sessionStore) Commit(token string, b []byte, expiry time.Time) error {
// 	second := time.Until(expiry)
// 	ctx, cancel := context.WithTimeout(context.TODO(), 2*time.Second)
// 	defer cancel()
// 	if err := s.cli.Set(ctx, token, b, second).Err(); err != nil {
// 		return err
// 	}
// 	return nil
// }

// // Delete implements scs.Store.
// func (s *sessionStore) Delete(token string) error {
// 	ctx, cancel := context.WithTimeout(context.TODO(), 2*time.Second)
// 	defer cancel()
// 	res := s.cli.Del(ctx, token)
// 	if res.Err() != nil {
// 		return res.Err()
// 	}
// 	return nil
// }

// // Find implements scs.Store.
// func (s *sessionStore) Find(token string) ([]byte, bool, error) {
// 	ctx, cancel := context.WithTimeout(context.TODO(), 2*time.Second)
// 	defer cancel()

// 	res := s.cli.Get(ctx, token)
// 	err := res.Err()
// 	if err != nil {
// 		if err == redis.Nil {
// 			return nil, false, nil
// 		}
// 		return nil, false, err
// 	}

// 	b, err := res.Bytes()
// 	if err != nil {
// 		return nil, false, err
// 	}

// 	return b, true, nil

// }
