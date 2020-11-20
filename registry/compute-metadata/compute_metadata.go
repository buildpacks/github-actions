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

package metadata

import (
	"encoding/json"
	"strings"

	"github.com/google/go-github/v32/github"
	"github.com/pelletier/go-toml"

	"github.com/buildpacks/github-actions/registry"
	"github.com/buildpacks/github-actions/toolkit"
)

func ComputeMetadata(tk toolkit.Toolkit) error {
	i, ok := tk.GetInput("issue")
	if !ok {
		return toolkit.FailedError("issue must be set")
	}

	var issue github.Issue
	if err := json.Unmarshal([]byte(i), &issue); err != nil {
		return toolkit.FailedErrorf("unable to unmarshal issue\n%w", err)
	}

	var request registry.IndexRequest
	if err := toml.Unmarshal([]byte(strings.ReplaceAll(*issue.Body, "```", "")), &request); err != nil {
		return toolkit.FailedErrorf("unable to unmarshal body\n%w", err)
	}

	var (
		namespace string
		name      string
	)
	if g := registry.ValidId.FindStringSubmatch(request.ID); g == nil {
		return toolkit.FailedErrorf("invalid id '%s'", request.ID)
	} else {
		namespace = g[1]
		name = g[2]
	}

	if !registry.ValidVersion.MatchString(request.Version) {
		return toolkit.FailedErrorf("invalid version '%s'", request.Version)
	}

	if !registry.ValidAddress.MatchString(request.Address) {
		return toolkit.FailedErrorf("invalid address '%s'", request.Address)
	}

	tk.SetOutput("id", request.ID)
	tk.SetOutput("version", request.Version)
	tk.SetOutput("address", request.Address)
	tk.SetOutput("namespace", namespace)
	tk.SetOutput("name", name)

	return nil
}
