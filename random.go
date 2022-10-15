package omw

import (
  "fmt"
  "math/rand"
)

func RandomInt(start, end int, random *rand.Rand) (int, error) {
  if start >= end {
    return 0, fmt.Errorf("start < end でなければならない")
  }
  return random.Intn(end - start) + start, nil
}

func RandomBool(random *rand.Rand) bool {
  return random.Intn(2) == 0
}
