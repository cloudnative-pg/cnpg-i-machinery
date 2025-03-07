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
		"Decode Functions",
		func(clusterJSON []byte, succeeds bool) {
			cluster, err := DecodeClusterStrict(clusterJSON)
			if !succeeds {
				Expect(err).To(HaveOccurred())
				return
			}

			Expect(cluster).NotTo(BeNil())
			Expect(cluster.GroupVersionKind()).To(Equal(getClusterGVK()))
		},
		Entry(
			"should decode valid cluster JSON",
			[]byte(`{"apiVersion":"postgresql.cnpg.io/v1","kind":"Cluster"}`),
			true,
		),
		Entry(
			"should return error for invalid cluster JSON",
			[]byte(`{"apiVersion":"v1","kind":}`),
			false,
		),
		Entry(
			"should fail when the JSON is valid but doesn't represent a Cluster",
			[]byte(`{"apiVersion":"postgresql.cnpg.io/v1","kind":"Backup"}`),
			false,
		),
	)
})
