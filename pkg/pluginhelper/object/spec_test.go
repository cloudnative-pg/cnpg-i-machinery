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

var _ = Describe("InjectIntoPostgresPluginVolume", func() {
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
			InjectIntoPostgresPluginVolume(pod)
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
			InjectIntoPostgresPluginVolume(pod)
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
			InjectIntoPostgresPluginVolume(pod)
			Expect(pod.Spec.Volumes).To(HaveLen(1))
			Expect(pod.Spec.Volumes[0].Name).To(Equal(pluginVolumeName))
			Expect(pod.Spec.Containers[0].VolumeMounts).To(HaveLen(1))
			Expect(pod.Spec.Containers[0].VolumeMounts[0].Name).To(Equal(pluginVolumeName))
			Expect(pod.Spec.Containers[0].VolumeMounts[0].MountPath).To(Equal(pluginMountPath))
			Expect(pod.Spec.Containers[1].VolumeMounts).To(BeEmpty())
		})
	})
})

var _ = Describe("InjectSidecarIntoPostgres", func() {
	var sidecar *corev1.Container

	BeforeEach(func() {
		sidecar = &corev1.Container{
			Name: "pluginname",
		}
	})

	When("when the passed Pod have no 'postgres' container", func() {
		pod := &corev1.Pod{
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name: "sidecar",
					},
				},
			},
		}

		It("will fail if we need to inject the volume mounts", func() {
			err := InjectSidecarIntoPostgres(pod, sidecar, true)
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(ErrNoMainContainerFound))
		})

		It("will fail if we don't need to inject the volume mounts", func() {
			err := InjectSidecarIntoPostgres(pod, sidecar, false)
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(ErrNoMainContainerFound))
		})
	})

	When("the passed Pod have a 'postgres' container", func() {
		var pod *corev1.Pod

		BeforeEach(func() {
			pod = &corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name: postgresContainerName,
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "pgdata",
									MountPath: "/pgdata",
								},
							},
						},
					},
				},
			}
		})

		When("the PG volume mounts injection is requested", func() {
			It("it will inherit the volume mounts and the plugin volume", func() {
				err := InjectSidecarIntoPostgres(pod, sidecar, true)
				Expect(err).ToNot(HaveOccurred())
				Expect(pod.Spec.Containers).To(HaveLen(2))
				Expect(pod.Spec.Containers[1].Name).To(Equal(sidecar.Name))

				// the plugin volume have been injected
				Expect(pod.Spec.Containers[0].VolumeMounts).To(HaveLen(2))

				// even in the sidecar
				Expect(pod.Spec.Containers[1].VolumeMounts).To(HaveLen(2))
			})
		})

		When("the PG volume mounts is set to not be inherited", func() {
			It("it will not inherit the volume mounts", func() {
				err := InjectSidecarIntoPostgres(pod, sidecar, false)
				Expect(err).ToNot(HaveOccurred())
				Expect(pod.Spec.Containers).To(HaveLen(2))
				Expect(pod.Spec.Containers[0].Name).To(Equal(postgresContainerName))
				Expect(pod.Spec.Containers[1].Name).To(Equal(sidecar.Name))

				// the plugin volume have been injected
				Expect(pod.Spec.Containers[0].VolumeMounts).To(HaveLen(2))

				// even in the sidecar
				Expect(pod.Spec.Containers[1].VolumeMounts).To(BeEmpty())
			})
		})
	})
})
