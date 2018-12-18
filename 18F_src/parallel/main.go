package main

import (
  "fmt"
  "time"
  //"runtime"
  "strconv"
)

func f(n int) {
  for i := 0; i < 10; i++ {
    fmt.Println(n, ":", i)
    time.Sleep(time.Millisecond * 100) // so that can see other function loaded in before this one finishes
  }
}

// A channel can store a value of a given type and allow functions to communicate values
// A private function
func pinger(c chan string) {
  for i:= 0;; i++ {
    c <- "ping " + strconv.Itoa(i) // c<-: put sth into the channel
  }
}

func printer(c chan string) {
  for {
    msg := <- c // <-c: retrive sth from channel
    fmt.Println(msg)
  }
}

// Exercise: write a funct that Parallelize the work of summing the element in an array of integers into to pieces


func main() {
  /*
  fmt.Println("Parallel programming in Go!")
  fmt.Println("We have", runtime.NumCPU(), "cores!") // number of cores available
  runtime.GOMAXPROCS(1) //limit the max number of processors to be 1
  //Multiple processors in fact slows down serial algorithm!!!

  for i := 0; i < 10; i++ {
    go f(i)
  }

  var input string
  fmt.Scanln(&input) // so that go wait for us for an input. "go" won't stop before return

  */

  /*
  c := make(chan string)
  go pinger(c)
  go printer(c)
  // printed in order, why?
  // channel blocking. When pinger() tries to place an element into the channel, it waits until the other end of the channel is ready to receive it.

  var input string
  fmt.Scanln(&input)
  */

  a := make([]int, 1000)
  for i  := range a {
    a[i] = 1*i*i + 3*i -1
  }


  c2:= make(chan int)
  //c3:= make(chan int)
  go SumParallel(a[:len(a)/2], c2)
  go SumParallel(a[len(a)/2:], c2)

  fmt.Println(<-c2 + <-c2)
/*
  start := time.Now()
  sumSerial := Sumserial(a)
  elapse := time.Since(start)
  // Devide summation task into numPros parts
  numProcs := runtime.NumCPU()
  start2 := time.Now()
  elapse2 := time.Since(start2)
  sumParallel := SumMultiProc(a, numProcs)

  fmt.Println(sumSerial, sumParallel, elapse, elapse2)


  n := 10
  c4 := make(chan int, n) // buffered channel with cap length == n
  for k := 0; k < n; k++ {
    go Push(k, c4)
  }
  for k := 0; k < n; k++ {
    fmt.Println(<-c4) // FIFO
  }

  */
  fmt.Println(SumMultiProc(a, 5))
  fmt.Println(SumMultiProc(a, 6))
  fmt.Println(SumMultiProc(a, 7))
}

// Example of buffered channel
func Push(k int, c chan int) {
  c <- k
}


func Sumserial(a []int) int {
  s:= 0
  for _, val := range a {
    s = s + val
  }
  return s
}

func SumParallel(a []int, c chan int) {
  s := 0
  for _, val := range a {
    s += val
  }
  c <- s
}

func SumMultiProc(a []int, numProcs int) int {
    n := len(a)
    s := 0
    c := make(chan int)

    for i:= 0; i < numProcs; i++ {
      if i == numProcs - 1 {
        go SumParallel(a[i*(n/numProcs):], c)
      } else {
        go SumParallel(a[i*(n/numProcs):(i+1)*(n/numProcs)], c)
      }
    }

    for i:= 0; i < numProcs; i++ { // !!Won't be parallel if put s += <- c into the loop above
      s += <-c
    }

    return s
}
