package timeutil

import "time"

// MilliSecond returns time.Time in millisecond
func MilliSecond(t time.Time) int64 {
	return t.UnixNano() / time.Millisecond.Nanoseconds()
}

// UnixMillis returns millisecond in time.Time
func UnixMillis(t int64) time.Time {
	return time.Unix(0, t*time.Millisecond.Nanoseconds())
}

// NowMS return now time in millisecond
func NowMS() int64 {
	return MilliSecond(time.Now())
}

// UTCNowMs return utc+0 time in millisecond
func UTCNowMs() int64 {
	return MilliSecond(time.Now().UTC())
}