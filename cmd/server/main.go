package main

import (
	"context"
	"errors"
	"expvar"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/odit-bit/cloudfs/internal/storage"
	"github.com/odit-bit/cloudfs/internal/user"
	"github.com/odit-bit/cloudfs/server"
	"github.com/odit-bit/cloudfs/server/apipb"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	var rpcPort, httpPort int
	flag.IntVar(&rpcPort, "rpc-port", 8181, "rpc-port")
	flag.IntVar(&httpPort, "http-port", 8282, "http-port")

	flag.Parse()

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	grpcAddr := fmt.Sprintf("%s:%d", "localhost", rpcPort)
	httpAddr := fmt.Sprintf("%s:%d", "localhost", httpPort)

	// setup logger
	logger := logrus.New()
	eg := errgroup.Group{}

	// setup listener
	l := tcpListener(grpcAddr)
	l2 := tcpListener(httpAddr)
	defer l.Close()
	defer l2.Close()

	// grpc module interceptor

	//grpc server
	o, _ := storage.NewWithMemory()
	u := user.NewWithMemory()
	ss := server.NewGrpcServer(o, u)
	srv := grpc.NewServer()
	reflection.Register(srv)
	apipb.RegisterStorageServiceServer(srv, ss)

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
	os.Exit(0)
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
