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

package common

import (
	apiv1 "github.com/cloudnative-pg/api/pkg/api/v1"
)

// Plugin represents a plugin with its associated cluster and parameters.
type Plugin struct {
	Cluster *apiv1.Cluster
	// Parameters are the configuration parameters of this plugin
	Parameters  map[string]string
	PluginIndex int
}

// NewPlugin creates a new Plugin instance for the given cluster and plugin name.
func NewPlugin(cluster apiv1.Cluster, pluginName string) *Plugin {
	result := &Plugin{Cluster: &cluster}

	result.PluginIndex = -1
	for idx, cfg := range result.Cluster.Spec.Plugins {
		if cfg.Name == pluginName {
			result.PluginIndex = idx
			result.Parameters = cfg.Parameters
		}
	}

	return result
}
