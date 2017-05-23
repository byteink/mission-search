// Copyright 2017 The mission-search Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package missionsearch

// Mission mission
type Mission struct {
	// CarID trello unique card id
	CardID string `json:"cid"`

	// Name mission name
	Name string `json:"name"`

	// DescData mission description data
	DescData string `json:"desc_data"`

	// CoverURL mission cover image url
	CoverURL string `json:"cover_url"`

	// Labels mission labels
	Labels []string
}
