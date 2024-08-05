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
	"fmt"
	"strconv"

	apiv1 "github.com/cloudnative-pg/cloudnative-pg/api/v1"
	"github.com/cloudnative-pg/cnpg-i/pkg/operator"
	"github.com/snorwin/jsonpatch"
	corev1 "k8s.io/api/core/v1"
)

const (
	pluginVolumeName = "plugins"
	pluginMountPath  = "/plugins"
)

// Data is an helper structure to be used by
// plugins wanting to enhance the CNPG validating webhooks.
type Data struct {
	// Parameters are the configuration parameters of this plugin
	Parameters map[string]string

	cluster     apiv1.Cluster
	pod         corev1.Pod
	pluginIndex int
}

// DataBuilder a fluent constructor for the Data struct.
type DataBuilder struct {
	pluginName  string
	clusterJSON []byte
	podJSON     []byte
}

// NewDataBuilder initializes a basic DataBuilder.
func NewDataBuilder(pluginName string, clusterJSON []byte) *DataBuilder {
	d := DataBuilder{clusterJSON: clusterJSON, pluginName: pluginName}
	d.clusterJSON = clusterJSON
	return &d
}

// WithPod adds Pod data to the DataBuilder.
func (d *DataBuilder) WithPod(podJSON []byte) *DataBuilder {
	d.podJSON = podJSON
	return d
}

// Build returns the constructed Data object and any errors encountered.
func (d *DataBuilder) Build() (*Data, error) {
	result := &Data{}

	if err := json.Unmarshal(d.clusterJSON, &result.cluster); err != nil {
		return nil, fmt.Errorf("error unmarshalling cluster JSON: %w", err)
	}

	if len(d.podJSON) > 0 {
		if err := json.Unmarshal(d.podJSON, &result.pod); err != nil {
			return nil, fmt.Errorf("error unmarshalling pod JSON: %w", err)
		}
	}

	result.pluginIndex = -1
	for idx, cfg := range result.cluster.Spec.Plugins {
		if cfg.Name == d.pluginName {
			result.pluginIndex = idx
			result.Parameters = cfg.Parameters
		}
	}

	return result, nil
}

// GetCluster gets the decoded cluster object.
func (helper *Data) GetCluster() *apiv1.Cluster {
	return &helper.cluster
}

// GetPod gets the decoded pod object.
func (helper *Data) GetPod() *corev1.Pod {
	return &helper.pod
}

// CreateClusterJSONPatch creates a JSON patch changing the cluster
// that was loaded into this helper into the cluster.
func (helper *Data) CreateClusterJSONPatch(newCluster apiv1.Cluster) ([]byte, error) {
	patch, err := jsonpatch.CreateJSONPatch(newCluster, helper.cluster)
	return []byte(patch.String()), err
}

// CreatePodJSONPatch creates a JSON patch changing the cluster
// that was loaded into this helper into the pod.
func (helper *Data) CreatePodJSONPatch(newPod corev1.Pod) ([]byte, error) {
	patch, err := jsonpatch.CreateJSONPatch(newPod, helper.pod)
	return []byte(patch.String()), err
}

// ValidationErrorForParameter creates a validation error for a certain plugin
// parameter.
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

// InjectPluginVolume injects the plugin volume into a CNPG Pod.
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

// DecodeBackup decodes a JSON representation of a backup.
func (*Data) DecodeBackup(backupDefinition []byte) (*apiv1.Backup, error) {
	var backup apiv1.Backup

	if err := json.Unmarshal(backupDefinition, &backup); err != nil {
		return nil, fmt.Errorf("error unmarshalling backup JSON: %w", err)
	}

	return &backup, nil
}

// GetKind gets the Kubernetes object kind from its JSON representation
func GetKind(definition []byte) (string, error) {
	var genericObject struct {
		Kind string `json:"kind"`
	}

	if err := json.Unmarshal(definition, &genericObject); err != nil {
		return "", err
	}

	return genericObject.Kind, nil
}
