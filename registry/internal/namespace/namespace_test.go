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

func TestNamespace(t *testing.T) {
	spec.Run(t, "namespace", func(t *testing.T, context spec.G, it spec.S) {
		var (
			Expect = NewWithT(t).Expect
		)

		it("identifies restricted namespaces", func() {
			Expect(namespace.IsRestricted("cnb")).To(BeTrue())
			Expect(namespace.IsRestricted("test-namespace")).To(BeFalse())
		})
	})
}
