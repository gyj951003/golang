package main

import (
  "fmt"
  "math"
)

type Circle struct {
  x, y float64
  radius float64
}

type Rectangle struct {
  x, y float64
  width, height, rotation float64
}

func (r Rectangle) Area() float64 {
  return r.width * r.height
}

func (c Circle) Area() float64 {
  return math.Pi * c.radius * c.radius
}

func (r *Rectangle) Translate(a, b float64) {
  r.x = r.x - a
  r.y = r.y - b
}

func (c *Circle) Translate(a, b float64) {
  c.x = c.x - a
  c.y = c.y - b
}

func main() {
  var c Circle
  c.radius = 3
  c.x = 0
  c.y = 0
  fmt.Println(c.Area())

  var ptrC *Circle
  ptrC = &c
  fmt.Println((*ptrC).x) // pointer dereference
  fmt.Println(ptrC.x) // no need in Go

  ptrC.Translate(1, 1)
  c.Translate(1, 1) // Also work!!
  fmt.Println(c.x, c.y)
}
