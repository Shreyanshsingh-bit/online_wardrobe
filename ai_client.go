package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

// SendToVirtualTryOn sends the user selfie and the garment photo to the Python AI service.
func SendToVirtualTryOn(userImagePath, garmentImagePath string) ([]byte, error) {
	
	targetURL := "http://127.0.0.1:8000/process-tryon"

	//Create a buffer to hold our multipart form data
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	//Open and attach the User Selfie file
	userFile, err := os.Open(userImagePath)
	if err != nil {
		return nil, fmt.Errorf("failed opening user image: %w", err)
	}
	defer userFile.Close()

	userPart, err := bodyWriter.CreateFormFile("user_image", userFile.Name())
	if err != nil {
		return nil, fmt.Errorf("failed creating user form file: %w", err)
	}
	if _, err = io.Copy(userPart, userFile); err != nil {
		return nil, fmt.Errorf("failed copying user file bytes: %w", err)
	}

	//Open and attach the Garment/Clothing file
	garmentFile, err := os.Open(garmentImagePath)
	if err != nil {
		return nil, fmt.Errorf("failed opening garment image: %w", err)
	}
	defer garmentFile.Close()

	garmentPart, err := bodyWriter.CreateFormFile("garment_image", garmentFile.Name())
	if err != nil {
		return nil, fmt.Errorf("failed creating garment form file: %w", err)
	}
	if _, err = io.Copy(garmentPart, garmentFile); err != nil {
		return nil, fmt.Errorf("failed copying garment file bytes: %w", err)
	}

	// Close the writer to finalize the multipart boundaries
	bodyWriter.Close()

	// 5. Build and fire the HTTP request
	req, err := http.NewRequest("POST", targetURL, bodyBuf)
	if err != nil {
		return nil, fmt.Errorf("failed creating request: %w", err)
	}
	
	// Set the crucial Content-Type header so Python knows it is a multi-part file stream
	req.Header.Set("Content-Type", bodyWriter.FormDataContentType())

	// Initialize an HTTP client with a safe timeout (AI processing can take time)
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed sending request to AI microservice: %w", err)
	}
	defer resp.Body.Close()

	// 6. Read the result back from Python
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("AI microservice returned bad status: %d, body: %s", resp.StatusCode, string(respBytes))
	}

	return respBytes, nil
}