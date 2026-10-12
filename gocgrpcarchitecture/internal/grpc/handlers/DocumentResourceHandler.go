package handlers

import (
	"gocgrpcarchitecture/internal/app/resources/user/resource"
	"gocgrpcarchitecture/internal/app/resources/user/service"
	"net/http"
	"time"

	"google.golang.org/grpc"
)

type DocumentResourceHandler struct {
	DocumentService        service.DocumentService
	DocumentServerResource resource.DocumentServerResource
	conn                   *grpc.ClientConn
}

func NewDocumentResourceHandler() *DocumentResourceHandler {
	return &DocumentResourceHandler{}
}

func (app *DocumentResourceHandler) DocumentResourceHandlerRoutes(mux *http.ServeMux, conn *grpc.ClientConn) {
	time.Sleep(2 * time.Second)
	mux.HandleFunc("POST /create", app.CreateDocumentHandler(conn))
	mux.HandleFunc("GET /docs/{uuid}", app.GetDocumentHandler(conn))
	return
}

func (app *DocumentResourceHandler) CreateDocumentHandler(conn *grpc.ClientConn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		return
	}
}

func (app *DocumentResourceHandler) GetDocumentHandler(conn *grpc.ClientConn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		return
	}
}
