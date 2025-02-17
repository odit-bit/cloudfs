package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/alexedwards/scs/v2"
	"github.com/odit-bit/cloudfs/handler/app"
	"github.com/odit-bit/cloudfs/internal/blob"
	"github.com/odit-bit/cloudfs/internal/user"
	repoUser "github.com/odit-bit/cloudfs/internal/user/repo"
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
	flag.BoolVar(&isProd, "production", false, "if true will use production env, default false")
	flag.StringVar(&host, "host", "localhost", "host name")
	flag.IntVar(&port, "port", 8181, "port")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server := new(app.App)
	switch isProd {
	case false:
		server = InitDevService(ctx)

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

func InitDevService(ctx context.Context) *app.App {
	// setup session token
	sess := scs.New() //session management

	// setup account service
	udb, _ := repoUser.NewInMemory()
	users, _ := user.NewStore(ctx, udb, udb)

	// setup blob service
	blobs, _ := blob.NewWithMemory()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

	server := app.New(users, blobs, sess, logger)
	return server
}

// func InitProductionService() service.Cloudfs {
// 	// init s3 constructor
// 	endpoint := os.Getenv(BLOB_MINIO_ENDPOINT)
// 	key := os.Getenv(BLOB_MINIO_ACCESS_KEY)
// 	secret := os.Getenv(BLOB_MINIO_SECRET_ACCESS_KEY)
// 	blobDriver, err := repoBlob.NewMinioBlob(endpoint, key, secret)
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
