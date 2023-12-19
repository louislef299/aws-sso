package aws_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/louislef299/aws-sso/pkg/v1/aws"
)

var _ = Describe("Sts", func() {
	Context("When gathering the region", func() {
		It("Should gather the proper region for all expected region inputs", func() {
			for _, r := range aws.AwsRegions {
				err := os.Setenv("AWS_REGION", r)
				Expect(err).NotTo(HaveOccurred())

				region, err := aws.GetRegion()
				Expect(err).NotTo(HaveOccurred())
				Expect(region).To(Equal(r))
			}
		})

		It("Should error for improper regions", func() {
			improperRegions := []string{"us-west-11", "dummy-region", "cn-south-10"}
			for _, r := range improperRegions {
				err := os.Setenv("AWS_REGION", r)
				Expect(err).NotTo(HaveOccurred())

				region, err := aws.GetRegion()
				Expect(region).To(BeEmpty())
				Expect(err).To(Equal(aws.ErrRegionInvalid))
			}
		})

		It("Should error when no region is given", func() {
			err := os.Unsetenv("AWS_REGION")
			Expect(err).NotTo(HaveOccurred())
			err = os.Unsetenv("AWS_DEFAULT_REGION")
			Expect(err).NotTo(HaveOccurred())

			region, err := aws.GetRegion()
			Expect(region).To(BeEmpty())
			Expect(err).To(Equal(aws.ErrRegionNotFound))
		})
	})

	Context("When gathering the AWS URL", func() {
		It("Should gather the global region for valid global regions", func() {
			globalRegions := []string{"us-east-2", "us-east-1", "us-west-1", "us-west-2", "af-south-1", "ap-east-1", "ap-south-1"}
			for _, r := range globalRegions {
				err := os.Setenv("AWS_REGION", r)
				Expect(err).NotTo(HaveOccurred())

				url, err := aws.GetURL()
				Expect(err).NotTo(HaveOccurred())
				Expect(url).To(Equal("amazonaws.com"))
			}
		})

		It("Should return the China url for China regions", func() {
			chinaRegions := []string{"cn-north-1", "cn-northwest-1"}
			for _, r := range chinaRegions {
				err := os.Setenv("AWS_REGION", r)
				Expect(err).NotTo(HaveOccurred())

				url, err := aws.GetURL()
				Expect(err).NotTo(HaveOccurred())
				Expect(url).To(Equal("amazonaws.com.cn"))
			}
		})

		It("Should error for invalid regions", func() {
			improperRegions := []string{"us-west-11", "dummy-region", "cn-south-10"}
			for _, r := range improperRegions {
				err := os.Setenv("AWS_REGION", r)
				Expect(err).NotTo(HaveOccurred())

				url, err := aws.GetURL()
				Expect(url).To(BeEmpty())
				Expect(err).To(Equal(aws.ErrRegionInvalid))
			}
		})
	})
})
