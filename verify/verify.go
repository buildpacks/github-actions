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

package verify

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

const MetadataLabel = "io.buildpacks.buildpackage.metadata"

type Verifier struct {
	Image func(name.Reference, ...remote.Option) (v1.Image, error)

	ID      string
	Version string
	Address string
}

func (v Verifier) Verify() error {
	ref, err := name.ParseReference(v.Address)
	if err != nil {
		return fmt.Errorf("unable to parse address %s as image reference\n%w", v.Address, err)
	}

	if _, ok := ref.(name.Digest); !ok {
		return fmt.Errorf("address %s must be in digest form <host>/<repository>@sh256:<digest>", v.Address)
	}

	image, err := v.Image(ref)
	if err != nil {
		return fmt.Errorf("unable to retrieve image %s\n%w", v.Address, err)
	}

	configFile, err := image.ConfigFile()
	if err != nil {
		return fmt.Errorf("unable to retrieve config file\n%w", err)
	}

	raw, ok := configFile.Config.Labels[MetadataLabel]
	if !ok {
		return fmt.Errorf("unable to retrieve %s label", MetadataLabel)
	}

	var m metadata
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		return fmt.Errorf("unable to unmarshal %s label", MetadataLabel)
	}

	if v.ID != m.ID {
		return fmt.Errorf("invalid id in buildpackage: expected '%s', found '%s'", v.ID, m.ID)
	}

	if v.Version != m.Version {
		return fmt.Errorf("invalid version in buildpackage: expected '%s', found '%s'", v.Version, m.Version)

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
`, v.Address, m.ID, m.Version, m.Homepage, strings.Join(stacks, ", "))

	return nil
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
