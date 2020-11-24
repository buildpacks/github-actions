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

	"github.com/buildpacks/github-actions/internal/toolkit"
	"github.com/buildpacks/github-actions/registry/internal/index"
	"github.com/buildpacks/github-actions/registry/internal/services"
)

func YankEntry(tk toolkit.Toolkit, repositories services.RepositoriesService) error {
	owner, ok := tk.GetInput("owner")
	if !ok {
		return toolkit.FailedError("owner must be set")
	}

	repository, ok := tk.GetInput("repository")
	if !ok {
		return toolkit.FailedError("repository must be set")
	}

	namespace, ok := tk.GetInput("namespace")
	if !ok {
		return toolkit.FailedError("namespace must be set")
	}

	name, ok := tk.GetInput("name")
	if !ok {
		return toolkit.FailedError("name must be set")
	}

	version, ok := tk.GetInput("version")
	if !ok {
		return toolkit.FailedError("version must be set")
	}

	file := index.Path(namespace, name)
	content, _, _, err := repositories.GetContents(context.Background(), owner, repository, file, nil)
	if err2, ok := err.(*github.ErrorResponse); ok && err2.Response.StatusCode == http.StatusNotFound {
		return toolkit.FailedErrorf("index %s does not exist", name)
	} else if err != nil {
		return toolkit.FailedErrorf("unable to read index %s\n%w", name, err)
	}

	s, err := content.GetContent()
	if err != nil {
		return toolkit.FailedErrorf("unable to get index content\n%w", err)
	}

	entries, err := index.UnmarshalEntries(s)
	if err != nil {
		return toolkit.FailedErrorf("unable to unmarshal entries\n%w", err)
	}

	i := indexOf(entries, namespace, version)
	if i == nil {
		return toolkit.FailedErrorf("index %s already does not have namespace %s and version %s", name, namespace, version)
	}

	entries[*i].Yanked = true

	s, err = index.MarshalEntries(entries)
	if err != nil {
		return toolkit.FailedErrorf("unable to marshal entries\n%w", err)
	}

	if _, _, err := repositories.CreateFile(context.Background(), owner, repository, file, &github.RepositoryContentFileOptions{
		Author: &github.CommitAuthor{
			Name:  github.String("buildpacks-bot"),
			Email: github.String("cncf-buildpacks-maintainers@lists.cncf.io"),
		},
		Message: github.String(fmt.Sprintf("ADD %s/%s@%s", namespace, name, version)),
		SHA:     content.SHA,
		Content: []byte(s),
	}); err != nil {
		return err
	}

	fmt.Printf("Yanked %s/%s@%s\n", namespace, name, version)
	return nil
}

func indexOf(entries []index.Entry, namespace string, version string) *int {
	for i, e := range entries {
		if e.Namespace == namespace && e.Version == version {
			return &i
		}
	}

	return nil
}
