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
	"github.com/minio/minio-go/v7"
)

func (s *Server) createObjectsHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Printf("ERROR: %s", err)
		webutils.WriteHTTPCode(w, http.StatusBadRequest)
		return
	}
	file, header, err := r.FormFile("file")
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
	user := r.Context().Value(keyTokenSubject).(string)
	if user == "" {
		log.Printf("Can't detect any user in request context. Is auth middleware working?")
		webutils.WriteHTTPCodeJSON(w, http.StatusInternalServerError, map[string]string{})
		return
	}
	msg := minioHandler.FileMessage{
		Header: header,
		File:   &buf,
		Owner:  user,
		Ctx:    context.Background(),
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
	if r.Context().Value(keyTokenGuest) != nil {
		webutils.WriteHTTPCodeJSON(w, http.StatusForbidden, map[string]string{})
		return
	}
	user := r.Context().Value(keyTokenSubject)
	if user == nil {
		log.Printf("Can't detect any user in request context. Is auth middleware working?")
		webutils.WriteHTTPCodeJSON(w, http.StatusInternalServerError, map[string]string{})
		return
	}
	objects := s.minioHandler.MinioListFiles(ctx, user.(string))
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
	guest := r.Context().Value(keyTokenGuest).(bool)
	ctx := context.Background()
	objectName := mux.Vars(r)["id"]
	var file *minio.Object
	var err error
	if guest {
		allowedKeys := r.Context().Value(keyTokenAllowedKeys).([]interface{})
		owner := r.Context().Value(keyTokenOwner).(string)
		if len(allowedKeys) == 0 || owner == "" {
			log.Printf("Bad token")
			webutils.WriteHTTPCodeJSON(w, http.StatusBadRequest, map[string]string{"authorization": "bad token"})
		}
		file, _, err = s.minioHandler.MinioGetSharedFile(ctx, objectName, allowedKeys, owner)
	} else {
		user := r.Context().Value(keyTokenSubject).(string)
		if user == "" {
			log.Printf("Can't detect any user in request context. Is auth middleware working?")
			webutils.WriteHTTPCodeJSON(w, http.StatusInternalServerError, map[string]string{})
			return
		}
		if objectName == "" {
			webutils.WriteHTTPCode(w, http.StatusBadRequest)
			return
		}
		file, _, err = s.minioHandler.MinioGetFile(ctx, objectName, user)
	}
	if err != nil {
		if err.Error() == "Not authorized" {
			webutils.WriteHTTPCode(w, http.StatusForbidden)
			return
		}
		if err.Error() == "The specified key does not exist." {
			webutils.WriteHTTPCode(w, http.StatusNotFound)
			return
		}
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
	user := r.Context().Value(keyTokenSubject).(string)
	if user == "" {
		log.Printf("Can't detect any user in request context. Is auth middleware working?")
		webutils.WriteHTTPCodeJSON(w, http.StatusInternalServerError, map[string]string{})
		return
	}
	ctx := context.Background()
	Name := mux.Vars(r)["id"]
	err := s.minioHandler.MinioDeleteFile(Name, ctx, user)
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
