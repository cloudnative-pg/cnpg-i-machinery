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

package decoder

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func getPodGVK() schema.GroupVersionKind {
	return schema.GroupVersionKind{
		Group:   corev1.SchemeGroupVersion.Group,
		Version: corev1.SchemeGroupVersion.Version,
		Kind:    "Pod",
	}
}

// DecodePodJSON decodes a JSON representation of a pod.
func DecodePodJSON(podJSON []byte) (*corev1.Pod, error) {
	var result corev1.Pod

	if err := DecodeObject(podJSON, &result, getPodGVK()); err != nil {
		return nil, err
	}

	return &result, nil
}
