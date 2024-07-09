package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/odit-bit/cloudfs/internal/user/pguser"
)

func connectMinio(endpoint, accessKeyID, secretAccessKey string, secure bool) (*minio.Client, error) {

	cli, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: secure,
	})
	if err != nil {
		return nil, err
	}

	cancel, err := cli.HealthCheck(2 * time.Second)
	if err != nil {
		return nil, err
	}
	defer cancel()

	if ok := cli.IsOnline(); !ok {
		return nil, fmt.Errorf("storage endpoint is offline, api-endpoint: %v", cli.EndpointURL().String())
	}

	if err := cli.MakeBucket(context.TODO(), "init-bucket", minio.MakeBucketOptions{}); err != nil {
		return nil, err
	}

	return cli, nil

}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var sessionTokenSecret, userDBEndpoint string
	var bsu, bsp, bse string
	flag.StringVar(&sessionTokenSecret, "token", "", "token secret")
	flag.StringVar(&userDBEndpoint, "user_db_addr", "", "host:port")
	flag.StringVar(&bse, "bse", "", "host:port")
	flag.StringVar(&bsu, "bsu", "", "blob storage user value")
	flag.StringVar(&bsp, "bsp", "", "blob storage password")

	flag.Parse()

	conf := Config{
		SessionTokenSecret: sessionTokenSecret,
		UserDBEndpoint:     userDBEndpoint,
		BSE:                bse,
		BSU:                bsu,
		BSP:                bsp,
	}

	if err := conf.validate(); err != nil {
		log.Fatal(err)
	}

	udb, err := pguser.ConnectDB(ctx, conf.UserDBEndpoint)
	if err != nil {
		log.Fatal(err)
	}

	mo, err := connectMinio(conf.BSE, conf.BSU, conf.BSP, false)
	if err != nil {
		log.Fatal(err)
	}

	cs := sessions.NewCookieStore([]byte(conf.SessionTokenSecret))

	app := api{
		blobStorage: mo,
		userDB:      udb,
		cs:          cs,
	}
	_ = app

	useEcho(&app)
}

func useEcho(app *api) {
	e := echo.New()

	// e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(5)))

	// html page endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.String(200, "health")
	})

	e.GET("/", func(c echo.Context) error {
		app.serveIndex(c.Response(), c.Request())
		return nil
	})

	e.GET(loginHTML, func(c echo.Context) error {
		app.serveLoginPage(loginAPIEndpoint).ServeHTTP(c.Response(), c.Request())
		return nil
	})

	e.GET(listHTML, func(c echo.Context) error {
		app.serveListPage(listAPIEndpoint).ServeHTTP(c.Response(), c.Request())
		return nil
	})

	e.GET(uploadHTML, func(c echo.Context) error {
		app.serveUploadPage(uploadAPIEndpoint).ServeHTTP(c.Response(), c.Request())
		return nil
	})

	e.GET(registerHTML, func(c echo.Context) error {
		app.serveRegisterPage(registerAPIEndpoint).ServeHTTP(c.Response(), c.Request())
		return nil
	})

	//=============================================================
	// service endpoint

	//user related

	e.POST(loginAPIEndpoint, func(c echo.Context) error {
		app.handleLogin(c.Response(), c.Request())
		return nil
	})

	e.POST(registerAPIEndpoint, func(c echo.Context) error {
		app.handleRegister(c.Response(), c.Request())
		return nil
	})

	//storage related

	e.GET(listAPIEndpoint, func(c echo.Context) error {
		app.HandleList(c.Response(), c.Request())
		return nil
	})

	e.POST(uploadAPIEndpoint, func(c echo.Context) error {
		app.handleUpload(c.Response(), c.Request())
		return nil
	})

	e.GET("/api/download", func(c echo.Context) error {
		app.handleDownload(c.Response(), c.Request())
		return nil
	})

	e.Start(":8181")
}
