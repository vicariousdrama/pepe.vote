package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

const pathToImages = "/home/bob/Pictures/pepe-vote"
const uploadUrl = "https://nostr.build/api/v2/upload/files"

func main() {
	fmt.Println(pathToImages)
	files, err := os.ReadDir(pathToImages)
	if err != nil {
		log.Fatal(err)
	}
	// make an auth event?
	hostedImageUrls := []string{}
	// upload each file
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fmt.Printf("processing file: %s\n", file.Name())
		fullPath := filepath.Join(pathToImages, file.Name())
		file, err := os.Open(fullPath)
		if err != nil {
			fmt.Printf("failed to open file: %s", err.Error())
			continue
		}
		defer file.Close()

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		part, err := writer.CreateFormFile("file", fullPath)
		if err != nil {
			fmt.Printf("failed to create form file: %s", err.Error())
			continue
		}

		_, err = io.Copy(part, file)
		if err != nil {
			fmt.Printf("failed to copy file contents: %s", err.Error())
			continue
		}

		err = writer.Close()
		if err != nil {
			fmt.Printf("failed to close writer: %s", err.Error())
			continue
		}

		req, err := http.NewRequest("POST", uploadUrl, &body)
		if err != nil {
			fmt.Printf("failed to create request: %s", err.Error())
			continue
		}

		req.Header.Set("Content-Type", writer.FormDataContentType())

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("failed to make request: %s", err.Error())
			continue
		}
		defer resp.Body.Close()
		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			fmt.Printf("failed to decode response: %s", err.Error())
			continue
		}
		if status, ok := result["status"].(string); ok && status == "success" {
			if data, ok := result["data"].([]interface{}); ok {
				if firstEntry, ok := data[0].(map[string]interface{}); ok {
					if url, ok := firstEntry["url"].(string); ok {
						fmt.Println("Uploaded file URL:", url)
						// grab the 720p image url
						if uploadUrl, ok := firstEntry["responsive"].(map[string]interface{})["720p"]; ok {
							hostedImageUrls = append(hostedImageUrls, uploadUrl.(string))
							fmt.Printf("hosted url: %s\n", uploadUrl)
						} else {
							fmt.Printf("no 720p url\n")
						}

					}
				}
			}
		} else {
			fmt.Println("Upload failed")
		}
	}
	for _, url := range hostedImageUrls {
		fmt.Printf("hosted url: %s\n", url)
	}
	fmt.Printf("num images with url: %d\n", len(hostedImageUrls))
}
