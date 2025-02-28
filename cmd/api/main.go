package main

import (
	"context"
	"database/sql"
	"errors"
	"expvar"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	blobModule "github.com/odit-bit/cloudfs/internal/blob/module"
	userModule "github.com/odit-bit/cloudfs/internal/user/module"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	var rpcPort, httpPort int
	var minioAddr, minioAccess, minioSecret, postgreDns string
	var minioSecure bool
	flag.IntVar(&rpcPort, "rpc-port", 8181, "rpc-port")
	flag.IntVar(&httpPort, "http-port", 8282, "http-port")

	flag.StringVar(&minioAddr, "minio-addr", "localhost:9000", "minio address")
	flag.StringVar(&minioAccess, "minio-access", "minioAdmin", "minio access")
	flag.StringVar(&minioSecret, "minio-secret", "minioAdmin", "minio secret")
	flag.BoolVar(&minioSecure, "minio-secure", false, "minio secure connection")

	flag.StringVar(&postgreDns, "pg-addr", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable", "postgre dns")

	flag.Parse()

	// setup logger
	logger := logrus.New()
	eg := errgroup.Group{}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//setup outbound connection
	mio, err := setupMinio(minioAddr, minioAccess, minioSecret, minioSecure)
	if err != nil {
		logger.Error(err)
		return
	}
	pg, err := setupPostgre(ctx, postgreDns)
	if err != nil {
		logger.Error(err)
		return
	}
	defer pg.Close()

	// setup listener
	grpcAddr := fmt.Sprintf("%s:%d", "localhost", rpcPort)
	httpAddr := fmt.Sprintf("%s:%d", "localhost", httpPort)
	l := tcpListener(grpcAddr)
	l2 := tcpListener(httpAddr)
	defer l.Close()
	defer l2.Close()

	// // OBJECTS
	// objects := objectRepo.NewMinioBlob(mio)
	// objectToken, err := objectRepo.NewPGShareToken(ctx, pg)
	// if err != nil {
	// 	logger.Error(err)
	// 	return
	// }
	// o, _ := blob.New(ctx, objectToken, objects)

	// // USERS
	// users, err := userRepo.NewUserPG(ctx, pg)
	// if err != nil {
	// 	logger.Error(err)
	// 	return
	// }
	// userToken, err := userRepo.NewUserTokenPG(ctx, pg)
	// if err != nil {
	// 	logger.Error(err)
	// 	return
	// }
	// u, _ := user.New(ctx, users, userToken)
	// ss := rpc.NewGrpcServer(o, u)
	// apipb.RegisterStorageServiceServer(srv, ss)

	//grpc server
	srv := grpc.NewServer()
	if err := userModule.Start(ctx, logger, srv, pg); err != nil {
		logger.Error(err)
		return
	}
	if err := blobModule.Start(ctx, logger, mio, pg, srv); err != nil {
		logger.Error(err)
		return
	}
	reflection.Register(srv)
	// grpc serve
	eg.Go(func() error {
		if err := srv.Serve(l); err != nil {
			return err
		}
		return nil
	})

	///http server
	mux := chi.NewMux()

	// http handler
	mux.Handle("/v1/metrics", expvar.Handler())
	httpSrv := http.Server{Handler: mux}

	// http serve
	eg.Go(func() error {
		if err := httpSrv.Serve(l2); err != nil {
			return err
		}
		return nil
	})

	// handling signal
	// setup signal
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	eg.Go(func() error {
		<-sig
		cancel()
		srv.Stop()
		err := httpSrv.Close()
		return errors.Join(err, l.Close(), l2.Close())
	})

	logger.Infof("serve grpc on port %v", l.Addr().String())
	logger.Infof("serve http on port %v", l2.Addr().String())

	if err := eg.Wait(); err != nil {
		logger.Infof("exit:%v", err)
	}
}

func tcpListener(addr string) *net.TCPListener {
	lAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		panic(err)
	}
	l, err := net.ListenTCP("tcp", lAddr)
	if err != nil {
		panic(err)
	}
	return l
}

func setupMinio(endpoint, access, secret string, secure bool) (*minio.Client, error) {

	creds := credentials.NewStaticV4(access, secret, "")
	cli, err := minio.New(endpoint, &minio.Options{
		Creds:  creds,
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
	if cli.IsOnline() {
		return cli, nil
	} else {
		return nil, fmt.Errorf("minio server is offline")
	}
}

func setupPostgre(ctx context.Context, dns string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dns)
	if err != nil {
		return nil, err
	}
	ctx2, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if err := db.PingContext(ctx2); err != nil {
		return nil, err
	}
	return db, nil
}
