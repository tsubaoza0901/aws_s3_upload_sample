package main

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func main() {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String("ap-northeast-1"),
			Credentials: credentials.NewEnvCredentials(),
			// Endpoint:         aws.String("http://localstack:4566"), // ★localstack利用時は必要
			// S3ForcePathStyle: aws.Bool(true),                       // ★localstack利用時は必要
		},
	})
	if err != nil {
		panic(err)
	}

	// ファイルを開く
	targetFilePath := "./sample.txt"
	file, err := os.Open(targetFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	bucketName := "teamup-t"
	objectKey := "user-image/sample.txt"

	// Uploaderを作成し、ローカルファイルをアップロード
	uploader := s3manager.NewUploader(sess)
	ret, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   file,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(ret.Location)
	log.Println("done")
}
