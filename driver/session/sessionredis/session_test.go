package sessionredis_test

// func Test_session(t *testing.T) {
// 	uri := "redis://:@localhost:6379/0" // password set
// 	sr := sessionredis.NewSessionRedis(uri)

// 	token := "token"
// 	value := []byte("token-value")
// 	expiry := time.Now().Add(2 * time.Second)
// 	if err := sr.Commit(token, value, expiry); err != nil {
// 		t.Fatal(err)
// 	}

// 	_, ok, err := sr.Find(token)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if !ok {
// 		t.Fatal("should ok")
// 	}
// }
