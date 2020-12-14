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
	"net/http"

	"github.com/google/go-github/v32/github"
	"gopkg.in/retry.v1"

	"github.com/buildpacks/github-actions/internal/toolkit"
	"github.com/buildpacks/github-actions/registry/internal/index"
	"github.com/buildpacks/github-actions/registry/internal/services"
)

func AddEntry(tk toolkit.Toolkit, repositories services.RepositoriesService, strategy retry.Strategy) error {
	c, err := parseConfig(tk)
	if err != nil {
		return err
	}

	file := index.Path(c.Namespace, c.Name)

	for a := retry.Start(strategy, nil); a.Next(); {
		content, _, resp, err := repositories.GetContents(context.Background(), c.Owner, c.Repository, file, nil)
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			fmt.Printf("New Index: %s\n", c.Name)
			content = &github.RepositoryContent{}
		} else if err != nil {
			return toolkit.FailedErrorf("unable to read index %s\n%w", c.Name, err)
		}

		s, err := content.GetContent()
		if err != nil {
			return toolkit.FailedErrorf("unable to get index content\n%w", err)
		}

		entries, err := index.UnmarshalEntries(s)
		if err != nil {
			return toolkit.FailedErrorf("unable to unmarshal entries\n%w", err)
		}

		if contains(entries, c.Namespace, c.Version) {
			return toolkit.FailedErrorf("index %s already has namespace %s and version %s", c.Name, c.Namespace, c.Version)
		}

		entries = append(entries, index.Entry{
			Namespace: c.Namespace,
			Name:      c.Name,
			Version:   c.Version,
			Address:   c.Address,
		})

		s, err = index.MarshalEntries(entries)
		if err != nil {
			return toolkit.FailedErrorf("unable to marshal entries\n%w", err)
		}

		if _, resp, err := repositories.CreateFile(context.Background(), c.Owner, c.Repository, file, &github.RepositoryContentFileOptions{
			Author: &github.CommitAuthor{
				Name:  github.String("buildpacks-bot"),
				Email: github.String("cncf-buildpacks-maintainers@lists.cncf.io"),
			},
			Message: github.String(fmt.Sprintf("ADD %s/%s@%s", c.Namespace, c.Name, c.Version)),
			SHA:     content.SHA,
			Content: []byte(s),
		}); resp != nil && resp.StatusCode == http.StatusConflict {
			tk.Warning("retrying index update after conflict")
			continue
		} else if err != nil {
			return toolkit.FailedErrorf("unable to create index\n%w", err)
		}

		fmt.Printf("Added %s/%s@%s\n", c.Namespace, c.Name, c.Version)
		return nil
	}

	return toolkit.FailedError("timed out")
}

type config struct {
	Owner      string
	Repository string
	Namespace  string
	Name       string
	Version    string
	Address    string
}

func parseConfig(tk toolkit.Toolkit) (config, error) {
	var (
		c  config
		ok bool
	)

	c.Owner, ok = tk.GetInput("owner")
	if !ok {
		return config{}, toolkit.FailedError("owner must be set")
	}

	c.Repository, ok = tk.GetInput("repository")
	if !ok {
		return config{}, toolkit.FailedError("repository must be set")
	}

	c.Namespace, ok = tk.GetInput("namespace")
	if !ok {
		return config{}, toolkit.FailedError("namespace must be set")
	}

	c.Name, ok = tk.GetInput("name")
	if !ok {
		return config{}, toolkit.FailedError("name must be set")
	}

	c.Version, ok = tk.GetInput("version")
	if !ok {
		return config{}, toolkit.FailedError("version must be set")
	}

	c.Address, ok = tk.GetInput("address")
	if !ok {
		return config{}, toolkit.FailedError("address must be set")
	}

	return c, nil
}

func contains(entries []index.Entry, namespace string, version string) bool {
	for _, e := range entries {
		if e.Namespace == namespace && e.Version == version {
			return true
		}
	}

	return false
}
