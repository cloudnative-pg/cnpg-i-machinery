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
	"encoding/json"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// WrongObjectTypeError is raised when the GVK of the passed JSON
// object is different from the expected one.
type WrongObjectTypeError struct {
	expectedGVK schema.GroupVersionKind
	receivedGVK schema.GroupVersionKind
}

// Error implements the error interface.
func (e *WrongObjectTypeError) Error() string {
	return fmt.Sprintf("received wrong GVK '%v' expected '%v'", e.receivedGVK.String(), e.expectedGVK.String())
}

// DecodeObject decodes a JSON representation of an object.
func DecodeObject(objectJSON []byte, object client.Object, expectedGVK schema.GroupVersionKind) error {
	if err := json.Unmarshal(objectJSON, object); err != nil {
		return fmt.Errorf("error unmarshalling object JSON: %w", err)
	}

	if object.GetObjectKind().GroupVersionKind() != expectedGVK {
		return &WrongObjectTypeError{
			expectedGVK: expectedGVK,
			receivedGVK: object.GetObjectKind().GroupVersionKind(),
		}
	}

	return nil
}
