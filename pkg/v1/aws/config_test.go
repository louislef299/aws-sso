package aws_test

import (
	"os"
	"path"

	"github.com/aws/aws-sdk-go-v2/config"
	laws "github.com/louislef299/aws-sso/pkg/v1/aws"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Account", Ordered, func() {
	BeforeAll(func() {
		// Create a temporary directory for testing usage information
		err := os.Mkdir(testFolder, 0744)
		Expect(err).ShouldNot(HaveOccurred())

		home := os.Getenv("HOME")
		pwd, err := os.Getwd()
		Expect(err).ShouldNot(HaveOccurred())

		err = os.Setenv("HOME", path.Join(pwd, testFolder))
		Expect(err).ShouldNot(HaveOccurred())

		DeferCleanup(func() {
			err := os.RemoveAll(testFolder)
			Expect(err).ShouldNot(HaveOccurred())

			err = os.Setenv("HOME", home)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Context("When gathering AWS configuration sections", func() {
		It("should throw ErrFileNotFound when config folder/file when information isn't found", func() {
			files := []string{
				config.DefaultSharedConfigFilename(),
				config.DefaultSharedCredentialsFilename(),
			}

			for _, f := range files {
				_, err := laws.GetAWSConfigSections(f)
				Expect(err).Should(BeEquivalentTo(laws.ErrFileNotFound))
			}
		})

		It("should return the sample profile in a local config example", func() {
			sampleProf := "sample-config-profile"
			err := laws.WriteAWSConfigFile(sampleProf, "sample-region", "json")
			Expect(err).ShouldNot(HaveOccurred())

			p, err := laws.GetAWSConfigSections(config.DefaultSharedConfigFilename())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(p).Should(ContainElements(sampleProf))
		})
	})

	Context("When gathering AWS profiles", func() {
		It("should create AWS config files when they don't exist", func() {
			_, err := laws.GetAWSProfiles()
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
