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
	apiv1 "github.com/cloudnative-pg/api/pkg/api/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func getClusterGVK() schema.GroupVersionKind {
	return schema.GroupVersionKind{
		Group:   apiv1.SchemeGroupVersion.Group,
		Version: apiv1.SchemeGroupVersion.Version,
		Kind:    apiv1.ClusterKind,
	}
}

// DecodeClusterLenient decodes a JSON representation of a cluster.
func DecodeClusterLenient(clusterJSON []byte) (*apiv1.Cluster, error) {
	var result apiv1.Cluster

	if err := DecodeObjectLenient(clusterJSON, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// DecodeClusterStrict decodes a JSON representation of a cluster.
func DecodeClusterStrict(clusterJSON []byte) (*apiv1.Cluster, error) {
	var result apiv1.Cluster

	if err := DecodeObjectStrict(clusterJSON, &result, getClusterGVK()); err != nil {
		return nil, err
	}

	return &result, nil
}
