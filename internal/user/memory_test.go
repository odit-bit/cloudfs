package user

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_account(t *testing.T) {
	inMem, _ := newInMemory()
	acc := CreateAccount("username", "password")
	if err := inMem.Insert(context.TODO(), acc); err != nil {
		t.Fatal(err)
	}

	//expect success
	acc2, err := inMem.FindUsername(context.TODO(), "username")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, acc, acc2)

	// expect error not nil
	if _, err := inMem.FindUsername(context.Background(), "not-exist"); err == nil {
		t.Fatal("fineUsername should not existed")
	} else {
		assert.EqualError(t, err, ErrAccountNotExist.Error())
	}

}

func Test_token(t *testing.T) {
	inMem, _ := newInMemory()
	tkn := NewToken("12345", 10*time.Minute)
	if err := inMem.PutToken(context.TODO(), tkn); err != nil {
		t.Fatal(err)
	}

	tkn2, err := inMem.GetToken(context.TODO(), tkn.Key)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, tkn, tkn2)

	tkn3, ok := inMem.GetTokenUserID(context.TODO(), tkn.UserID)
	if !ok {
		t.Fatal("token should exist finded by user id")
	}
	assert.Equal(t, tkn2, tkn3)

	inMem.Delete(context.TODO(), tkn.Key)

	// expect err not nil
	if _, err := inMem.GetToken(context.TODO(), tkn.Key); err == nil {
		t.Fatal(err)
	} else {
		assert.EqualError(t, err, ErrTokenNotExist.Error())
	}
}
