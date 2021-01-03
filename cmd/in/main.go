package main

import (
	"encoding/json"
	resource "github.com/PG2000/artifacthub-resource"
	"log"
	"os"
)

func main() {

	var request resource.GetRequest

	decoder := json.NewDecoder(os.Stdin)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		log.Fatalf("failed to unmarshal request: %s", err)
	}

	if len(os.Args) < 2 {
		log.Fatalf("missing arguments")
	}

	outputDir := os.Args[1]

	response, err := resource.Get(request, outputDir, resource.NewArtifactHubClient())

	if err != nil {
		log.Fatalf("get failed: %s", err)
	}

	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		log.Fatalf("failed to marshal response: %s", err)
	}
}
