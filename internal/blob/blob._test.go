package blob

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_iterator(t *testing.T) {

	tc := []struct {
		input []ObjectInfo
	}{
		{input: []ObjectInfo{
			{
				UserID:   "1",
				Filename: "obj1",
			},
			{
				UserID:   "2",
				Filename: "obj2",
			},
		}},
		{input: []ObjectInfo{}},
	}

	// //
	// infos := make([]*ObjectInfo, 2)
	// infos[0] = &ObjectInfo{
	// 	UserID:   "1",
	// 	Filename: "obj1",
	// }
	// infos[1] = &ObjectInfo{
	// 	UserID:   "2",
	// 	Filename: "obj2",
	// }

	for _, test := range tc {
		objC := make(chan ObjectInfo, 10)

		for _, v := range test.input {
			objC <- v
		}
		close(objC)

		it := Iterator{
			UserID: "",
			C:      objC,
			obj:    ObjectInfo{},
			// err:    nil,
		}

		for _, expected := range test.input {
			if !it.Next() {
				t.Fatal("should return info")
			}

			obj := it.Value()
			assert.Equal(t, expected.UserID, obj.UserID)
		}

	}
}

func Test_readBlob(t *testing.T) {

	data := []byte("this is data")

	buf := bytes.NewBuffer(data)
	rc := io.NopCloser(buf)
	obj := ObjectInfo{
		Size: int64(len(data)),
		Data: rc,
	}

	dst := bytes.Buffer{}
	defer obj.Data.Close()
	n, err := io.Copy(&dst, obj.Data)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, n, obj.Size)
	assert.Equal(t, data, dst.Bytes())
}
