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

package main

import (
	"fmt"
	"os"

	"github.com/google/go-containerregistry/pkg/v1/remote"

	"github.com/buildpacks/github-actions/verify"
)

func main() {
	var (
		v verify.Verifier

		err error
		ok  bool
	)

	v.Image = remote.Image

	if v.ID, ok = os.LookupEnv("INPUT_ID"); !ok {
		panic(fmt.Errorf("id must be specified"))
	}

	if v.Version, ok = os.LookupEnv("INPUT_VERSION"); !ok {
		panic(fmt.Errorf("version must be specified"))
	}

	if v.Address, ok = os.LookupEnv("INPUT_ADDRESS"); !ok {
		panic(fmt.Errorf("address must be specified"))
	}

	if err = v.Verify(); err != nil {
		panic(err)
	}
}
