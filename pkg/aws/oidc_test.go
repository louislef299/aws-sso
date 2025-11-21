package aws_test

import (
	"os"
	"path"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/louislef299/knot/pkg/aws"
)

var _ = Describe("Oidc", func() {
	Context("When getting AWS client configurations", func() {
		It("Should error there are too many configuration locations", func() {
			_, err := aws.ClientInfoFileDestination("oneLocation", "twoLocations")
			Expect(err).To(Equal(aws.ErrMoreThanOneLocation))
		})

		It("Shouldn't error with sample local configuration", func() {
			configDir, err := os.Getwd()
			Expect(err).NotTo(HaveOccurred())
			returnedLocation, err := aws.ClientInfoFileDestination(configDir)
			Expect(err).NotTo(HaveOccurred())

			DeferCleanup(func() {
				err := os.RemoveAll(path.Join(configDir, ".aws"))
				Expect(err).NotTo(HaveOccurred())
			})

			atok, err := aws.GetTokenHash(aws.GetAccessToken())
			Expect(err).NotTo(HaveOccurred())
			Expect(returnedLocation).To(Equal(path.Join(configDir, aws.AWS_TOKEN_PATH, atok)))
		})
	})
})
