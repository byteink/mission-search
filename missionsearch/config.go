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
	"errors"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

const defaultDataDir = "./data"

// Config the config
type Config struct {
	DataDir string `yaml:"data_dir"`

	TrelloAppKey string `yaml:"trello_appkey"`

	TelegramBotToken string `yaml:"telegram_bot_token"`

	SegoDicts string `yaml:"sego_dicts"`
}

// Validate check config validate
func (c *Config) Validate() error {
	if c == nil {
		return errors.New("config is nil")
	}
	if c.DataDir == "" {
		return errors.New("invalid data directory")
	}
	if c.TrelloAppKey == "" {
		return errors.New("invalid trello appkey")
	}
	if c.TelegramBotToken == "" {
		return errors.New("invalid telegram bot token")
	}
	if c.SegoDicts == "" {
		return errors.New("invalid sego dictionarys")
	}

	return nil
}

// ParseConfigFile parse config from file
func ParseConfigFile(conf string) (*Config, error) {
	data, err := ioutil.ReadFile(conf)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	if cfg.DataDir == "" {
		cfg.DataDir = defaultDataDir
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}
