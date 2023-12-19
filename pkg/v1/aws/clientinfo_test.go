package aws_test

import (
	"encoding/json"
	"os"
	"path"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/louislef299/aws-sso/pkg/v1/aws"
)

var _ = Describe("Clientinfo", func() {
	Context("When reading client information", func() {
		It("Should work with a local test and not be expired", func() {
			testClient := aws.ClientInformation{
				AccessTokenExpiresAt:    time.Now().Add(time.Hour),
				AccessToken:             "dummy",
				ClientId:                "123456",
				ClientSecret:            "S3cr3t!",
				ClientSecretExpiresAt:   "tomorrow",
				DeviceCode:              "devicecode",
				VerificationUriComplete: "yup",
				StartUrl:                "start.com",
			}

			data, err := json.Marshal(testClient)
			Expect(err).NotTo(HaveOccurred())
			err = os.WriteFile(aws.GetAccessToken(), data, 0744)
			Expect(err).NotTo(HaveOccurred())

			DeferCleanup(func() {
				err := os.Remove(aws.GetAccessToken())
				Expect(err).NotTo(HaveOccurred())
			})

			pwd, err := os.Getwd()
			Expect(err).NotTo(HaveOccurred())
			clientInfo, err := aws.ReadClientInformation(path.Join(pwd, aws.GetAccessToken()))
			Expect(err).NotTo(HaveOccurred())

			Expect(areClientsEqual(*clientInfo, testClient)).To(Equal(true))
			Expect(clientInfo.IsExpired()).To(Equal(false))
		})

		It("Should have expired when access token expiry is an hour earlier", func() {
			testClient := aws.ClientInformation{
				AccessTokenExpiresAt:    time.Now().Add(-time.Hour),
				AccessToken:             "dummy",
				ClientId:                "123456",
				ClientSecret:            "S3cr3t!",
				ClientSecretExpiresAt:   "tomorrow",
				DeviceCode:              "devicecode",
				VerificationUriComplete: "yup",
				StartUrl:                "start.com",
			}

			data, err := json.Marshal(testClient)
			Expect(err).NotTo(HaveOccurred())
			err = os.WriteFile(aws.GetAccessToken(), data, 0744)
			Expect(err).NotTo(HaveOccurred())

			DeferCleanup(func() {
				err := os.Remove(aws.GetAccessToken())
				Expect(err).NotTo(HaveOccurred())
			})

			pwd, err := os.Getwd()
			Expect(err).NotTo(HaveOccurred())
			clientInfo, err := aws.ReadClientInformation(path.Join(pwd, aws.GetAccessToken()))
			Expect(err).NotTo(HaveOccurred())

			Expect(clientInfo.IsExpired()).To(Equal(true))
		})
	})
})

func areClientsEqual(client1, client2 aws.ClientInformation) bool {
	if !client1.AccessTokenExpiresAt.Equal(client2.AccessTokenExpiresAt) {
		return false
	}
	if strings.Compare(client1.AccessToken, client2.AccessToken) != 0 {
		return false
	}
	if strings.Compare(client1.ClientId, client2.ClientId) != 0 {
		return false
	}
	if strings.Compare(client1.ClientSecret, client2.ClientSecret) != 0 {
		return false
	}
	if strings.Compare(client1.ClientSecretExpiresAt, client2.ClientSecretExpiresAt) != 0 {
		return false
	}
	if strings.Compare(client1.DeviceCode, client2.DeviceCode) != 0 {
		return false
	}
	if strings.Compare(client1.VerificationUriComplete, client2.VerificationUriComplete) != 0 {
		return false
	}
	if strings.Compare(client1.StartUrl, client2.StartUrl) != 0 {
		return false
	}
	return true
}
