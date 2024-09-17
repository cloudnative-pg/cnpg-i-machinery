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
	"k8s.io/apimachinery/pkg/runtime"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CreatePatch", func() {
	DescribeTable("should create patches for different objects",
		func(newObject, oldObject runtime.Object, expectedPatch []byte, succeeds bool) {
			patch, err := CreatePatch(newObject, oldObject)
			if succeeds {
				Expect(err).NotTo(HaveOccurred())
			} else {
				Expect(err).To(HaveOccurred())
			}
			Expect(patch).To(Equal(expectedPatch))
		},
		Entry(
			"valid patch on Pod",
			&corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{Name: "nginx", Image: "nginx:1.16.0"},
					},
				},
			},
			&corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{Name: "nginx", Image: "nginx:1.14.2"},
					},
				},
			},
			[]byte(`[{"op":"replace","path":"/spec/containers/0/image","value":"nginx:1.16.0"}]`),
			true,
		),
		Entry(
			"should return an empty patch for identical objects",
			&corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{Name: "nginx", Image: "nginx:1.14.2"},
					},
				},
			},
			&corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{Name: "nginx", Image: "nginx:1.14.2"},
					},
				},
			},
			[]byte(``),
			true,
		),
		Entry(
			"should return an error for nil oldObject",
			&corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{Name: "nginx", Image: "nginx:1.16.0"},
					},
				},
			},
			nil,
			nil,
			false,
		),
		Entry(
			"should return an error for nil newObject",
			nil,
			&corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{Name: "nginx", Image: "nginx:1.16.0"},
					},
				},
			},
			nil,
			false,
		),
	)
})
