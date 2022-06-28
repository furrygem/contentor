package minioHandler

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioHandler struct {
	client           *minio.Client
	workers          int
	channelCapacity  int
	bucketName       string
	fileMessagesChan chan FileMessage
}

type FileMessage struct {
	Header *multipart.FileHeader
	File   io.Reader
	Owner  string
}

func New(fileMessagesChan chan FileMessage, config *Config) *MinioHandler {
	minioClient := ConnectMinio(config.Endpoint,
		config.AccessKeyID,
		config.SecretKeyID,
		config.UseSSL,
	)
	return &MinioHandler{
		client:           minioClient,
		workers:          config.Workers,
		channelCapacity:  config.ChannelCapacity,
		bucketName:       config.BucketName,
		fileMessagesChan: fileMessagesChan,
	}
}

func (mh *MinioHandler) MinioUploadObject(messages <-chan FileMessage) {
	for msg := range messages {
		objectUUID, _ := uuid.NewRandom()
		objectName := fmt.Sprintf("%s-%s", objectUUID.String(), msg.Header.Filename)
		log.Printf("Uploading %s\n", objectName)
		info, err := mh.client.PutObject(context.Background(), mh.bucketName, objectName, msg.File, msg.Header.Size, minio.PutObjectOptions{
			UserMetadata: map[string]string{"owner": msg.Owner},
			ContentType:  msg.Header.Header.Get("content-type"),
		})
		if err != nil {
			log.Printf("ERROR: %s: %v\n", objectName, err)
		}
		log.Printf("Done uploading %s\n", info.Key)
	}
}

func (mh *MinioHandler) MinioGetFile(key string) (*minio.Object, minio.ObjectInfo, error) {
	object, err := mh.client.GetObject(context.Background(), mh.bucketName, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, minio.ObjectInfo{}, err
	}
	info, err := object.Stat()
	return object, info, nil
}

func (mh *MinioHandler) MinioListFiles() []*Object {
	objects := mh.client.ListObjects(context.Background(), mh.bucketName, minio.ListObjectsOptions{
		WithMetadata: true,
	})
	objectList := []*Object{}
	for o := range objects {
		object := &Object{
			Key:  o.Key,
			Size: o.Size,
		}
		objectList = append(objectList, object)
	}
	return objectList
}

func (mh *MinioHandler) InitGoroutines() {
	for i := 0; i < mh.workers; i++ {
		go mh.MinioUploadObject(mh.fileMessagesChan)
	}
}

func ConnectMinio(endpoint, accessKeyID, secretKeyID string, useSSL bool) *minio.Client {
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
}
