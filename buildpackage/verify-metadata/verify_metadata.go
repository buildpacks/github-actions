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
	"fmt"
	"strings"

	"github.com/google/go-containerregistry/pkg/name"

	"github.com/buildpacks/github-actions/internal/toolkit"
)

const MetadataLabel = "io.buildpacks.buildpackage.metadata"

func VerifyMetadata(tk toolkit.Toolkit, imageFn ImageFunction) error {
	c, err := parseConfig(tk)
	if err != nil {
		return err
	}

	ref, err := name.ParseReference(c.Address)
	if err != nil {
		return toolkit.FailedErrorf("unable to parse address %s as image reference", c.Address)
	}

	if _, ok := ref.(name.Digest); !ok {
		return toolkit.FailedErrorf("address %s must be in digest form <host>/<repository>@sh256:<digest>", c.Address)
	}

	image, err := imageFn(ref)
	if err != nil {
		return toolkit.FailedErrorf("unable to retrieve image %s", c.Address)
	}

	configFile, err := image.ConfigFile()
	if err != nil {
		return toolkit.FailedErrorf("unable to retrieve config file\n%w", err)
	}

	raw, ok := configFile.Config.Labels[MetadataLabel]
	if !ok {
		return toolkit.FailedErrorf("unable to retrieve %s label", MetadataLabel)
	}

	var m metadata
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		return toolkit.FailedErrorf("unable to unmarshal %s label", MetadataLabel)
	}

	if c.ID != m.ID {
		return toolkit.FailedErrorf("invalid id in buildpackage: expected %s, found %s", c.ID, m.ID)
	}

	if c.Version != m.Version {
		return toolkit.FailedErrorf("invalid version in buildpackage: expected %s, found %s", c.Version, m.Version)

	}

	var stacks []string
	for _, s := range m.Stacks {
		stacks = append(stacks, s.ID)
	}

	fmt.Printf(`Verified %s
  ID:       %s
  Version:  %s
  Homepage: %s
  Stacks:   %s
`, c.Address, m.ID, m.Version, m.Homepage, strings.Join(stacks, ", "))

	return nil
}

type config struct {
	ID      string
	Version string
	Address string
}

func parseConfig(tk toolkit.Toolkit) (config, error) {
	var (
		c  config
		ok bool
	)

	c.ID, ok = tk.GetInput("id")
	if !ok {
		return config{}, toolkit.FailedError("id must be specified")
	}

	c.Version, ok = tk.GetInput("version")
	if !ok {
		return config{}, toolkit.FailedError("version must be specified")
	}

	c.Address, ok = tk.GetInput("address")
	if !ok {
		return config{}, toolkit.FailedError("address must be specified")
	}

	return c, nil
}

type metadata struct {
	ID       string
	Version  string
	Homepage string
	Stacks   []stack
}

type stack struct {
	ID string
}
