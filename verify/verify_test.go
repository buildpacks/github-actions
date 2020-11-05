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

package verify_test

import (
	"testing"

	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/fake"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	"github.com/buildpacks/github-actions/verify"
)

func TestVerify(t *testing.T) {
	spec.Run(t, "verify", func(t *testing.T, when spec.G, it spec.S) {
		var (
			Expect = NewWithT(t).Expect

			r = &TestRemote{Labels: make(map[string]string)}

			v = verify.Verifier{
				Image: r.Image,

				ID:      "test-id",
				Version: "test-version",
				Address: "host/repository@sha256:04ba2d17480910bd340f0305d846b007148dafd64bc6fc2626870c174b7c7de7",
			}
		)

		it("fails if address is not a digest image reference", func() {
			v.Address = "test-host/test-repository:test-version"

			Expect(v.Verify()).To(MatchError("address test-host/test-repository:test-version must be in digest form <host>/<repository>@sh256:<digest>"))
		})

		it("fails if io.buildpacks.buildpackage.metadata is not on image", func() {
			Expect(v.Verify()).To(MatchError("unable to retrieve io.buildpacks.buildpackage.metadata label"))
		})

		it("fails if id does not match", func() {
			r.Labels[verify.MetadataLabel] = `{ "id": "another-id", "version": "test-version" }`

			Expect(v.Verify()).To(MatchError("invalid id in buildpackage: expected 'test-id', found 'another-id'"))
		})

		it("fails if version does not match", func() {
			r.Labels[verify.MetadataLabel] = `{ "id": "test-id", "version": "another-version" }`

			Expect(v.Verify()).To(MatchError("invalid version in buildpackage: expected 'test-version', found 'another-version'"))
		})

		it("passes if version and id match", func() {
			r.Labels[verify.MetadataLabel] = `{ "id": "test-id", "version": "test-version" }`

			Expect(v.Verify()).To(Succeed())
		})

	}, spec.Report(report.Terminal{}))
}

type TestRemote struct {
	Labels map[string]string
}

func (t *TestRemote) Image(_ name.Reference, _ ...remote.Option) (v1.Image, error) {
	i := &fake.FakeImage{}
	i.ConfigFileReturns(&v1.ConfigFile{
		Config: v1.Config{
			Labels: t.Labels,
		},
	}, nil)

	return i, nil
}
