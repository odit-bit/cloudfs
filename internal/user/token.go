package user

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

var (
	Default_Token_Expire = time.Duration(24 * 7 * time.Hour) // a week
)

type Token struct {
	Key    string
	UserID string
	Expire time.Time
}

func NewToken(userID string, expire time.Duration) *Token {
	//make token from this enc
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	key := hex.EncodeToString(b)

	if expire <= 1*time.Hour {
		expire = Default_Token_Expire
	}
	t := Token{
		Key:    key,
		UserID: userID,
		Expire: time.Now().UTC().Add(expire).Round(1 * time.Microsecond),
	}
	return &t
}

// func TokenFromStore(key, userID string, expire time.Time) *Token {
// 	return &Token{
// 		Key:    key,
// 		UserID: userID,
// 		Expire: expire,
// 	}
// }

func (t *Token) IsNotExpire() bool {
	return time.Now().Before(t.Expire)
}

// func (t *Token) Key() string {
// 	return t.key
// }

// func (t *Token) UserID() string {
// 	return t.userID
// }

func (t *Token) ValidUntil() time.Time {
	return t.Expire
}

func (t *Token) RefreshExpire(dur time.Duration) error {
	t.Expire = time.Now().UTC().Add(Default_Token_Expire)
	return nil
}

type TokenOption struct {
	Expire time.Duration
}
