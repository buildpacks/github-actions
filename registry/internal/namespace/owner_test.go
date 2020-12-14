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

package namespace_test

import (
	"testing"

	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"

	"github.com/buildpacks/github-actions/registry/internal/namespace"
)

func TestOwner(t *testing.T) {
	spec.Run(t, "owner", func(t *testing.T, context spec.G, it spec.S) {
		var (
			Expect = NewWithT(t).Expect
		)

		it("identifies owner by predicate", func() {
			Expect(namespace.IsOwner([]namespace.Owner{{}}, func(_ namespace.Owner) bool {
				return false
			})).To(BeFalse())

			Expect(namespace.IsOwner([]namespace.Owner{{}}, func(_ namespace.Owner) bool {
				return true
			})).To(BeTrue())
		})

		it("identifies owner by user", func() {
			Expect(namespace.ByUser(1)(namespace.Owner{ID: 2, Type: namespace.OrganizationType})).To(BeFalse())
			Expect(namespace.ByUser(1)(namespace.Owner{ID: 2, Type: namespace.UserType})).To(BeFalse())
			Expect(namespace.ByUser(1)(namespace.Owner{ID: 1, Type: namespace.OrganizationType})).To(BeFalse())
			Expect(namespace.ByUser(1)(namespace.Owner{ID: 1, Type: namespace.UserType})).To(BeTrue())
		})

		it("identifies owner by organizations", func() {
			Expect(namespace.ByOrganizations([]int64{1})(namespace.Owner{ID: 2, Type: namespace.UserType})).To(BeFalse())
			Expect(namespace.ByOrganizations([]int64{1})(namespace.Owner{ID: 2, Type: namespace.OrganizationType})).To(BeFalse())
			Expect(namespace.ByOrganizations([]int64{1})(namespace.Owner{ID: 1, Type: namespace.UserType})).To(BeFalse())
			Expect(namespace.ByOrganizations([]int64{1})(namespace.Owner{ID: 1, Type: namespace.OrganizationType})).To(BeTrue())
		})
	})
}
