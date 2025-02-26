package blob

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

// represent shareFileToken
// every shared file has this unique token
type Token struct {
	Key      string
	Bucket   string
	Filename string
	Expire   time.Time
}

func NewShareToken(bucket, filename string, expire time.Duration) *Token {
	//make token from this enc
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	key := hex.EncodeToString(b)

	if expire <= 1*time.Hour {
		expire = 1 * time.Hour
	}
	t := Token{
		Key:      key,
		Bucket:   bucket,
		Filename: filename,
		Expire:   time.Now().Add(expire).UTC(),
	}
	return &t
}

// func FromStore(key, userID, filename string, expire time.Time) *Token {
// 	return &Token{
// 		key:      key,
// 		bucket:   userID,
// 		filename: filename,
// 		expire:   expire,
// 	}
// }

func (t *Token) IsNotExpire() bool {
	return time.Now().Before(t.Expire)
}

// func (t *Token) Key() string {
// 	return t.key
// }

// func (t *Token) Bucket() string {
// 	return t.bucket
// }

// func (t *Token) Filename() string {
// 	return t.filename
// }

func (t *Token) ValidUntil() time.Time {
	return t.Expire
}

func (t *Token) IsFilenameEqualTo(filename string) bool {
	return t.Filename == filename
}

// func (t *Token) IsUserEqualTo(userID string) bool {
// 	return t.userID == userID
// }
