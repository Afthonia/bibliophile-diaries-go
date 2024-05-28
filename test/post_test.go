package test

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var baseURLPost string = "http://127.0.0.1:8081/post"

func TestGetPosts(t *testing.T) {

	// Create a request with the request body
	req, err := http.Get(fmt.Sprintf("%s/list", baseURLPost))
	if err != nil {
		t.Fatal(err)
	}

	// Check the status code
	assert.Equal(t, http.StatusOK, req.StatusCode)
}

func TestCreatePostUnauthorized(t *testing.T) {
	// Create a request body for the test
	requestBody := []byte(`{"title": "Test Title", "content": "Test Content", "bookTitle": "Test Book", "vote": 5}`)

	// Create a request with the request body
	req, err := http.Post(baseURLPost, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer randomString")

	// Check the status code
	assert.Equal(t, http.StatusUnauthorized, req.StatusCode)
}

func TestCreatePostAuthorized(t *testing.T) {
	client := &http.Client{}

	// Create a request body for the test
	requestBody := []byte(`{"title": "Test Title", "content": "Test Content", "bookTitle": "Test Book", "vote": 6}`)

	// Create a request with the request body
	req, err := http.NewRequest("POST", baseURLPost, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFpc2VzZW43NjVAZ21haWwuY29tIiwiZXhwIjoxNzI0NjYyNDA5LCJuYW1lIjoiYXlzZSIsInVzZXJfaWQiOiIxIn0.LYbqO5njEFJS9-ZU6p_HH9D_9QtlvdoVy_TIHDZ92F0"

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	// Check the status code
	assert.Equal(t, http.StatusCreated, res.StatusCode)
}

func TestToggleLikeAuthorized(t *testing.T) {
	client := &http.Client{}

	// Create a request body for the test
	requestBody := []byte(`{"title": "Test Title", "content": "Test Content", "bookTitle": "Test Book", "vote": 6}`)

	// Create a request with the request body
	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/like?id=1", baseURLPost), bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFpc2VzZW43NjVAZ21haWwuY29tIiwiZXhwIjoxNzI0NjYyNDA5LCJuYW1lIjoiYXlzZSIsInVzZXJfaWQiOiIxIn0.LYbqO5njEFJS9-ZU6p_HH9D_9QtlvdoVy_TIHDZ92F0"

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	// Check the status code
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestDeletePostAuthorizedValid(t *testing.T) {
	client := &http.Client{}

	// Create a request with the request body
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s?id=2", baseURLPost), nil)
	if err != nil {
		t.Fatal(err)
	}

	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFpc2VzZW43NjVAZ21haWwuY29tIiwiZXhwIjoxNzI0NjYyNDA5LCJuYW1lIjoiYXlzZSIsInVzZXJfaWQiOiIxIn0.LYbqO5njEFJS9-ZU6p_HH9D_9QtlvdoVy_TIHDZ92F0"

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))


	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	// Check the status code
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestDeletePostAuthorizedInvalid(t *testing.T) {
	client := &http.Client{}

	// Create a request with the request body
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s?id=0", baseURLPost), nil)
	if err != nil {
		t.Fatal(err)
	}

	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFpc2VzZW43NjVAZ21haWwuY29tIiwiZXhwIjoxNzI0NjYyNDA5LCJuYW1lIjoiYXlzZSIsInVzZXJfaWQiOiIxIn0.LYbqO5njEFJS9-ZU6p_HH9D_9QtlvdoVy_TIHDZ92F0"

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))


	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	// Check the status code
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}
