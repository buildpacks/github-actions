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
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
	"gopkg.in/retry.v1"

	"github.com/buildpacks/github-actions/internal/toolkit"
	"github.com/buildpacks/github-actions/registry/request-yank-entry"
)

func main() {
	tk := &toolkit.DefaultToolkit{}

	t, ok := tk.GetInput("token")
	if !ok {
		fmt.Println(toolkit.FailedError("token must be specified"))
		os.Exit(1)
	}

	gh := github.NewClient(oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(&oauth2.Token{AccessToken: t})))

	strategy := retry.LimitTime(8*time.Minute,
		retry.Exponential{
			Initial:  time.Second,
			MaxDelay: 30 * time.Second,
		},
	)

	if err := entry.RequestYankEntry(tk, gh.Issues, strategy); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
