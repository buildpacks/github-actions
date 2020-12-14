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

package index_test

import (
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"

	"github.com/buildpacks/github-actions/registry/internal/index"
)

func TestPath(t *testing.T) {
	spec.Run(t, "path", func(t *testing.T, context spec.G, it spec.S) {
		var (
			Expect = NewWithT(t).Expect
		)

		it("returns path for 1 character name", func() {
			Expect(index.Path("test-namespace", "a")).To(Equal(filepath.Join("1", "test-namespace_a")))
		})

		it("returns path for 2 character name", func() {
			Expect(index.Path("test-namespace", "ab")).To(Equal(filepath.Join("2", "test-namespace_ab")))
		})

		it("returns path for 3 character name", func() {
			Expect(index.Path("test-namespace", "abc")).To(Equal(filepath.Join("3", "ab", "test-namespace_abc")))
		})

		it("returns path for 4+ character name", func() {
			Expect(index.Path("test-namespace", "abcd")).To(Equal(filepath.Join("ab", "cd", "test-namespace_abcd")))
		})
	})
}
