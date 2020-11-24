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
	"encoding/json"
	"fmt"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"

	"github.com/buildpacks/github-actions/registry/internal/index"
)

func TestEntry(t *testing.T) {
	spec.Run(t, "entry", func(t *testing.T, context spec.G, it spec.S) {
		var (
			Expect           = NewWithT(t).Expect
			ExpectWithOffset = NewWithT(t).ExpectWithOffset
		)

		asJSONString := func(v interface{}) string {
			b, err := json.Marshal(v)
			ExpectWithOffset(1, err).NotTo(HaveOccurred())

			return string(b)
		}

		it("marshals entries", func() {
			Expect(index.MarshalEntries([]index.Entry{
				{
					Namespace: "test-namespace-1",
					Name:      "test-name-1",
					Version:   "test-version-1",
					Address:   "test-address-1",
				},
				{
					Namespace: "test-namespace-2",
					Name:      "test-name-2",
					Version:   "test-version-2",
					Address:   "test-address-2",
				},
			})).To(Equal(fmt.Sprintf("%s\n%s\n",
				asJSONString(index.Entry{
					Namespace: "test-namespace-1",
					Name:      "test-name-1",
					Version:   "test-version-1",
					Address:   "test-address-1",
				}),
				asJSONString(index.Entry{
					Namespace: "test-namespace-2",
					Name:      "test-name-2",
					Version:   "test-version-2",
					Address:   "test-address-2",
				}),
			)))
		})

		it("unmarshals entries", func() {
			Expect(index.UnmarshalEntries(fmt.Sprintf("%s\n%s\n",
				asJSONString(index.Entry{
					Namespace: "test-namespace-1",
					Name:      "test-name-1",
					Version:   "test-version-1",
					Address:   "test-address-1",
				}),
				asJSONString(index.Entry{
					Namespace: "test-namespace-2",
					Name:      "test-name-2",
					Version:   "test-version-2",
					Address:   "test-address-2",
				}),
			))).To(Equal([]index.Entry{
				{
					Namespace: "test-namespace-1",
					Name:      "test-name-1",
					Version:   "test-version-1",
					Address:   "test-address-1",
				},
				{
					Namespace: "test-namespace-2",
					Name:      "test-name-2",
					Version:   "test-version-2",
					Address:   "test-address-2",
				},
			}))
		})
	})
}
