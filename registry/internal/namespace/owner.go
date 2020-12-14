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

const (
	OrganizationType = "github_org"
	UserType         = "github_user"
)

type Owner struct {
	ID   int64  `json:"id"`
	Type string `json:"type"`
}

type OwnerPredicate func(Owner) bool

func IsOwner(owners []Owner, predicate OwnerPredicate) bool {
	for _, o := range owners {
		if predicate(o) {
			return true
		}
	}

	return false
}

func ByUser(id int64) OwnerPredicate {
	return func(owner Owner) bool {
		return owner.Type == UserType && owner.ID == id
	}
}

func ByOrganizations(ids []int64) OwnerPredicate {
	return func(owner Owner) bool {
		for _, id := range ids {
			if owner.Type == OrganizationType && owner.ID == id {
				return true
			}
		}
		return false
	}
}
