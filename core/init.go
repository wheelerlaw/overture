// Copyright (c) 2016 shawn1m. All rights reserved.
// Use of this source code is governed by The MIT License (MIT) that can be
// found in the LICENSE file.

// Package core implements the essential features.
package core

// Initiate the server with config file
func InitServer(configFilePath string) {

	config := NewConfig(configFilePath)


	s := &Server{
		BindAddress: config.BindAddress,
		Protocols:   config.Protocols,
		Cache:       config.Cache,
		Forwarders:  config.Forwarders,
	}

	s.Run()
}
