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

package info

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/buildpacks/libcnb"
	"github.com/pelletier/go-toml"
)

type BuildpackInfo struct {
	Path   string
	Writer io.Writer
}

func (b BuildpackInfo) Inform() error {
	c, err := ioutil.ReadFile(b.Path)
	if err != nil {
		return fmt.Errorf("unable to read %s\n%w", b.Path, err)
	}

	var bp libcnb.Buildpack
	if err := toml.Unmarshal(c, &bp); err != nil {
		return fmt.Errorf("unable to unmarhal %s\n%w", b.Path, err)
	}

	_, _ = fmt.Fprintf(b.Writer, "::set-output name=id::%s\n", bp.Info.ID)
	_, _ = fmt.Fprintf(b.Writer, "::set-output name=name::%s\n", bp.Info.Name)
	_, _ = fmt.Fprintf(b.Writer, "::set-output name=version::%s\n", bp.Info.Version)
	_, _ = fmt.Fprintf(b.Writer, "::set-output name=homepage::%s\n", bp.Info.Homepage)

	return nil
}
