package myfs

import (
	"bytes"
	"crypto/sha256"
	"io"
	"os"
	"testing"
)

func Test_Write(t *testing.T) {
	data := []byte("test-input")

	var buf bytes.Buffer
	sha := sha256.New()

	buf.Write(data)
	sum, _ := sha.Write(data)

	f, err := os.OpenFile("test-file.txt", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		t.Fatal(err)
	}

	defer f.Close()
	defer func() {
		os.Remove("test-file.txt")
	}()

	_, err = f.Write(buf.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	f.Seek(0, 0)

	b, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	sha.Reset()
	sum2, _ := sha.Write(b)

	if sum != sum2 {
		t.Fatalf("not same got:%v, expect %v \n ", sum, sum2)
	}

}


//put object into fs 
func Put(filename string, )