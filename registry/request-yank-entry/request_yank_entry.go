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

	"github.com/google/go-github/v32/github"
	"github.com/pelletier/go-toml"
	"gopkg.in/retry.v1"

	"github.com/buildpacks/github-actions/internal/toolkit"
	"github.com/buildpacks/github-actions/registry/internal/index"
	"github.com/buildpacks/github-actions/registry/internal/services"
)

func RequestYankEntry(tk toolkit.Toolkit, issues services.IssuesService, strategy retry.Strategy) error {
	c, err := parseConfig(tk)
	if err != nil {
		return err
	}

	body, err := toml.Marshal(index.Request{
		ID:      c.ID,
		Version: c.Version,
		Yank:    true,
	})
	if err != nil {
		return toolkit.FailedErrorf("unable to marshal to TOML\n%w", err)
	}

	req := &github.IssueRequest{
		Title: github.String(fmt.Sprintf("YANK %s@%s", c.ID, c.Version)),
		Body:  github.String(fmt.Sprintf("```\n%s\n```", string(body))),
	}

	issue, _, err := issues.Create(context.Background(), "buildpacks", "registry-index", req)
	if err != nil {
		return toolkit.FailedErrorf("unable to create issue\n%w", err)
	}

	url := *issue.HTMLURL
	number := *issue.Number

	fmt.Printf("Created issue %s\n", url)

	for a := retry.Start(strategy, nil); a.Next(); {
		issue, _, err = issues.Get(context.Background(), "buildpacks", "registry-index", number)
		if err != nil {
			tk.Warningf("unable to get state for %s", url)
			continue
		}

		for _, l := range issue.Labels {
			if *l.Name == index.RequestFailureLabel {
				return toolkit.FailedErrorf("Registry request %s failed", url)
			} else if *l.Name == index.RequestSuccessLabel {
				fmt.Printf("Registry request %s succeeded\n", url)
				return nil
			}
		}
	}

	return toolkit.FailedError("timed out waiting for request to be processed")
}

type config struct {
	ID      string
	Version string
}

func parseConfig(tk toolkit.Toolkit) (config, error) {
	var (
		c  config
		ok bool
	)

	c.ID, ok = tk.GetInput("id")
	if !ok {
		return config{}, toolkit.FailedError("id must be set")
	}

	c.Version, ok = tk.GetInput("version")
	if !ok {
		return config{}, toolkit.FailedError("version must be set")
	}

	return c, nil
}
