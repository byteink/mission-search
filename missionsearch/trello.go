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
	"log"
	"time"

	"golang.org/x/time/rate"

	"strings"

	"fmt"

	gtrello "github.com/VojtechVitek/go-trello"
)

const (
	missionsBoardID = "LvwOjrYP"
)

// TrelloSyncTask  sync data from trello
type TrelloSyncTask struct {
	c  *Config
	db *DB

	stopCh <-chan struct{}
}

// NewTrelloSyncTask start sync with trello
// outputDir: new blotdb path
func NewTrelloSyncTask(c *Config, db *DB, stopCh <-chan struct{}) *TrelloSyncTask {
	return &TrelloSyncTask{
		c:      c,
		db:     db,
		stopCh: stopCh,
	}
}

// Run run task
func (t *TrelloSyncTask) Run() (err error) {
	start := time.Now()
	count := 0

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

	for _, l := range lists {
		if strings.HasPrefix(l.Name, "通知") {
			fmt.Println(l.Name)
			continue
		}

		checkRateLimit()
		cards, err := l.Cards()
		if err != nil {
			log.Printf("\n\nget card failed: %v, list: %v", err, l.Name)
			return err
		}
		for _, c := range cards {
			var labels []string
			for _, l := range c.Labels {
				labels = append(labels, l.Name)
			}

			// get cover image url
			var coverURL string
			if c.IdAttachmentCover != "" {
				checkRateLimit()
				coverAttachment, err := c.Attachment(c.IdAttachmentCover)
				if err != nil {
					log.Printf("\n\nget cover attachment failed: %v, card: %v, cover: %v, url: %v", err, c.Id, c.IdAttachmentCover, c.Url)
					return err
				}
				if coverAttachment != nil {
					coverURL = coverAttachment.Url
				}
			}

			m := NewMission(c.Id, c.Name, c.Url, c.Desc, coverURL, labels)
			if err := t.db.Put(m); err != nil {
				log.Printf("\n\nput mission to db failed: %v, card: %v", err, c)
				return err
			}
			count++
		}
	}

	log.Printf("trello sync finished. num: %d, take: %v", count, time.Since(start))

	return nil
}
