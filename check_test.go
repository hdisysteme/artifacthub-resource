package resource_test

import (
	"fmt"
	resource "github.com/PG2000/artifacthub-resource"
	"github.com/PG2000/artifacthub-resource/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("ArtifacthubResource Check", func() {

	var (
		artifacthub  *fakes.FakeArtifactHub
		checkRequest resource.CheckRequest
	)

	BeforeEach(func() {
		artifacthub = new(fakes.FakeArtifactHub)
		checkRequest = createCheckRequest("acme-charts", "my-package-name", "some-fake-api-key")
	})

	When("check is called with missing parameters", func() {

		testdata := []struct {
			description    string
			repositoryName string
			packageName    string
			apiKey         string
		}{
			{description: "should return an error when package name and repository name are empty", repositoryName: "", packageName: "", apiKey: ""},
			{description: "should return an error when package name is empty", repositoryName: "acme-charts", packageName: "", apiKey: ""},
			{description: "should return an error when repository name is empty", repositoryName: "", packageName: "my-package-name", apiKey: ""},
		}

		for _, data := range testdata {
			It(data.description, func() {
				test(createCheckRequest(data.repositoryName, data.packageName, data.apiKey), artifacthub)
			})
		}

	})

	When("check is called with valid source", func() {

		It("should call list versions with expected parameters", func() {
			var packageVersions []resource.Version

			packageVersions = append(packageVersions, resource.Version{
				Version: "9.2.4",
				TS:      time.Now(),
			})

			artifacthub.ListHelmVersionsReturns(packageVersions, nil)

			check, err := resource.Check(checkRequest, artifacthub)

			Expect(artifacthub.ListHelmVersionsCallCount()).To(Equal(1))
			Expect(artifacthub.ListHelmVersionsArgsForCall(0)).To(Equal(resource.Package{
				RepositoryName: "acme-charts",
				PackageName:    "my-package-name",
				ApiKey:         "some-fake-api-key",
			}))
			Expect(check).To(HaveLen(1))
			Expect(err).ToNot(HaveOccurred())

		})

	})
	When("list versions fails", func() {

		It("should return an error when call to list versions failed", func() {
			artifacthub.ListHelmVersionsReturns(nil, fmt.Errorf("some error occured"))
			check, err := resource.Check(checkRequest, artifacthub)
			Expect(err).To(HaveOccurred())
			Expect(check).To(BeNil())
		})

	})
})

func test(request resource.CheckRequest, artifacthub *fakes.FakeArtifactHub) {
	check, err := resource.Check(request, artifacthub)
	Expect(check).To(BeNil())
	Expect(err).To(HaveOccurred())
}

func createCheckRequest(repositoryName string, packageName string, apiKey string) resource.CheckRequest {
	return resource.CheckRequest{
		Source: resource.Source{
			RepositoryName: repositoryName,
			PackageName:    packageName,
			ApiKey:         apiKey,
		},
	}
}
