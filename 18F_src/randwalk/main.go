package main

import (
  "fmt"
  "math/rand"
  "time"
)

//Check if a step will go outof the board
//Return true if (x, y) is within the board
func InField(x, y, boardWidth int) bool {
  if x <= boardWidth && y <= boardWidth && x > 0 && y > 0{
    return true
  } else {
    return false
  }
}

// return the direction of movement
func RandomDelta() int {
  rand.Seed(time.Now().UTC().UnixNano())
  randomDelta := rand.Intn(3) - 1
  return randomDelta
}

//TakeRandomStep takes the current step and return the next
func TakeRandomStep(x, y, boardWidth int) (int, int) {
  a := x
  b := y
  for ((a == x && b == y) || !InField(a, b, boardWidth)) {
    a = a + RandomDelta()
    b = b + RandomDelta()
  }
  return a, b
}

//TakeRandomWalk takes the number of steps and the size of board
func TakeRandomWalk(numSteps, boardWidth int) [][2]int {
  stepArray := make([][2]int, numSteps + 1)
  x, y := boardWidth / 2, boardWidth / 2
  stepArray[0][0], stepArray[0][1] = x, y
  for i := 1; i <= numSteps; i++ {
    x, y = TakeRandomStep(x, y, boardWidth)
    stepArray[i][0], stepArray[i][1] = x, y
  }
  return stepArray
}

func PrintRandonWalk(coordinates [][2]int) {
  for r := range coordinates {
    fmt.Print(coordinates[r][0])
    fmt.Print(" ")
    fmt.Print(coordinates[r][1])
    fmt.Println()
  }
}

func main() {
  numSteps := 28
  boardWidth := 18
  coords := TakeRandomWalk(numSteps, boardWidth)
  PrintRandonWalk(coords)
}
