// +build e2e

package e2e_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hdisysteme/artifacthub-resource/internal/pkg/resource"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/ghttp"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var _ = Describe("E2E In Resource", func() {

	var (
		apiToken                    string
		server                      *ghttp.Server
		execPath                    string
		session                     *Session
		tmpDir                      string
		fakeArtifactHubJsonResponse = `
{
  "package_id": "be378d3f-d6c5-47ac-a2ae-0cb5d9f6d8f5",
  "name": "some-package",
  "normalized_name": "some-package",
  "logo_image_id": "d77ccb10-b972-4494-a813-f9b0f303bcde",
  "is_operator": false,
  "description": "SomePackage is an open sourced code quality scanning tool",
  "keywords": [
    "coverage",
    "security",
    "code",
    "quality"
  ],
  "home_url": "https://www.example.local/",
  "readme": "# README",
  "links": [
    {
      "url": "https://git.local/SomePackage/docker-some-package",
      "name": "source"
    }
  ],
  "security_report_created_at": 1608740109,
  "data": {
    "dependencies": []
  },
  "version": "9.2.4",
  "available_versions": [
    {
      "version": "9.2.0",
      "ts": 1605806528
    },
    {
      "version": "9.2.4",
      "ts": 1606316622
    },
    {
      "version": "9.1.2",
      "ts": 1604507915
    }
  ],
  "app_version": "8.5.1-community",
  "digest": "d0a4a8230cd5e23beff38131e3aad706bcab97a79c6bf26a94318957271a8c6e",
  "deprecated": false,
  "signed": false,
  "content_url": "https://git.local/acme/charts/releases/download/some-package-9.2.4/some-package-9.2.4.tgz",
  "has_values_schema": false,
  "has_changelog": false,
  "ts": 1606316622,
  "maintainers": [
    {
      "name": "acme",
      "email": "acme@gmail.com"
    }
  ],
  "repository": {
    "repository_id": "534a9dcb-0942-4ebb-b1d8-3a716e80f17e",
    "name": "acme-charts",
    "display_name": "Acme Charts",
    "url": "https://acme.github.io/charts",
    "private": false,
    "kind": 0,
    "verified_publisher": false,
    "official": false,
    "organization_name": "acme",
    "organization_display_name": "Acme"
  }
}`
	)

	BeforeEach(func() {
		execPath = buildExec("github.com/hdisysteme/artifacthub-resource/cmd/in")
		server = ghttp.NewServer()
		apiToken = "MY_SECRET_TOKEN"

		var err error
		tmpDir, err = ioutil.TempDir("", "resource-test-")
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		CleanupBuildArtifacts()
		server.Close()
		os.RemoveAll(tmpDir)
	})

	When("in is executed with api key", func() {

		BeforeEach(func() {
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/api/v1/packages/helm/acme-charts/some-package/9.2.4"),
				ghttp.VerifyHeader(http.Header{
					"Authorization": []string{"Bearer " + apiToken},
				}),
				ghttp.RespondWith(http.StatusOK, fakeArtifactHubJsonResponse),
			))

			session = executeCheckCommand(
				execPath,
				fmt.Sprintf("{ \"source\": {\"repository_name\": \"acme-charts\", \"package_name\": \"some-package\", \"api_key\": \"%s\"}, \"version\": {\"created_at\":\"2020-11-25T16:03:42+01:00\",\"version\":\"9.2.4\"} }", apiToken),
				[]string{"/opt/resource/in", tmpDir},
				"ARTIFACTHUB_BASE_URL=http://"+server.Addr(),
			)

			Eventually(session).Should(Exit(0))

		})

		It("it should return the expected response and metadata", func() {
			var response = resource.GetResponse{}
			err := json.NewDecoder(bytes.NewBuffer(session.Out.Contents())).Decode(&response)

			Expect(err).ToNot(HaveOccurred())
			Expect(response.Version).To(Equal(
				resource.Version{
					CreatedAt: time.Date(2020, 11, 25, 15, 3, 42, 0, time.UTC),
					Version:   "9.2.4",
				}))

			Expect(response.Metadata).To(ConsistOf(resource.Metadata{
				{Name: "app_version", Value: "8.5.1-community"},
				{Name: "charts_url", Value: "https://acme.github.io/charts"},
				{Name: "chart_download_url", Value: "https://git.local/acme/charts/releases/download/some-package-9.2.4/some-package-9.2.4.tgz"},
				{Name: "name", Value: "some-package"},
				{Name: "organization_name", Value: "Acme"},
				{Name: "repository_name", Value: "Acme Charts"},
				{Name: "version", Value: "9.2.4"},
			}))
		})

		It("should create files with details about the given version", func() {
			testFileContainsExpectedText(tmpDir, "app_version", "8.5.1-community")
			testFileContainsExpectedText(tmpDir, "version", "9.2.4")
			testFileContainsExpectedText(tmpDir, "name", "some-package")
			testFileContainsExpectedText(tmpDir, "chart_download_url",
				"https://git.local/acme/charts/releases/download/some-package-9.2.4/some-package-9.2.4.tgz")
		})

	})
})

func testFileContainsExpectedText(dir string, filename string, expectedText string) {
	file, err := ioutil.ReadFile(dir + "/" + filename)
	Expect(err).ToNot(HaveOccurred())
	Expect(string(file)).To(Equal(expectedText))
}
