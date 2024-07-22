package blob

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

// represent shareFileToken
// every shared file has this unique token
type ShareToken struct {
	Key      string
	UserID   string
	Filename string
	Expire   time.Time
}

func NewShareToken(userID, filename string, expire time.Duration) *ShareToken {
	//make token from this enc
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	key := hex.EncodeToString(b)

	t := ShareToken{
		Key:      key,
		UserID:   userID,
		Filename: filename,
		Expire:   time.Now().Add(expire),
	}
	return &t
}

func (t *ShareToken) IsNotExpire() bool {
	return time.Now().Before(t.Expire)
}
