package main

import (
	"errors"
	"gocgrpcarchitecture/internal/grpc/client"
	"log"
	"net/http"
	"os"
)

func main() {
	httpClientServer := client.NewHttpClientServer(":8080")
	if err := httpClientServer.Start(); err != nil && errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
		os.Exit(1)
	}
}
