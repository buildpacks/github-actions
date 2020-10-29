/*
 * Copyright 2020 the original author or authors.
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
	"strconv"

	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"

	"github.com/buildpacks/github-actions/registry"
)

func main() {
	var (
		r registry.Registry

		err error
		ok  bool
	)

	var token oauth2.Token
	if token.AccessToken, ok = os.LookupEnv("INPUT_TOKEN"); !ok {
		panic(fmt.Errorf("token must be specified"))
	}
	r.Issues = github.NewClient(oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(&token))).Issues

	if r.ID, ok = os.LookupEnv("INPUT_ID"); !ok {
		panic(fmt.Errorf("id must be specified"))
	}

	if r.Version, ok = os.LookupEnv("INPUT_VERSION"); !ok {
		panic(fmt.Errorf("id must be specified"))
	}

	yank := false
	if s, ok := os.LookupEnv("INPUT_YANK"); ok {
		if yank, err = strconv.ParseBool(s); err != nil {
			panic(fmt.Errorf("unable to parse INPUT_YANK='%s' as a bool", s))
		}
	}

	if yank {
		if err := r.Yank(); err != nil {
			panic(err)
		}
	} else {
		address, ok := os.LookupEnv("INPUT_ADDRESS")
		if !ok {
			panic(fmt.Errorf("address must be specified"))
		}

		if err := r.Add(address); err != nil {
			panic(err)
		}
	}
}
