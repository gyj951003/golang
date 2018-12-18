package main

import (
  "fmt"
  "math/rand"
  "time"
)

func DieRoll() int {
  return rand.Intn(6) + 1
}

func WeightedDie() int {
  rand.Seed(time.Now().UTC().UnixNano())
  n := rand.Intn(10)
  if (n == 0 || n == 7 || n == 8 || n == 9) {
    return 3
  } else {
    return n
  }
}

func main() {
  rand.Seed(time.Now().UTC().UnixNano()) // nanosecond scale, real-time
/*
  for i := 0; i < 10; i++ {
    fmt.Println(DieRoll())
  }
  */
  fmt.Println(WeightedDie())
}
