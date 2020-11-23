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

package metadata_test

import (
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	metadata "github.com/buildpacks/github-actions/buildpack/compute-metadata"
	"github.com/buildpacks/github-actions/toolkit/mocks"
)

func TestComputeMetadata(t *testing.T) {
	spec.Run(t, "compute-metadata", func(t *testing.T, context spec.G, it spec.S) {
		var (
			Expect = NewWithT(t).Expect

			tk = &mocks.Toolkit{}
		)

		it.Before(func() {
			tk.On("GetInput", "path").Return(filepath.Join("testdata", "buildpack.toml"), true)
		})

		it("computes metadata", func() {
			tk.On("SetOutput", "id", "test-id")
			tk.On("SetOutput", "name", "test-name")
			tk.On("SetOutput", "version", "test-version")
			tk.On("SetOutput", "homepage", "test-homepage")

			Expect(metadata.ComputeMetadata(tk)).To(Succeed())
		})
	}, spec.Report(report.Terminal{}))
}
