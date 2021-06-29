package ntp

import "time"

var timeOffset time.Duration

func GetTimeOffset() time.Duration {
	return timeOffset
}

func SetTimeOffset(off time.Duration) error {
	timeOffset = off
	return nil
}