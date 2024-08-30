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

package object

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Decode Functions", func() {
	DescribeTable(
		"GetKind",
		func(definition []byte, expectedKind string, succeeds bool) {
			kind, err := GetKind(definition)
			if !succeeds {
				Expect(err).To(HaveOccurred())
				return
			}

			Expect(err).NotTo(HaveOccurred())
			Expect(kind).To(Equal(expectedKind))
		},
		Entry("should get kind from valid JSON", []byte(`{"kind":"Pod"}`), "Pod", true),
		Entry("should return error for invalid JSON", []byte(`{"kind":}`), "", false),
	)
})
