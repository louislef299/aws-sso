package cmd

import (
	"os"
	"path"

	"github.com/aws/aws-sdk-go-v2/config"
	laws "github.com/louislef299/aws-sso/pkg/v1/aws"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const testFolder = "testing_files"

var _ = Describe("Account", func() {
	Context("When gathering configuration information", Ordered, func() {
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

		It("getAWSConfigSections should throw ErrFileNotFound when config folder/file when information isn't found", func() {
			files := []string{
				config.DefaultSharedConfigFilename(),
				config.DefaultSharedCredentialsFilename(),
			}

			for _, f := range files {
				_, err := getAWSConfigSections(f)
				Expect(err).Should(BeEquivalentTo(ErrFileNotFound))
			}
		})

		It("getAWSConfigSections should return the sample profile in a local config example", func() {
			sampleProf := "sample-config-profile"
			err := laws.WriteAWSConfigFile(sampleProf, "sample-region", "json")
			Expect(err).ShouldNot(HaveOccurred())

			p, err := getAWSConfigSections(config.DefaultSharedConfigFilename())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(p).Should(ContainElements(sampleProf))
		})
	})
})
