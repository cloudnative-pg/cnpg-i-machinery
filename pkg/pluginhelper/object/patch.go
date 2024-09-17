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
	"fmt"

	"github.com/snorwin/jsonpatch"
	"k8s.io/apimachinery/pkg/runtime"
)

// CreatePatch creates a JSON patch from the diff between the old and new object.
func CreatePatch(newObject, oldObject runtime.Object) ([]byte, error) {
	ptc, err := jsonpatch.CreateJSONPatch(newObject, oldObject)
	if err != nil {
		return nil, fmt.Errorf("while creating JSON patch: %w", err)
	}

	return []byte(ptc.String()), nil
}
