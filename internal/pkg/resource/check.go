package resource

import "fmt"

// Check for CheckRequest will fetch all versions of a given helm chart
func Check(request CheckRequest, repository ArtifactHub) (*[]Version, error) {

	err := request.validate()

	if err != nil {
		return nil, err
	}

	var versions []Version
	versions, err = repository.ListHelmVersions(Package{
		RepositoryName: request.Source.RepositoryName,
		PackageName:    request.Source.PackageName,
		ApiKey:         request.Source.ApiKey,
	})

	if err != nil {
		return nil, err
	}

	return &versions, nil
}

func (c CheckRequest) validate() error {
	if len(c.Source.PackageName) == 0 || len(c.Source.RepositoryName) == 0 {
		return fmt.Errorf(
			"package name: %s or repository name: %s should not be empty",
			c.Source.PackageName,
			c.Source.RepositoryName,
		)
	}
	return nil
}

// CheckRequest contains the information for the desired Source and Version
type CheckRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

// Source contains information for the helm repository and chart package
type Source struct {
	RepositoryName string `json:"repository_name"`
	PackageName    string `json:"package_name"`
	ApiKey         string `json:"api_key"`
}
