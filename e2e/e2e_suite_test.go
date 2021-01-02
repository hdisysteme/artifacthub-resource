// +build e2e

package e2e_test

import (
	"io"
	"log"
	"os"
	"os/exec"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

func TestE2e(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "E2e Suite")
}

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
