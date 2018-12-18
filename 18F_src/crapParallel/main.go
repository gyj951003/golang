package main

import (
  "fmt"
  "math/rand"
)

func CrapsHouseEdgeMultiProc(n int, numProcs int) float64 {
  count := 0
  c := make(chan int)
  for i := 1; i < numProcs; i++ {
    go TotalWinOneProc(n/numProcs, c)
  }

  go TotalWinOneProc(n/numProcs + n%numProcs, c)

  for i:= 0; i < numProcs; i++ {
    count += <- c
  }

  return float64(n)/float64(count)
}

func TotalWinOneProc(n int, c chan int) {
  count := 0
  for i := 0; i < n; i++ {
    count = count + PlayCrapsOnce()
  }

  c <- count
  return
}

func DieRoll() int {
	return rand.Intn(6) + 1
}

//SumTwoDice takes no inputs and returns a random number between 2 and 12
//simulating the sum of two dice. It uses DieRoll as a subroutine.
func SumTwoDice() int {
	return DieRoll() + DieRoll()
}

func SumMultipleDice(numDice int) int {
	s := 0
	for i := 0; i < numDice; i++ {
		s += DieRoll()
	}
	return s
}

func PlayCrapsOnce() int {
	roll := SumTwoDice()
  result := 0
	if roll == 2 || roll == 3 || roll == 12 {
		result = 0 // loser on first roll
	} else if roll == 7 || roll == 11 {
		result = 1 // winner on first roll
	} else { // roll until we hit either a 7 or our first roll
		for true { // in practice this won't be an infinite loop
			roll2 := SumTwoDice()
			if roll2 == roll { // winner!
        break
				result = 1
			} else if roll2 == 7 { // loser :(
        break
				result = 0
			}
		}
	}
	return result
}

func main() {
  fmt.Println(CrapsHouseEdgeMultiProc(100, 4))
}
