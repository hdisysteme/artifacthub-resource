package resource_test

import (
	"github.com/hdisysteme/artifacthub-resource/internal/pkg/resource"
	"github.com/hdisysteme/artifacthub-resource/internal/pkg/resource/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"time"
)

var _ = Describe("Artifacthub Resource In", func() {

	var (
		artifacthub     *fakes.FakeArtifactHub
		getRequest      resource.GetRequest
		fixedTime       time.Time
		testHelmVersion *resource.HelmVersion
	)

	BeforeEach(func() {
		artifacthub = new(fakes.FakeArtifactHub)
		fixedTime = time.Now().UTC()
		getRequest = resource.GetRequest{
			Source: resource.Source{
				RepositoryName: "acme-charts",
				PackageName:    "my-package-name",
				ApiKey:         "some-fake-api-key",
			},
			Version: resource.Version{
				Version:   "9.2.4",
				CreatedAt: fixedTime,
			},
		}
		testHelmVersion = &resource.HelmVersion{
			AppVersion:        "8.2.1",
			ContentUrl:        "https://git.local/",
			TS:                resource.Epoch(fixedTime),
			Name:              "some-package",
			Version:           "9.2.4",
			AvailableVersions: nil,
			Repository: resource.Repository{
				Url:                     "https://git.local/some-package/",
				DisplayName:             "Some Package",
				Name:                    "some-package",
				OrganizationDisplayName: "Acme Charts",
			},
		}
	})

	When("in is called with valid arguments", func() {
		It("should call list helm version with specific version", func() {

			artifacthub.ListHelmVersionReturns(testHelmVersion, nil)

			_, err := resource.Get(getRequest, os.TempDir(), artifacthub)
			Expect(err).ToNot(HaveOccurred())

			pkg, version := artifacthub.ListHelmVersionArgsForCall(0)
			Expect(pkg).To(Equal(resource.Package{
				RepositoryName: "acme-charts",
				PackageName:    "my-package-name",
				ApiKey:         "some-fake-api-key",
			}))
			Expect(artifacthub.ListHelmVersionCallCount()).To(Equal(1))
			Expect(version).To(Equal("9.2.4"))

		})

		It("should return a response with expected version and metadata", func() {

			artifacthub.ListHelmVersionReturns(testHelmVersion, nil)

			response, err := resource.Get(getRequest, os.TempDir(), artifacthub)

			Expect(err).ToNot(HaveOccurred())
			Expect(response.Version).To(Equal(resource.Version{
				Version:   "9.2.4",
				CreatedAt: fixedTime,
			}))

			Expect(response.Metadata).To(ConsistOf(resource.Metadata{
				{
					Name:  "app_version",
					Value: "8.2.1",
				},
				{
					Name:  "charts_url",
					Value: "https://git.local/some-package/",
				},
				{
					Name:  "chart_download_url",
					Value: "https://git.local/",
				},
				{
					Name:  "name",
					Value: "some-package",
				},
				{
					Name:  "organization_name",
					Value: "Acme Charts",
				},
				{
					Name:  "repository_name",
					Value: "some-package",
				},
				{
					Name:  "repository_display_name",
					Value: "Some Package",
				},
				{
					Name:  "version",
					Value: "9.2.4",
				},
			}))

		})
	})

})
