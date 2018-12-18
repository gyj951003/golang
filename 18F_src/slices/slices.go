package main

import "fmt"

func main() {
  fmt.Print(FindSquares([]int{100, -4, 1, 13, 0, -47, 16, 8}))
  return
}


func IsSquare(n int) bool {
  if n >= 0 {
    for i := 0; i <= n; i++ {
      if i * i == n {
        return true
      }
    }
  }
  return false
}

func FindSquares(intSlice []int) []int {
  var squaresSlice []int
  for i := 0; i < len(intSlice); i++ {
    if IsSquare(intSlice[i]) == true {
      squaresSlice = append(squaresSlice, intSlice[i])
    }
  }
  return squaresSlice
}
