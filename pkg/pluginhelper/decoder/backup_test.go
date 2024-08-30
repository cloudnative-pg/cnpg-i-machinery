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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DecodeBackup", func() {
	DescribeTable(
		"DecodeBackup",
		func(backupJSON []byte, succeeds bool) {
			backup, err := DecodeBackup(backupJSON)
			if err != nil {
				Expect(succeeds).To(BeFalse())
				return
			}

			Expect(succeeds).To(BeTrue())
			Expect(backup.GetObjectKind().GroupVersionKind()).To(Equal(getBackupGVK()))
		},
		Entry(
			"when the backup JSON is valid",
			[]byte(`{"apiVersion":"postgresql.cnpg.io/v1","kind":"Backup"}`),
			true,
		),
		Entry(
			"when the backup JSON is valid but the Kind is wrong",
			[]byte(`{"apiVersion":"postgresql.cnpg.io/v1","kind":"Pooler"}`),
			false,
		),
		Entry(
			"when the backup JSON is valid but the object type is wrong",
			[]byte(`{"apiVersion":"apps/v1","kind":"Backup"}`),
			false,
		),
		Entry(
			"when the backup JSON is invalid",
			[]byte(`{"apiVersion":"v1","kind":}`),
			false,
		),
		Entry(
			"when the backup JSON is empty",
			[]byte(``),
			false,
		),
	)
})
