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

// Server server
type Server struct {
	c      *Config
	stopCh chan struct{}
}

// NewServer create new server
func NewServer(c *Config) *Server {
	s := &Server{
		c:      c,
		stopCh: make(chan struct{}),
	}
	return s
}

// Start start server
func (s *Server) Start() error {
	// TODO:
	return nil
}

// Stop stop server
func (s *Server) Stop() {
	select {
	case <-s.stopCh:
		return
	default:
	}
	close(s.stopCh)
}
