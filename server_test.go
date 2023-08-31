package main

import (
	"httpserver/database"
	"net/http"
	"os/exec"
	"testing"

	"github.com/gin-gonic/gin"
)

func runCurlCommand(t *testing.T, method, url, requestData, expectedResponse string) {
	cmd := exec.Command("curl", "-X", method, "-H", "Content-Type: application/json", "-d", requestData, url)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to execute cURL command: %v\nOutput: %s", err, output)
	}
	actualResponse := string(output)

	if actualResponse != expectedResponse {
		t.Errorf("Expected response:\n%s\nActual response:\n%s", expectedResponse, actualResponse)
	}
}

func TestCreateSegment(t *testing.T) {
	db, _ := database.New()
	s := server{db: db}

	req := &gin.Context{
		Request: &http.Request{
			// Body: ,
		},
	}

	s.createSegment(req)

	url := "http://localhost:8080/segment/create"
	requestData := `{"seg_name": "NewSegment"}`
	expectedResponse := "Segment created successfully"
	runCurlCommand(t, "POST", url, requestData, expectedResponse)
}

