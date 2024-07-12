package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_config(t *testing.T) {
	conf, err := loadConfig("../../cloudfs.yaml")
	if err != nil {
		t.Fatal(err)
	}

	expected := "my-access-key"
	assert.Equal(t, expected, conf.Storage.Blob.AccessKey)

	expected = "postgres://user:password@host/db"
	assert.Equal(t, expected, conf.Storage.User.URI)

}

func Test_LoadConfig(t *testing.T) {
	expected := "my-secret-key-modified"
	os.Setenv(BLOB_STORAGE_SECRET_KEY, expected)

	conf, err := LoadConfig("../../cloudfs.yaml")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, expected, conf.Storage.Blob.SecretKey)
	if !assert.Equal(t, "my-token-secret", conf.Session.Token) {
		t.FailNow()
	}

}
