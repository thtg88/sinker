package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/fsnotify/fsnotify"
	"github.com/joho/godotenv"
)

var watcher *fsnotify.Watcher
var sinkerApiDeviceId string

// main
func main() {
	err := godotenv.Load()
	if err != nil {
		panic(errors.New("could not load .env file"))
	}

	sinkerApiDeviceId, err = registerDevice()
	if err != nil {
		panic(err)
	}

	// TODO sync before watching

	// creates a new file watcher
	watcher, _ = fsnotify.NewWatcher()
	defer watcher.Close()

	go watchPeriodically(os.Getenv("SINKER_BASE_PATH"), 5)

	done := make(chan bool)

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				handleFsEvent(event)

			case err := <-watcher.Errors:
				fmt.Println("ERROR", err)
			}
		}
	}()

	<-done
}

// watchDir gets run as a walk func, searching for directories to add watchers to
func watchDir(path string, fi os.FileInfo, err error) error {
	// since fsnotify can watch all the files in a directory, watchers only need
	// to be added to each nested directory
	if fi.Mode().IsDir() {
		return watcher.Add(path)
	}

	return nil
}

// watchPeriodically adds sub directories peridically to watch, with the help
// of fsnotify which maintains a directory map rather than slice.
func watchPeriodically(directory string, interval int) {
	done := make(chan struct{})
	go func() {
		done <- struct{}{}
	}()

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for ; ; <-ticker.C {
		<-done

		if err := filepath.Walk(directory, watchDir); err != nil {
			fmt.Println(err)
		}

		go func() {
			done <- struct{}{}
		}()
	}
}

// handleFsEvent handles a file system event, uploading a file to S3,
// and updates the state backend
func handleFsEvent(event fsnotify.Event) {
	var err error

	// Skip CHMOD event as macOS sends 2 for every WRITE event (before and after)
	if event.Op.String() == "CHMOD" {
		return
	}

	fmt.Printf("EVENT! %#v\n", event.String())

	if event.Op.String() == "CREATE" {
		_, err = uploadFile(event.Name)
	}

	if event.Op.String() == "REMOVE" {
		_, err = removeFile(event.Name)
	}

	if event.Op.String() == "WRITE" {
		_, err = uploadFile(event.Name)
	}

	if err != nil {
		fmt.Println("ERROR", err, event.Name)
		return
	}

	_, err = updateState(event)
	if err != nil {
		fmt.Println("ERROR", err, event.Name)
	}

	fmt.Println("file updated", event.Name)
}

// updateState updates the state backend
func updateState(event fsnotify.Event) ([]byte, error) {
	values := map[string]string{
		"path": relativePath(event.Name),
		"type": event.Op.String(),
	}
	jsonValue, _ := json.Marshal(values)

	return sinkerApiRequest("POST", os.Getenv("SINKER_API_STORE_EVENT_PATH"), jsonValue)
}

// removeFile removes a file from a given absolute path from the S3 bucket
// specified by the AWS_BUCKET env variable
func removeFile(path string) (*s3.DeleteObjectOutput, error) {
	return s3Client().DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET")),
		Key:    aws.String(relativePath(path)),
	})
}

// uploadFile uploads a file from a given absolute path to the S3 bucket
// specified by the AWS_BUCKET env variable
func uploadFile(path string) (*s3.PutObjectOutput, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.New("could not open path")
	}

	return s3Client().PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET")),
		Key:    aws.String(relativePath(path)),
		Body:   file,
	})
}

// relativePath returns the relative path of a file from a given aboslute path string
func relativePath(path string) string {
	return strings.Trim(strings.Replace(path, os.Getenv("SINKER_BASE_PATH"), "", 1), "/")
}

// s3Client returns a new S3 client
func s3Client() *s3.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	return s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.Region = "eu-west-1"
	})
}

func registerDevice() (string, error) {
	type Device struct {
		Uuid string `json:"uuid"`
		Name string `json:"name"`
	}
	type DeviceResponse struct {
		Data Device `json:"data"`
	}
	var response DeviceResponse

	values := map[string]string{"name": "Sinker macOS client"}
	jsonValue, _ := json.Marshal(values)

	body, err := sinkerApiRequest("POST", os.Getenv("SINKER_API_STORE_DEVICE_PATH"), jsonValue)
	if err != nil {
		return "", err
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}
	if response.Data.Uuid == "" {
		return "", errors.New(fmt.Sprint("Could not register device. ERROR ", string(body)))
	}

	return string(response.Data.Uuid), nil
}

func sinkerApiRequest(method string, uri string, requestBody []byte) ([]byte, error) {
	url := fmt.Sprint(os.Getenv("SINKER_API_BASE_URL"), uri)
	fmt.Println("URL:>", url)

	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set(os.Getenv("SINKER_API_KEY_HEADER_NAME"), os.Getenv("SINKER_API_KEY_HEADER_VALUE"))
	req.Header.Set(os.Getenv("SINKER_API_USER_ID_HEADER_NAME"), os.Getenv("SINKER_API_USER_ID_HEADER_VALUE"))
	req.Header.Set(os.Getenv("SINKER_API_DEVICE_ID_HEADER_NAME"), sinkerApiDeviceId)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// fmt.Println("request Headers:", req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	return body, nil
}
