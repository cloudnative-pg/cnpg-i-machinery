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

// ErrNoPostgresContainerFound is raised when there's no PostgreSQL container
// in the passed instance Pod.
var ErrNoPostgresContainerFound = errors.New("no postgres container into instance Pod")

// InjectPluginVolume injects the plugin volume into a CNPG Pod.
func InjectPluginVolume(pod *corev1.Pod) {
	foundPluginVolume := false
	for i := range pod.Spec.Volumes {
		if pod.Spec.Volumes[i].Name == pluginVolumeName {
			foundPluginVolume = true
		}
	}

	if foundPluginVolume {
		return
	}

	pod.Spec.Volumes = append(pod.Spec.Volumes, corev1.Volume{
		Name: pluginVolumeName,
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		},
	})

	for i := range pod.Spec.Containers {
		if pod.Spec.Containers[i].Name == postgresContainerName {
			pod.Spec.Containers[i].VolumeMounts = append(
				pod.Spec.Containers[i].VolumeMounts,
				corev1.VolumeMount{
					Name:      pluginVolumeName,
					MountPath: pluginMountPath,
				},
			)
		}
	}
}

// InjectPluginSidecar injects a plugin sidecar into a CNPG Pod.
//
// If the "injectPostgresVolumeMount" flag is true, this will append all the volume
// mounts that are used in the instance manager Pod to the passed sidecar
// container, granting it superuser access to the PostgreSQL instance.
//
// Besides the value of "injectPostgresVolumeMount", the plugin volume
// will always be injected in the PostgreSQL container.
func InjectPluginSidecar(pod *corev1.Pod, sidecar *corev1.Container, injectPostgresEnvironment bool) error {
	sidecar = sidecar.DeepCopy()
	InjectPluginVolume(pod)

	var volumeMounts []corev1.VolumeMount
	var envs []corev1.EnvVar
	sidecarContainerFound := false
	postgresContainerFound := false
	for i := range pod.Spec.Containers {
		if pod.Spec.Containers[i].Name == postgresContainerName {
			volumeMounts = pod.Spec.Containers[i].VolumeMounts
			envs = pod.Spec.Containers[i].Env
			postgresContainerFound = true
		} else if pod.Spec.Containers[i].Name == sidecar.Name {
			sidecarContainerFound = true
		}
	}

	if sidecarContainerFound {
		// The sidecar container was already added
		return nil
	}

	if !postgresContainerFound {
		return ErrNoPostgresContainerFound
	}

	// Do not modify the passed sidecar definition
	if injectPostgresEnvironment {
		sidecar.VolumeMounts = append(sidecar.VolumeMounts, volumeMounts...)
		sidecar.Env = append(sidecar.Env, envs...)
	}
	pod.Spec.Containers = append(pod.Spec.Containers, *sidecar)

	return nil
}
