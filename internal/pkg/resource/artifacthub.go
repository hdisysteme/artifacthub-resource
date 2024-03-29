// Package resource provides functions for obtaining artifacthub.io Helm Chart versions
package resource

import (
	"encoding/json"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"
)

// NewArtifactHubClient returns an ArtifactHubClient that contains a HTTP client that is already preconfigured.
//
// The contained http.Client is configured as follows.
// http.Timeout = 10sec
// http.Transport = http.ProxyFromEnvironment
//
// The Base URL is https://artifacthub.io and can be overwritten by the Environment Variable ARTIFACTHUB_BASE_URL
func NewArtifactHubClient() ArtifactHubClient {
	return ArtifactHubClient{
		client: &http.Client{
			Timeout:   10 * time.Second,
			Transport: &http.Transport{Proxy: http.ProxyFromEnvironment},
		},
		baseUrl: baseUrl(),
	}
}

// ListHelmVersion returns a specific HelmVersion of the given Package
func (a ArtifactHubClient) ListHelmVersion(p Package, version string) (*HelmVersion, error) {
	url := fmt.Sprintf("%s/api/v1/packages/helm/%s/%s/%s", a.baseUrl, p.RepositoryName, p.PackageName, version)
	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, fmt.Errorf("build new artifacthub http request failed: %s", err)
	}

	prepareHttpHeader(p, request)

	response, err := a.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error while requesting artifacthub: %w", err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"artifacthub http request returned status code: %d with message: %w",
			response.StatusCode,
			err,
		)
	}

	var target HelmVersion
	err = json.NewDecoder(response.Body).Decode(&target)

	if err != nil {
		return nil, fmt.Errorf("could not marshal JSON: %s", err)
	}

	return &target, nil
}

// ListHelmVersions lists all available versions for the given Package
// The []Version is returned in descending order of the Version
func (a ArtifactHubClient) ListHelmVersions(p Package) ([]Version, error) {

	url := fmt.Sprintf("%s/api/v1/packages/helm/%s/%s", a.baseUrl, p.RepositoryName, p.PackageName)
	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, fmt.Errorf("build new artifacthub http request failed: %s", err)
	}

	prepareHttpHeader(p, request)

	response, err := a.client.Do(request)

	if err != nil {
		return nil, fmt.Errorf("error while requesting artifacthub: %w", err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"artifacthub http request returned status code: %d with message: %w",
			response.StatusCode,
			err,
		)
	}

	var target HelmVersion
	err = json.NewDecoder(response.Body).Decode(&target)

	if err != nil {
		return nil, fmt.Errorf("could not marshal JSON: %s", err)
	}

	sort.Slice(target.AvailableVersions, func(i, j int) bool {
		version, err := semver.NewVersion(target.AvailableVersions[i].Version)

		if err != nil {
			printError(target.Name, target.AvailableVersions[i], nil)
		}

		otherVersion, err := semver.NewVersion(target.AvailableVersions[j].Version)

		if err != nil {
			printError(target.Name, target.AvailableVersions[j], err)
		}

		return version.LessThan(otherVersion)
	})

	var versions []Version

	for _, version := range target.AvailableVersions {
		versions = append(versions, Version{
			CreatedAt: time.Time(version.TS).UTC(),
			Version:   version.Version,
		})
	}

	return versions, nil

}

// MarshalJSON marshals an Epoch into a formatted time.RFC3339 representation
func (t Epoch) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", time.Time(t).Format(time.RFC3339))), nil
}

// UnmarshalJSON unmarshals the given time.RFC3339 formatted string to an Epoch representation
func (t *Epoch) UnmarshalJSON(s []byte) (err error) {
	q, err := strconv.ParseInt(string(s), 10, 64)

	if err != nil {
		return err
	}
	*(*time.Time)(t) = time.Unix(q, 0)
	return
}

// String transforms an Epoch to a time.Time string representation
func (t Epoch) String() string { return time.Time(t).String() }

// ArtifactHub is the interface implemented by
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o ./fakes/fake_artifacthub.go . ArtifactHub
type ArtifactHub interface {
	ListHelmVersions(p Package) ([]Version, error)
	ListHelmVersion(p Package, version string) (*HelmVersion, error)
}

// ArtifactHubClient is used to query the artifacthub.io endpoint.
type ArtifactHubClient struct {
	client  *http.Client
	baseUrl string
}

func printError(name string, target AvailableVersion, err error) {
	fmt.Println(fmt.Sprintf(
		"Error while getting semver version for package %s and version %s with error: %v",
		name,
		target.Version,
		err,
	))
}

func baseUrl() string {
	var baseUrl string
	baseUrl, ok := os.LookupEnv("ARTIFACTHUB_BASE_URL")
	if !ok {
		baseUrl = "https://artifacthub.io"
	}
	return baseUrl
}

func prepareHttpHeader(p Package, request *http.Request) {
	request.Header.Add("User-Agent", "artifacthub-resource/0.1")
	request.Header.Add("Accept", "application/json")

	if len(p.ApiKey) > 0 {
		request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", p.ApiKey))
	}
}

// Package represents an artifacthub Helm Package
type Package struct {
	RepositoryName string
	PackageName    string
	ApiKey         string
}

// Epoch is an alias for time.Time
type Epoch time.Time

// AvailableVersion represents a version and version timestamp for a HelmVersion
type AvailableVersion struct {
	Version string `json:"version"`
	TS      Epoch  `json:"ts"`
}

// Repository represents information about the repository of a HelmVersion
type Repository struct {
	Url                     string `json:"url"`
	DisplayName             string `json:"display_name"`
	Name                    string `json:"name"`
	OrganizationDisplayName string `json:"organization_display_name"`
}

// HelmVersion represents a helm chart package version
type HelmVersion struct {
	AppVersion        string             `json:"app_version"`
	ContentUrl        string             `json:"content_url"`
	TS                Epoch              `json:"ts"`
	Name              string             `json:"name"`
	Version           string             `json:"version"`
	AvailableVersions []AvailableVersion `json:"available_versions"`
	Repository        Repository         `json:"repository"`
}

// Version represents a specific version for a HelmVersion
type Version struct {
	CreatedAt time.Time `json:"created_at"`
	Version   string    `json:"version"`
}
