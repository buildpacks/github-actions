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

package entry_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/google/go-github/v32/github"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
	"github.com/stretchr/testify/mock"
	"gopkg.in/retry.v1"

	"github.com/buildpacks/github-actions/internal/toolkit"
	"github.com/buildpacks/github-actions/registry/internal/index"
	"github.com/buildpacks/github-actions/registry/internal/services"
	entry "github.com/buildpacks/github-actions/registry/yank-entry"
)

func TestYankEntry(t *testing.T) {
	spec.Run(t, "yank-entry", func(t *testing.T, context spec.G, it spec.S) {
		var (
			Expect           = NewWithT(t).Expect
			ExpectWithOffset = NewWithT(t).ExpectWithOffset

			r     = &services.MockRepositoriesService{}
			rOpts *github.RepositoryContentGetOptions
			s     = retry.LimitCount(2, retry.Regular{Min: 2})
			tk    = &toolkit.MockToolkit{}
		)

		asJSONString := func(v interface{}) string {
			b, err := json.Marshal(v)
			ExpectWithOffset(1, err).NotTo(HaveOccurred())

			return string(b)
		}

		it.Before(func() {
			tk.On("GetInput", "owner").Return("test-owner", true)
			tk.On("GetInput", "repository").Return("test-repository", true)
			tk.On("GetInput", "namespace").Return("test-namespace", true)
			tk.On("GetInput", "name").Return("test-name", true)
			tk.On("GetInput", "version").Return("test-version", true)

		})

		context("index does not exist", func() {
			it.Before(func() {
				r.On("GetContents", mock.Anything, "test-owner", "test-repository", filepath.Join("te", "st", "test-namespace_test-name"), rOpts).
					Return(nil, nil, &github.Response{Response: &http.Response{StatusCode: http.StatusNotFound}}, nil)
			})

			it("fails if index does not exist", func() {
				Expect(entry.YankEntry(tk, r, s)).
					To(MatchError("::error ::index test-name does not exist"))
			})
		})

		context("index does exist", func() {

			it("fails if version does not exist", func() {
				r.On("GetContents", mock.Anything, "test-owner", "test-repository", filepath.Join("te", "st", "test-namespace_test-name"), rOpts).
					Return(&github.RepositoryContent{
						Content: github.String(asJSONString(index.Entry{
							Namespace: "test-namespace",
							Name:      "test-name",
							Version:   "another-version",
							Address:   "test-address",
						})),
						SHA: github.String("test-sha"),
					}, nil, nil, nil)

				Expect(entry.YankEntry(tk, r, s)).
					To(MatchError("::error ::index test-name already does not have namespace test-namespace and version test-version"))
			})

			it("adds entry to index", func() {
				r.On("GetContents", mock.Anything, "test-owner", "test-repository", filepath.Join("te", "st", "test-namespace_test-name"), rOpts).
					Return(&github.RepositoryContent{
						Content: github.String(asJSONString(index.Entry{
							Namespace: "test-namespace",
							Name:      "test-name",
							Version:   "test-version",
							Address:   "test-address",
						})),
						SHA: github.String("test-sha"),
					}, nil, nil, nil)

				r.On("CreateFile", mock.Anything, "test-owner", "test-repository", filepath.Join("te", "st", "test-namespace_test-name"), &github.RepositoryContentFileOptions{
					Author: &github.CommitAuthor{
						Name:  github.String("buildpacks-bot"),
						Email: github.String("cncf-buildpacks-maintainers@lists.cncf.io"),
					},
					Message: github.String("YANK test-namespace/test-name@test-version"),
					Content: []byte(fmt.Sprintf("%s\n",
						asJSONString(index.Entry{
							Namespace: "test-namespace",
							Name:      "test-name",
							Version:   "test-version",
							Address:   "test-address",
							Yanked:    true,
						}),
					)),
					SHA: github.String("test-sha"),
				}).
					Return(nil, nil, nil)

				Expect(entry.YankEntry(tk, r, s)).To(Succeed())
			})

		})

	}, spec.Report(report.Terminal{}))
}
