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

package pluginhelper

import (
	"encoding/json"
	"strconv"

	apiv1 "github.com/cloudnative-pg/cloudnative-pg/api/v1"
	"github.com/cloudnative-pg/cnpg-i/pkg/operator"
	jsonpatch "github.com/evanphx/json-patch/v5"
	corev1 "k8s.io/api/core/v1"
)

const (
	pluginVolumeName = "plugins"
	pluginMountPath  = "/plugins"
)

// Data is an helper structure to be used by
// plugins wanting to enhance the CNPG validating webhooks
type Data struct {
	// Parameters are the configuration parameters of this plugin
	Parameters map[string]string

	cluster     apiv1.Cluster
	pod         corev1.Pod
	pluginIndex int
}

// NewFromCluster creates a new validation helper loading
// a cluster definition
func NewFromCluster(
	pluginName string,
	clusterDefinition []byte,
) (*Data, error) {
	result := &Data{}

	if err := json.Unmarshal(clusterDefinition, &result.cluster); err != nil {
		return nil, err
	}

	result.pluginIndex = -1
	for idx, cfg := range result.cluster.Spec.Plugins {
		if cfg.Name == pluginName {
			result.pluginIndex = idx
			result.Parameters = cfg.Parameters
		}
	}

	return result, nil
}

// NewFromClusterAndPod creates a new validation helper loading
// a cluster and a Pod definition
func NewFromClusterAndPod(
	pluginName string,
	clusterDefinition []byte,
	podDefinition []byte,
) (*Data, error) {
	result, err := NewFromCluster(pluginName, clusterDefinition)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(podDefinition, &result.pod); err != nil {
		return nil, err
	}

	return result, nil
}

// GetCluster gets the decoded cluster object
func (helper *Data) GetCluster() *apiv1.Cluster {
	return &helper.cluster
}

// GetPod gets the decoded pod object
func (helper *Data) GetPod() *corev1.Pod {
	return &helper.pod
}

// CreateClusterJSONPatch creates a JSON patch changing the cluster
// that was loaded into this helper into the
func (helper *Data) CreateClusterJSONPatch(newCluster apiv1.Cluster) ([]byte, error) {
	originalCluster, err := json.Marshal(helper.cluster)
	if err != nil {
		return nil, err
	}

	currentCluster, err := json.Marshal(newCluster)
	if err != nil {
		return nil, err
	}

	return jsonpatch.CreateMergePatch(originalCluster, currentCluster)
}

// CreatePodJSONPatch creates a JSON patch changing the cluster
// that was loaded into this helper into the
func (helper *Data) CreatePodJSONPatch(newPod corev1.Pod) ([]byte, error) {
	originalPod, err := json.Marshal(helper.pod)
	if err != nil {
		return nil, err
	}

	currentPod, err := json.Marshal(newPod)
	if err != nil {
		return nil, err
	}

	return jsonpatch.CreateMergePatch(originalPod, currentPod)
}

// ValidationErrorForParameter creates a validation error for a certain plugin
// parameter
func (helper *Data) ValidationErrorForParameter(name, message string) *operator.ValidationError {
	if helper.pluginIndex == -1 {
		return &operator.ValidationError{
			PathComponents: []string{
				"spec",
				"plugins",
				name,
			},
			Message: message,
		}
	}

	return &operator.ValidationError{
		PathComponents: []string{
			"spec",
			"plugins",
			strconv.Itoa(helper.pluginIndex),
			name,
		},
		Message: message,
		Value:   helper.Parameters[name],
	}
}

// InjectPluginVolume injects the plugin volume into a CNPG Pod
func (*Data) InjectPluginVolume(pod *corev1.Pod) {
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
		if pod.Spec.Containers[i].Name == "postgres" {
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
