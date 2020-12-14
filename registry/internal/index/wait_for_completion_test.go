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
)

func TestWaitForCompletion(t *testing.T) {
	spec.Run(t, "wait-for-completion", func(t *testing.T, context spec.G, it spec.S) {
		var (
			Expect = NewWithT(t).Expect

			i  = &services.MockIssuesService{}
			s  = retry.LimitCount(2, retry.Regular{Min: 2})
			tk = &toolkit.MockToolkit{}
		)

		it("it handles success", func() {
			i.On("Get", mock.Anything, "buildpacks", "registry-index", 1).
				Return(&github.Issue{Labels: []*github.Label{{Name: github.String(index.RequestSuccessLabel)}}}, nil, nil)

			Expect(index.WaitForCompletion(1, "test-url", tk, i, s)).To(Succeed())
		})

		it("it handles failure", func() {
			i.On("Get", mock.Anything, "buildpacks", "registry-index", 1).
				Return(&github.Issue{Labels: []*github.Label{{Name: github.String(index.RequestFailureLabel)}}}, nil, nil)

			Expect(index.WaitForCompletion(1, "test-url", tk, i, s)).
				To(MatchError("::error ::Registry request test-url failed"))
		})

		it("retries", func() {
			i.On("Get", mock.Anything, "buildpacks", "registry-index", 1).
				Return(&github.Issue{}, nil, nil).
				Once()
			i.On("Get", mock.Anything, "buildpacks", "registry-index", 1).
				Return(&github.Issue{Labels: []*github.Label{{Name: github.String(index.RequestSuccessLabel)}}}, nil, nil)

			Expect(index.WaitForCompletion(1, "test-url", tk, i, s)).To(Succeed())
		})

	}, spec.Report(report.Terminal{}))
}
