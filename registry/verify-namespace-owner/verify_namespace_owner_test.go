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

package owner_test

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/google/go-github/v32/github"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
	"github.com/stretchr/testify/mock"

	"github.com/buildpacks/github-actions/internal/toolkit"
	"github.com/buildpacks/github-actions/registry/internal/namespace"
	"github.com/buildpacks/github-actions/registry/internal/services"
	"github.com/buildpacks/github-actions/registry/verify-namespace-owner"
)

func TestVerifyNamespaceOwner(t *testing.T) {
	spec.Run(t, "verify-namespace-owner", func(t *testing.T, context spec.G, it spec.S) {
		var (
			Expect           = NewWithT(t).Expect
			ExpectWithOffset = NewWithT(t).ExpectWithOffset

			o     = &services.MockOrganizationsService{}
			r     = &services.MockRepositoriesService{}
			rOpts *github.RepositoryContentGetOptions
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
			tk.On("GetInput", "user").
				Return(asJSONString(github.User{ID: github.Int64(1), Login: github.String("test-user")}), true)
		})

		context("unknown namespace", func() {
			it.Before(func() {

				r.On("GetContents", mock.Anything, "test-owner", "test-repository", filepath.Join("v1", "test-namespace.json"), rOpts).
					Return(nil, nil, nil, &github.ErrorResponse{Response: &http.Response{StatusCode: http.StatusNotFound}})
			})

			it("fails if add-if-missing is false", func() {
				tk.On("GetInput", "add-if-missing").Return("", false)

				Expect(owner.VerifyNamespaceOwner(tk, o, r)).
					To(MatchError("::error ::invalid namespace test-namespace"))
			})

			it("succeeds if add-if-missing is true", func() {
				tk.On("GetInput", "add-if-missing").Return("true", true)

				c := asJSONString(namespace.Namespace{Owners: []namespace.Owner{{ID: 1, Type: namespace.UserType}}})

				r.On("CreateFile", mock.Anything, "test-owner", "test-repository", filepath.Join("v1", "test-namespace.json"), &github.RepositoryContentFileOptions{
					Author: &github.CommitAuthor{
						Name:  github.String("buildpacks-bot"),
						Email: github.String("cncf-buildpacks-maintainers@lists.cncf.io"),
					},
					Message: github.String("New Namespace: test-namespace"),
					Content: []byte(c),
				}).Return(&github.RepositoryContentResponse{
					Content: &github.RepositoryContent{Content: github.String(c)},
				}, nil, nil)

				Expect(owner.VerifyNamespaceOwner(tk, o, r)).To(Succeed())
			})
		})

		context("user-owned namespace", func() {

			it.Before(func() {
				o.On("List", mock.Anything, "test-user", mock.Anything).
					Return([]*github.Organization{}, &github.Response{}, nil)
			})

			it("fails if user does not own", func() {
				r.On("GetContents", mock.Anything, "test-owner", "test-repository", filepath.Join("v1", "test-namespace.json"), rOpts).
					Return(&github.RepositoryContent{
						Content: github.String(asJSONString(namespace.Namespace{Owners: []namespace.Owner{{ID: 2, Type: namespace.UserType}}})),
					}, nil, nil, nil)

				Expect(owner.VerifyNamespaceOwner(tk, o, r)).
					To(MatchError("::error ::test-user is not an owner of test-namespace"))
			})

			it("succeeds if user does own", func() {
				r.On("GetContents", mock.Anything, "test-owner", "test-repository", filepath.Join("v1", "test-namespace.json"), rOpts).
					Return(&github.RepositoryContent{
						Content: github.String(asJSONString(namespace.Namespace{Owners: []namespace.Owner{{ID: 1, Type: namespace.UserType}}})),
					}, nil, nil, nil)

				Expect(owner.VerifyNamespaceOwner(tk, o, r)).To(Succeed())
			})
		})

		context("organization-owned namespace", func() {

			it.Before(func() {
				r.On("GetContents", mock.Anything, "test-owner", "test-repository", filepath.Join("v1", "test-namespace.json"), rOpts).
					Return(&github.RepositoryContent{
						Content: github.String(asJSONString(namespace.Namespace{Owners: []namespace.Owner{{ID: 1, Type: namespace.OrganizationType}}})),
					}, nil, nil, nil)
			})

			it("fails if user does not own", func() {
				o.On("List", mock.Anything, "test-user", mock.Anything).
					Return([]*github.Organization{}, &github.Response{}, nil)

				Expect(owner.VerifyNamespaceOwner(tk, o, r)).
					To(MatchError("::error ::test-user is not an owner of test-namespace"))
			})

			it("succeeds if user does own", func() {
				o.On("List", mock.Anything, "test-user", mock.Anything).
					Return([]*github.Organization{{ID: github.Int64(1)}}, &github.Response{}, nil)

				Expect(owner.VerifyNamespaceOwner(tk, o, r)).To(Succeed())
			})

		})

	}, spec.Report(report.Terminal{}))
}
