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
	"fmt"
	"testing"
	"time"

	"github.com/google/go-github/v32/github"
	. "github.com/onsi/gomega"
	"github.com/pelletier/go-toml"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
	"github.com/stretchr/testify/mock"

	"github.com/buildpacks/github-actions/internal/toolkit"
	"github.com/buildpacks/github-actions/registry/internal/index"
	"github.com/buildpacks/github-actions/registry/internal/services"
	"github.com/buildpacks/github-actions/registry/request-yank-entry"
)

func TestRequestYankEntry(t *testing.T) {
	spec.Run(t, "request-yank-entry", func(t *testing.T, context spec.G, it spec.S) {
		var (
			Expect = NewWithT(t).Expect

			i  = &services.MockIssuesService{}
			tk = &toolkit.MockToolkit{}
		)

		it.Before(func() {
			tk.On("GetInput", "id").Return("test-namespace/test-name", true)
			tk.On("GetInput", "version").Return("test-version", true)
			tk.On("GetInput", "address").Return("test-address", true)

			b, err := toml.Marshal(index.Request{
				ID:      "test-namespace/test-name",
				Version: "test-version",
				Address: "test-address",
				Yank:    true,
			})
			Expect(err).NotTo(HaveOccurred())

			i.On("Create", mock.Anything, "buildpacks", "registry-index", &github.IssueRequest{
				Title: github.String("YANK test-namespace/test-name@test-version"),
				Body:  github.String(fmt.Sprintf("```\n%s\n```", string(b))),
			}).Return(&github.Issue{
				Number:  github.Int(1),
				HTMLURL: github.String("test-html-url"),
			}, nil, nil)
		})

		it("yank entry succeeds", func() {
			i.On("Get", mock.Anything, "buildpacks", "registry-index", 1).Return(&github.Issue{
				Labels: []*github.Label{{Name: github.String(index.RequestSuccessLabel)}},
			}, nil, nil)

			timeout := time.NewTimer(1 * time.Second)
			defer timeout.Stop()
			interval := time.NewTicker(1)
			defer interval.Stop()

			Expect(entry.RequestYankEntry(tk, i, timeout, interval)).To(Succeed())
		})

		it("yank entry fails", func() {
			i.On("Get", mock.Anything, "buildpacks", "registry-index", 1).Return(&github.Issue{
				Labels: []*github.Label{{Name: github.String(index.RequestFailureLabel)}},
			}, nil, nil)

			timeout := time.NewTimer(1 * time.Second)
			defer timeout.Stop()
			interval := time.NewTicker(1)
			defer interval.Stop()

			Expect(entry.RequestYankEntry(tk, i, timeout, interval)).
				To(MatchError("::error ::Registry request test-html-url failed"))
		})

	}, spec.Report(report.Terminal{}))
}
