package server

import (
	"log"
	"net/http"

	"github.com/furrygem/contentor/content-service/internal/minioHandler"
	"github.com/gorilla/mux"
)

type Server struct {
	router           *mux.Router
	fileMessagesChan chan minioHandler.FileMessage
	minioHandler     *minioHandler.MinioHandler
}

func New(fileMessagesChan chan minioHandler.FileMessage, minioHandler *minioHandler.MinioHandler) *Server {
	return &Server{
		router:           mux.NewRouter(),
		fileMessagesChan: fileMessagesChan,
		minioHandler:     minioHandler,
	}
}

func (s *Server) Start() {
	handler := AuthenticationMiddleware(s.router)
	s.router.HandleFunc("/api/objects", s.objectsHandler)
	s.router.HandleFunc("/api/objects/{id}", s.objectHandler)
	server := http.Server{
		Handler: handler,
		Addr:    ":8000",
	}
	log.Fatal(server.ListenAndServe())
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s: %s %s", r.RemoteAddr, r.Method, r.RequestURI)
	s.router.ServeHTTP(w, r)
}

// func (s *Server) deleteObjectHandler(w http.ResponseWriter, r *http.Request) {
// 	object Name := mux.Vars(r)
// }
