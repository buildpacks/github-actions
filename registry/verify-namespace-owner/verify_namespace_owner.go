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

package owner

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/google/go-github/v32/github"

	"github.com/buildpacks/github-actions/registry"
	"github.com/buildpacks/github-actions/toolkit"
)

func VerifyNamespaceOwner(tk toolkit.Toolkit, organizations registry.OrganizationsService, repositories registry.RepositoriesService) error {
	u, ok := tk.GetInput("user")
	if !ok {
		return toolkit.FailedError("user must be set")
	}

	var user github.User
	if err := json.Unmarshal([]byte(u), &user); err != nil {
		return toolkit.FailedErrorf("unable to unmarshal user\n%w", err)
	}

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

	file := filepath.Join(registry.Version, fmt.Sprintf("%s.json", namespace))
	content, _, _, err := repositories.GetContents(context.Background(), owner, repository, file, nil)
	if err2, ok := err.(*github.ErrorResponse); ok && err2.Response.StatusCode == http.StatusNotFound {
		if !resolveBool("add-if-missing", tk) {
			return toolkit.FailedErrorf("invalid namespace %s", namespace)
		}

		message := fmt.Sprintf("New Namespace: %s", namespace)
		if content, err = addNamespace(user, repositories, owner, repository, file, message); err != nil {
			return toolkit.FailedErrorf("unable to add namespace %s\n%w", namespace, err)
		}

		fmt.Println(message)
	} else if err != nil {
		return toolkit.FailedErrorf("unable to read namespace %s\n%w", namespace, err)
	}

	s, err := content.GetContent()
	if err != nil {
		return toolkit.FailedError("unable to get namespace content")
	}

	var n registry.Namespace
	if err := json.Unmarshal([]byte(s), &n); err != nil {
		return toolkit.FailedErrorf("unable to unmarshal owners\n%w", err)
	}

	if registry.IsOwner(n.Owners, registry.ByUser(*user.ID)) {
		fmt.Printf("Verified %s is an owner of %s\n", *user.Login, namespace)
		return nil
	}

	ids, err := listOrganizations(*user.Login, organizations)
	if err != nil {
		return toolkit.FailedErrorf("unable to list organizations for %s\n%w", *user.Login, err)
	}

	if registry.IsOwner(n.Owners, registry.ByOrganizations(ids)) {
		fmt.Printf("Verified %s is an owner of %s\n", *user.Login, namespace)
		return nil
	}

	return toolkit.FailedErrorf("%s is not an owner of %s", *user.Login, namespace)
}

func addNamespace(user github.User, repositories registry.RepositoriesService, owner string, repository string, path string, message string) (*github.RepositoryContent, error) {
	b, err := json.Marshal(registry.Namespace{Owners: []registry.Owner{{ID: *user.ID, Type: registry.UserType}}})
	if err != nil {
		return nil, err
	}

	createFile, _, err := repositories.CreateFile(context.Background(), owner, repository, path, &github.RepositoryContentFileOptions{
		Author: &github.CommitAuthor{
			Name:  github.String("buildpacks-bot"),
			Email: github.String("cncf-buildpacks-maintainers@lists.cncf.io"),
		},
		Message: github.String(message),
		Content: b,
	})
	if err != nil {
		return nil, err
	}

	return createFile.Content, nil
}

func listOrganizations(user string, organizations registry.OrganizationsService) ([]int64, error) {
	var ids []int64

	opt := &github.ListOptions{PerPage: 100}

	for {
		orgs, rsp, err := organizations.List(context.Background(), user, opt)
		if err != nil {
			return nil, err
		}

		for _, o := range orgs {
			ids = append(ids, *o.ID)
		}

		if rsp.NextPage == 0 {
			break
		}
		opt.Page = rsp.NextPage
	}

	return ids, nil
}

func resolveBool(name string, tk toolkit.Toolkit) bool {
	s, _ := tk.GetInput(name)
	t, err := strconv.ParseBool(s)
	if err != nil {
		return false
	}

	return t
}
