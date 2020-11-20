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

package entry

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v32/github"
	"github.com/pelletier/go-toml"

	"github.com/buildpacks/github-actions/registry"
	"github.com/buildpacks/github-actions/toolkit"
)

func YankEntry(tk toolkit.Toolkit, issues registry.IssuesService, timeout *time.Timer, interval *time.Ticker) error {
	id, ok := tk.GetInput("id")
	if !ok {
		return toolkit.FailedError("id must be set")
	}

	version, ok := tk.GetInput("version")
	if !ok {
		return toolkit.FailedError("version must be set")
	}

	address, ok := tk.GetInput("address")
	if !ok {
		return toolkit.FailedError("address must be set")
	}

	body, err := toml.Marshal(registry.Request{
		ID:      id,
		Version: version,
		Address: address,
		Yank:    true,
	})
	if err != nil {
		return toolkit.FailedErrorf("unable to marshal to TOML\n%w", err)
	}

	req := &github.IssueRequest{
		Title: github.String(fmt.Sprintf("YANK %s@%s", id, version)),
		Body:  github.String(string(body)),
	}

	issue, _, err := issues.Create(context.Background(), "buildpacks", "registry-index", req)
	if err != nil {
		return toolkit.FailedErrorf("unable to create issue\n%w", err)
	}

	url := *issue.HTMLURL
	number := *issue.Number

	fmt.Printf("Created issue %s\n", url)

	for {
		select {
		case <-timeout.C:
			return toolkit.FailedError("timed out waiting for request to be processed")
		case <-interval.C:
			issue, _, err = issues.Get(context.Background(), "buildpacks", "registry-index", number)
			if err != nil {
				tk.Warningf("unable to get state for %s", url)
			}

			for _, l := range issue.Labels {
				if *l.Name == registry.FailureLabel {
					return toolkit.FailedErrorf("Registry request %s failed", url)
				} else if *l.Name == registry.SuccessLabel {
					fmt.Printf("Registry request %s succeeded\n", url)
					return nil
				}
			}
		}
	}
}
