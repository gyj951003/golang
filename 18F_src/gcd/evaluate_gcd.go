package main
import (
  "fmt"
  "math"
  "math/rand"
  "time"
  "os"
  "log"
  "strconv"
)

func main () {
  f, err := os.OpenFile("table.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatal(err)
    }
  line := "TrivialGCD time cost, EuclidGCD time cost, Speed Up\n"
  if _, err := f.WriteString(line); err != nil {
    log.Fatal(err)
  }

  for i := 3; i <= 6; i++ {
    lowerbound := 1 * int(math.Pow(10.0, float64(i)))
    upperbound := 2 * int(math.Pow(10.0, float64(i)))
    numberPairs := 10
    elapsedT := 0 * time.Nanosecond
    elapsedE := 0 * time.Nanosecond
    for j := 0; j < numberPairs; j++ {
      m := RandomNumber (lowerbound, upperbound)
      n := RandomNumber (lowerbound, upperbound)
      startT := time.Now()
      gcdT := TrivialGCD(m, n)
      startE := time.Now()
      elapsedT = startE.Sub(startT) + elapsedT
      gcdE := EuclidGCD(m, n)
      endE := time.Now()
      elapsedE = endE.Sub(startE) + elapsedE
      fmt.Println(m, n, gcdT, gcdE, startE.Sub(startT), endE.Sub(startE))
    }
    speedUp := int(elapsedT / elapsedE)
    // Write to a table
    line = elapsedT.String() + "," + elapsedE.String() + "," + strconv.Itoa(speedUp) + "\n"
    if _, err := f.WriteString(line); err != nil {
      log.Fatal(err)
    }
  }

  if err := f.Close(); err != nil {
      log.Fatal(err)
  }
}

func RandomNumber (lowerbound, upperbound int) int {
  rand.Seed(time.Now().UTC().UnixNano())
  return rand.Intn(upperbound - lowerbound + 1) + lowerbound
}

func Max(a, b int) int {
  if a < b {
    a = b
  }
  return a
}

func Min(a, b int) int {
  if a > b {
    a = b
  }
  return a
}

// This function takes two integers and
// output their gcd by trivial method and steps taken to get gcd.
func TrivialGCD(m, n int) (int) {
  min := Min(m, n)
  gcd := 1
  for i := min; i >= 1; i-- {
    if m % i == 0 && n % i == 0 {
      gcd = i
      break
    }
  }
  return gcd
}

// This function takes two integers and
// output their gcd by euclid method and steps taken to get gcd.
func EuclidGCD(m, n int) (int) {
  for m != n {
    nextm := Min(m, n)
    nextn := Max(m, n) - Min(m, n)
    m = nextm
    n = nextn
  }
  return m
}
