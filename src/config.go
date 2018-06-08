// Copyright (c) 2016 shawn1m. All rights reserved.
// Use of this source code is governed by The MIT License (MIT) that can be
// found in the LICENSE file.

package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	log "github.com/Sirupsen/logrus"
)

const DEFAULT_CONFIG = `
{
  "BindAddress": "127.0.0.1:5053",
  "Protocols": ["udp", "tcp"],
  "Forwarders": [
    {
      "Name": "Local",
      "Address": "127.0.0.1:53",
      "Timeout": 10
    }
  ],
  "CacheSize" : 3000
}
`

type Config struct {
	BindAddress        string `json:"BindAddress"`
	Protocols          []string
	Forwarders         []*Forwarder
	CacheSize          int

	Cache           *Cache
}

func (config *Config) UnmarshalJSON(b []byte) error {

	type ConfigAlias Config

	configAlias := new(ConfigAlias)

	if err := json.Unmarshal(b, &configAlias); err != nil {
		return err
	}

	if len(configAlias.Protocols) == 0 {
		configAlias.Protocols = []string{"tcp", "udp"}
	}

	*config = Config(*configAlias)



	return nil
}

// New config with json file and do some other initiate works
func NewConfig(configFile string) *Config {

	config := parseJson(configFile)

	config.Cache = New(config.CacheSize)
	if config.CacheSize > 0 {
		log.Info("CacheSize is " + strconv.Itoa(config.CacheSize))
	} else {
		log.Info("Cache is disabled")
	}

	return config
}

func parseJson(path string) *Config {

	config := new(Config)

	f, err := os.Open(path)
	if err != nil {
		log.Warn("Open config file failed: ", err)
		_ = json.Unmarshal([]byte(DEFAULT_CONFIG), config)

	} else {
		defer f.Close()

		rawJson, err := ioutil.ReadAll(f)
		if err != nil {
			log.Fatal("Read config file failed: ", err)
			os.Exit(1)
		}

		err = json.Unmarshal(rawJson, config)
		if err != nil {
			log.Fatal("Json syntex error: ", err)
			os.Exit(1)
		}
	}

	return config
}

