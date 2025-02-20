package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/alexedwards/scs/v2"
	"github.com/odit-bit/cloudfs/web"
	"github.com/odit-bit/cloudfs/lib/xhttp"
)

const (
	BLOB_MINIO_ENDPOINT          = "BLOB_MINIO_ENDPOINT"
	BLOB_MINIO_ACCESS_KEY        = "BLOB_MINIO_ACCESS_KEY"
	BLOB_MINIO_SECRET_ACCESS_KEY = "BLOB_MINIO_SECRET_ACCESS_KEY"

	USER_PG_URI = "USER_PG_URI"

	TOKEN_REDIS_URI = "TOKEN_REDIS_URI"

	SESSION_REDIS_URI = "SESSION_REDIS_URI"
)

func main() {

	var isProd bool
	var port int
	var host string
	var backendAddr string
	flag.BoolVar(&isProd, "production", false, "if true will use production env, default false")
	flag.StringVar(&host, "host", "localhost", "host name")
	flag.IntVar(&port, "port", 8181, "port")
	flag.StringVar(&backendAddr, "backend-addr", "", "backend address")

	flag.Parse()

	if backendAddr == "" {
		log.Println("missing backend address")
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server := new(web.App)
	var err error
	switch isProd {
	case false:
		server, err = InitDevService(ctx, backendAddr)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
	default:
		panic("production code is not implemented")
		// sess.Store = sessionredis.NewSessionRedis(os.Getenv(SESSION_REDIS_URI))
		// cfs = InitProductionService()
	}

	address := fmt.Sprintf("%s:%d", host, port)
	if err := server.Run(
		ctx,
		address,
		// xhttp.LogRequest(logger),
		xhttp.CorsDefault(),
	); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}

func InitDevService(ctx context.Context, backendAddr string) (*web.App, error) {
	// setup session token
	sess := scs.New() //session management

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

	server, err := web.New(backendAddr, sess, logger)
	if err != nil {
		return nil, err
	}
	return server, nil
}

// func InitProductionService() service.Cloudfs {
// 	// init s3 constructor
// 	endpoint := os.Getenv(BLOB_MINIO_ENDPOINT)
// 	key := os.Getenv(BLOB_MINIO_ACCESS_KEY)
// 	secret := os.Getenv(BLOB_MINIO_SECRET_ACCESS_KEY)
// 	blobDriver, err := repostorage.NewMinioBlob(endpoint, key, secret)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// init PG
// 	pgURI := os.Getenv(USER_PG_URI)
// 	userDriver, err := repoUser.NewUserPG(context.TODO(), pgURI)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// init token
// 	redURI := os.Getenv(TOKEN_REDIS_URI)
// 	tokenDriver := repoToken.NewRedisToken(redURI)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	svc, err := service.NewCloudfs(&tokenDriver, blobDriver, userDriver)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return svc
// }
