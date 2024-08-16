package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	apiKey := os.Getenv("NFT_STORAGE_API_KEY") // Your NFT.Storage API key
	fmt.Println("Printing the api key : ", apiKey)
	filePath := "images/linkdine.png" // Replace with your file path

	cid, err := uploadToNFTStorage(apiKey, filePath)
	if err != nil {
		log.Fatalf("Error uploading file: %v", err)
	}

	fmt.Printf("File uploaded successfully: %s\n", cid)
}

func uploadToNFTStorage(apiKey, filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Prepare the form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		return "", err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return "", err
	}
	writer.Close()

	// Create the request
	req, err := http.NewRequest("POST", "https://api.nft.storage/upload", body)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to upload file: %s", resp.Status)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Here, you should parse the response JSON to extract the CID.
	// For simplicity, we're returning the raw response.
	return string(responseBody), nil
}
