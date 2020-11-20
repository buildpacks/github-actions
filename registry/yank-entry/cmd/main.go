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

	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"

	"github.com/buildpacks/github-actions/registry/yank-entry"
	"github.com/buildpacks/github-actions/toolkit"
)

func main() {
	tk := &toolkit.DefaultToolkit{}

	t, ok := tk.GetInput("token")
	if !ok {
		fmt.Println(toolkit.FailedError("token must be specified"))
		os.Exit(1)
	}

	gh := github.NewClient(oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(&oauth2.Token{AccessToken: t})))

	timeout := time.NewTimer(5 * time.Minute)
	defer timeout.Stop()
	interval := time.NewTicker(30 * time.Second)
	defer interval.Stop()

	if err := entry.YankEntry(tk, gh.Issues, timeout, interval); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
