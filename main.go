package main

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func NewSession() (*session.Session, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String("ap-northeast-1"),
			Credentials: credentials.NewStaticCredentials(
				"your accsess id", // AWS ACCESS KEY
				"your sercret id", // AWS SECRET KEY
				"token",           // Token ※Tokenを使用していない場合は空文字を設定
			),
		},
	})
	return sess, err
}

func main() {
	// sessionの作成
	sess, err := NewSession()
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
