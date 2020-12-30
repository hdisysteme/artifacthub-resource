package resource_test

import (
	"fmt"
	resource "github.com/PG2000/artifacthub-resource"
	"github.com/PG2000/artifacthub-resource/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("ArtifacthubResource", func() {

	var (
		artifacthub  *fakes.FakeArtifactHub
		checkRequest resource.CheckRequest
	)

	BeforeEach(func() {
		artifacthub = new(fakes.FakeArtifactHub)
		checkRequest = createCheckRequest("acme-charts", "my-package-name", "some-fake-api-key")
	})

	When("check is called with missing parameters", func() {

		It("should return an error when package name and repository name are empty", func() {
			test(createCheckRequest("", "", ""), artifacthub)
		})

		It("should return an error when package name is empty", func() {
			test(createCheckRequest("acme-charts", "", ""), artifacthub)
		})

		It("should return an error when repository name is empty", func() {
			test(createCheckRequest("", "my-package-name", ""), artifacthub)
		})

		It("should return no error if only apiKey is empty", func() {
			var packageVersions []resource.Version
			packageVersions = append(packageVersions, resource.Version{})

			artifacthub.ListVersionsReturns(packageVersions, nil)
			check, err := resource.Check(createCheckRequest("acme-charts", "my-package-name", ""), artifacthub)

			Expect(check).ToNot(BeNil())
			Expect(err).ToNot(HaveOccurred())
		})
	})

	When("check is called", func() {

		It("should call list versions with expected parameters", func() {
			var packageVersions []resource.Version

			packageVersions = append(packageVersions, resource.Version{
				Version:   "9.2.4",
				CreatedAt: time.Now(),
			})

			artifacthub.ListVersionsReturns(packageVersions, nil)

			check, err := resource.Check(checkRequest, artifacthub)

			Expect(artifacthub.ListVersionsCallCount()).To(Equal(1))
			Expect(artifacthub.ListVersionsArgsForCall(0)).To(Equal(resource.Package{
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
			artifacthub.ListVersionsReturns(nil, fmt.Errorf("some error occured"))
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
