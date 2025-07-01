package os_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	pluginos "github.com/louislef299/aws-sso/pkg/os"
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

	Context("When checking for managed profiles", func() {
		It("Should return true for managed profiles", func() {
			Expect(pluginos.IsManagedProfile("profile-aws-sso")).To(BeTrue())
		})

		It("Should return false for unmanaged profiles", func() {
			Expect(pluginos.IsManagedProfile("profile")).To(BeFalse())
		})
	})
})
