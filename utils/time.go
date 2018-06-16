package utils

import (
  "time"
)

func GetTime() int64 {
  return time.Now().UnixNano() / int64(time.Millisecond)
}