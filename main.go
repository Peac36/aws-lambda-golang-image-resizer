package main

import (
	"bytes"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/disintegration/imaging"
)

var s3Service *s3.S3
var sizes []ImageSize

type ImageSize struct {
	OutputDirectory string
	SizeWidth       int32
	SizeHeight      int32
}

func getS3Service() *s3.S3 {
	return s3.New(session.New())
}

func handler(event events.SQSEvent) {

	outputBucket := os.Getenv("OUTPUT_BUCKET")

	var sizes []ImageSize = []ImageSize{
		{OutputDirectory: "/thumbnail/", SizeWidth: 600, SizeHeight: 300},
		{OutputDirectory: "/medium/", SizeWidth: 1920, SizeHeight: 1020},
	}

	for _, record := range event.Records {
		var result events.S3Event
		log.Println(record.Body)
		json.Unmarshal([]byte(record.Body), &result)

		for _, file := range result.Records {
			for _, size := range sizes {
				ResizeImage(size, file, outputBucket)
			}
		}
	}
}

func ResizeImage(size ImageSize, file events.S3EventRecord, outputBucket string) {
	var key string = file.S3.Object.Key
	var bucket string = file.S3.Bucket.Name
	var outputExtension = ".png"
	var outputContentType = "image/png"

	log.Printf("Input key %s \n", key)
	log.Printf("Input Bucket %s \n", bucket)
	log.Printf("Size output directory %s \n", size.OutputDirectory)
	log.Printf("Size width %d \n", size.SizeWidth)
	log.Printf("Size height %d \n", size.SizeHeight)

	input := s3.GetObjectInput{
		Key:    &key,
		Bucket: &bucket,
	}
	fileContent, err := s3Service.GetObject(&input)
	if err != nil {
		log.Printf("File can not be found: %s \n", err)
		return
	}
	defer fileContent.Body.Close()

	newFileName := file.S3.Object.ETag + outputExtension
	localFileName := "/tmp/" + newFileName
	removeFileName := size.OutputDirectory + newFileName

	log.Printf("Local file:  %s \n", localFileName)
	log.Printf("Remote File %s \n", removeFileName)

	imageSrc, err := imaging.Decode(fileContent.Body)
	if err != nil {
		log.Printf("The image can not be decoded %s \n", err)
	}
	image := imaging.Fit(imageSrc, int(size.SizeWidth), int(size.SizeHeight), imaging.Gaussian)
	imaging.Save(image, localFileName)

	defer os.Remove(localFileName)

	fileBody, err := os.ReadFile(localFileName)
	if err != nil {
		log.Printf("Can not open resized file: %s \n", err)
	}

	initialPermission := s3.ObjectCannedACLPublicRead
	_, err = s3Service.PutObject(&s3.PutObjectInput{
		Body:        bytes.NewReader(fileBody),
		Key:         &removeFileName,
		Bucket:      &outputBucket,
		ACL:         &initialPermission,
		ContentType: &outputContentType,
	})

	if err != nil {
		log.Printf("Error uploading the file %s", err)
		return
	}
	log.Println("File was successfully uploaded.")
}

func init() {
	s3Service = getS3Service()
}

func main() {
	lambda.Start(handler)
}
