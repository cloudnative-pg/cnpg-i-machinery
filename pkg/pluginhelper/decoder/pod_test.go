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

package decoder

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Decode Functions", func() {
	DescribeTable(
		"DecodePodJSON",
		func(podJSON []byte, succeeds bool) {
			pod, err := DecodePodJSON(podJSON)
			if !succeeds {
				Expect(err).To(HaveOccurred())
				return
			}

			Expect(pod).ToNot(BeNil())
			Expect(pod.GroupVersionKind()).To(Equal(getPodGVK()))
		},
		Entry(
			"should decode valid pod JSON",
			[]byte(`{"apiVersion":"v1","kind":"Pod"}`),
			true,
		),
		Entry(
			"should return error for invalid pod JSON",
			[]byte(`{"apiVersion":"v1","kind":}`),
			false,
		),
		Entry(
			"should return error for invalid pod JSON",
			[]byte(`{"apiVersion":"v1","kind":"Node"}`),
			false,
		),
	)
})
