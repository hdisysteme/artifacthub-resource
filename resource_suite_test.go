package resource_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestArtifacthubResource(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ArtifacthubResource Suite")
}
