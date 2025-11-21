package os_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	pluginos "github.com/louislef299/knot/pkg/os"
)

var _ = Describe("Encrypt", func() {
	Context("When testing plugin encryptions", func() {
		It("Should be able to decode base64", func() {
			secret := "password"
			encoded := pluginos.Encode(secret)

			r, err := pluginos.Decode(encoded)
			Expect(err).NotTo(HaveOccurred())

			Expect(r).To(Equal(secret))
		})
	})
})
