// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ntp

import (
	"log"
	"runtime/debug"
	"time"

	cli "github.com/beevik/ntp"
)

var ticker *time.Ticker
var clockOffset time.Duration

// GetClockOffset returns current clock offset with the remote ntp server
func GetClockOffset() time.Duration {
	return clockOffset
}

// Config is the relevant ntp config
type Config struct {
	host         string
	interval     time.Duration
	queryOptions cli.QueryOptions
	verbose      bool
}

// NewConfig returns a new default config with host specified and 1 minute interval
func NewConfig(host string) *Config {
	return &Config{
		host:     host,
		interval: time.Minute,
	}
}

// WithInterval setting querying interval
func (cfg *Config) WithInterval(d time.Duration) *Config {
	cfg.interval = d
	return cfg
}

// WithVerbose setting if verbose is enabled
func (cfg *Config) WithVerbose(b bool) *Config {
	cfg.verbose = b
	return cfg
}

// WithQueryOptions settings provides access to more query options
func (cfg *Config) WithQueryOptions(opts cli.QueryOptions) *Config {
	cfg.queryOptions = opts
	return cfg
}

// ShouldStart determine if we should start
func ShouldStart(config *Config) bool {
	if config == nil {
		return false
	}
	return config.host != ""
}

func updateClockOffset(config *Config) {
	defer func() {
		if panicInfo := recover(); panicInfo != nil {
			log.Printf("%v, %s", panicInfo, string(debug.Stack()))
		}
	}()

	resp, err := cli.QueryWithOptions(config.host, config.queryOptions)
	if err != nil {
		log.Printf("error querying %s: %v\n", config.host, err)
		return
	}
	if config.verbose {
		log.Println("got ClockOffset: ", resp.ClockOffset)
	}
	clockOffset = resp.ClockOffset
}

// StartNTPWorker start querying goroutine
func StartNTPWorker(config *Config) {
	if config.verbose {
		log.Println("ntp worker starting")
	}
	// do one update first
	updateClockOffset(config)
	ticker = time.NewTicker(config.interval)
	go func() {
		for range ticker.C {
			updateClockOffset(config)
		}
		if config.verbose {
			log.Println("ntp worker exiting")
		}
	}()
}

// StopNTPWorker stop the goroutine
func StopNTPWorker(config *Config) {
	if ticker == nil {
		return
	}
	if config != nil && config.verbose {
		log.Println("stopping ntp worker")
	}
	ticker.Stop()
}
