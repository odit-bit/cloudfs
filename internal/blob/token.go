package blob

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)




// represent shareFileToken
// every shared file has this unique token
type Token struct {
	key      string
	userID   string
	filename string
	expire   time.Time
}

func NewShareToken(userID, filename string, expire time.Duration) *Token {
	//make token from this enc
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	key := hex.EncodeToString(b)

	if expire <= 1*time.Hour {
		expire = 1 * time.Hour
	}
	t := Token{
		key:      key,
		userID:   userID,
		filename: filename,
		expire:   time.Now().Add(expire),
	}
	return &t
}

func FromStore(key, userID, filename string, expire time.Time) *Token {
	return &Token{
		key:      key,
		userID:   userID,
		filename: filename,
		expire:   expire,
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

func (t *Token) Filename() string {
	return t.filename
}

func (t *Token) ValidUntil() time.Time {
	return t.expire
}

func (t *Token) IsFilenameEqualTo(filename string) bool {
	return t.filename == filename
}

func (t *Token) IsUserEqualTo(userID string) bool {
	return t.userID == userID
}
