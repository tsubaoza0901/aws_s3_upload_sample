package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// --------
// Others↓
// --------

func NewSession() (*session.Session, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String("ap-northeast-1"),
			Credentials: credentials.NewStaticCredentials(
				os.Getenv("AWS_ACCESS_KEY_ID"),     // AWS ACCESS KEY
				os.Getenv("AWS_SECRET_ACCESS_KEY"), // AWS SECRET KEY
				"",
			),
		},
	})
	return sess, err
}

func Decode(targetData string) ([]byte, error) {
	idx := strings.Index(targetData, ",")
	dec, err := base64.StdEncoding.DecodeString(targetData[idx+1:])
	if err != nil {
		return nil, err
	}
	return dec, nil
}

// --------
// model↓
// --------

// Image ...
type Image struct {
	ID              uint   `json:"id"`
	EncodedURL      string `json:"encoded_url"`
	FileName        string `json:"file_name"`
	StorageLocation string `json:"storage_location"`
}

func (i *Image) Upload() (string, error) {
	// AWSアクセス用のsession作成
	sess, err := NewSession()
	if err != nil {
		panic(err)
	}

	decodedImageURL, err := Decode(i.EncodedURL)
	if err != nil {
		log.Println("Decode Image failed")
		return "", err
	}

	log.Printf("BUCKET_NAME:%v", os.Getenv("BUCKET_NAME"))

	uploader := s3manager.NewUploader(sess)
	ret, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(os.Getenv("BUCKET_NAME")),
		Key:    aws.String(fmt.Sprintf("%v/%v", "test", i.FileName)),
		Body:   ioutil.NopCloser(bytes.NewReader(decodedImageURL)),
	})
	if err != nil {
		log.Println("Upload Image failed")
		return "", err
	}

	log.Println("Success Upload Image ")

	return ret.Location, err
}

// --------
// RDB model↓
// --------

// ImageLocation ...
type ImageLocation struct {
	gorm.Model
	StorageLocation string `gorm:"storage_location"`
}

// --------
// router↓
// --------

// InitRouting ...
func InitRouting(e *echo.Echo, u *Image) {
	e.POST("image", u.CreateImage)
	e.GET("image/:id", u.GetImage)
}

// CreateImage ...
func (u *Image) CreateImage(c echo.Context) error {
	image := &Image{}

	if err := c.Bind(image); err != nil {
		return err
	}

	loc, err := image.Upload()
	if err != nil {
		return err
	}

	imageLocation := ImageLocation{
		StorageLocation: loc,
	}
	err = db.Debug().Create(&imageLocation).Error
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, imageLocation)
}

// GetImage ...
func (u *Image) GetImage(c echo.Context) error {
	id := c.Param("id")

	imageLocation := ImageLocation{}
	err := db.Debug().Where("id = ?", id).First(&imageLocation).Error
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, imageLocation)
}

// --------
// infrastructure↓
// --------

var db *gorm.DB

// InitDB ...
func InitDB() *gorm.DB {
	dsn := "root:root@tcp(db:3306)/awss3uploadsample?parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	return db
}

// --------
// main.go↓
// --------

func main() {
	db = InitDB()

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	e := echo.New()

	u := new(Image)
	InitRouting(e, u)

	e.Start(":9111")
}
