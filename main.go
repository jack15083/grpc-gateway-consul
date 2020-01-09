package main

import (
	"flag"
	"net/http"

	"./consul"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	gw "./proto/helloworld"
)

func run() error {
	r := consul.NewResolver("Helloworld")
	b := grpc.RoundRobin(r)

	conn, err := grpc.Dial("192.168.126.128:8500", grpc.WithInsecure(), grpc.WithBalancer(b), grpc.WithBlock())
	if err != nil {
		return err
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	errReg := gw.RegisterGreeterHandler(ctx, mux, conn)
	if errReg != nil {
		return errReg
	}

	return http.ListenAndServe(":8080", mux)
}

func main() {
	flag.Parse()
	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}

}
