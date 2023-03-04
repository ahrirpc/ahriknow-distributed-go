package service

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/ahrirpc/ahriknow-distributed-go/registry"
)

func Start(ctx context.Context, host, port string, reg registry.Registration, registerHandlersFunc func()) (context.Context, error) {
	registerHandlersFunc()
	ctx = startService(ctx, reg, host, port)
	err := registry.RegisterService(reg)
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}

func startService(ctx context.Context, reg registry.Registration, host, port string) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	var srv http.Server
	srv.Addr = ":" + port

	go func() {
		log.Println(srv.ListenAndServe())
		err := registry.ShutdownService(reg)
		fmt.Println("111=================")
		if err != nil {
			log.Println(err)
		}
		cancel()
	}()

	go func() {
		fmt.Printf("%v started. Press any key to stop.", reg.ServiceName)
		var s string
		fmt.Scanln(&s)
		err := registry.ShutdownService(reg)
		fmt.Println("222=================")
		if err != nil {
			log.Println(err)
		}
		srv.Shutdown(ctx)
		cancel()
	}()

	return ctx
}
