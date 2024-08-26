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
	apiv1 "github.com/cloudnative-pg/cloudnative-pg/api/v1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("NewPlugin", func() {
	Context("when the plugin exists in the cluster", func() {
		It("should initialize the plugin with the correct parameters", func() {
			cluster := apiv1.Cluster{
				Spec: apiv1.ClusterSpec{
					Plugins: apiv1.PluginConfigurationList{
						{
							Name: "test-plugin",
							Parameters: map[string]string{
								"param1": "value1",
							},
						},
					},
				},
			}
			plugin := NewPlugin(cluster, "test-plugin")
			Expect(plugin.PluginIndex).To(Equal(0))
			Expect(plugin.Parameters).To(HaveKeyWithValue("param1", "value1"))
		})
	})

	Context("when the plugin does not exist in the cluster", func() {
		It("should initialize the plugin with PluginIndex set to -1", func() {
			cluster := apiv1.Cluster{
				Spec: apiv1.ClusterSpec{
					Plugins: apiv1.PluginConfigurationList{},
				},
			}
			plugin := NewPlugin(cluster, "non-existent-plugin")
			Expect(plugin.PluginIndex).To(Equal(-1))
			Expect(plugin.Parameters).To(BeNil())
		})
	})

	Context("when the cluster has multiple plugins", func() {
		It("should initialize the correct plugin based on the name", func() {
			cluster := apiv1.Cluster{
				Spec: apiv1.ClusterSpec{
					Plugins: apiv1.PluginConfigurationList{
						{
							Name: "plugin1",
							Parameters: map[string]string{
								"param1": "value1",
							},
						},
						{
							Name: "plugin2",
							Parameters: map[string]string{
								"param2": "value2",
							},
						},
					},
				},
			}
			plugin := NewPlugin(cluster, "plugin2")
			Expect(plugin.PluginIndex).To(Equal(1))
			Expect(plugin.Parameters).To(HaveKeyWithValue("param2", "value2"))
		})
	})
})
