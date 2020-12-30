package resource

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

func Get(request GetRequest, dir string, repository ArtifactHub) (GetResponse, error) {

	version, err := repository.ListHelmVersion(Package{
		RepositoryName: request.Source.RepositoryName,
		PackageName:    request.Source.PackageName,
		ApiKey:         request.Source.ApiKey,
	}, request.Version.Version)

	emptyResponse := GetResponse{}

	if err != nil {
		return emptyResponse, err
	}

	var metadata Metadata = Metadata{}
	metadata.append("app_version", version.AppVersion)
	metadata.append("charts_url", version.Repository.Url)
	metadata.append("chart_download_url", version.ContentUrl)
	metadata.append("name", version.Name)
	metadata.append("organization_name", version.Repository.OrganizationDisplayName)
	metadata.append("repository_name", version.Repository.DisplayName)

	path := dir
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return emptyResponse, fmt.Errorf("failed to create output directory: %s", err)
	}

	if err := ioutil.WriteFile(filepath.Join(path, "app_version"), []byte(version.AppVersion), 0644); err != nil {
		return emptyResponse, fmt.Errorf("failed to write version: %s", err)
	}

	if err := ioutil.WriteFile(filepath.Join(path, "version"), []byte(version.Version), 0644); err != nil {
		return emptyResponse, fmt.Errorf("failed to write version: %s", err)
	}

	if err := ioutil.WriteFile(filepath.Join(path, "name"), []byte(version.Name), 0644); err != nil {
		return emptyResponse, fmt.Errorf("failed to write version: %s", err)
	}

	if err := ioutil.WriteFile(filepath.Join(path, "chart_download_url"), []byte(version.ContentUrl), 0644); err != nil {
		return emptyResponse, fmt.Errorf("failed to write version: %s", err)
	}

	return GetResponse{
		Version: Version{
			CreatedAt: time.Time(version.CreatedAt).UTC(),
			Version:   version.Version,
		},
		Metadata: metadata,
	}, nil
}

type GetRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

type Metadata []*MetadataPair

type MetadataPair struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
type GetResponse struct {
	Version  Version  `json:"version"`
	Metadata Metadata `json:"metadata,omitempty"`
}

func (m *Metadata) append(name string, value string) {
	*m = append(*m, &MetadataPair{
		Name:  name,
		Value: value,
	})
}
