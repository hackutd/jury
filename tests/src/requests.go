package src

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
)

type H map[string]any

// WaitForBackend will wait for the backend to load, checking the URL every 5 seconds
func WaitForBackend(logger *Logger) {
	logger.Log(Info, "Waiting for backend to load...\n")

	url := getBaseUrl()

	for {
		// Send a GET request to the backend
		res, err := http.Get(url)
		if err != nil {
			logger.Log(Info, "Error sending GET request to %s, waiting 5 seconds: %s\n", url, err.Error())
			time.Sleep(5 * time.Second)
			continue
		}

		if res.StatusCode == http.StatusOK {
			logger.LogLn(Info, "Backend loaded")
			break
		}
	}
}

func PostRequest(logger *Logger, url string, body H, authHeader string) string {
	logger.Log(Verbose, "Sending POST request to %s\n", url)

	jsonBody, err := json.Marshal(body)
	if err != nil {
		logger.Log(Error, "Error marshalling body of POST request to %s: %s\n", url, err.Error())
		return err.Error()
	}
	fullUrl := getBaseUrl() + url

	logger.Log(Verbose, "Request Body: %s\n", jsonBody)

	// Send the POST request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(jsonBody))
	if err != nil {
		logger.Log(Error, "Error creating POST request to %s: %s\n", url, err.Error())
		return err.Error()
	}
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Log(Error, "Error sending POST request to %s: %s\n", url, err.Error())
		return err.Error()
	}

	// Read the body
	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Log(Error, "Error reading response body: %s\n", err.Error())
		return err.Error()
	}

	logger.Log(Verbose, "Response: %s\n", string(resBody))

	return string(resBody)
}

func getBaseUrl() string {
	return GetEnv("API_URL")
}

func DefaultAuth() string {
	return ""
}

func JudgeAuth(token string) string {
	return "Bearer " + token
}

func AdminAuth() string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte("admin:"+GetEnv("ADMIN_PASSWORD")))
}

// Checks if body is okay
func IsOk(body string) bool {
	return strings.Contains(body, "\"ok\":1")
}
