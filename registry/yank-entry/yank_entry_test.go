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
	"testing"
	"time"

	"github.com/google/go-github/v32/github"
	. "github.com/onsi/gomega"
	"github.com/pelletier/go-toml"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
	"github.com/stretchr/testify/mock"

	"github.com/buildpacks/github-actions/registry"
	mocks2 "github.com/buildpacks/github-actions/registry/mocks"
	"github.com/buildpacks/github-actions/registry/yank-entry"
	"github.com/buildpacks/github-actions/toolkit/mocks"
)

func TestYankEntry(t *testing.T) {
	spec.Run(t, "yank-entry", func(t *testing.T, when spec.G, it spec.S) {
		var (
			Expect = NewWithT(t).Expect

			i  = &mocks2.IssuesService{}
			tk = &mocks.Toolkit{}
		)

		it.Before(func() {
			tk.On("GetInput", "id").Return("test-namespace/test-name", true)
			tk.On("GetInput", "version").Return("test-version", true)
			tk.On("GetInput", "address").Return("test-address", true)

			b, err := toml.Marshal(registry.Request{
				ID:      "test-namespace/test-name",
				Version: "test-version",
				Address: "test-address",
				Yank:    true,
			})
			Expect(err).NotTo(HaveOccurred())

			i.On("Create", mock.Anything, "buildpacks", "registry-index", &github.IssueRequest{
				Title: github.String("YANK test-namespace/test-name@test-version"),
				Body:  github.String(string(b)),
			}).Return(&github.Issue{
				Number:  github.Int(1),
				HTMLURL: github.String("test-html-url"),
			}, nil, nil)
		})

		it("yank entry succeeds", func() {
			i.On("Get", mock.Anything, "buildpacks", "registry-index", 1).Return(&github.Issue{
				Labels: []*github.Label{{Name: github.String(registry.SuccessLabel)}},
			}, nil, nil)

			timeout := time.NewTimer(1 * time.Second)
			defer timeout.Stop()
			interval := time.NewTicker(1)
			defer interval.Stop()

			Expect(entry.YankEntry(tk, i, timeout, interval)).To(Succeed())
		})

		it("yank entry fails", func() {
			i.On("Get", mock.Anything, "buildpacks", "registry-index", 1).Return(&github.Issue{
				Labels: []*github.Label{{Name: github.String(registry.FailureLabel)}},
			}, nil, nil)

			timeout := time.NewTimer(1 * time.Second)
			defer timeout.Stop()
			interval := time.NewTicker(1)
			defer interval.Stop()

			Expect(entry.YankEntry(tk, i, timeout, interval)).
				To(MatchError("::error ::Registry request test-html-url failed"))
		})

	}, spec.Report(report.Terminal{}))
}
