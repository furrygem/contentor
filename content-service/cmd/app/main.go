package main

import (
	"log"

	"github.com/furrygem/contentor/content-service/internal/authHandler"
	"github.com/furrygem/contentor/content-service/internal/minioHandler"
	"github.com/furrygem/contentor/content-service/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	fileMessagesChan := make(chan minioHandler.FileMessage, 1)

	minioHandlerConfig := minioHandler.NewConfig()
	// minioHandlerConfig.AccessKeyID = utils.GetEnvDefault("MINIO_ACCESS_KEY", "", true)
	// minioHandlerConfig.SecretKeyID = utils.GetEnvDefault("MINIO_SECRET_KEY", "", true)
	// minioHandlerConfig.Endpoint = utils.GetEnvDefault("MINIO_ENDPOINT", minioHandlerConfig.Endpoint, false)
	// minioHandlerConfig.Endpoint = utils.GetEnvDefault("MINIO_ENDPOINT", minioHandlerConfig.Endpoint, false)
	if err := minioHandlerConfig.FromENV(); err != nil {
		log.Fatal(err)
	}
	mh := minioHandler.New(fileMessagesChan, minioHandlerConfig)

	// mh := MinioHandler{
	// 	client:           *connectMinio(),
	// 	workers:          1,
	// 	channelCapacity:  1,
	// 	bucketName:       "files",
	// 	fileMessagesChan: fileMessagesChan,
	// }
	mh.InitGoroutines()

	ahConfig := authHandler.NewConfig()
	if err := ahConfig.FromEnv(); err != nil {
		log.Fatal(err)
	}

	ah, err := authHandler.New(ahConfig)
	if err != nil {
		log.Fatal(err)
	}
	// s := Server{
	// 	router:           mux.NewRouter(),
	// 	fileMessagesChan: fileMessagesChan,
	// 	minioHandler:     mh,
	// }
	// handler := authenticatetionMiddleware(s.router)
	// s.router.HandleFunc("/api/objects", s.objectsHandler)
	// s.router.HandleFunc("/api/objects/{id}", s.objectHandler)
	// server := http.Server{
	// 	Handler: handler,
	// 	Addr:    ":8000",
	// }
	// log.Fatal(server.ListenAndServe())
	sc := server.NewConfig()
	sc.FromENV()
	s := server.New(fileMessagesChan, mh, ah)
	s.Start(sc)
}

func init() {
	godotenv.Load()
}
