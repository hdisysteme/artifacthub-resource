package resource

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

func Get(request GetRequest, path string, repository ArtifactHub) (GetResponse, error) {

	version, err := repository.ListHelmVersion(Package{
		RepositoryName: request.Source.RepositoryName,
		PackageName:    request.Source.PackageName,
		ApiKey:         request.Source.ApiKey,
	}, request.Version.Version)

	emptyResponse := GetResponse{}

	if err != nil {
		return emptyResponse, err
	}

	var metadata = Metadata{}
	metadata.append("app_version", version.AppVersion)
	metadata.append("charts_url", version.Repository.Url)
	metadata.append("chart_download_url", version.ContentUrl)
	metadata.append("name", version.Name)
	metadata.append("organization_name", version.Repository.OrganizationDisplayName)
	metadata.append("repository_name", version.Repository.DisplayName)
	metadata.append("version", version.Version)

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return emptyResponse, fmt.Errorf("failed to create output directory: %s", err)
	}
	for _, metadatum := range metadata {
		if err := ioutil.WriteFile(filepath.Join(path, metadatum.Name), []byte(metadatum.Value), 0600); err != nil {
			return emptyResponse, fmt.Errorf("failed to write %s: %s", metadatum.Name, err)
		}
	}

	return GetResponse{
		Version: Version{
			CreatedAt: time.Time(version.TS).UTC(),
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
