package main

import "fmt"
import "math/rand"
import "time"

func DieRoll() int {
  rand.Seed(time.Now().UTC().UnixNano())
  return rand.Intn(6) + 1
}

func  SumTwoDice() int {
  return DieRoll() + DieRoll()
}

func main() {
  for i := 1; i < 100; i++{
    fmt.Println(SumTwoDice())
  }
}
