/*
 * Copyright 2020 the original author or authors.
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

package info_test

import (
	"bytes"
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	"github.com/buildpacks/github-actions/info"
)

func TestBuildpackInfo(t *testing.T) {
	spec.Run(t, "buildpack-info", func(t *testing.T, when spec.G, it spec.S) {
		var (
			Expect = NewWithT(t).Expect

			b = &bytes.Buffer{}

			i = info.BuildpackInfo{
				Path:   filepath.Join("testdata", "buildpack.toml"),
				Writer: b,
			}
		)

		it("informs", func() {
			Expect(i.Inform()).To(Succeed())

			Expect(b.String()).To(Equal(`::set-output name=id::test-id
::set-output name=name::test-name
::set-output name=version::test-version
::set-output name=homepage::test-homepage
`,
			))
		})
	}, spec.Report(report.Terminal{}))
}
