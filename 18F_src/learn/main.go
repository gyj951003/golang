package main

import (
  "fmt"
)


func TestDefaltVariable() {
  var a int
  var b uint
  var c byte
  var d float64
  var e string
  var f bool
  fmt.Println(a) //0
  fmt.Println(b) //0
  fmt.Println(c) //0
  fmt.Println(d) //0
  fmt.Println(e) //""
  fmt.Println(f) //false
}

func TestSlice() {
  s1 := []int{1, 2, 3, 4}
  s2 := []int{0, 0, 0}
  s2 = append(s2, s1 ...) //use ... to append multiple items
  fmt.Println(s2)
}


func TestScope(x uint, y uint, s []int, a [5]int) {
  x = x - 1
  y = y - 1
  s[1] = 5
  a[1] = 5
}

func TestStrings() {
  s1 := "hi! there"
  //s2 := "hi!"
  //s1[3] = 't'
  s1 = s1[:3] + "t" + s1[4:]
  fmt.Println(s1)
}

func SliceField(m, n int) {
  field := make([][]int, n)
  for i := range field {
    field[i] = make([]int, m)
  }
}

func main () {
  //TestDefaltVariable()
  //TestSlice()
  var x uint = 1
  var y uint = 1


  s := []int{1, 2, 3, 4, 5} // original when pass
  a := [5]int{1, 2, 3, 4, 5} // copy when pass to the func
  TestScope(x, y, s, a)

  fmt.Println(x, y)
  fmt.Println(s, a)

  TestStrings()
  SliceField(3,5)
}
