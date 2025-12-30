package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/fsnotify/fsnotify"
	"github.com/joho/godotenv"

	"github.com/thtg88/sinker/internal/config"
	"github.com/thtg88/sinker/internal/handlers"
	"github.com/thtg88/sinker/internal/uploaders"
	"github.com/thtg88/sinker/internal/watchers"
	"github.com/thtg88/sinker/pkg/sinker"
)

var watcher *fsnotify.Watcher
var sinkerAPIDeviceID string

func main() {
	if err := run(); err != nil {
		log.Printf("[ERROR] %v", err)
		os.Exit(1)
	}
}

func run() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("godotenv load: %v", err)
	}

	cfg := config.Load()
	httpClient := &http.Client{}
	sinkerAPIClient := sinker.NewAPIClient(httpClient, cfg.SinkerAPI)

	s3Config, err := awsconfig.LoadDefaultConfig(context.TODO())
	if err != nil {
		return fmt.Errorf("awsconfig loaddefaultconfig: %v", err)
	}

	s3Client := s3.NewFromConfig(s3Config, func(o *s3.Options) { o.Region = "eu-west-1" })
	fileUploader := uploaders.NewS3FileUploader(s3Client)
	handler := handlers.NewFSEventHandler(fileUploader, sinkerAPIClient)

	sinkerAPIDeviceID, err = sinkerAPIClient.RegisterDevice()
	if err != nil {
		return fmt.Errorf("sinker register device: %v",err)
	}

	// TODO sync before watching

	fsNotifyWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("fsnotify newwatcher: %v", err)
	}
	defer fsNotifyWatcher.Close()

	watcher := watchers.NewWatcher(fsNotifyWatcher)

	go watcher.WatchPeriodically(cfg.Sinker.BasePath, 5)

	done := make(chan bool)

	go func() {
		for {
			select {
			case event := <-fsNotifyWatcher.Events:
				handler.Handle(event, sinkerAPIDeviceID)

			case err := <-fsNotifyWatcher.Errors:
				fmt.Println("ERROR", err)
			}
		}
	}()

	<-done

	return nil
}
