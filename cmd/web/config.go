package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/odit-bit/cloudfs/internal/blob"
	"github.com/odit-bit/cloudfs/internal/user/pguser"
)

type Config struct {
	SessionTokenSecret string
	UserDBEndpoint     string

	//blob storage endpoint
	BSE string
	//blob storage username
	BSU string
	//blob storage password
	BSP string
}

func (conf *Config) validate() error {
	if conf.BSU == "" {
		conf.BSU = os.Getenv("BLOB_STORAGE_USER")
		if conf.BSU == "" {
			return fmt.Errorf("blob storage user not set")
		}
	}

	if conf.BSP == "" {
		conf.BSP = os.Getenv("BLOB_STORAGE_PASSWORD")
		if conf.BSP == "" {
			return fmt.Errorf("blob storage password not set")
		}
	}

	if conf.BSE == "" {
		conf.BSE = os.Getenv("BLOB_STORAGE_ENDPOINT")

		if conf.BSE == "" {
			conf.BSE = blob.DefaultEndpoint
			// return fmt.Errorf("storage end point not set")
		}
	}

	if conf.SessionTokenSecret == "" {
		ts := os.Getenv("SESSION_TOKEN_SECRET")
		if ts == "" {
			// Define the length of the random string in bytes
			length := 8

			// Create a byte slice to store the random bytes
			bytes := make([]byte, length)

			// Read random bytes from crypto/rand
			_, err := rand.Read(bytes)
			if err != nil {
				panic(err)
			}

			// Encode the bytes to a hexadecimal string
			randomString := hex.EncodeToString(bytes)
			fmt.Printf("config create random token %v\n", randomString)
			ts = randomString
		}

		conf.SessionTokenSecret = ts
	}

	if conf.UserDBEndpoint == "" {
		se := os.Getenv("USER_DB_ENDPOINT")
		if se == "" {
			se = pguser.DefaultEndpoint
			// return fmt.Errorf("storage end point not set")
		}
		conf.UserDBEndpoint = se
	}

	return nil
}
