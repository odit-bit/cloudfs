package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/spf13/viper"
)

const (
	BLOB_STORAGE_ENDPOINT   = "BLOB_STORAGE_ENDPOINT"
	BLOB_STORAGE_ACCESS_KEY = "BLOB_STORAGE_ACCESS_KEY"
	BLOB_STORAGE_SECRET_KEY = "BLOB_STORAGE_SECRET_KEY"

	USER_DATABASE_URI = "USER_DATABASE_URI"

	SESSION_TOKEN_SECRET = "SESSION_TOKEN_SECRET"
)

type config struct {
	HTTP struct {
		HOST string
		PORT string
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
	}
	Session struct {
		Token string
	}
}

func (c *config) validate() error {
	if c.Session.Token == "" {
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
		c.Session.Token = randomString
	}

	return nil
}

func (conf *config) Address() string {
	return conf.HTTP.HOST + ":" + conf.HTTP.PORT
}

type xConfig struct {
	v *viper.Viper
	config
}

func LoadConfig(path string) (*xConfig, error) {
	c, err := loadConfig(path)
	if err != nil {
		return nil, err
	}
	if err := c.loadEnv(); err != nil {
		return nil, err
	}
	if err := c.validate(); err != nil {
		return nil, err
	}
	return c, nil
}

func loadConfig(path string) (*xConfig, error) {
	v := viper.New()
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var conf config
	if err := v.Unmarshal(&conf); err != nil {
		return nil, err
	}
	c := xConfig{
		v:      v,
		config: conf,
	}
	return &c, nil
}

func (c *xConfig) loadEnv() error {
	c.v.BindEnv("storage.blob.endpoint", BLOB_STORAGE_ENDPOINT)
	c.v.BindEnv("storage.blob.accesskey", BLOB_STORAGE_ACCESS_KEY)
	c.v.BindEnv("storage.blob.secretkey", BLOB_STORAGE_SECRET_KEY)
	c.v.BindEnv("storage.user.uri", USER_DATABASE_URI)

	var conf config
	if err := c.v.Unmarshal(&conf); err != nil {
		return err
	}
	c.config = conf
	return nil
}

// func (c *xConfig) bindFlag() error {
// 	var sessionTokenSecret, userDBEndpoint string
// 	var bsu, bsp, bse string
// 	flag.StringVar(&sessionTokenSecret, "token", "", "token secret")
// 	flag.StringVar(&userDBEndpoint, "user_db_addr", "", "host:port")
// 	flag.StringVar(&bse, "bse", "", "host:port")
// 	flag.StringVar(&bsu, "bsu", "", "blob storage user value")
// 	flag.StringVar(&bsp, "bsp", "", "blob storage password")
// 	flag.Parse()

// 	c.v.BindPFlag("storage.blob.accesskey", pflag.Lookup(""))
// 	c.v.Set("storage.blob.accesskey", bsu)
// 	c.v.Set("storage.blob.secretkey", bsp)
// 	c.v.Set("storage.user.uri", userDBEndpoint)

// 	var conf config
// 	if err := c.v.Unmarshal(&conf); err != nil {
// 		return err
// 	}
// 	c.config = conf
// 	return nil
// }

// /////

// type Config struct {
// 	SessionTokenSecret string
// 	UserDBEndpoint     string

// 	//blob storage endpoint
// 	BSE string
// 	//blob storage username
// 	BSU string
// 	//blob storage password
// 	BSP string
// }

// func (conf *Config) validate() error {
// 	if conf.BSU == "" {
// 		conf.BSU = os.Getenv("BLOB_STORAGE_USER")
// 		if conf.BSU == "" {
// 			return fmt.Errorf("blob storage user not set")
// 		}
// 	}

// 	if conf.BSP == "" {
// 		conf.BSP = os.Getenv("BLOB_STORAGE_PASSWORD")
// 		if conf.BSP == "" {
// 			return fmt.Errorf("blob storage password not set")
// 		}
// 	}

// 	if conf.BSE == "" {
// 		conf.BSE = os.Getenv("BLOB_STORAGE_ENDPOINT")

// 		if conf.BSE == "" {
// 			conf.BSE = blob.DefaultEndpoint
// 			// return fmt.Errorf("storage end point not set")
// 		}
// 	}

// 	if conf.SessionTokenSecret == "" {
// 		ts := os.Getenv("SESSION_TOKEN_SECRET")
// 		if ts == "" {
// 			// Define the length of the random string in bytes
// 			length := 8

// 			// Create a byte slice to store the random bytes
// 			bytes := make([]byte, length)

// 			// Read random bytes from crypto/rand
// 			_, err := rand.Read(bytes)
// 			if err != nil {
// 				panic(err)
// 			}

// 			// Encode the bytes to a hexadecimal string
// 			randomString := hex.EncodeToString(bytes)
// 			fmt.Printf("config create random token %v\n", randomString)
// 			ts = randomString
// 		}

// 		conf.SessionTokenSecret = ts
// 	}

// 	if conf.UserDBEndpoint == "" {
// 		se := os.Getenv("USER_DB_ENDPOINT")
// 		if se == "" {
// 			se = pguser.DefaultEndpoint
// 			// return fmt.Errorf("storage end point not set")
// 		}
// 		conf.UserDBEndpoint = se
// 	}

// 	return nil
// }
