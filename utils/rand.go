package utils

import (
  "math/rand"
  "time"
)

var r *rand.Rand // Rand for this package.

func init() {
  r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RandomString(strlen int) string {
  const chars = "pvwxyz012345abcqrstudef6789ghijklmno"
  result := make([]byte, strlen)
  for i := range result {
    result[i] = chars[r.Intn(len(chars))]
  }
  return string(result)
}
