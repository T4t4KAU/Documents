package minio

import (
	"context"
	"douyin/pkg/constants"
	"fmt"
	"github.com/minio/minio-go/v7"
	"testing"
)

func TestBuketExist(t *testing.T) {
	Init()
	ctx := context.Background()
	exists, err := Client.BucketExists(ctx, constants.MinioVideoBucketName)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	if exists {
		fmt.Printf("%v found\n", constants.MinioVideoBucketName)
	} else {
		fmt.Printf("%v not found\n", constants.MinioVideoBucketName)
	}
}

func TestMakeBucket(t *testing.T) {
	Init()
	ctx := context.Background()
	exists, err := Client.BucketExists(ctx, constants.MinioVideoBucketName)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	if exists {
		fmt.Printf("%v found", constants.MinioVideoBucketName)
	} else {
		err = Client.MakeBucket(ctx, constants.MinioVideoBucketName, minio.MakeBucketOptions{})
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		fmt.Printf("Successfully created bucket %v\n", constants.MinioVideoBucketName)
	}
}
