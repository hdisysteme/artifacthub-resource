package main

import (
	"encoding/json"
	resource "github.com/hdisysteme/artifacthub-resource"
	"log"
	"os"
)

func main() {

	var request resource.CheckRequest

	decoder := json.NewDecoder(os.Stdin)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		log.Fatalf("failed to unmarshal request: %s", err)
	}

	response, err := resource.Check(request, resource.NewArtifactHubClient())

	if err != nil {
		log.Fatalf("resource check failed with: %s", err)
	}

	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		log.Fatalf("failed to marshal response: %s", err)
	}

}
