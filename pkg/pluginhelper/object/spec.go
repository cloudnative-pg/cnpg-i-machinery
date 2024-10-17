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
)

const (
	pluginVolumeName = "plugins"
	pluginMountPath  = "/plugins"

	postgresContainerName = "postgres"
)

// ErrNoMainContainerFound is raised when there's no main container.
var ErrNoMainContainerFound = errors.New("no main container found into the Pod")

// InjectIntoPostgresPluginVolume injects the plugin volume into a CNPG Pod.
func InjectIntoPostgresPluginVolume(pod *corev1.Pod) {
	InjectPluginVolumePodSpec(&pod.Spec, postgresContainerName)
}

// InjectPluginVolumePodSpec injects the plugin volume into a CNPG Pod spec.
func InjectPluginVolumePodSpec(spec *corev1.PodSpec, mainContainerName string) {
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
		if spec.Containers[i].Name == mainContainerName {
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

// InjectSidecarIntoPostgres refer to InjectPluginSidecarPodSpec.
func InjectSidecarIntoPostgres(pod *corev1.Pod, sidecar *corev1.Container, injectPostgresVolumeMounts bool) error {
	return InjectPluginSidecarPodSpec(&pod.Spec, sidecar, postgresContainerName, injectPostgresVolumeMounts)
}

// InjectPluginSidecarPodSpec injects a plugin sidecar into a CNPG Pod spec.
//
// If the "injectMainContainerVolumes" flag is true, this will append all the volume
// mounts that are used in the instance manager Pod to the passed sidecar
// container, granting it superuser access to the PostgreSQL instance.
func InjectPluginSidecarPodSpec(
	spec *corev1.PodSpec,
	sidecar *corev1.Container,
	mainContainerName string,
	injectMainContainerVolumes bool,
) error {
	sidecar = sidecar.DeepCopy()
	InjectPluginVolumePodSpec(spec, mainContainerName)

	var volumeMounts []corev1.VolumeMount
	sidecarContainerFound := false
	mainContainerFound := false
	for i := range spec.Containers {
		if spec.Containers[i].Name == mainContainerName {
			volumeMounts = spec.Containers[i].VolumeMounts
			mainContainerFound = true
		} else if spec.Containers[i].Name == sidecar.Name {
			sidecarContainerFound = true
		}
	}

	if sidecarContainerFound {
		// The sidecar container was already added
		return nil
	}

	if !mainContainerFound {
		return ErrNoMainContainerFound
	}

	// Do not modify the passed sidecar definition
	if injectMainContainerVolumes {
		sidecar.VolumeMounts = append(sidecar.VolumeMounts, volumeMounts...)
	}
	spec.Containers = append(spec.Containers, *sidecar)

	return nil
}
