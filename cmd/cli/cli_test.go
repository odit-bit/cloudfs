package main

import (
	"os"
	"testing"
)

func Test_(t *testing.T) {
	toConf := "/mnt/d/wsl"
	expectPath := "/mndt/d/wsl/.cloudfs/config.json"
	if err := createConfigFile(toConf); err != nil {
		t.Fatal(err)
	}
	if _, err := loadConfig(toConf); err != nil {
		t.Fatal(err)
	}

	os.Remove(expectPath)

}
