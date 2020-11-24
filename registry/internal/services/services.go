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

package services

import (
	"context"

	"github.com/google/go-github/v32/github"
)

//go:generate mockery --all  --inpackage --case=underscore

type IssuesService interface {
	Create(ctx context.Context, owner string, repo string, issue *github.IssueRequest) (*github.Issue, *github.Response, error)
	Get(ctx context.Context, owner string, repo string, number int) (*github.Issue, *github.Response, error)
}

type OrganizationsService interface {
	List(ctx context.Context, user string, opts *github.ListOptions) ([]*github.Organization, *github.Response, error)
}

type RepositoriesService interface {
	CreateFile(ctx context.Context, owner, repo, path string, opts *github.RepositoryContentFileOptions) (*github.RepositoryContentResponse, *github.Response, error)
	GetContents(ctx context.Context, owner, repo, path string, opts *github.RepositoryContentGetOptions) (fileContent *github.RepositoryContent, directoryContent []*github.RepositoryContent, resp *github.Response, err error)
}
