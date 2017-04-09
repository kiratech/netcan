package integration_test

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Basecommand", func() {
	Describe("Extracting base network info", func() {
		Context("Without any specific container created", func() {
			It("should print only host interfaces and docker0", func() {
				var err error
				cmd := exec.Command("docker", "run", "--privileged", "--rm", "-v", "/home/fntlnz/go/src/github.com/fntlnz/netcan:/go/src/github.com/fntlnz/netcan", "-w", "/go/src/github.com/fntlnz/netcan", "alpine:latest", "./netcan")

				session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
				Ω(err).ShouldNot(HaveOccurred())
				Eventually(session).Should(gexec.Exit(0))
			})
		})
	})
})
