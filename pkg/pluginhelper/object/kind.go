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
	"encoding/json"
	"fmt"
)

// GetKind gets the Kubernetes object kind from its JSON representation.
func GetKind(definition []byte) (string, error) {
	var genericObject struct {
		Kind string `json:"kind"`
	}

	if err := json.Unmarshal(definition, &genericObject); err != nil {
		return "", fmt.Errorf("while unmarshalling resource definition: %w", err)
	}

	return genericObject.Kind, nil
}
