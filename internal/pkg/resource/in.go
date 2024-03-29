package resource

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

// Get metadata for GetRequest will fetch meta information for the given helm chart version
func Get(request GetRequest, path string, repository ArtifactHub) (*GetResponse, error) {

	version, err := repository.ListHelmVersion(Package{
		RepositoryName: request.Source.RepositoryName,
		PackageName:    request.Source.PackageName,
		ApiKey:         request.Source.ApiKey,
	}, request.Version.Version)

	if err != nil {
		return nil, err
	}

	var metadata = &Metadata{}
	metadata.append("app_version", version.AppVersion)
	metadata.append("charts_url", version.Repository.Url)
	metadata.append("chart_download_url", version.ContentUrl)
	metadata.append("name", version.Name)
	metadata.append("organization_name", version.Repository.OrganizationDisplayName)
	metadata.append("repository_name", version.Repository.Name)
	metadata.append("repository_display_name", version.Repository.DisplayName)
	metadata.append("version", version.Version)

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %s", err)
	}
	for _, metadatum := range *metadata {
		if err := ioutil.WriteFile(filepath.Join(path, metadatum.Name), []byte(metadatum.Value), 0600); err != nil {
			return nil, fmt.Errorf("failed to write %s: %s", metadatum.Name, err)
		}
	}

	return &GetResponse{
		Version: Version{
			CreatedAt: time.Time(version.TS).UTC(),
			Version:   version.Version,
		},
		Metadata: *metadata,
	}, nil
}

// GetRequest contains the information for a specific Source and Version
type GetRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

// GetResponse contains a Version and Metadata for a Version
type GetResponse struct {
	Version  Version  `json:"version"`
	Metadata Metadata `json:"metadata,omitempty"`
}

// Metadata is a store for key, value information
type Metadata []*metadataPair

type metadataPair struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (m *Metadata) append(name string, value string) {
	*m = append(*m, &metadataPair{
		Name:  name,
		Value: value,
	})
}
