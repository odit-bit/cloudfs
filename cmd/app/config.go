package main

import (
	"time"
)

const (
	BLOB_STORAGE_ENDPOINT   = "BLOB_STORAGE_ENDPOINT"
	BLOB_STORAGE_ACCESS_KEY = "BLOB_STORAGE_ACCESS_KEY"
	BLOB_STORAGE_SECRET_KEY = "BLOB_STORAGE_SECRET_KEY"

	USER_DB_URI = "USER_DB_URI"

	TOKEN_DB_URI = "TOKEN_DB_URI"

	SESSION_TOKEN_SECRET = "SESSION_TOKEN_SECRET"

	HOST = "HOST"
	PORT = "PORT"
)

type config struct {
	HTTP struct {
		Host  string
		Port  string
		Limit int
		Burst int
	}
	Storage struct {
		Blob struct {
			Endpoint  string
			AccessKey string
			SecretKey string
		}
		User struct {
			URI string
		}
		Token struct {
			URI string
		}
	}
	Session struct {
		Token   string
		Expired time.Duration
	}
}

func (c *config) validate() error {
	// if c.Session.Token == "" {
	// 	// Define the length of the random string in bytes
	// 	length := 8

	// 	// Create a byte slice to store the random bytes
	// 	bytes := make([]byte, length)

	// 	// Read random bytes from crypto/rand
	// 	_, err := rand.Read(bytes)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	// Encode the bytes to a hexadecimal string
	// 	randomString := hex.EncodeToString(bytes)
	// 	fmt.Printf("config create random token %v\n", randomString)
	// 	c.Session.Token = randomString
	// }

	if c.Session.Expired < 1*time.Hour {
		c.Session.Expired = 4 * time.Hour
	}
	return nil
}

func (conf *config) Address() string {
	return conf.HTTP.Host + ":" + conf.HTTP.Port
}

// type xConfig struct {
// 	v *viper.Viper
// 	config
// }

// func LoadConfig(path string) (*xConfig, error) {
// 	c, err := loadConfig(path)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if err := c.loadEnv(); err != nil {
// 		return nil, err
// 	}
// 	if err := c.validate(); err != nil {
// 		return nil, err
// 	}
// 	return c, nil
// }

// func loadConfig(path string) (*xConfig, error) {
// 	v := viper.New()
// 	v.SetConfigFile(path)
// 	if err := v.ReadInConfig(); err != nil {
// 		return nil, err
// 	}

// 	var conf config
// 	if err := v.Unmarshal(&conf); err != nil {
// 		return nil, err
// 	}
// 	c := xConfig{
// 		v:      v,
// 		config: conf,
// 	}
// 	return &c, nil
// }

// func (c *xConfig) loadEnv() error {
// 	c.v.BindEnv("storage.blob.endpoint", BLOB_STORAGE_ENDPOINT)
// 	c.v.BindEnv("storage.blob.accesskey", BLOB_STORAGE_ACCESS_KEY)
// 	c.v.BindEnv("storage.blob.secretkey", BLOB_STORAGE_SECRET_KEY)
// 	c.v.BindEnv("storage.user.uri", USER_DB_URI)

// 	var conf config
// 	if err := c.v.Unmarshal(&conf); err != nil {
// 		return err
// 	}
// 	c.config = conf
// 	return nil
// }
