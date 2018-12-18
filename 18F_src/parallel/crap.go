package main

import (
  //"fmt"
  "craps"
  "runtime"
  "time"
  "log"
  "math/rand"
)

func CrapsHouseEdgeMultiProc(numTrials, numProcs int) float64 {
  count := 0
  c := make(chan int) // synch. channel OK

  // play the game in parallel
  for i := 1; i < numProcs; i++ {
    // create a PRNG for each chunk
    source := rand.NewSource(time.Now().UnixNano())
    generator := rand.New(source)

    go TotalWinOneProc(numTrials/numProcs, generator, c)
  }
  source := rand.NewSource(time.Now().UnixNano())
  generator := rand.New(source)
  // get the remaining trials
  go TotalWinOneProc(numTrials/numProcs + numTrials%numProcs, generator, c)

  // grab from channel
  for i:=0; i<numProcs; i++ {
    count += <- c
  }

  return float64(count)/float64(numTrials)
}

func TotalWinOneProc(numTrials int, generator *(rand.Rand), c chan int) {
  count := 0
  for i := 0; i<numTrials; i++ {
    play := PlayCrapsOnce(generator)
    if play { //WINNER
      count++
    } else { //LOSER
      count--
    }
  }

  c <- count

}

func PlayCrapsOnce(generator *(rand.Rand)) bool {
	roll := SumTwoDice(generator)
	if roll == 2 || roll == 3 || roll == 12 {
		return false // loser on first roll
	} else if roll == 7 || roll == 11 {
		return true // winner on first roll
	} else { // roll until we hit either a 7 or our first roll
		for true { // in practice this won't be an infinite loop
			roll2 := SumTwoDice(generator)
			if roll2 == roll { // winner!
				return true
			} else if roll2 == 7 { // loser :(
				return false
			}
		}
	}
	// there is no way we will ever make it here -- (Go doesn't know that)
	return false
}

func SumTwoDice(generator *(rand.Rand)) int {
	return DieRoll(generator) + DieRoll(generator)
}

func DieRoll(generator *(rand.Rand)) int {
  return generator.Intn(6) + 1
}

func main() {
  numTrials := 1000000

  // time the serial part
  start := time.Now()
  craps.ComputeHouseEdgeCraps(numTrials)
  elapsed := time.Since(start)
  log.Printf("Serial Craps took %s", elapsed)


  // then, time the parallel part
  start2 := time.Now()
  numProcs := runtime.NumCPU()
  CrapsHouseEdgeMultiProc(numTrials, numProcs)
  elapsed2 := time.Since(start2)
  log.Printf("Parallel Craps took %s", elapsed2)

}
