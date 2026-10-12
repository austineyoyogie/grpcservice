package main

import (
	"errors"
	"gocgrpcarchitecture/internal/app/resources/user/models"
	"gocgrpcarchitecture/internal/app/server"
	"gocgrpcarchitecture/internal/pkg/configs/database"
	"log"
	"net/http"
	"os"
)

var (
	DC = database.CDriver()
)

func main() {
	models.AutoMigration(DC)
	S := server.NewHttpServer(":9090")
	if err := S.Start(); err != nil && errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
		os.Exit(1)
	}
}
