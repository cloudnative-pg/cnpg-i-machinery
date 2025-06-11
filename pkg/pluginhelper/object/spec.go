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
	"errors"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/ptr"
)

const (
	pluginVolumeName = "plugins"
	pluginMountPath  = "/plugins"

	postgresContainerName = "postgres"
)

var (
	// ErrNoPostgresContainerFound is raised when there's no PostgreSQL container
	// in the passed instance Pod.
	ErrNoPostgresContainerFound = errors.New("no postgres container into instance Pod")

	// ErrNilPodPassed is raised when a nil Pod is passed to a function requiring it.
	ErrNilPodPassed = errors.New("nil pod passed")
)

// InjectPluginVolume injects the plugin volume into a CNPG Pod.
func InjectPluginVolume(pod *corev1.Pod) {
	InjectPluginVolumeSpec(&pod.Spec)
}

// InjectPluginVolumeSpec injects the plugin volume into a CNPG Pod spec.
func InjectPluginVolumeSpec(spec *corev1.PodSpec) {
	foundPluginVolume := false
	for i := range spec.Volumes {
		if spec.Volumes[i].Name == pluginVolumeName {
			foundPluginVolume = true
		}
	}

	if foundPluginVolume {
		return
	}

	spec.Volumes = append(spec.Volumes, corev1.Volume{
		Name: pluginVolumeName,
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		},
	})

	for i := range spec.Containers {
		if spec.Containers[i].Name == postgresContainerName {
			spec.Containers[i].VolumeMounts = append(
				spec.Containers[i].VolumeMounts,
				corev1.VolumeMount{
					Name:      pluginVolumeName,
					MountPath: pluginMountPath,
				},
			)
		}
	}
}

// InjectPluginSidecar refer to InjectPluginSidecarSpec.
//
// Deprecated: Kubernetes versions >= 1.29 support sidecars as InitContainers by default,
// so this function should not be used anymore. Use InjectPluginInitContainerSidecarSpec instead.
func InjectPluginSidecar(pod *corev1.Pod, sidecar *corev1.Container, injectPostgresVolumeMounts bool) error {
	if pod == nil {
		return ErrNilPodPassed
	}

	return InjectPluginSidecarSpec(&pod.Spec, sidecar, injectPostgresVolumeMounts)
}

// InjectPluginSidecarInitContainer refer to InjectPluginSidecarInitContainerSpec.
func InjectPluginSidecarInitContainer(pod *corev1.Pod,
	sidecar *corev1.Container,
	injectPostgresVolumeMounts bool,
) error {
	if pod == nil {
		return ErrNilPodPassed
	}

	return InjectPluginInitContainerSidecarSpec(&pod.Spec, sidecar, injectPostgresVolumeMounts)
}

// InjectPluginSidecarSpec injects a plugin sidecar into a CNPG Pod spec.
//
// If the "injectPostgresVolumeMount" flag is true, this will append all the volume
// mounts that are used in the instance manager Pod to the passed sidecar
// container, granting it superuser access to the PostgreSQL instance.
//
// Besides the value of "injectPostgresVolumeMount", the plugin volume
// will always be injected in the PostgreSQL container.
//
// Deprecated: Kubernetes versions >= 1.29 support sidecars as InitContainers by default,
// so this function should not be used anymore. Use InjectPluginInitContainerSidecarSpec instead.
func InjectPluginSidecarSpec(spec *corev1.PodSpec, sidecar *corev1.Container, injectPostgresVolumeMounts bool) error {
	fetcher := func(spec *corev1.PodSpec, _ *corev1.Container) *[]corev1.Container {
		return &spec.Containers
	}

	return injectSidecar(spec, sidecar, injectPostgresVolumeMounts, fetcher)
}

// InjectPluginInitContainerSidecarSpec injects a plugin sidecar into a CNPG Pod spec as
// an InitContainer. This requires the SidecarContainers feature gate to be enabled, which is
// the default for kubernetes versions >= 1.29.
//
// If the "injectPostgresVolumeMount" flag is true, this will append all the volume
// mounts that are used in the instance manager Pod to the passed sidecar
// container, granting it superuser access to the PostgreSQL instance.
//
// Besides the value of "injectPostgresVolumeMount", the plugin volume
// will always be injected in the PostgreSQL container.
func InjectPluginInitContainerSidecarSpec(
	spec *corev1.PodSpec, sidecar *corev1.Container, injectPostgresVolumeMounts bool,
) error {
	fetcher := func(spec *corev1.PodSpec, sidecar *corev1.Container) *[]corev1.Container {
		// ensure the sidecar has the correct restartPolicy
		sidecar.RestartPolicy = ptr.To(corev1.ContainerRestartPolicyAlways)

		if spec.InitContainers == nil {
			spec.InitContainers = []corev1.Container{}
		}

		return &spec.InitContainers
	}

	return injectSidecar(spec, sidecar, injectPostgresVolumeMounts, fetcher)
}

func injectSidecar(
	spec *corev1.PodSpec,
	sidecar *corev1.Container,
	injectPostgresVolumeMounts bool,
	containerFetcher func(spec *corev1.PodSpec, sidecar *corev1.Container) *[]corev1.Container,
) error {
	if spec == nil || sidecar == nil {
		return nil
	}

	modifiedSidecar := sidecar.DeepCopy()
	InjectPluginVolumeSpec(spec)

	// Find PostgreSQL container and its volume mounts
	volumeMounts, err := getPostgresVolumeMounts(spec)
	if err != nil {
		return err
	}

	targetContainers := containerFetcher(spec, modifiedSidecar)

	for _, container := range *targetContainers {
		if container.Name == modifiedSidecar.Name {
			return nil
		}
	}

	if injectPostgresVolumeMounts {
		modifiedSidecar.VolumeMounts = append(modifiedSidecar.VolumeMounts, volumeMounts...)
	}

	*targetContainers = append(*targetContainers, *modifiedSidecar)

	return nil
}

func getPostgresVolumeMounts(spec *corev1.PodSpec) ([]corev1.VolumeMount, error) {
	for _, container := range spec.Containers {
		if container.Name == postgresContainerName {
			return container.VolumeMounts, nil
		}
	}

	return nil, ErrNoPostgresContainerFound
}
