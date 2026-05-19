package util

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/valyala/fastjson"
)

type H map[string]any

// WaitForBackend will wait for the backend to load, checking the URL every 5 seconds
func WaitForBackend(logger *Logger) error {
	logger.Log(Info, "Waiting for backend to load...\n")

	url := getBaseUrl()
	retry := 3 // Retry 3 times before giving up

	for {
		res, err := http.Get(url)
		if err != nil {
			if retry < 0 {
				logger.Log(Error, "Error: failed to connect to backend. Aborting.")
				return errors.New("failed to connect to backend")
			}

			logger.Log(Info, "Error sending GET request to %s, waiting 5 seconds: %s\n", url, err.Error())
			time.Sleep(5 * time.Second)
			retry--
			continue
		}

		if res.StatusCode == http.StatusOK {
			logger.LogLn(Info, "Backend loaded")
			break
		}
	}
	return nil
}

func GetRequest(logger *Logger, url string, authHeader string) string {
	_, body := GetRequestWithStatus(logger, url, authHeader)
	return body
}

// GetRequestWithStatus sends a GET request and returns both the status code and body
func GetRequestWithStatus(logger *Logger, url string, authHeader string) (int, string) {
	logger.Log(Verbose, "Sending GET request to %s\n", url)

	fullUrl := getBaseUrl() + url

	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		logger.Log(Error, "Error creating GET request to %s: %s\n", url, err.Error())
		return 0, err.Error()
	}
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Log(Error, "Error sending GET request to %s: %s\n", url, err.Error())
		return 0, err.Error()
	}

	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Log(Error, "Error reading response body: %s\n", err.Error())
		return res.StatusCode, err.Error()
	}

	logger.Log(Verbose, "Response (%d): %s\n", res.StatusCode, string(resBody))

	return res.StatusCode, string(resBody)
}

func PostRequest(logger *Logger, url string, body H, authHeader string) string {
	_, resBody := PostRequestWithStatus(logger, url, body, authHeader)
	return resBody
}

// PostRequestWithStatus sends a POST request and returns both the status code and body
func PostRequestWithStatus(logger *Logger, url string, body H, authHeader string) (int, string) {
	logger.Log(Verbose, "Sending POST request to %s\n", url)

	jsonBody, err := json.Marshal(body)
	if err != nil {
		logger.Log(Error, "Error marshalling body of POST request to %s: %s\n", url, err.Error())
		return 0, err.Error()
	}
	fullUrl := getBaseUrl() + url

	logger.Log(Verbose, "Request Body: %s\n", jsonBody)

	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(jsonBody))
	if err != nil {
		logger.Log(Error, "Error creating POST request to %s: %s\n", url, err.Error())
		return 0, err.Error()
	}
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Log(Error, "Error sending POST request to %s: %s\n", url, err.Error())
		return 0, err.Error()
	}

	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Log(Error, "Error reading response body: %s\n", err.Error())
		return res.StatusCode, err.Error()
	}

	logger.Log(Verbose, "Response (%d): %s\n", res.StatusCode, string(resBody))

	return res.StatusCode, string(resBody)
}

// PutRequest sends a PUT request and returns the response body
func PutRequest(logger *Logger, url string, body H, authHeader string) string {
	_, resBody := PutRequestWithStatus(logger, url, body, authHeader)
	return resBody
}

// PutRequestWithStatus sends a PUT request and returns both the status code and body
func PutRequestWithStatus(logger *Logger, url string, body H, authHeader string) (int, string) {
	logger.Log(Verbose, "Sending PUT request to %s\n", url)

	jsonBody, err := json.Marshal(body)
	if err != nil {
		logger.Log(Error, "Error marshalling body of PUT request to %s: %s\n", url, err.Error())
		return 0, err.Error()
	}
	fullUrl := getBaseUrl() + url

	logger.Log(Verbose, "Request Body: %s\n", jsonBody)

	req, err := http.NewRequest("PUT", fullUrl, bytes.NewBuffer(jsonBody))
	if err != nil {
		logger.Log(Error, "Error creating PUT request to %s: %s\n", url, err.Error())
		return 0, err.Error()
	}
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Log(Error, "Error sending PUT request to %s: %s\n", url, err.Error())
		return 0, err.Error()
	}

	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Log(Error, "Error reading response body: %s\n", err.Error())
		return res.StatusCode, err.Error()
	}

	logger.Log(Verbose, "Response (%d): %s\n", res.StatusCode, string(resBody))

	return res.StatusCode, string(resBody)
}

// DeleteRequest sends a DELETE request and returns the response body
func DeleteRequest(logger *Logger, url string, authHeader string) string {
	_, resBody := DeleteRequestWithStatus(logger, url, authHeader)
	return resBody
}

// DeleteRequestWithStatus sends a DELETE request and returns both the status code and body
func DeleteRequestWithStatus(logger *Logger, url string, authHeader string) (int, string) {
	logger.Log(Verbose, "Sending DELETE request to %s\n", url)

	fullUrl := getBaseUrl() + url

	req, err := http.NewRequest("DELETE", fullUrl, nil)
	if err != nil {
		logger.Log(Error, "Error creating DELETE request to %s: %s\n", url, err.Error())
		return 0, err.Error()
	}
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Log(Error, "Error sending DELETE request to %s: %s\n", url, err.Error())
		return 0, err.Error()
	}

	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Log(Error, "Error reading response body: %s\n", err.Error())
		return res.StatusCode, err.Error()
	}

	logger.Log(Verbose, "Response (%d): %s\n", res.StatusCode, string(resBody))

	return res.StatusCode, string(resBody)
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

// IsOk checks if body is okay
func IsOk(body string) bool {
	return strings.Contains(body, "\"ok\":1")
}

// AssertOk checks if body is okay and returns a Result
func AssertOk(body string, err string) Result {
	if IsOk(body) {
		return ResultOk()
	}
	return NewResult(false, err)
}

// AssertNotOk checks if body is not okay and returns a Result
func AssertNotOk(body string, err string) Result {
	if !IsOk(body) {
		return ResultOk()
	}
	return NewResult(false, err)
}

// AssertStatus checks that the given status code matches expected and returns a Result
func AssertStatus(actual int, expected int, err string) Result {
	if actual == expected {
		return ResultOk()
	}
	return NewResult(false, err)
}

type JsonType int

const (
	StringType JsonType = iota
	IntType
	BoolType
	Float64Type
)

// IsValue checks if body contains key with value
func IsValue(body string, key string, valueType JsonType, value any) bool {
	data := []byte(body)

	switch valueType {
	case StringType:
		return fastjson.GetString(data, key) == value
	case IntType:
		return fastjson.GetInt(data, key) == value
	case BoolType:
		return fastjson.GetBool(data, key) == value
	case Float64Type:
		return fastjson.GetFloat64(data, key) == value
	default:
		return false
	}
}

// ExtractString extracts a string field from a JSON response body
func ExtractString(body string, key string) string {
	return fastjson.GetString([]byte(body), key)
}

// ExtractInt extracts an int field from a JSON response body
func ExtractInt(body string, key string) int {
	return fastjson.GetInt([]byte(body), key)
}