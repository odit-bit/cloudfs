package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"

	"github.com/fsnotify/fsnotify"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/odit-bit/cloudfs/internal/ui"
	"github.com/spf13/viper"
)

//this package still WIP

type config struct {
	Api struct {
		Scheme string
		Host   string
		List   struct {
			Path string
		}
		Upload struct {
			Path string
		}
		Login struct {
			Path string
		}
		Register struct {
			Path string
		}
	}
	baseURL url.URL
}

func (c *config) Validate() error {
	if c.Api.Scheme == "" {
		return fmt.Errorf("config: api scheme cannot be nil")
	}
	if c.Api.Host == "" {
		return fmt.Errorf("config: api Host cannot be nil")
	}

	if c.Api.List.Path == "" {
		return fmt.Errorf("config: api 'list' path cannot be nil, you mean '/' ? ")
	}

	if c.Api.Upload.Path == "" {
		return fmt.Errorf("config: api 'upload' path cannot be nil, you mean '/' ? ")
	}

	if c.Api.Login.Path == "" {
		return fmt.Errorf("config: api 'login' path cannot be nil, you mean '/' ? ")
	}

	if c.Api.Register.Path == "" {
		return fmt.Errorf("config: api 'register' path cannot be nil, you mean '/' ? ")
	}

	c.baseURL = url.URL{
		Scheme: c.Api.Scheme,
		Host:   c.Api.Host,
	}

	return nil
}

func (c *config) ListURL() string {
	u := c.baseURL
	u.Path = c.Api.List.Path
	return u.String()
}

func (c *config) UploadURL() string {
	u := c.baseURL
	u.Path = c.Api.Upload.Path
	return u.String()
}

func main() {
	var isWatchConfig = flag.Bool("watch", false, "watch config file for change")
	flag.Parse()

	viper.SetConfigFile("./ui-config.yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	var c config
	if err := viper.Unmarshal(&c); err != nil {
		log.Fatal(err)
	}
	if err := c.Validate(); err != nil {
		log.Fatal(err)
	}

	if *isWatchConfig {
		viper.WatchConfig()
		viper.OnConfigChange(func(in fsnotify.Event) {
			if err := viper.Unmarshal(&c); err != nil {
				log.Fatal(err)
			}
			if err := c.Validate(); err != nil {
				log.Fatal(err)
			}

			log.Println("got config change", c)
		})
		log.Println("watchimg config file")
	}

	// sigC := make(chan os.Signal, 1)
	// signal.Notify(sigC, os.Interrupt)

	// <-sigC

	e := echo.New()

	// e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(5)))

	index := ui.NewIndexPage()
	index.AddMenu("list", c.ListURL())
	index.AddMenu("upload", c.UploadURL())

	e.GET("/", func(c echo.Context) error {
		index.Render(c.Response())
		return nil
	})

	e.Start(":8080")
}
