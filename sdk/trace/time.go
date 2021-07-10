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

package trace

import (
	"time"
)

// Clock is the entrypoint for providing time to span's start/end timestamp.
// By default the standard "time" package will be used. User can replace
// it with customized clock implementation (e.g. has additional clock
// synchronization logic) by using the `WithClock` option.
type Clock interface {
	Start() (time.Time, Stopwatch)
}

type Stopwatch interface {
	Stop(time.Time) time.Time
}

func defaultClock() Clock {
	return standardClock{}
}

func defaultStopwatch() Stopwatch {
	return standardStopwatch{}
}

type standardClock struct{}
type standardStopwatch struct{}

func (standardStopwatch) Stop(t time.Time) time.Time {
	return t.Add(time.Since(t))
}

func (standardClock) Start() (time.Time, Stopwatch) {
	return time.Now(), standardStopwatch{}
}