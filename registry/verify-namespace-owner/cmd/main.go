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
	"time"

	"github.com/google/go-github/v89/github"
	"gopkg.in/retry.v1"

	"github.com/buildpacks/github-actions/internal/toolkit"
	owner "github.com/buildpacks/github-actions/registry/verify-namespace-owner"
)

func main() {
	tk := &toolkit.DefaultToolkit{}

	t, ok := tk.GetInput("token")
	if !ok {
		fmt.Println(toolkit.FailedError("token must be specified"))
		os.Exit(1)
	}

	gh, err := github.NewClient(github.WithAuthToken(t))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	strategy := retry.LimitTime(
		2*time.Minute,
		retry.Exponential{
			Initial: time.Second,
			Jitter:  true,
		},
	)

	if err := owner.VerifyNamespaceOwner(tk, gh.Organizations, gh.Repositories, strategy); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
