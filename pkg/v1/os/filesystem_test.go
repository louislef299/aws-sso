package os_test

import (
	"encoding/json"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	pluginos "github.com/louislef299/aws-sso/pkg/v1/os"
)

var _ = Describe("Filesystem", func() {
	Context("When testing if a file exists", func() {
		It("Should succeed when creating a sample file", func() {
			filename := "tempfile.txt"
			f, err := os.Create(filename)
			Expect(err).NotTo(HaveOccurred())
			DeferCleanup(func() {
				f.Close()
				err := os.Remove(filename)
				Expect(err).NotTo(HaveOccurred())
			})

			exists, err := pluginos.IsFileOrFolderExisting(filename)
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(Equal(true))
		})

		It("Should fail when file doesn't exist", func() {
			exists, err := pluginos.IsFileOrFolderExisting("dummy-file")
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(Equal(false))
		})
	})

	Context("When writing to a file", func() {
		It("Should write to a file successfully", func() {
			filename := "dummy.txt"
			name := "Louis"
			number := 500

			type examplePayload struct {
				Name   string `json:"name"`
				Number int    `json:"number"`
			}

			err := pluginos.WriteStructToFile(examplePayload{
				Name:   name,
				Number: number,
			}, filename)
			Expect(err).NotTo(HaveOccurred())

			DeferCleanup(func() {
				err := os.Remove(filename)
				Expect(err).NotTo(HaveOccurred())
			})

			data, err := os.ReadFile(filename)
			Expect(err).NotTo(HaveOccurred())

			returnedJson := examplePayload{}
			expectedJson := examplePayload{
				Name:   name,
				Number: number,
			}
			err = json.Unmarshal(data, &returnedJson)
			Expect(err).NotTo(HaveOccurred())
			Expect(returnedJson).To(Equal(expectedJson))
		})
	})
})
