package pluginhelper

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("BuildSetStatusResponse", func() {
	type test struct {
		Name string `json:"string"`
	}

	It("should properly form a response with an object, allowing the plugins to set the status", func() {
		jsonBody := test{Name: "test"}
		b, err := NewSetStatusResponseBuilder().JsonStatusResponse(&jsonBody)
		Expect(err).NotTo(HaveOccurred())
		Expect(b.JsonStatus).ToNot(BeEmpty())
	})

	It("should properly form a response for a 'nil' value, allowing the plugins to do a 'noop'", func() {
		b := NewSetStatusResponseBuilder().NoOpResponse()
		Expect(b.JsonStatus).To(BeNil())
	})

	It("should serialize an empty JSONStatus, allowing the plugins to reset its status", func() {
		b := NewSetStatusResponseBuilder().SetEmptyStatusResponse()
		Expect(b.JsonStatus).ToNot(BeEmpty())
	})

	It("should return an error if it is an invalid JSON object", func() {
		wrongType := 4
		_, err := NewSetStatusResponseBuilder().JsonStatusResponse(&wrongType)
		Expect(err).To(HaveOccurred())
	})
})
