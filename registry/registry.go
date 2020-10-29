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

package registry

import (
	"context"
	"fmt"

	"github.com/google/go-github/v32/github"
	"github.com/pelletier/go-toml"
)

type IssuesService interface {
	Create(ctx context.Context, owner string, repo string, issue *github.IssueRequest) (*github.Issue, *github.Response, error)
}

type Registry struct {
	Issues  IssuesService
	ID      string
	Version string
}

func (r Registry) Add(address string) error {
	body, err := toml.Marshal(map[string]string{
		"id":      r.ID,
		"version": r.Version,
		"addr":    address,
	})
	if err != nil {
		return fmt.Errorf("unable to marshal to TOML\n%w", err)
	}

	req := &github.IssueRequest{
		Title: github.String(fmt.Sprintf("ADD %s@%s", r.ID, r.Version)),
		Body:  github.String(string(body)),
	}

	issue, _, err := r.Issues.Create(context.Background(), "buildpacks", "registry-index", req)
	if err != nil {
		return fmt.Errorf("unable to create issue\n%w", err)
	}

	fmt.Printf("Created issue %s\n", *issue.HTMLURL)
	return nil
}

func (r Registry) Yank() error {
	body, err := toml.Marshal(map[string]string{
		"id":      r.ID,
		"version": r.Version,
	})
	if err != nil {
		return fmt.Errorf("unable to marshal to TOML\n%w", err)
	}

	req := &github.IssueRequest{
		Title: github.String(fmt.Sprintf("YANK %s@%s", r.ID, r.Version)),
		Body:  github.String(string(body)),
	}

	issue, _, err := r.Issues.Create(context.Background(), "buildpacks", "registry-index", req)
	if err != nil {
		return fmt.Errorf("unable to create issue\n%w", err)
	}

	fmt.Printf("Created issue %s\n", *issue.HTMLURL)
	return nil
}
