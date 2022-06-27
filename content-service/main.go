package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Server struct {
	router           *mux.Router
	fileMessagesChan chan FileMessage
	minioHandler     MinioHandler
}

type MinioHandler struct {
	client           minio.Client
	workers          int
	channelCapacity  int
	bucketName       string
	fileMessagesChan chan FileMessage
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s: %s %s", r.RemoteAddr, r.Method, r.RequestURI)
	s.router.ServeHTTP(w, r)
}

func writeHTTPCode(w http.ResponseWriter, err int) {
	w.WriteHeader(err)
	w.Write([]byte(http.StatusText(err)))
	return
}

func writeHTTPCodeJSON(w http.ResponseWriter, err int, messages map[string]string) {
	w.WriteHeader(err)
	if messages == nil {
		messages = make(map[string]string)
	}
	messages["error"] = http.StatusText(err)
	result, _ := json.Marshal(messages)
	w.Write(result)
}

func getenv(key, fallback string, required bool) string {
	value := os.Getenv(key)
	if required && value == "" {
		log.Fatalf("%s environment variable must be set\n", key)
		os.Exit(1)
	}
	if value == "" {
		return fallback
	}
	return value
}

func (mh *MinioHandler) minioUploadObject(messages <-chan FileMessage) {
	for msg := range messages {
		objectUUID, _ := uuid.NewRandom()
		objectName := fmt.Sprintf("%s-%s", objectUUID.String(), msg.header.Filename)
		log.Printf("Uploading %s\n", objectName)
		info, err := mh.client.PutObject(context.Background(), mh.bucketName, objectName, msg.file, msg.header.Size, minio.PutObjectOptions{
			UserMetadata: map[string]string{"owner": msg.owner},
			ContentType:  msg.header.Header.Get("content-type"),
		})
		if err != nil {
			log.Printf("ERROR: %s: %v\n", objectName, err)
		}
		log.Printf("Done uploading %s\n", info.Key)
	}
}

func (mh *MinioHandler) minioGetFile(key string) (*minio.Object, minio.ObjectInfo, error) {
	object, err := mh.client.GetObject(context.Background(), mh.bucketName, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, minio.ObjectInfo{}, err
	}
	info, err := object.Stat()
	return object, info, nil
}

func (mh *MinioHandler) minioListFiles() []*Object {
	objects := mh.client.ListObjects(context.Background(), mh.bucketName, minio.ListObjectsOptions{
		WithMetadata: true,
	})
	objectList := []*Object{}
	for o := range objects {
		object := &Object{
			Key:  o.Key,
			Size: o.Size,
			Type: o.Metadata.Get("Content-Type"),
		}
		objectList = append(objectList, object)
	}
	return objectList
}

func (mh *MinioHandler) initGoroutines() {
	for i := 0; i < mh.workers; i++ {
		go mh.minioUploadObject(mh.fileMessagesChan)
	}
}

func connectMinio() *minio.Client {
	endpoint := getenv("MINIO_ENDPOINT", "127.0.0.1:9000", false)
	accessKeyID := getenv("MINIO_ACCESS_KEY", "", true)
	secretKeyID := getenv("MINIO_SECRET_KEY", "", true)
	useSSL := getenv("MINIO_USE_SSL", "false", false) == "true"
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretKeyID, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatal(err)
	}
	return client
}

type Object struct {
	Key  string `json:"key"`
	Size int64  `json:"size"`
	Type string `json:"type"`
}

type FileMessage struct {
	header *multipart.FileHeader
	file   io.Reader
	owner  string
}

func (s *Server) createObjectsHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Printf("ERROR: %s", err)
		writeHTTPCode(w, http.StatusBadRequest)
		return
	}
	file, header, err := r.FormFile("file")
	owner := r.FormValue("owner")
	if err != nil {
		log.Printf("ERROR: %v", err)
		writeHTTPCode(w, http.StatusBadRequest)
		return
	}
	defer file.Close()
	var buf bytes.Buffer
	written, err := io.Copy(&buf, file)
	if err != nil {
		log.Printf("ERROR: %v", err)
		writeHTTPCode(w, http.StatusInternalServerError)
		return
	}
	if written != header.Size {
		log.Printf("ERROR: Couldn't write required amount of data. %d < %d", written, header.Size)
		writeHTTPCode(w, http.StatusInternalServerError)
		return
	}
	msg := FileMessage{
		header: header,
		file:   &buf,
		owner:  owner,
	}
	file.Seek(0, 0)
	probe := make([]byte, 5*1024)

	file.Read(probe)
	log.Println(http.DetectContentType(probe))

	msg.header.Header.Set("content-type", http.DetectContentType(probe))
	s.fileMessagesChan <- msg
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Object is being created\n"))
	return
}

func (s *Server) listObjectsHandler(w http.ResponseWriter, r *http.Request) {
	objects := s.minioHandler.minioListFiles()
	response, err := json.Marshal(objects)
	if err != nil {
		log.Printf("ERROR: %v", err)
		writeHTTPCode(w, http.StatusInternalServerError)
		return
	}
	w.Write(response)
	return
}

func (s *Server) objectsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		s.createObjectsHandler(w, r)
		return
	}
	if r.Method == http.MethodGet {
		s.listObjectsHandler(w, r)
		return
	}
	writeHTTPCode(w, http.StatusMethodNotAllowed)
}

func (s *Server) downloadObjectHandler(w http.ResponseWriter, r *http.Request) {
	objectName := mux.Vars(r)["id"]
	if objectName == "" {
		writeHTTPCode(w, http.StatusBadRequest)
		return
	}
	file, _, err := s.minioHandler.minioGetFile(objectName)
	if err != nil {
		log.Printf("ERROR: %v", err)
		writeHTTPCode(w, http.StatusInternalServerError)
		return
	}
	if file == nil {
		writeHTTPCode(w, http.StatusNotFound)
		return
	}
	io.Copy(w, file)
}

// func (s *Server) deleteObjectHandler(w http.ResponseWriter, r *http.Request) {
// 	object Name := mux.Vars(r)
// }

func (s *Server) objectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		s.downloadObjectHandler(w, r)
		return
	}
}

func authenticatetionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			writeHTTPCode(w, http.StatusUnauthorized)
			return
		}
		if auth == "creds" {
			next.ServeHTTP(w, r)
			return
		}
		writeHTTPCode(w, http.StatusForbidden)
		return
	})
}

func main() {
	godotenv.Load()
	fileMessagesChan := make(chan FileMessage, 1)
	mh := MinioHandler{
		client:           *connectMinio(),
		workers:          1,
		channelCapacity:  1,
		bucketName:       "files",
		fileMessagesChan: fileMessagesChan,
	}
	mh.initGoroutines()
	s := Server{
		router:           mux.NewRouter(),
		fileMessagesChan: fileMessagesChan,
		minioHandler:     mh,
	}
	handler := authenticatetionMiddleware(s.router)
	s.router.HandleFunc("/api/objects", s.objectsHandler)
	s.router.HandleFunc("/api/objects/{id}", s.objectHandler)
	server := http.Server{
		Handler: handler,
		Addr:    ":8000",
	}
	log.Fatal(server.ListenAndServe())
}
