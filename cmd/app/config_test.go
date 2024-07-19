package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_config(t *testing.T) {
	conf, err := loadConfig("../../cloudfs.yaml")
	if err != nil {
		t.Fatal(err)
	}

	expected := "admin"
	assert.Equal(t, expected, conf.Storage.Blob.AccessKey)

	expected = "postgres://admin:admin@localhost/postgres"
	assert.Equal(t, expected, conf.Storage.User.URI)

}
