package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/ahrirpc/ahriknow-distributed-go/registry"
)

func main() {
	http.Handle("/services", &registry.RegistryService{})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var srv http.Server
	srv.Addr = ":3000"

	go func() {
		log.Println(srv.ListenAndServe())
		cancel()
	}()

	go func() {
		var s string
		fmt.Println("Registry service started. Press any key to stop.")
		fmt.Scanln(&s)
		srv.Shutdown(ctx)
		cancel()
	}()

	<-ctx.Done()

	fmt.Println("Shutting down registry service")
}
