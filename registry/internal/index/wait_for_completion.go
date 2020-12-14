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
	"context"
	"fmt"

	"gopkg.in/retry.v1"

	"github.com/buildpacks/github-actions/internal/toolkit"
	"github.com/buildpacks/github-actions/registry/internal/services"
)

const (
	RequestFailureLabel = "failure"
	RequestSuccessLabel = "succeeded"
)

func WaitForCompletion(number int, url string, tk toolkit.Toolkit, issues services.IssuesService, strategy retry.Strategy) error {
	for a := retry.Start(strategy, nil); a.Next(); {
		issue, _, err := issues.Get(context.Background(), "buildpacks", "registry-index", number)
		if err != nil {
			tk.Warningf("unable to get state for %s", url)
			continue
		}

		for _, l := range issue.Labels {
			if *l.Name == RequestFailureLabel {
				return toolkit.FailedErrorf("Registry request %s failed", url)
			} else if *l.Name == RequestSuccessLabel {
				fmt.Printf("Registry request %s succeeded\n", url)
				return nil
			}
		}
	}

	return toolkit.FailedError("timed out waiting for request to be processed")
}
