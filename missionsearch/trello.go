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
	"time"

	"golang.org/x/time/rate"

	gtrello "github.com/VojtechVitek/go-trello"
	"github.com/boltdb/bolt"
)

const (
	missionsBoardID = "LvwOjrYP"
)

// TrelloSyncTask  sync data from trello
type TrelloSyncTask struct {
	c         *Config
	outputDir string
	stopCh    <-chan struct{}

	docID int
}

// NewTrelloSyncTask start sync with trello
// outputDir: new blotdb path
func NewTrelloSyncTask(c *Config, outputDir string, stopCh <-chan struct{}) *TrelloSyncTask {
	return &TrelloSyncTask{
		c:         c,
		outputDir: outputDir,
	}
}

// Run run task
func (t *TrelloSyncTask) Run() (db *bolt.DB, err error) {
	db, err = bolt.Open(t.outputDir, 0644, nil)
	if err != nil {
		return
	}

	var token string
	var cli *gtrello.Client
	cli, err = gtrello.NewAuthClient(t.c.TrelloAppKey, &token)
	if err != nil {
		return
	}

	rateLimiter := rate.NewLimiter(rate.Every(10*time.Second), 150)
	checkRateLimit := func() {
		r := rateLimiter.ReserveN(time.Now(), 1)
		if !r.OK() {
			time.Sleep(r.Delay())
		}
	}

	var board *gtrello.Board
	checkRateLimit()
	board, err = cli.Board(missionsBoardID)
	if err != nil {
		return
	}

	var lists []gtrello.List
	checkRateLimit()
	lists, err = board.Lists()
	if err != nil {
		return
	}

	_ = lists

	return nil, nil
}
