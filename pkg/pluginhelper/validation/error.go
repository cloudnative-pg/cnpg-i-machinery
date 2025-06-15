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

package validation

import (
	"strconv"

	"github.com/cloudnative-pg/cnpg-i/pkg/operator"

	"github.com/hh24k/cnpg-i-machinery/pkg/pluginhelper/common"
)

// BuildErrorForParameter creates a validation error for a certain plugin
// parameter.
func BuildErrorForParameter(plugin *common.Plugin, name, message string) *operator.ValidationError {
	if plugin.PluginIndex == -1 {
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
			strconv.Itoa(plugin.PluginIndex),
			name,
		},
		Message: message,
		Value:   plugin.Parameters[name],
	}
}
