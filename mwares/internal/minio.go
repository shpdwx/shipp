package internal

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/shpdwx/mwares/conf"
)

var (
	newMC      *minio.Client
	bucketName string
	rootPath   string
)

func UploadMinio(ctx context.Context, cfg *conf.CfgMinio, content string) bool {

	var (
		accessKey = cfg.AccessKey
		secretKey = cfg.SecretKey
		endpoint  = cfg.Endpoint
	)

	bucketName = cfg.Bucket
	rootPath = cfg.RootPath

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})

	if err != nil {
		log.Fatalf("new minio client failed: %v", err)
	}

	found, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		log.Fatalf("get bucket failed: %v", err)
	}

	if !found {
		log.Fatal(bucketName + " is not found")
	}
	newMC = minioClient

	resp, err := TempFile(content, putFile)
	if err != nil {
		log.Fatalf("upload minio failed: %v", err)
	}

	v, ok := resp.(minio.UploadInfo)
	if !ok {
		log.Fatal("parse resp failed")
	}

	fmt.Println(v.ETag)
	fmt.Println(v.VersionID)

	return true
}

func putFile(tmpf string) interface{} {

	file, err := os.Open(tmpf)
	if err != nil {
		log.Fatalf("open file failed: %v", err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		log.Fatalf("get file info failed: %v", err)
	}

	opts := minio.PutObjectOptions{
		ContentType: "application/text",
	}

	filename := fmt.Sprintf("%s/%d.txt", rootPath, time.Now().Unix())

	upload, err := newMC.PutObject(context.Background(), bucketName, filename, file, stat.Size(), opts)
	if err != nil {
		log.Fatalf("upload failed: %v", err)
	}

	return upload
}
