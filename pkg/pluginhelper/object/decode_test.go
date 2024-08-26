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
	corev1 "k8s.io/api/core/v1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Decode Functions", func() {
	Context("DecodeClusterJSON", func() {
		It("should decode valid cluster JSON", func() {
			clusterJSON := []byte(`{"apiVersion":"v1","kind":"Cluster"}`)
			cluster, err := DecodeClusterJSON(clusterJSON)
			Expect(err).NotTo(HaveOccurred())
			Expect(cluster).NotTo(BeNil())
			Expect(cluster.Kind).To(Equal("Cluster"))
		})

		It("should return error for invalid cluster JSON", func() {
			clusterJSON := []byte(`{"apiVersion":"v1","kind":}`)
			cluster, err := DecodeClusterJSON(clusterJSON)
			Expect(err).To(HaveOccurred())
			Expect(cluster).To(BeNil())
		})
	})

	Context("DecodePodJSON", func() {
		It("should decode valid pod JSON", func() {
			podJSON := []byte(`{"apiVersion":"v1","kind":"Pod"}`)
			pod, err := DecodePodJSON(podJSON)
			Expect(err).NotTo(HaveOccurred())
			Expect(pod).NotTo(BeNil())
			Expect(pod.Kind).To(Equal("Pod"))
		})

		It("should return error for invalid pod JSON", func() {
			podJSON := []byte(`{"apiVersion":"v1","kind":}`)
			pod, err := DecodePodJSON(podJSON)
			Expect(err).To(HaveOccurred())
			Expect(pod).To(BeNil())
		})
	})

	Context("DecodeObject", func() {
		It("should decode valid object JSON", func() {
			objectJSON := []byte(`{"apiVersion":"v1","kind":"Pod"}`)
			var pod corev1.Pod
			err := DecodeObject(objectJSON, &pod)
			Expect(err).NotTo(HaveOccurred())
			Expect(pod.Kind).To(Equal("Pod"))
		})

		It("should return error for invalid object JSON", func() {
			objectJSON := []byte(`{"apiVersion":"v1","kind":}`)
			var pod corev1.Pod
			err := DecodeObject(objectJSON, &pod)
			Expect(err).To(HaveOccurred())
		})
	})
})
