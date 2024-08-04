package token

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

// represent shareFileToken
// every shared file has this unique token
type ShareToken struct {
	key      string
	userID   string
	filename string
	expire   time.Time
}

func NewShareToken(userID, filename string, expire time.Duration) *ShareToken {
	//make token from this enc
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	key := hex.EncodeToString(b)

	if expire <= 1*time.Hour {
		expire = 1 * time.Hour
	}
	t := ShareToken{
		key:      key,
		userID:   userID,
		filename: filename,
		expire:   time.Now().Add(expire),
	}
	return &t
}

func FromStore(key, userID, filename string, expire time.Time) *ShareToken {
	return &ShareToken{
		key:      key,
		userID:   userID,
		filename: filename,
		expire:   expire,
	}
}

func (t *ShareToken) IsNotExpire() bool {
	return time.Now().Before(t.expire)
}

func (t *ShareToken) Key() string {
	return t.key
}

func (t *ShareToken) UserID() string {
	return t.userID
}

func (t *ShareToken) Filename() string {
	return t.filename
}

func (t *ShareToken) ValidUntil() time.Time {
	return t.expire
}

func (t *ShareToken) IsFilenameEqualTo(filename string) bool {
	return t.filename == filename
}

func (t *ShareToken) IsUserEqualTo(userID string) bool {
	return t.userID == userID
}
