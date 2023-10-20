// MIT License
// Copyright (c) 2023 Zhang Xiaofeng
// Email: i@zxf.io

// S3 Parallel Downloader
//
// This program downloads an S3 object in parallel chunks.
// It uses the AWS SDK for Go V2.

package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// downloadPart downloads a part of the S3 object and writes it to a local file.
// It uses goroutines for parallel downloading.
func downloadPart(svc *s3.Client, bucket, key string, partNum int32, start, end int64, wg *sync.WaitGroup, localFilePath string) {
	defer wg.Done() // Decrement the counter when the goroutine completes.

	// Define the byte range to download.
	rangeStr := fmt.Sprintf("bytes=%d-%d", start, end)

	// Prepare the download request.
	input := &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Range:  &rangeStr,
	}

	// Execute the download request.
	resp, err := svc.GetObject(context.TODO(), input)
	if err != nil {
		fmt.Println("Error downloading part:", err)
		return
	}
	defer resp.Body.Close()

	// Open the local file for writing.
	outFile, err := os.OpenFile(localFilePath, os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer outFile.Close()

	// Seek to the correct position in the file.
	outFile.Seek(start, 0)

	// Write the downloaded part to the file.
	io.Copy(outFile, resp.Body)
}

func main() {
	// Load AWS SDK configuration.
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}

	// Initialize S3 client.
	client := s3.NewFromConfig(cfg)

	// Check if enough arguments are provided.
	if len(os.Args) < 3 {
		fmt.Println("Usage: <program> <s3uri> <localpath>")
		return
	}

	// Get S3 URI and local path from command-line arguments.
	s3URI := os.Args[1]
	localPath := os.Args[2]

	// Parse S3 URI to extract bucket and key.
	parts := strings.Split(s3URI, "/")
	bucket := parts[2]
	key := strings.Join(parts[3:], "/")

	// Get object metadata to find out the total size.
	headInput := &s3.HeadObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}
	headResp, err := client.HeadObject(context.TODO(), headInput)
	if err != nil {
		panic("failed to get object metadata: " + err.Error())
	}
	totalSize := headResp.ContentLength

	// Define the size of each part.
	partSize := int64(5 * 1024 * 1024) // 5MB

	// Prepare the local file for writing.
	fileName := filepath.Base(key)
	var localFilePath string

	// Check if localPath is a directory.
	fi, err := os.Stat(localPath)
	if err == nil && fi.IsDir() {
		// localPath is a directory; append the S3 object's filename to it.
		localFilePath = filepath.Join(localPath, fileName)
	} else if os.IsNotExist(err) {
		// localPath does not exist; use it as the file path.
		localFilePath = localPath
	} else {
		// An error occurred trying to get the file info.
		panic("Error determining local file path: " + err.Error())
	}

	outFile, err := os.Create(localFilePath)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	// Pre-allocate space for the file.
	err = outFile.Truncate(totalSize)
	if err != nil {
		panic("failed to truncate file: " + err.Error())
	}

	// Initialize wait group for goroutines.
	var wg sync.WaitGroup

	// Download each part in parallel.
	for i := int64(0); i < totalSize; i += partSize {
		wg.Add(1)
		end := i + partSize - 1
		if end > totalSize {
			end = totalSize - 1
		}
		go downloadPart(client, bucket, key, int32(i/partSize), i, end, &wg, localFilePath)
	}

	// Wait for all goroutines to complete.
	wg.Wait()
}
