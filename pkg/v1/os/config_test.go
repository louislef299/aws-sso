package os_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	pluginos "github.com/louislef299/aws-sso/pkg/v1/os"
)

var _ = Describe("Config", func() {
	Context("When gathering the config path", func() {
		It("Should get the proper home directory", func() {
			home, err := os.UserHomeDir()
			Expect(err).NotTo(HaveOccurred())

			configPath, err := pluginos.GetConfigPath()
			Expect(err).NotTo(HaveOccurred())
			Expect(configPath).To(Equal(home + pluginos.AWS_LOGIN_PATH))
		})
	})
})
