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
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/go-github/v32/github"
	. "github.com/onsi/gomega"
	"github.com/pelletier/go-toml"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	"github.com/buildpacks/github-actions/registry"
	metadata "github.com/buildpacks/github-actions/registry/compute-metadata"
	"github.com/buildpacks/github-actions/toolkit/mocks"
)

func TestComputeMetadata(t *testing.T) {
	spec.Run(t, "compute-metadata", func(t *testing.T, context spec.G, it spec.S) {
		var (
			Expect           = NewWithT(t).Expect
			ExpectWithOffset = NewWithT(t).ExpectWithOffset

			tk = &mocks.Toolkit{}
		)

		asJSONString := func(v interface{}) string {
			b, err := json.Marshal(v)
			ExpectWithOffset(1, err).NotTo(HaveOccurred())

			return string(b)
		}

		asTOMLString := func(v interface{}) string {
			b, err := toml.Marshal(v)
			ExpectWithOffset(1, err).NotTo(HaveOccurred())

			return string(b)
		}

		it("returns error when id is invalid", func() {
			tk.On("GetInput", "issue").Return(asJSONString(github.Issue{
				Body: github.String(fmt.Sprintf("```\n%s\n```", asTOMLString(registry.IndexRequest{
					ID: "test@namespace/test-name",
				}))),
			}), true)

			Expect(metadata.ComputeMetadata(tk)).To(MatchError("::error ::invalid id test@namespace/test-name"))
		})

		it("returns error if namespace is restricted", func() {
			tk.On("GetInput", "issue").Return(asJSONString(github.Issue{
				Body: github.String(asTOMLString(registry.IndexRequest{
					ID: "cnb/test-name",
				})),
			}), true)

			Expect(metadata.ComputeMetadata(tk)).To(MatchError("::error ::restricted namespace cnb"))
		})

		it("returns error when version is invalid", func() {
			tk.On("GetInput", "issue").Return(asJSONString(github.Issue{
				Body: github.String(fmt.Sprintf("```\n%s\n```", asTOMLString(registry.IndexRequest{
					ID:      "test-namespace/test-name",
					Version: "test-version",
				}))),
			}), true)

			Expect(metadata.ComputeMetadata(tk)).To(MatchError("::error ::invalid version test-version"))
		})

		it("returns error when address is invalid", func() {
			tk.On("GetInput", "issue").Return(asJSONString(github.Issue{
				Body: github.String(asTOMLString(registry.IndexRequest{
					ID:      "test-namespace/test-name",
					Version: "0.0.0",
					Address: "host.com:443/repository/image:tag",
				})),
			}), true)

			Expect(metadata.ComputeMetadata(tk)).To(MatchError("::error ::invalid address host.com:443/repository/image:tag"))
		})

		it("computes metadata", func() {
			tk.On("GetInput", "issue").Return(asJSONString(github.Issue{
				Body: github.String(asTOMLString(registry.IndexRequest{
					ID:      "test-namespace/test-name",
					Version: "0.0.0",
					Address: "host.com:443/repository/image@sha256:133f2117e15569ca59645eddad78f4a6a675c435f9614e4b137364274f3a7614",
				})),
			}), true)
			tk.On("SetOutput", "id", "test-namespace/test-name")
			tk.On("SetOutput", "version", "0.0.0")
			tk.On("SetOutput", "address", "host.com:443/repository/image@sha256:133f2117e15569ca59645eddad78f4a6a675c435f9614e4b137364274f3a7614")
			tk.On("SetOutput", "namespace", "test-namespace")
			tk.On("SetOutput", "name", "test-name")

			Expect(metadata.ComputeMetadata(tk)).To(Succeed())
		})

	}, spec.Report(report.Terminal{}))
}
