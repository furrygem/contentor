package server

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/furrygem/contentor/content-service/internal/minioHandler"
	"github.com/furrygem/contentor/content-service/pkg/webutils"
	"github.com/gorilla/mux"
)

func (s *Server) createObjectsHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Printf("ERROR: %s", err)
		webutils.WriteHTTPCode(w, http.StatusBadRequest)
		return
	}
	file, header, err := r.FormFile("file")
	owner := r.FormValue("owner")
	if err != nil {
		log.Printf("ERROR: %v", err)
		webutils.WriteHTTPCode(w, http.StatusBadRequest)
		return
	}
	defer file.Close()
	var buf bytes.Buffer
	written, err := io.Copy(&buf, file)
	if err != nil {
		log.Printf("ERROR: %v", err)
		webutils.WriteHTTPCode(w, http.StatusInternalServerError)
		return
	}
	if written != header.Size {
		log.Printf("ERROR: Couldn't write required amount of data. %d < %d", written, header.Size)
		webutils.WriteHTTPCode(w, http.StatusInternalServerError)
		return
	}
	msg := minioHandler.FileMessage{
		Header: header,
		File:   &buf,
		Owner:  owner,
	}
	file.Seek(0, 0)
	probe := make([]byte, 5*1024)

	file.Read(probe)
	log.Println(http.DetectContentType(probe))

	msg.Header.Header.Set("content-type", http.DetectContentType(probe))
	s.fileMessagesChan <- msg
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Object is being created\n"))
	return
}

func (s *Server) listObjectsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	objects := s.minioHandler.MinioListFiles(ctx)
	response, err := json.Marshal(objects)
	if err != nil {
		log.Printf("ERROR: %v", err)
		webutils.WriteHTTPCode(w, http.StatusInternalServerError)
		return
	}
	w.Write(response)
	return
}

func (s *Server) downloadObjectHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	objectName := mux.Vars(r)["id"]
	if objectName == "" {
		webutils.WriteHTTPCode(w, http.StatusBadRequest)
		return
	}
	file, _, err := s.minioHandler.MinioGetFile(ctx, objectName)
	if err != nil {
		log.Printf("ERROR: %v", err)
		webutils.WriteHTTPCode(w, http.StatusInternalServerError)
		return
	}
	if file == nil {
		webutils.WriteHTTPCode(w, http.StatusNotFound)
		return
	}
	io.Copy(w, file)
}

func (s *Server) deleteObjectHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	Name := mux.Vars(r)["id"]
	err := s.minioHandler.MinioDeleteFile(Name, ctx)
	if err != nil {
		log.Printf("ERROR: %v", err)
		webutils.WriteHTTPCodeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Couldn't delete file"})
		return
	}
	webutils.WriteHTTPCodeJSON(w, http.StatusNoContent, map[string]string{})
	return
}

// func (s *Server) objectsHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == http.MethodPost {
// 		s.createObjectsHandler(w, r)
// 		return
// 	}
// 	if r.Method == http.MethodGet {
// 		s.listObjectsHandler(w, r)
// 		return
// 	}
// 	webutils.WriteHTTPCode(w, http.StatusMethodNotAllowed)
// }

// func (s *Server) objectHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == http.MethodGet {
// 		s.downloadObjectHandler(w, r)
// 		return
// 	}
// }
