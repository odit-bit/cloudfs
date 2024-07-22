package blob

import (
	"testing"
	"time"
)

func Test_shareToken(t *testing.T) {
	token := NewShareToken("id-123", "myfile", time.Duration(1*time.Minute))
	if ok := token.IsNotExpire(); !ok {
		t.Fatal("should ok")
	}

	token = NewShareToken("id-123", "myfile", time.Duration(1*time.Nanosecond))
	time.Sleep(2 * time.Nanosecond)
	if ok := token.IsNotExpire(); ok {
		t.Fatal("should not ok")
	}
}
