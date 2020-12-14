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
	"strconv"

	"github.com/google/go-github/v32/github"
	"gopkg.in/retry.v1"

	"github.com/buildpacks/github-actions/internal/toolkit"
	"github.com/buildpacks/github-actions/registry/internal/namespace"
	"github.com/buildpacks/github-actions/registry/internal/services"
)

func VerifyNamespaceOwner(tk toolkit.Toolkit, organizations services.OrganizationsService, repositories services.RepositoriesService, strategy retry.Strategy) error {
	c, err := parseConfig(tk)
	if err != nil {
		return err
	}

	var user github.User
	if err := json.Unmarshal([]byte(c.User), &user); err != nil {
		return toolkit.FailedErrorf("unable to unmarshal user\n%w", err)
	}

	n, err := getNamespace(tk, c, user, repositories, strategy)
	if err != nil {
		return err
	}

	if namespace.IsOwner(n.Owners, namespace.ByUser(*user.ID)) {
		fmt.Printf("Verified %s is an owner of %s\n", *user.Login, c.Namespace)
		return nil
	}

	ids, err := listOrganizations(*user.Login, organizations)
	if err != nil {
		return toolkit.FailedErrorf("unable to list organizations for %s\n%w", *user.Login, err)
	}

	if namespace.IsOwner(n.Owners, namespace.ByOrganizations(ids)) {
		fmt.Printf("Verified %s is an owner of %s\n", *user.Login, c.Namespace)
		return nil
	}

	return toolkit.FailedErrorf("%s is not an owner of %s", *user.Login, c.Namespace)
}

type config struct {
	User         string
	Owner        string
	Repository   string
	Namespace    string
	AddIfMissing bool
}

func parseConfig(tk toolkit.Toolkit) (config, error) {
	var (
		c  = config{AddIfMissing: false}
		ok bool
	)

	c.User, ok = tk.GetInput("user")
	if !ok {
		return config{}, toolkit.FailedError("user must be set")
	}

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

	if s, ok := tk.GetInput("add-if-missing"); ok {
		if t, err := strconv.ParseBool(s); err == nil {
			c.AddIfMissing = t
		}
	}

	return c, nil
}

func getNamespace(tk toolkit.Toolkit, c config, user github.User, repositories services.RepositoriesService, strategy retry.Strategy) (namespace.Namespace, error) {
	file := namespace.Path(c.Namespace)

	for a := retry.Start(strategy, nil); a.Next(); {
		content, _, resp, err := repositories.GetContents(context.Background(), c.Owner, c.Repository, file, nil)
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			if !c.AddIfMissing {
				return namespace.Namespace{}, toolkit.FailedErrorf("invalid namespace %s", c.Namespace)
			}

			b, err := json.Marshal(namespace.Namespace{Owners: []namespace.Owner{{ID: *user.ID, Type: namespace.UserType}}})
			if err != nil {
				return namespace.Namespace{}, toolkit.FailedErrorf("unable to decode namespace\n%w", err)
			}

			if _, resp, err := repositories.CreateFile(context.Background(), c.Owner, c.Repository, file, &github.RepositoryContentFileOptions{
				Author: &github.CommitAuthor{
					Name:  github.String("buildpacks-bot"),
					Email: github.String("cncf-buildpacks-maintainers@lists.cncf.io"),
				},
				Message: github.String(fmt.Sprintf("New Namespace: %s", c.Namespace)),
				Content: b,
			}); resp != nil && resp.StatusCode == http.StatusConflict {
				tk.Warning("retrying namespace update after conflict")
				continue
			} else if err != nil {
				return namespace.Namespace{}, toolkit.FailedErrorf("unable to create namespace\n%w", err)
			}

			fmt.Printf("New Namespace: %s\n", c.Namespace)
			continue
		} else if err != nil {
			return namespace.Namespace{}, toolkit.FailedErrorf("unable to read namespace %s\n%w", c.Namespace, err)
		}

		s, err := content.GetContent()
		if err != nil {
			return namespace.Namespace{}, toolkit.FailedErrorf("unable to get namespace content\n%w", err)
		}

		var n namespace.Namespace
		if err := json.Unmarshal([]byte(s), &n); err != nil {
			return namespace.Namespace{}, toolkit.FailedErrorf("unable to unmarshal owners\n%w", err)
		}

		return n, nil
	}

	return namespace.Namespace{}, toolkit.FailedError("timed out")
}

func listOrganizations(user string, organizations services.OrganizationsService) ([]int64, error) {
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
