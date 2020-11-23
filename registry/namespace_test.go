/*
 * Copyright 2018-2020 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package registry_test

import (
	"testing"

	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"

	"github.com/buildpacks/github-actions/registry"
)

func TestVerifyMetadata(t *testing.T) {
	spec.Run(t, "verify-metadata", func(t *testing.T, context spec.G, it spec.S) {
		var (
			Expect = NewWithT(t).Expect
		)

		it("identifies owner by predicate", func() {
			Expect(registry.IsOwner([]registry.Owner{{}}, func(_ registry.Owner) bool {
				return false
			})).To(BeFalse())

			Expect(registry.IsOwner([]registry.Owner{{}}, func(_ registry.Owner) bool {
				return true
			})).To(BeTrue())
		})

		it("identifies owner by user", func() {
			Expect(registry.ByUser(1)(registry.Owner{ID: 2, Type: registry.OrganizationType})).To(BeFalse())
			Expect(registry.ByUser(1)(registry.Owner{ID: 2, Type: registry.UserType})).To(BeFalse())
			Expect(registry.ByUser(1)(registry.Owner{ID: 1, Type: registry.OrganizationType})).To(BeFalse())
			Expect(registry.ByUser(1)(registry.Owner{ID: 1, Type: registry.UserType})).To(BeTrue())
		})

		it("identifies owner by organizations", func() {
			Expect(registry.ByOrganizations([]int64{1})(registry.Owner{ID: 2, Type: registry.UserType})).To(BeFalse())
			Expect(registry.ByOrganizations([]int64{1})(registry.Owner{ID: 2, Type: registry.OrganizationType})).To(BeFalse())
			Expect(registry.ByOrganizations([]int64{1})(registry.Owner{ID: 1, Type: registry.UserType})).To(BeFalse())
			Expect(registry.ByOrganizations([]int64{1})(registry.Owner{ID: 1, Type: registry.OrganizationType})).To(BeTrue())
		})
	})
}
