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

package registry_test

import (
	"context"
	"testing"

	"github.com/google/go-github/v32/github"
	. "github.com/onsi/gomega"
	"github.com/pelletier/go-toml"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	"github.com/buildpacks/github-actions/registry"
)

func TestRegistry(t *testing.T) {
	spec.Run(t, "registry", func(t *testing.T, when spec.G, it spec.S) {
		var (
			Expect = NewWithT(t).Expect

			i = TestIssues{Response: &github.Issue{HTMLURL: github.String("test-html-url")}}

			r = registry.Registry{
				Issues:  &i,
				ID:      "test-namespace/test-name",
				Version: "test-version",
			}
		)

		it("creates add issue", func() {
			Expect(r.Add("test-address")).To(Succeed())

			Expect(i.Owner).To(Equal("buildpacks"))
			Expect(i.Repo).To(Equal("registry-index"))
			Expect(*i.Request.Title).To(Equal("ADD test-namespace/test-name@test-version"))

			var body map[string]string
			Expect(toml.Unmarshal([]byte(*i.Request.Body), &body)).To(Succeed())
			Expect(body).To(Equal(map[string]string{
				"id":      "test-namespace/test-name",
				"version": "test-version",
				"addr":    "test-address",
			}))
		})

		it("creates yank issue", func() {
			Expect(r.Yank()).To(Succeed())

			Expect(i.Owner).To(Equal("buildpacks"))
			Expect(i.Repo).To(Equal("registry-index"))
			Expect(*i.Request.Title).To(Equal("YANK test-namespace/test-name@test-version"))

			var body map[string]string
			Expect(toml.Unmarshal([]byte(*i.Request.Body), &body)).To(Succeed())
			Expect(body).To(Equal(map[string]string{
				"id":      "test-namespace/test-name",
				"version": "test-version",
			}))
		})

	}, spec.Report(report.Terminal{}))
}

type TestIssues struct {
	Owner   string
	Repo    string
	Request *github.IssueRequest

	Response *github.Issue
}

func (t *TestIssues) Create(_ context.Context, owner string, repo string, issue *github.IssueRequest) (*github.Issue, *github.Response, error) {
	t.Owner = owner
	t.Repo = repo
	t.Request = issue

	return t.Response, nil, nil
}
