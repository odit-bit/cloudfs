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
	key    string
	userID string
	expire time.Time
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
		key:    key,
		userID: userID,
		expire: time.Now().UTC().Add(expire),
	}
	return &t
}

func FromStore(key, userID, filename string, expire time.Time) *Token {
	return &Token{
		key:    key,
		userID: userID,
		expire: expire,
	}
}

func (t *Token) IsNotExpire() bool {
	return time.Now().Before(t.expire)
}

func (t *Token) Key() string {
	return t.key
}

func (t *Token) UserID() string {
	return t.userID
}

func (t *Token) ValidUntil() time.Time {
	return t.expire
}

func (t *Token) RefreshExpire(dur time.Duration) error {
	t.expire = time.Now().UTC().Add(Default_Token_Expire)
	return nil
}

type TokenOption struct {
	Expire time.Duration
}
