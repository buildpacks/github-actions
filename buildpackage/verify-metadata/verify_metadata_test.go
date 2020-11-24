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
	"testing"

	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/fake"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
	"github.com/stretchr/testify/mock"

	"github.com/buildpacks/github-actions/buildpackage/verify-metadata"
	"github.com/buildpacks/github-actions/buildpackage/verify-metadata/internal/mocks"
	"github.com/buildpacks/github-actions/internal/toolkit"
)

func TestVerifyMetadata(t *testing.T) {
	spec.Run(t, "verify-metadata", func(t *testing.T, context spec.G, it spec.S) {
		var (
			Expect = NewWithT(t).Expect

			f = &mocks.ImageFunction{}
			i = &fake.FakeImage{}

			tk = &toolkit.MockToolkit{}
		)

		it.Before(func() {
			tk.On("GetInput", "id").Return("test-id", true)
			tk.On("GetInput", "version").Return("test-version", true)
		})

		it("fails if address is not a digest image reference", func() {
			tk.On("GetInput", "address").Return("test-host/test-repository:test-version", true)

			Expect(metadata.VerifyMetadata(tk, f.Execute)).
				To(MatchError("::error ::address test-host/test-repository:test-version must be in digest form <host>/<repository>@sh256:<digest>"))
		})

		context("valid address", func() {
			it.Before(func() {
				tk.On("GetInput", "address").Return("host/repository@sha256:04ba2d17480910bd340f0305d846b007148dafd64bc6fc2626870c174b7c7de7", true)
				f.On("Execute", mock.Anything).Return(i, nil)
			})

			it("fails if io.buildpacks.buildpackage.metadata is not on image", func() {
				i.ConfigFileReturns(&v1.ConfigFile{
					Config: v1.Config{
						Labels: map[string]string{},
					},
				}, nil)

				Expect(metadata.VerifyMetadata(tk, f.Execute)).
					To(MatchError("::error ::unable to retrieve io.buildpacks.buildpackage.metadata label"))
			})

			it("fails if id does not match", func() {
				i.ConfigFileReturns(&v1.ConfigFile{
					Config: v1.Config{
						Labels: map[string]string{metadata.MetadataLabel: `{ "id": "another-id", "version": "test-version" }`},
					},
				}, nil)

				Expect(metadata.VerifyMetadata(tk, f.Execute)).
					To(MatchError("::error ::invalid id in buildpackage: expected test-id, found another-id"))
			})

			it("fails if version does not match", func() {
				i.ConfigFileReturns(&v1.ConfigFile{
					Config: v1.Config{
						Labels: map[string]string{metadata.MetadataLabel: `{ "id": "test-id", "version": "another-version" }`},
					},
				}, nil)

				Expect(metadata.VerifyMetadata(tk, f.Execute)).
					To(MatchError("::error ::invalid version in buildpackage: expected test-version, found another-version"))
			})

			it("passes if version and id match", func() {
				i.ConfigFileReturns(&v1.ConfigFile{
					Config: v1.Config{
						Labels: map[string]string{metadata.MetadataLabel: `{ "id": "test-id", "version": "test-version" }`},
					},
				}, nil)

				Expect(metadata.VerifyMetadata(tk, f.Execute)).To(Succeed())
			})

		})

	}, spec.Report(report.Terminal{}))
}
