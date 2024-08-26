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

var _ = Describe("InjectPluginVolume", func() {
	Context("when the pod does not have the plugin volume", func() {
		It("should inject the plugin volume and mount", func() {
			pod := &corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name: "postgres",
						},
					},
				},
			}
			InjectPluginVolume(pod)
			Expect(pod.Spec.Volumes).To(HaveLen(1))
			Expect(pod.Spec.Volumes[0].Name).To(Equal(pluginVolumeName))
			Expect(pod.Spec.Containers[0].VolumeMounts).To(HaveLen(1))
			Expect(pod.Spec.Containers[0].VolumeMounts[0].Name).To(Equal(pluginVolumeName))
			Expect(pod.Spec.Containers[0].VolumeMounts[0].MountPath).To(Equal(pluginMountPath))
		})
	})

	Context("when the pod already has the plugin volume", func() {
		It("should not inject the plugin volume again", func() {
			pod := &corev1.Pod{
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						{
							Name: pluginVolumeName,
						},
					},
					Containers: []corev1.Container{
						{
							Name: "postgres",
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      pluginVolumeName,
									MountPath: pluginMountPath,
								},
							},
						},
					},
				},
			}
			InjectPluginVolume(pod)
			Expect(pod.Spec.Volumes).To(HaveLen(1))
			Expect(pod.Spec.Containers[0].VolumeMounts).To(HaveLen(1))
		})
	})

	Context("when the pod has multiple containers", func() {
		It("should inject the plugin volume and mount only into the postgres container", func() {
			pod := &corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name: "postgres",
						},
						{
							Name: "sidecar",
						},
					},
				},
			}
			InjectPluginVolume(pod)
			Expect(pod.Spec.Volumes).To(HaveLen(1))
			Expect(pod.Spec.Volumes[0].Name).To(Equal(pluginVolumeName))
			Expect(pod.Spec.Containers[0].VolumeMounts).To(HaveLen(1))
			Expect(pod.Spec.Containers[0].VolumeMounts[0].Name).To(Equal(pluginVolumeName))
			Expect(pod.Spec.Containers[0].VolumeMounts[0].MountPath).To(Equal(pluginMountPath))
			Expect(pod.Spec.Containers[1].VolumeMounts).To(BeEmpty())
		})
	})
})
