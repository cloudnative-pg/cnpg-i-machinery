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

var _ = Describe("CreatePatch", func() {
	It("should create a patch for different objects", func() {
		oldObject := &corev1.Pod{
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{Name: "nginx", Image: "nginx:1.14.2"},
				},
			},
		}
		newObject := &corev1.Pod{
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{Name: "nginx", Image: "nginx:1.16.0"},
				},
			},
		}
		patch, err := CreatePatch(oldObject, newObject)
		Expect(err).NotTo(HaveOccurred())
		Expect(patch).NotTo(BeNil())
		Expect(string(patch)).To(ContainSubstring(`"op":"replace"`))
	})

	It("should return an empty patch for identical objects", func() {
		oldObject := &corev1.Pod{
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{Name: "nginx", Image: "nginx:1.14.2"},
				},
			},
		}
		newObject := &corev1.Pod{
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{Name: "nginx", Image: "nginx:1.14.2"},
				},
			},
		}
		patch, err := CreatePatch(oldObject, newObject)
		Expect(err).NotTo(HaveOccurred())
		Expect(patch).To(BeEmpty())
	})

	It("should return an error for nil oldObject", func() {
		newObject := &corev1.Pod{
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{Name: "nginx", Image: "nginx:1.16.0"},
				},
			},
		}
		patch, err := CreatePatch(nil, newObject)
		Expect(err).To(HaveOccurred())
		Expect(patch).To(BeNil())
	})

	It("should return an error for nil newObject", func() {
		oldObject := &corev1.Pod{
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{Name: "nginx", Image: "nginx:1.14.2"},
				},
			},
		}
		patch, err := CreatePatch(oldObject, nil)
		Expect(err).To(HaveOccurred())
		Expect(patch).To(BeNil())
	})
})
