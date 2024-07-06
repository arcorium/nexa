package util

import "time"

func RoundTimeToSecond(times time.Time) time.Time {
  sec := times.Unix()
  return time.Unix(sec, 0)
}
