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

package toolkit_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	"github.com/buildpacks/github-actions/internal/toolkit"
)

func TestToolkit(t *testing.T) {
	spec.Run(t, "Toolkit", func(t *testing.T, context spec.G, it spec.S) {
		var (
			Expect = NewWithT(t).Expect
		)

		context("MessageContext", func() {

			it("renders string", func() {
				mc := toolkit.MessageContext{Message: "test-message-1\ntest-message-2"}
				Expect(mc.String()).To(Equal("::test-message-1%0Atest-message-2"))

				mc.File = "test-file"
				Expect(mc.String()).To(Equal("file=test-file::test-message-1%0Atest-message-2"))

				mc.Line = "test-line"
				Expect(mc.String()).To(Equal("file=test-file,line=test-line::test-message-1%0Atest-message-2"))

				mc.Column = "test-column"
				Expect(mc.String()).To(Equal("file=test-file,line=test-line,col=test-column::test-message-1%0Atest-message-2"))
			})

		})

		context("FailedError", func() {

			it("returns failed error", func() {
				Expect(toolkit.FailedError("test-message-1", "test-message-2\ntest-message-3")).
					To(MatchError("::error ::test-message-1test-message-2%0Atest-message-3"))
			})

			it("returns failed errorc", func() {
				Expect(toolkit.FailedErrorc(toolkit.MessageContext{
					File:    "test-file",
					Line:    "test-line",
					Column:  "test-column",
					Message: "test-message-1 test-message-2\ntest-message-3",
				})).
					To(MatchError("::error file=test-file,line=test-line,col=test-column::test-message-1 test-message-2%0Atest-message-3"))
			})

			it("returns failed errorf", func() {
				Expect(toolkit.FailedErrorf("%s %s\n%s", "test-message-1", "test-message-2", "test-message-3")).
					To(MatchError("::error ::test-message-1 test-message-2%0Atest-message-3"))
			})

		})

		context("DefaultToolkit", func() {

			var (
				b  = &bytes.Buffer{}
				tk = toolkit.DefaultToolkit{Writer: b}
			)

			it("adds path", func() {
				f, err := ioutil.TempFile("", "github-path")
				Expect(err).NotTo(HaveOccurred())
				_, err = fmt.Fprintln(f, "test-value")
				Expect(err).NotTo(HaveOccurred())
				Expect(f.Close()).To(Succeed())

				tk.Environment = map[string]string{"GITHUB_PATH": f.Name()}

				Expect(tk.AddPath("test-path-1", "test-path-2")).To(Succeed())

				b, err := ioutil.ReadFile(f.Name())
				Expect(string(b)).To(Equal("test-value\ntest-path-1\ntest-path-2\n"))
			})

			it("exports variable", func() {
				f, err := ioutil.TempFile("", "github-env")
				Expect(err).NotTo(HaveOccurred())
				_, err = fmt.Fprintln(f, "TEST_KEY=test-value")
				Expect(err).NotTo(HaveOccurred())
				Expect(f.Close()).To(Succeed())

				tk.Environment = map[string]string{"GITHUB_ENV": f.Name()}

				Expect(tk.ExportVariable("TEST_NAME_1", "test-value-1")).To(Succeed())
				Expect(tk.ExportVariable("TEST_NAME_2", "test-value-2\ntest-value-3")).To(Succeed())

				b, err := ioutil.ReadFile(f.Name())
				Expect(string(b)).To(Equal("TEST_KEY=test-value\nTEST_NAME_1=test-value-1\nTEST_NAME_2<<EOF\ntest-value-2\ntest-value-3\nEOF\n"))
			})

			it("gets input", func() {
				tk.Environment = map[string]string{"INPUT_TEST-1": "test-value"}

				s, ok := tk.GetInput("test-1")
				Expect(ok).To(BeTrue())
				Expect(s).To(Equal("test-value"))

				_, ok = tk.GetInput("test-2")
				Expect(ok).To(BeFalse())
			})

			it("sets output", func() {
				tk.SetOutput("test-name", "test-value-1\ntest-value-2")

				Expect(b.String()).To(Equal("::set-output name=test-name::test-value-1%0Atest-value-2\n"))
			})

			it("gets state", func() {
				tk.Environment = map[string]string{"STATE_TEST-1": "test-value"}

				s, ok := tk.GetState("test-1")
				Expect(ok).To(BeTrue())
				Expect(s).To(Equal("test-value"))

				_, ok = tk.GetState("test-2")
				Expect(ok).To(BeFalse())
			})

			it("sets state", func() {
				tk.SetState("test-name", "test-value-1\ntest-value-2")

				Expect(b.String()).To(Equal("::save-state name=test-name::test-value-1%0Atest-value-2\n"))
			})

			it("adds mask", func() {
				tk.AddMask("test-mask")

				Expect(b.String()).To(Equal("::add-mask::test-mask\n"))
			})

			it("starts group", func() {
				tk.StartGroup("test-title")

				Expect(b.String()).To(Equal("::group::test-title\n"))
			})

			it("ends group", func() {
				tk.EndGroup()

				Expect(b.String()).To(Equal("::endgroup::\n"))
			})

			it("returns isDebug", func() {
				tk.Environment = map[string]string{}

				Expect(tk.IsDebug()).To(BeFalse())

				tk.Environment["RUNNER_DEBUG"] = "true"

				Expect(tk.IsDebug()).To(BeTrue())
			})

			it("writes debug", func() {
				tk.Debug("test-message-1", "test-message-2\ntest-message-3")

				Expect(b.String()).To(Equal("::debug::test-message-1test-message-2%0Atest-message-3\n"))
			})

			it("writes debugf", func() {
				tk.Debugf("%s %s\n%s", "test-message-1", "test-message-2", "test-message-3")

				Expect(b.String()).To(Equal("::debug::test-message-1 test-message-2%0Atest-message-3\n"))
			})

			it("writes warning", func() {
				tk.Warning("test-message-1", "test-message-2\ntest-message-3")

				Expect(b.String()).To(Equal("::warning ::test-message-1test-message-2%0Atest-message-3\n"))
			})

			it("writes warningc", func() {
				tk.Warningc(toolkit.MessageContext{
					File:    "test-file",
					Line:    "test-line",
					Column:  "test-column",
					Message: "test-message-1 test-message-2\ntest-message-3",
				})

				Expect(b.String()).To(Equal("::warning file=test-file,line=test-line,col=test-column::test-message-1 test-message-2%0Atest-message-3\n"))
			})

			it("writes warningf", func() {
				tk.Warningf("%s %s\n%s", "test-message-1", "test-message-2", "test-message-3")

				Expect(b.String()).To(Equal("::warning ::test-message-1 test-message-2%0Atest-message-3\n"))
			})

			it("writes error", func() {
				tk.Error("test-message-1", "test-message-2\ntest-message-3")

				Expect(b.String()).To(Equal("::error ::test-message-1test-message-2%0Atest-message-3\n"))
			})

			it("writes errorc", func() {
				tk.Errorc(toolkit.MessageContext{
					File:    "test-file",
					Line:    "test-line",
					Column:  "test-column",
					Message: "test-message-1 test-message-2\ntest-message-3",
				})

				Expect(b.String()).To(Equal("::error file=test-file,line=test-line,col=test-column::test-message-1 test-message-2%0Atest-message-3\n"))
			})

			it("writes errorf", func() {
				tk.Errorf("%s %s\n%s", "test-message-1", "test-message-2", "test-message-3")

				Expect(b.String()).To(Equal("::error ::test-message-1 test-message-2%0Atest-message-3\n"))
			})

		})

	}, spec.Report(report.Terminal{}))
}
