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

import (
	"sync/atomic"
	"unsafe"
)

// Store missions data storage
type Store struct {
	db unsafe.Pointer

	nextDBNum uint64
}

// OpenStore open DB
func OpenStore(c *Config) (*Store, error) {
	// TODO:
	return nil, nil
}

// Update
func (s *Store) Update(newDB *DB) {
	old := atomic.SwapPointer(&s.db, unsafe.Pointer(newDB))
	if old != nil {
		(*DB)(old).Close()
	}
}
