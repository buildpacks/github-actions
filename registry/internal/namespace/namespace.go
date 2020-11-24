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

package namespace

var restrictedNamespaces = []string{
	"buildpack",
	"buildpack-io",
	"buildpack.io",
	"buildpackio",
	"buildpacks",
	"buildpacks-io",
	"buildpacks.io",
	"buildpacksio",
	"cnb",
	"cnbs",
	"cncf",
	"cncf-cnb",
	"cncf-cnbs",
	"example",
	"examples",
	"official",
	"pack",
	"sample",
	"samples",
}

type Namespace struct {
	Owners []Owner
}

func IsRestricted(namespace string) bool {
	for _, n := range restrictedNamespaces {
		if n == namespace {
			return true
		}
	}

	return false
}
