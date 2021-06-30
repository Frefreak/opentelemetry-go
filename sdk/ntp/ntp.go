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
	"fmt"
	"time"

	cli "github.com/beevik/ntp"
)

var ticker *time.Ticker
var clockOffset time.Duration

func GetClockOffset() time.Duration {
	return clockOffset
}

type Config struct {
	host         string
	interval     time.Duration
	queryOptions cli.QueryOptions
	verbose      bool
}

func NewConfig(host string) *Config {
	return &Config{
		host:     host,
		interval: time.Minute,
	}
}

func (cfg *Config) WithInterval(d time.Duration) *Config {
	cfg.interval = d
	return cfg
}
func (cfg *Config) WithVerbose(b bool) *Config {
	cfg.verbose = b
	return cfg
}

func (cfg *Config) WithQueryOptions(opts cli.QueryOptions) *Config {
	cfg.queryOptions = opts
	return cfg
}

func ShouldStart(config *Config) bool {
	if config == nil {
		return false
	}
	return config.host != ""
}

func updateClockOffset(config *Config) {
	resp, err := cli.QueryWithOptions(config.host, config.queryOptions)
	if err != nil {
		fmt.Printf("error querying %s: %v\n", config.host, err)
		return
	}
	if config.verbose {
		fmt.Println("got ClockOffset: ", resp.ClockOffset)
	}
	clockOffset = resp.ClockOffset
}

func StartNTPWorker(config *Config) {
	if config.verbose {
		fmt.Println("ntp worker starting")
	}
	// do one update first
	updateClockOffset(config)
	ticker = time.NewTicker(config.interval)
	go func() {
		for range ticker.C {
			updateClockOffset(config)
		}
		if config.verbose {
			fmt.Println("ntp worker exiting")
		}
	}()
}

func StopNTPWorker(config *Config) {
	if ticker == nil {
		return
	}
	if config.verbose {
		fmt.Println("stopping ntp worker")
	}
	ticker.Stop()
}
