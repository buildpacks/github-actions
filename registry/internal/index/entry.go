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

package index

import (
	"bufio"
	"bytes"
	"encoding/json"
	"strings"
)

type Entry struct {
	Namespace string `json:"ns"`
	Name      string `json:"name"`
	Version   string `json:"version"`
	Yanked    bool   `json:"yanked"`
	Address   string `json:"addr"`
}

func MarshalEntries(entries []Entry) (string, error) {
	b := &bytes.Buffer{}
	j := json.NewEncoder(b)

	for _, e := range entries {
		if err := j.Encode(e); err != nil {
			return "", err
		}
	}

	return b.String(), nil
}

func UnmarshalEntries(content string) ([]Entry, error) {
	var entries []Entry

	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		var e Entry
		if err := json.Unmarshal(scanner.Bytes(), &e); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}

	return entries, scanner.Err()
}
