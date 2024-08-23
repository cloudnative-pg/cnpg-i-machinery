package object

import (
	"encoding/json"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	apiv1 "github.com/cloudnative-pg/cloudnative-pg/api/v1"
)

// DecodeClusterJSON decodes a JSON representation of a cluster.
func DecodeClusterJSON(clusterJSON []byte) (*apiv1.Cluster, error) {
	cluster := &apiv1.Cluster{}
	if err := DecodeObject(clusterJSON, cluster); err != nil {
		return nil, err
	}

	return cluster, nil
}

// DecodePodJSON decodes a JSON representation of a pod.
func DecodePodJSON(podJSON []byte) (*corev1.Pod, error) {
	pod := &corev1.Pod{}
	if err := DecodeObject(podJSON, pod); err != nil {
		return nil, err
	}

	return pod, nil
}

// DecodeObject decodes a JSON representation of an object.
func DecodeObject[T client.Object](objectJSON []byte, obj T) error {
	if err := json.Unmarshal(objectJSON, &obj); err != nil {
		return fmt.Errorf("error unmarshalling object JSON: %w", err)
	}
	return nil
}
