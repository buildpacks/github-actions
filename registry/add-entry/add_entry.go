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
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/go-github/v32/github"

	"github.com/buildpacks/github-actions/internal/toolkit"
	"github.com/buildpacks/github-actions/registry/internal/index"
	"github.com/buildpacks/github-actions/registry/internal/services"
)

func AddEntry(tk toolkit.Toolkit, repositories services.RepositoriesService) error {
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

	address, ok := tk.GetInput("address")
	if !ok {
		return toolkit.FailedError("address must be set")
	}

	file := index.Path(namespace, name)
	content, _, _, err := repositories.GetContents(context.Background(), owner, repository, file, nil)
	if err2, ok := err.(*github.ErrorResponse); ok && err2.Response.StatusCode == http.StatusNotFound {
		fmt.Printf("New Index: %s\n", name)
		content = &github.RepositoryContent{}
	} else if err != nil {
		return toolkit.FailedErrorf("unable to read index %s\n%w", name, err)
	}

	s, err := content.GetContent()
	if err != nil {
		return toolkit.FailedErrorf("unable to get index content\n%w", err)
	}

	entries, err := unmarshalEntries(s)
	if err != nil {
		return toolkit.FailedErrorf("unable to unmarshal entries\n%w", err)
	}

	if contains(entries, namespace, version) {
		return toolkit.FailedErrorf("index %s already has namespace %s and version %s", name, namespace, version)
	}

	entries = append(entries, index.Entry{
		Namespace: namespace,
		Name:      name,
		Version:   version,
		Address:   address,
	})

	s, err = marshalEntries(entries)
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

	fmt.Printf("Added %s/%s@%s\nx", namespace, name, version)

	return nil
}

func contains(entries []index.Entry, namespace string, version string) bool {
	for _, e := range entries {
		if e.Namespace == namespace && e.Version == version {
			return true
		}
	}

	return false
}

func marshalEntries(entries []index.Entry) (string, error) {
	b := &bytes.Buffer{}
	j := json.NewEncoder(b)

	for _, e := range entries {
		if err := j.Encode(e); err != nil {
			return "", err
		}
	}

	return b.String(), nil
}

func unmarshalEntries(content string) ([]index.Entry, error) {
	var entries []index.Entry

	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		var e index.Entry
		if err := json.Unmarshal(scanner.Bytes(), &e); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}

	return entries, scanner.Err()
}
