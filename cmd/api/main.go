package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	storagemodule "github.com/odit-bit/cloudfs/internal/blob/storagemod"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8181, "port")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	address := fmt.Sprintf("%s:%d", "localhost", port)

	// setup listener
	lAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		panic(err)
	}
	l, err := net.ListenTCP("tcp", lAddr)
	if err != nil {
		panic(err)
	}
	defer l.Close()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	srv := grpc.NewServer()
	reflection.Register(srv)
	storagemodule.Run(ctx, nil, srv)

	// blobs, _ := blob.NewWithMemory()

	// udb, _ := userRepo.NewInMemory()
	// accounts, _ := user.NewStore(ctx, udb, udb)

	logger := logrus.New()
	// app := api.New(accounts, blobs, logger)

	// eg.Go(func() error {
	// 	err := http.Serve(l, app.Route())
	// 	if err != nil {
	// 		logger.Error(err)
	// 	}
	// 	return err
	// })

	eg := errgroup.Group{}
	eg.Go(func() error {
		if err := srv.Serve(l); err != nil {
			return err
		}
		return nil
	})

	eg.Go(func() error {
		<-sig
		srv.Stop()
		return l.Close()
	})

	logger.Infof("serve grpc on port %v", l.Addr().String())
	eg.Wait()
}
