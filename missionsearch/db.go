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
	"encoding/binary"
	"encoding/json"
	"errors"

	"github.com/boltdb/bolt"
	"github.com/huichen/sego"
)

var bucketMission = []byte{'m'}
var bucketIndex = []byte{'i'}

type DB struct {
	b *bolt.DB

	segmenter sego.Segmenter

	docID   uint64
	scratch []byte
}

// OpenDB open or create new mission db
// path: db path
// dicts: segmenter dictionary file paths(separated by comma)
func OpenDB(path, dicts string) (*DB, error) {
	// init bolt db
	b, err := bolt.Open(path, 0644, nil)
	if err != nil {
		return nil, err
	}
	err = b.Batch(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucket(bucketMission); err != nil {
			return err
		}

		if _, err := tx.CreateBucket(bucketIndex); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	db := &DB{
		b:       b,
		docID:   1,
		scratch: make([]byte, 8),
	}

	db.segmenter.LoadDictionary(dicts)

	return db, nil
}

func (db *DB) Put(m *Mission) error {
	if m == nil || m.Title == "" {
		return errors.New("invalid mission")
	}

	// store mission
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	binary.BigEndian.PutUint64(db.scratch, db.docID)
	err = db.b.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketMission)
		return bucket.Put(db.scratch, data)
	})
	if err != nil {
		return err
	}

	// create index
	text := []byte(m.Title) // TODO: unsafe convert
	segments := db.segmenter.Segment(text)
	if len(segments) == 0 {
		return errors.New("segment result is empty.")
	}
	for _, seg := range segments {
		term := seg.Token().Text()
		if isIndexTerm(term) {
			if err := db.updateIndex(term, db.scratch); err != nil {
				return err
			}
		}
	}

	db.docID++
	return nil
}

func (db *DB) Search(query string) ([]*Mission, error) {
	segments := db.segmenter.Segment([]byte(query))
	if len(segments) == 0 {
		return nil, nil
	}
	for _, seg := range segments {
		term := seg.Token().Text()
		if isIndexTerm(term) {
			if err := db.updateIndex(term, db.scratch); err != nil {
				return nil, err
			}
		}
	}
	// TODO:
	return nil, nil
}

func (db *DB) updateIndex(term string, docID []byte) error {
	err := db.b.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketIndex)
		docs := bucket.Get([]byte(term))
		docs = append(docs, docID...)
		return bucket.Put([]byte(term), docs)
	})
	return err
}

func (db *DB) queryIndex(term string) (docs []uint64, err error) {
	var docsData []byte
	err = db.b.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketIndex)
		docsData = bucket.Get([]byte(term))
		return nil
	})
	if err != nil {
		return
	}
	for i := 0; i < len(docsData)/8; i++ {
		docs = append(docs, binary.BigEndian.Uint64(docsData[i*8:]))
	}
	return
}

func (db *DB) Close() {
	if err := db.b.Close(); err != nil {
		// TODO: log
	}
}

func isIndexTerm(term string) bool {
	return term != "[" && term != "]" && term != " "
}
