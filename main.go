package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"log"
	"net/http"
	"time"
)

func makeAWSSession(region string) *session.Session {
	// Initialize a session in us-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		log.Fatalf("Something failed while creating the session")
	}
	return sess
}

func getObjectPresignedURL(region string, bucket string, filename string) string {
	sess := makeAWSSession(region)
	// Create S3 service client
	svc := s3.New(sess)
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})
	urlStr, err := req.Presign(15 * time.Minute)
	if err != nil {
		log.Println("Failed to sign request", err)
	}
	log.Println(urlStr)

	return urlStr
}

func putObjectPresignedURL(region string, bucket string, key string) string {
	sess := makeAWSSession(region)
	// Create S3 service client
	svc := s3.New(sess)
	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	urlStr, err := req.Presign(15 * time.Minute)
	if err != nil {
		log.Println("Failed to sign request", err)
	}
	log.Println(urlStr)

	return urlStr
}

func makeAPutRequest(url string, data io.Reader) *http.Response {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, url, data)
	if err != nil {
		// handle error
		log.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		// handle error
		log.Fatal(err)
	}

	return resp
}

type SomeStruct struct {
	Info1 string `json:"info1"`
	Info2 string `json:"info2"`
}

func main() {
	url := putObjectPresignedURL("us-east-2", "liveon-sangam-partiql", "xyz.json")
	object := &SomeStruct{Info1: "value1", Info2: "value2"}
	jsonifiedObject, err := json.Marshal(object)
	if err != nil {
		log.Fatalf("Something went wrong when jsonifying the object")
	}
	bufferWriter := new(bytes.Buffer)
	log.Println(object)
	log.Println(string(jsonifiedObject))
	bufferWriter.Write(jsonifiedObject)
	log.Println(bufferWriter)
	fmt.Println(makeAPutRequest(url, bufferWriter))
}
