package minio

import (
	"bytes"
	"context"
	"douyin/pkg/constants"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"mime/multipart"
	"net/url"
	"time"
)

var (
	Client *minio.Client
)

func MakeBucket(ctx context.Context, bucketName string) error {
	exists, err := Client.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}
	if !exists {
		err = Client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
	}
	fmt.Printf("Successfully created bucket %v\n", bucketName)
	return nil
}

func PutToBucket(ctx context.Context, bucketName string,
	file *multipart.FileHeader) (minio.UploadInfo, error) {
	fileObj, _ := file.Open()
	defer fileObj.Close()
	info, err := Client.PutObject(ctx, bucketName, file.Filename,
		fileObj, file.Size, minio.PutObjectOptions{})
	return info, err
}

func GetObjectURL(ctx context.Context, bucketName, filename string) (*url.URL, error) {
	expires := time.Hour * 24
	reqParams := make(url.Values)
	return Client.PresignedGetObject(ctx, bucketName, filename, expires, reqParams)
}

func PutToBucketByBuffer(ctx context.Context, bucketName,
	filename string, buf *bytes.Buffer) (minio.UploadInfo, error) {
	return Client.PutObject(ctx, bucketName, filename, buf, int64(buf.Len()), minio.PutObjectOptions{})
}

func Init() {
	ctx := context.Background()
	var err error

	Client, err = minio.New(constants.MinioEndPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(constants.MinioAccessKeyId, constants.MinioSecretAccessKey, ""),
		Secure: constants.MinioUseSSL,
	})
	if err != nil {
		log.Fatalln("Minio connect error: ", err)
	}

	err = MakeBucket(ctx, constants.MinioVideoBucketName)
	err = MakeBucket(ctx, constants.MinioImageBucketName)
	if err != nil {
		log.Println("make bucket error: ", err)
	}
}
