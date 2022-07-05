package server

import (
	"log"
	"net/http"

	"github.com/furrygem/contentor/content-service/internal/authHandler"
	"github.com/furrygem/contentor/content-service/internal/minioHandler"
	"github.com/gorilla/mux"
)

type Server struct {
	router           *mux.Router
	fileMessagesChan chan minioHandler.FileMessage
	minioHandler     *minioHandler.MinioHandler
	authHandler      *authHandler.AuthHandler
}

func New(fileMessagesChan chan minioHandler.FileMessage,
	minioHandler *minioHandler.MinioHandler,
	authHandler *authHandler.AuthHandler) *Server {

	return &Server{
		router:           mux.NewRouter(),
		fileMessagesChan: fileMessagesChan,
		minioHandler:     minioHandler,
		authHandler:      authHandler,
	}
}

func (s *Server) Start(c *Config) {
	handler := s.AuthenticationMiddleware(s.router)
	s.router.HandleFunc("/api/objects", s.listObjectsHandler).Methods("GET")
	s.router.HandleFunc("/api/objects", s.createObjectsHandler).Methods("POST")
	s.router.HandleFunc("/api/objects/{id}", s.downloadObjectHandler).Methods("GET")
	s.router.HandleFunc("/api/objects/{id}", s.deleteObjectHandler).Methods("DELETE")
	server := http.Server{
		Handler: handler,
		Addr:    c.ListenAddr,
	}
	log.Fatal(server.ListenAndServe())
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s: %s %s", r.RemoteAddr, r.Method, r.RequestURI)
	s.router.ServeHTTP(w, r)
}
