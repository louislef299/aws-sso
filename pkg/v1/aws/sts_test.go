package aws_test

import (
	"log"
	"os"
	"path"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sso"
	"github.com/aws/aws-sdk-go-v2/service/sso/types"
	laws "github.com/louislef299/aws-sso/pkg/v1/aws"
	los "github.com/louislef299/aws-sso/pkg/v1/os"
)

const testFolder = "testing_files"

var _ = Describe("Sts", func() {
	Context("When saving usage information", Ordered, func() {
		BeforeAll(func() {
			// Create a temporary directory for testing usage information
			err := os.Mkdir(testFolder, 0744)
			Expect(err).ShouldNot(HaveOccurred())

			pwd, err := os.Getwd()
			Expect(err).ShouldNot(HaveOccurred())
			err = os.Setenv("HOME", path.Join(pwd, testFolder))
			Expect(err).ShouldNot(HaveOccurred())

			DeferCleanup(func() {
				err := os.RemoveAll(testFolder)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		It("Should create aws sso cache folders and files if they don't exist", func() {
			err := laws.SaveUsageInformation(&types.AccountInfo{
				AccountId:    aws.String("0123456789"),
				AccountName:  aws.String("temp"),
				EmailAddress: aws.String("john.smith@email.com"),
			}, &types.RoleInfo{
				AccountId: aws.String("0123456789"),
				RoleName:  aws.String("sample-role"),
			})
			Expect(err).ShouldNot(HaveOccurred())

			exists, err := los.IsFileOrFolderExisting(path.Join(testFolder, laws.LastUsageLocation))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(exists).Should(BeTrue())
		})

		It("Should create aws credentials folder and files if they don't exist", func() {
			err := laws.WriteAWSCredentialsFile("sample-profile", &sso.GetRoleCredentialsOutput{
				RoleCredentials: &types.RoleCredentials{
					AccessKeyId:     aws.String("dummyKey"),
					Expiration:      121212,
					SecretAccessKey: aws.String("dummySecret"),
					SessionToken:    aws.String("dummyToken"),
				},
			})
			Expect(err).ShouldNot(HaveOccurred())

			credentialsFile := config.DefaultSharedCredentialsFilename()
			log.Println("checking if", credentialsFile, "exists")
			exists, err := los.IsFileOrFolderExisting(credentialsFile)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(exists).Should(BeTrue())
		})

		It("Should create aws config folder and files if they don't exist", func() {
			err := laws.WriteAWSConfigFile("sample-profile", "sample-region", "json")
			Expect(err).ShouldNot(HaveOccurred())

			configFile := config.DefaultSharedConfigFilename()
			log.Println("checking if", configFile, "exists")
			exists, err := los.IsFileOrFolderExisting(configFile)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(exists).Should(BeTrue())
		})
	})

	Context("When gathering the region", func() {
		It("Should gather the proper region for all expected region inputs", func() {
			for _, r := range laws.AwsRegions {
				err := os.Setenv("AWS_REGION", r)
				Expect(err).NotTo(HaveOccurred())

				region, err := laws.GetRegion()
				Expect(err).NotTo(HaveOccurred())
				Expect(region).To(Equal(r))
			}
		})

		It("Should error for improper regions", func() {
			improperRegions := []string{"us-west-11", "dummy-region", "cn-south-10"}
			for _, r := range improperRegions {
				err := os.Setenv("AWS_REGION", r)
				Expect(err).NotTo(HaveOccurred())

				region, err := laws.GetRegion()
				Expect(region).To(BeEmpty())
				Expect(err).To(Equal(laws.ErrRegionInvalid))
			}
		})

		It("Should error when no region is given", func() {
			err := os.Unsetenv("AWS_REGION")
			Expect(err).NotTo(HaveOccurred())
			err = os.Unsetenv("AWS_DEFAULT_REGION")
			Expect(err).NotTo(HaveOccurred())

			region, err := laws.GetRegion()
			Expect(region).To(BeEmpty())
			Expect(err).To(Equal(laws.ErrRegionNotFound))
		})
	})

	Context("When gathering the AWS URL", func() {
		It("Should gather the global region for valid global regions", func() {
			globalRegions := []string{"us-east-2", "us-east-1", "us-west-1", "us-west-2", "af-south-1", "ap-east-1", "ap-south-1"}
			for _, r := range globalRegions {
				err := os.Setenv("AWS_REGION", r)
				Expect(err).NotTo(HaveOccurred())

				url, err := laws.GetURL()
				Expect(err).NotTo(HaveOccurred())
				Expect(url).To(Equal("amazonaws.com"))
			}
		})

		It("Should return the China url for China regions", func() {
			chinaRegions := []string{"cn-north-1", "cn-northwest-1"}
			for _, r := range chinaRegions {
				err := os.Setenv("AWS_REGION", r)
				Expect(err).NotTo(HaveOccurred())

				url, err := laws.GetURL()
				Expect(err).NotTo(HaveOccurred())
				Expect(url).To(Equal("amazonaws.com.cn"))
			}
		})

		It("Should error for invalid regions", func() {
			improperRegions := []string{"us-west-11", "dummy-region", "cn-south-10"}
			for _, r := range improperRegions {
				err := os.Setenv("AWS_REGION", r)
				Expect(err).NotTo(HaveOccurred())

				url, err := laws.GetURL()
				Expect(url).To(BeEmpty())
				Expect(err).To(Equal(laws.ErrRegionInvalid))
			}
		})
	})
})
