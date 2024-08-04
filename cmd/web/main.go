package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/alexedwards/scs/v2"
	"github.com/odit-bit/cloudfs/driver/session/sessionredis"
	"github.com/odit-bit/cloudfs/handler/app"
	repoBlob "github.com/odit-bit/cloudfs/internal/blob/repo"
	repoToken "github.com/odit-bit/cloudfs/internal/token/repo"
	repoUser "github.com/odit-bit/cloudfs/internal/user/repo"
	"github.com/odit-bit/cloudfs/lib/xhttp"
	"github.com/odit-bit/cloudfs/service"
)

func main() {

	var isProd bool
	var port int
	var host string
	flag.BoolVar(&isProd, "production", false, "if true will use production env, default false")
	flag.StringVar(&host, "host", "0.0.0.0", "host name")
	flag.IntVar(&port, "port", 8181, "port")
	flag.Parse()

	ctx := context.TODO()

	var cfs service.Cloudfs
	switch isProd {
	case false:
		cfs = InitDevService()
	default:
		cfs = InitProductionService()
	}

	sess := scs.New() //session management
	sess.Store = sessionredis.NewSessionRedis(os.Getenv(SESSION_REDIS_URI))

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

	address := fmt.Sprintf("%s:%d", host, port)

	server := app.New(&cfs, sess, logger)
	if err := server.Run(
		ctx,
		address,
		xhttp.LogRequest(logger),
		xhttp.CorsDefault(),
	); err != nil {
		log.Fatal(err)
	}

}

func InitDevService() service.Cloudfs {
	//token
	ts, err := repoToken.NewInMemToken("")
	if err != nil {
		log.Fatal(err)
	}

	// blob
	bs, err := repoBlob.NewInMemBlob()
	if err != nil {
		log.Fatal(err)
	}

	// user
	us, err := repoUser.NewSQLiteDB("")
	if err != nil {
		log.Fatal(err)
	}

	svc, err := service.NewCloudfs(&ts, &bs, &us)
	if err != nil {
		log.Fatal(err)
	}

	return svc
}

func InitProductionService() service.Cloudfs {
	// init s3 constructor
	endpoint := os.Getenv(BLOB_MINIO_ENDPOINT)
	key := os.Getenv(BLOB_MINIO_ACCESS_KEY)
	secret := os.Getenv(BLOB_MINIO_SECRET_ACCESS_KEY)
	blobDriver, err := repoBlob.NewMinioBlob(endpoint, key, secret)
	if err != nil {
		log.Fatal(err)
	}

	// init PG
	pgURI := os.Getenv(USER_PG_URI)
	userDriver, err := repoUser.NewUserPG(context.TODO(), pgURI)
	if err != nil {
		log.Fatal(err)
	}

	// init token
	redURI := os.Getenv(TOKEN_REDIS_URI)
	tokenDriver := repoToken.NewRedisToken(redURI)
	if err != nil {
		log.Fatal(err)
	}

	svc, err := service.NewCloudfs(&tokenDriver, blobDriver, userDriver)
	if err != nil {
		log.Fatal(err)
	}
	return svc
}
