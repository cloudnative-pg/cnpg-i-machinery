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

package backup

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DecodeBackup", func() {
	Context("when the backup JSON is valid", func() {
		It("should decode the backup JSON successfully", func() {
			backupJSON := []byte(`{"apiVersion":"v1","kind":"Backup"}`)
			backup, err := DecodeBackup(backupJSON)
			Expect(err).NotTo(HaveOccurred())
			Expect(backup).NotTo(BeNil())
			Expect(backup.Kind).To(Equal("Backup"))
		})
	})

	Context("when the backup JSON is invalid", func() {
		It("should return an error for invalid JSON", func() {
			backupJSON := []byte(`{"apiVersion":"v1","kind":}`)
			backup, err := DecodeBackup(backupJSON)
			Expect(err).To(HaveOccurred())
			Expect(backup).To(BeNil())
		})
	})

	Context("when the backup JSON is empty", func() {
		It("should return an error for empty JSON", func() {
			backupJSON := []byte(``)
			backup, err := DecodeBackup(backupJSON)
			Expect(err).To(HaveOccurred())
			Expect(backup).To(BeNil())
		})
	})
})
