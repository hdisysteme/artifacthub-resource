// +build e2e

package e2e_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	resource "github.com/PG2000/artifacthub-resource"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/ghttp"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

var _ = Describe("E2E Check Resource", func() {

	var (
		token        string
		server       *ghttp.Server
		execPath     string
		session      *Session
		jsonResponse = `
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
      "created_at": 1605806528
    },
    {
      "version": "9.2.4",
      "created_at": 1606316622
    }
  ],
  "app_version": "8.5.1-community",
  "digest": "d0a4a8230cd5e23beff38131e3aad706bcab97a79c6bf26a94318957271a8c6e",
  "deprecated": false,
  "signed": false,
  "content_url": "https://git.local/acme/charts/releases/download/some-package-9.2.4/some-package-9.2.4.tgz",
  "has_values_schema": false,
  "has_changelog": false,
  "created_at": 1606316622,
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
		execPath = buildExec("github.com/PG2000/artifacthub-resource/cmd/check")
		server = ghttp.NewServer()
		token = "MY_SECRET_TOKEN"

	})

	AfterEach(func() {
		CleanupBuildArtifacts()
		server.Close()
	})

	When("check is executed with a repository and a package name", func() {
		It("it should send the api token and return the expected output", func() {

			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/api/v1/packages/helm/acme-charts/some-package"),
				ghttp.VerifyHeader(http.Header{
					"Authorization": []string{"Bearer " + token},
				}),
				ghttp.RespondWith(http.StatusOK, jsonResponse),
			))

			session = executeCheckCommand(
				execPath,
				fmt.Sprintf("{ \"source\": {\"repository_name\": \"acme-charts\", \"package_name\": \"some-package\", \"api_key\": \"%s\"} }", token),
				nil,
				"ARTIFACTHUB_BASE_URL=http://"+server.Addr(),
			)

			Eventually(session).Should(Exit(0))

			var result = resource.CheckResponse{}
			err := json.NewDecoder(bytes.NewBuffer(session.Out.Contents())).Decode(&result)

			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal(resource.CheckResponse{
				{
					CreatedAt: time.Date(2020, 11, 19, 17, 22, 8, 0, time.UTC),
					Version:   "9.2.0",
				},
				{
					CreatedAt: time.Date(2020, 11, 25, 15, 3, 42, 0, time.UTC),
					Version:   "9.2.4",
				},
			}))

		})

	})
})

func buildExec(packagePath string) string {
	compiledPath, err := Build(packagePath)
	Expect(err).NotTo(HaveOccurred())
	return compiledPath
}

func executeCheckCommand(compilePath string, stdinInput string, args []string, envs ...string) *Session {
	cmd := exec.Command(compilePath)

	if len(args) > 0 {
		cmd.Args = args
	}

	cmd.Env = os.Environ()

	for _, env := range envs {
		cmd.Env = append(cmd.Env, env)
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, stdinInput)
	}()

	session, err := Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	return session
}
