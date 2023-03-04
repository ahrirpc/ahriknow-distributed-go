package main

import (
	"context"
	"fmt"
	stlog "log"

	"github.com/ahrirpc/ahriknow-distributed-go/log"
	"github.com/ahrirpc/ahriknow-distributed-go/registry"
	"github.com/ahrirpc/ahriknow-distributed-go/service"
)

func main() {
	log.Run("distributed.log")
	host, port := "localhost", "4000"

	r := registry.Registration{
		ServiceName: "Log Service",
		ServiceURL:  fmt.Sprintf("http://%s:%s", host, port),
	}

	ctx, err := service.Start(context.Background(), host, port, r, log.RegisterHandlers)

	if err != nil {
		stlog.Fatalln(err)
	}
	<-ctx.Done()

	fmt.Println("Shutting down log service")
}
