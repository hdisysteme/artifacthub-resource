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

func NewArtifactHubClient() ArtifactHubClient {
	return ArtifactHubClient{
		client:  &http.Client{Timeout: 10 * time.Second},
		baseUrl: baseUrl(),
	}
}

func (a ArtifactHubClient) ListHelmVersion(p Package, version string) (HelmVersion, error) {
	url := fmt.Sprintf("%s/api/v1/packages/helm/%s/%s/%s", a.baseUrl, p.RepositoryName, p.PackageName, version)
	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return HelmVersion{}, fmt.Errorf("build new artifacthub http request failed: %s", err)
	}

	prepareHttpHeader(p, request)

	response, err := a.client.Do(request)
	if err != nil {
		return HelmVersion{}, fmt.Errorf("error while requesting artifacthub: %w", err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return HelmVersion{}, fmt.Errorf(
			"artifacthub http request returned status code: %d with message: %w",
			response.StatusCode,
			err,
		)
	}

	var target HelmVersion
	err = json.NewDecoder(response.Body).Decode(&target)

	if err != nil {
		return HelmVersion{}, fmt.Errorf("could not marshal JSON: %s", err)
	}

	return target, nil
}

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
			printError(target.Name, target.AvailableVersions[i])
		}

		otherVersion, err := semver.NewVersion(target.AvailableVersions[j].Version)

		if err != nil {
			printError(target.Name, target.AvailableVersions[j])
		}

		return version.LessThan(otherVersion)
	})

	var versions []Version

	for _, version := range target.AvailableVersions {
		versions = append(versions, Version{
			TS:      time.Time(version.TS).UTC(),
			Version: version.Version,
		})
	}

	return versions, nil

}

func printError(name string, target AvailableVersion) (int, error) {
	return fmt.Println(fmt.Printf(fmt.Sprintf(
		"Error while getting semver version for package %s and version %s",
		name,
		target.Version,
	)))
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

func (t Epoch) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", time.Time(t).Format(time.RFC3339))), nil
}

func (t *Epoch) UnmarshalJSON(s []byte) (err error) {
	q, err := strconv.ParseInt(string(s), 10, 64)

	if err != nil {
		return err
	}
	*(*time.Time)(t) = time.Unix(q, 0)
	return
}

func (t Epoch) String() string { return time.Time(t).String() }

// ArtifactHub for testing purposes.
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o fakes/fake_artifacthub.go . ArtifactHub
type ArtifactHub interface {
	ListHelmVersions(p Package) ([]Version, error)
	ListHelmVersion(p Package, version string) (HelmVersion, error)
}

type ArtifactHubClient struct {
	client  *http.Client
	baseUrl string
}

type Package struct {
	RepositoryName string
	PackageName    string
	ApiKey         string
}

type Epoch time.Time

type AvailableVersion struct {
	Version string `json:"version"`
	TS      Epoch  `json:"ts"`
}

type Repository struct {
	Url                     string `json:"url"`
	DisplayName             string `json:"display_name"`
	OrganizationDisplayName string `json:"organization_display_name"`
}

type HelmVersion struct {
	AppVersion        string             `json:"app_version"`
	ContentUrl        string             `json:"content_url"`
	TS                Epoch              `json:"ts"`
	Name              string             `json:"name"`
	Version           string             `json:"version"`
	AvailableVersions []AvailableVersion `json:"available_versions"`
	Repository        Repository         `json:"repository"`
}

type Version struct {
	TS      time.Time `json:"ts"`
	Version string    `json:"version"`
}
