package minioHandler

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"strings"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Object struct {
	Key   string `json:"key"`
	Size  int64  `json:"size"`
	Owner string `json:"owner"`
}

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
	Ctx    context.Context
}

func isAuthorizedOwner(ctx context.Context, objectInfo *minio.ObjectInfo, key string, subject string) (bool, string) {
	if objectInfo.Metadata.Get("X-Amz-Meta-Owner") == subject || strings.HasPrefix(key, fmt.Sprintf("%s/", subject)) {
		return true, ""
	}

	return false, "Not authorized"
}

func isAuthorizedGuest(ctx context.Context, key string, allowedKeys []interface{}) (bool, string) {
	for _, k := range allowedKeys {
		if key == k {
			return true, ""
		}
	}
	return false, "Not authorized"
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
		objectName := fmt.Sprintf("%s/%s-%s", msg.Owner, objectUUID.String(), msg.Header.Filename)
		log.Printf("Uploading %s\n", objectName)
		info, err := mh.client.PutObject(msg.Ctx, mh.bucketName, objectName, msg.File, msg.Header.Size, minio.PutObjectOptions{
			UserMetadata: map[string]string{"owner": msg.Owner},
			ContentType:  msg.Header.Header.Get("content-type"),
		})
		if err != nil {
			log.Printf("ERROR: %s: %v\n", objectName, err)
		}
		log.Printf("Done uploading %s\n", info.Key)
	}
}

func (mh *MinioHandler) MinioGetSharedFile(ctx context.Context, key string, allowedKeys []interface{}, owner string) (*minio.Object, *minio.ObjectInfo, error) {
	fmt.Sprintf("%v, %v, %v", ctx, key, allowedKeys)
	key = fmt.Sprintf("%s/%s", owner, key)
	if authorized, msg := isAuthorizedGuest(ctx, key, allowedKeys); !authorized {
		return nil, nil, fmt.Errorf("%s", msg)
	}
	object, err := mh.client.GetObject(ctx, mh.bucketName, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, nil, err
	}
	info, err := mh.client.GetObjectACL(ctx, mh.bucketName, key)
	if err != nil {
		return nil, nil, err
	}
	return object, info, nil
}

func (mh *MinioHandler) MinioGetFile(ctx context.Context, key string, owner string) (*minio.Object, *minio.ObjectInfo, error) {
	fmt.Sprintf("%v, %v, %v", ctx, key, owner)
	key = fmt.Sprintf("%s/%s", owner, key)
	info, err := mh.client.GetObjectACL(ctx, mh.bucketName, key)
	if err != nil {
		return nil, nil, err
	}
	if authorized, msg := isAuthorizedOwner(ctx, info, key, owner); !authorized {
		return nil, nil, fmt.Errorf("%s", msg)
	}

	object, err := mh.client.GetObject(ctx, mh.bucketName, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, nil, err
	}
	return object, info, nil
}

func (mh *MinioHandler) MinioListFiles(ctx context.Context, prefix string) []*Object {
	objects := mh.client.ListObjects(ctx, mh.bucketName, minio.ListObjectsOptions{
		WithMetadata: true,
		Prefix:       prefix,
		Recursive:    true,
	})
	objectList := []*Object{}
	for o := range objects {
		object := &Object{
			Key:   o.Key,
			Size:  o.Size,
			Owner: o.UserMetadata["X-Amz-Meta-Owner"],
		}
		objectList = append(objectList, object)
	}
	return objectList
}

func (mh *MinioHandler) MinioDeleteFile(key string, ctx context.Context, owner string) error {
	key = fmt.Sprintf("%s/%s", owner, key)
	info, err := mh.client.GetObjectACL(ctx, mh.bucketName, key)
	if err != nil {
		return err
	}
	if info.Metadata.Get("X-Amz-Meta-Owner") != owner || !strings.HasPrefix(key, fmt.Sprintf("%s/", owner)) {
		return fmt.Errorf("Not authorized")
	}
	err = mh.client.RemoveObject(ctx, mh.bucketName, key, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	return nil
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
