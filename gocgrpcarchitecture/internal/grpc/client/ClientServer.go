package client

import (
	"context"
	"gocgrpcarchitecture/internal/grpc/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/rs/cors"
	"google.golang.org/grpc"
)

type HttpClientServer struct {
	addr string
}

func NewHttpClientServer(addr string) *HttpClientServer {
	return &HttpClientServer{addr: addr}
}

func (c *HttpClientServer) Start() error {
	mux := http.NewServeMux()
	mux.Handle("/v2/", http.StripPrefix("/v2", mux))

	conn := NewClientConnection(":9090")
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatal("Failed to close connection.")
		}
	}(conn)
	// -------------------------------------------------------
	userResource := handlers.NewUserResourceHandler()
	userResource.UserResourceHandlerRoutes(mux, conn)
	// ------------------------------------------------------
	documentResource := handlers.NewDocumentResourceHandler()
	documentResource.DocumentResourceHandlerRoutes(mux, conn)
	// ------------------------------------------------------
	logger := CustomRouterLogger("RouterLoggerLogs.txt")
	server := &http.Server{
		Addr:         c.addr,
		Handler:      loadCors(LoggerMiddleware(logger, mux)),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	shutdownChain := make(chan os.Signal, 1)
	signal.Notify(shutdownChain, os.Interrupt)

	go func() {
		<-shutdownChain
		log.Println("HTTP [CLIENT] shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Fatal("HTTP [CLIENT] failed to shutdown: ", err)
		}
	}()

	log.Println("HTTP [CLIENT] running smoothly on port", c.addr)
	return server.ListenAndServe()
	//return server.ListenAndServeTLS("server.crt", "server.key")
}

func loadCors(r http.Handler) http.Handler {
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200", "http://localhost:3000"},
		AllowedHeaders:   []string{"Origin", "Access-Control-Allow-Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With", "X-CSRF-Token", "Access-Control-Request-Method", "Access-Control-Request-Headers", "Location", "Entity"},
		ExposedHeaders:   []string{"Origin", "Content-Type", "Accept", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Credentials", "true"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowCredentials: true,
	})
	return corsMiddleware.Handler(r)
}

func LoggerMiddleware(logger *log.Logger, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("Request received: Method %s, Path: %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	}
}

func CustomRouterLogger(filePath string) *log.Logger {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}
	routerLogger := log.New(file, "Router Logger: ", log.LstdFlags|log.Lshortfile)
	return routerLogger
}
