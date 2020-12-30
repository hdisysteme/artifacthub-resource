package main

import (
	"encoding/json"
	resource "github.com/PG2000/artifacthub-resource"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	var request resource.CheckRequest

	decoder := json.NewDecoder(os.Stdin)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		log.Fatalf("failed to unmarshal request: %s", err)
	}

	var httpClient = &http.Client{Timeout: 10 * time.Second}

	baseUrl := baseUrl()

	response, err := resource.Check(
		request,
		resource.NewArtifactHubClient(httpClient, baseUrl),
	)

	if err != nil {
		log.Fatalf("resource check failed with: %s", err)
	}

	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		log.Fatalf("failed to marshal response: %s", err)
	}

}

func baseUrl() string {
	var baseUrl string

	env, ok := os.LookupEnv("ARTIFACTHUB_BASE_URL")

	if ok {
		baseUrl = env
	} else {
		baseUrl = "https://artifacthub.io"
	}
	return baseUrl
}
