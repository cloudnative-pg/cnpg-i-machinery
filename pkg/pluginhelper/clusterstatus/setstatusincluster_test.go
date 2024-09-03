/*
Copyright The CloudNativePG Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package clusterstatus

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
		b, err := NewSetStatusInClusterResponseBuilder().JSONStatusResponse(&jsonBody)
		Expect(err).NotTo(HaveOccurred())
		Expect(b.GetJsonStatus()).To(Equal([]byte(`{"string":"test"}`)))
	})

	It("should properly form a response for a 'nil' value, allowing the plugins to do a 'noop'", func() {
		b := NewSetStatusInClusterResponseBuilder().NoOpResponse()
		Expect(b.GetJsonStatus()).To(BeNil())
	})

	It("should serialize an empty JSONStatus, allowing the plugins to reset its status", func() {
		b := NewSetStatusInClusterResponseBuilder().SetEmptyStatusResponse()
		Expect(b.GetJsonStatus()).ToNot(BeEmpty())
	})

	It("should return an error if it is an invalid JSON object", func() {
		wrongType := 4
		_, err := NewSetStatusInClusterResponseBuilder().JSONStatusResponse(&wrongType)
		Expect(err).To(HaveOccurred())
	})
})
