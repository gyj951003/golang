package main

import (
  "os"
  "fmt"
  "strconv"
  "time"
  "math/rand"
  "runtime"
)

type GameBoard [][]int

func InitializeCentralBoard(size, pile int) GameBoard {
  board := make(GameBoard, size)

  // Allocate and initialze each entry to 0
  for i := range(board) {
    board[i] = make([]int, size)
    for j := range(board[i]) {
      board[i][j] = 0
    }
  }

  board[size/2][size/2] = pile
  return board
}

func InitializeRandomBoard(size, pile int) GameBoard {
  board := make(GameBoard, size)

  for i := range(board) {
    board[i] = make([]int, size)
    for j := range(board[i]) {
      board[i][j] = 0
    }
  }

  blockList := make([][2]int, 100)

  for i := 0; i < 100; i++ {
    r, c := SelectRandomBlock(size)
    blockList[i] = [2]int{r, c}
  }

  b := 0

  for pile > 0 {
    pos := blockList[b % 100]
    amount := rand.Intn(pile * 2 / 100 + 1) + 1
    board[pos[0]][pos[1]] += amount
    pile -= amount
    b++

    if pile < 0 {
      board[pos[0]][pos[1]] += pile
      pile = 0
    }
  }
  return board
}

func SelectRandomBlock(size int) (int, int) {
  num := rand.Intn(size * size)
  row := num / size
  col := num % size

  return row, col
}

func SerialSandPile(board GameBoard) GameBoard {
  for true {
    hasChange := (&board).SerialUpdateBoard()
    if hasChange == false {
      break
    }
  }

  return board
}

func (board *GameBoard) SerialUpdateBoard() bool {
  hasChange := false

  for r := range(*board) {
    for c := range((*board)[0]) {
      hasChangeCell := board.UpdateCell(r, c)
      if hasChangeCell == true {
        hasChange = true
      }
    }
  }
  return hasChange
}

func (board *GameBoard)UpdateCell(r, c int) bool {
  hasChange := false

  if (*board)[r][c] >= 4 {
     board.Topple(r, c)
     hasChange = true
  }
  return hasChange
}

func CopyBoard(board GameBoard) GameBoard {
  new := make(GameBoard, len(board))

  for i:=0; i < len(board); i++ {
    new[i] = make([]int, len(board[i]))
    copy(new[i], board[i])
  }

  return new
}

func (board *GameBoard)Topple(r, c int) {
  posArray := make([][2]int, 4)

  posArray[0] = [2]int{r-1, c}
  posArray[1] = [2]int{r+1, c}
  posArray[2] = [2]int{r, c-1}
  posArray[3] = [2]int{r, c+1}

  n := (*board)[r][c] / 4
  for _,  pos := range posArray {
    row := pos[0]
    col := pos[1]
    if board.WithinBoard(row, col) == true {
      (*board)[row][col] += n
    }
  }

  (*board)[r][c] = (*board)[r][c] % 4
  return
}

func (board *GameBoard)WithinBoard(r, c int) bool {
  rLen := len(*board)
  cLen := len((*board)[0])

  if r < 0 || r >= rLen  || c < 0 || c >= cLen {
    return false
  }

  return true
}

func (board *GameBoard)UpdateSingleCell(r, c int, position [][2]int) ([][2]int, bool){
  hasChange := false

  if (*board)[r][c] >= 4 {
     board.Topple(r, c)
     hasChange = true
     if board.WithinBoard(r-1, c) == true {
       position = append(position, [2]int{r-1, c})
     }
     if board.WithinBoard(r+1, c) == true {
       position = append(position, [2]int{r+1, c})
     }
     if board.WithinBoard(r, c-1) == true {
       position = append(position, [2]int{r, c-1})
     }
     if board.WithinBoard(r, c+1) == true {
       position = append(position, [2]int{r, c+1})
     }
  }
  return position, hasChange
}

func ParallelSandPile(board GameBoard, n int) GameBoard {
  hasChange := false

  partition := make([][2]int, 0)
  rLen := len(board) / n

  for i := 0; i < n; i++ {
    s := i * rLen + 1
    e := (i + 1) * rLen - 2
    if i == 0 {
      s = 0
    }
    if i == n-1 {
      e = len(board) - 1
    }
    partition = append(partition, [2]int{s, e})
  }



  for true {
    hasChange = false
    hasChangeBlock := make(chan bool, n)

    for i := 1; i < n; i++ {
      go (&board).ParallelUpdataBlock(i * rLen - 1, i * rLen, hasChangeBlock)
    }

    for i := 1; i < n; i++ {
      if <-hasChangeBlock == true {
        hasChange = true
      }
    }

    for _, part := range partition {
      go (&board).ParallelUpdataBlock(part[0], part[1], hasChangeBlock)
    }

    for i := 0; i < n; i++ {
      if <-hasChangeBlock == true {
        hasChange = true
        continue
      }
    }

    if hasChange == false {
      break
    }
  }

  return board
}

func (board *GameBoard) ParallelUpdataBlock(s, e int, hasChangeBlock chan bool) {
  hasChange := false

  for r := s; r <= e; r++ {
    for c := 0; c < len(*board); c++ {
      hasChangeCell := board.UpdateCell(r, c)
      if hasChangeCell == true {
        hasChange = true
      }
    }
  }
  hasChangeBlock <- hasChange
  return
}

func DrawGridLines(pic Canvas, cellWidth int) {
	w, h := pic.width, pic.height
	// first, draw vertical lines
	for i := 1; i < pic.width/cellWidth; i++ {
		y := i * cellWidth
		pic.MoveTo(0.0, float64(y))
		pic.LineTo(float64(w), float64(y))
	}
	// next, draw horizontal lines
	for j := 1; j < pic.height/cellWidth; j++ {
		x := j * cellWidth
		pic.MoveTo(float64(x), 0.0)
		pic.LineTo(float64(x), float64(h))
	}
	pic.Stroke()
}

func DrawGameBoard(board GameBoard, cellWidth int) Canvas {
	height := len(board) * cellWidth
	width := len(board) * cellWidth
	c := CreateNewCanvas(width, height)

	// declare colors
	color0 := MakeColor(0, 0, 0)
	color1 := MakeColor(85, 85, 85)
	color2 := MakeColor(170, 170, 170)
	color3 := MakeColor(255, 255, 255)
	red := MakeColor(255, 0, 0)
	// fill in colored squares
	for i := range board {
		for j := range board[i] {
			if board[i][j] == 0 {
				c.SetFillColor(color0)
			} else if board[i][j] == 1{
        c.SetFillColor(color1)
      } else if board[i][j] == 2{
        c.SetFillColor(color2)
      } else if board[i][j] == 3{
        c.SetFillColor(color3)
      } else {
				c.SetFillColor(red)
			}
			x := j * cellWidth
			y := i * cellWidth
			c.ClearRect(x, y, x+cellWidth, y+cellWidth)
			c.Fill()
		}
	}

	return c
}

func main() {
  size, error1 := strconv.Atoi(os.Args[1])
  if error1 != nil {
    fmt.Println("Error in reading size!")
    os.Exit(1)
  }

  pile, error2 := strconv.Atoi(os.Args[2])
  if error2 != nil {
    fmt.Println("Error in reading pile number!")
    os.Exit(1)
  }

  placement := os.Args[3]

  initialBoard := make([][]int, 0)
  cellWidth := 1

  n := 0
  if size <= 10 {
    n = 1
  } else if size <= 50 {
    n = 2
  } else {
    n = runtime.NumCPU()
  }

  if placement == "compete" {
    initialBoard = InitializeCentralBoard(size, pile)
    start := time.Now()
    boards := ParallelSandPile(initialBoard, n)
    elapse := time.Since(start)
    fmt.Println("Parallel sand pile takes ", elapse)
    finalBoard := DrawGameBoard(boards, cellWidth)
    finalBoard.SaveToPNG("parallel.png")
    return
  }

  if placement == "central" {
    initialBoard = InitializeCentralBoard(size, pile)
  } else if placement == "random" {
    initialBoard = InitializeRandomBoard(size, pile)
  }

  initialBoard2 := CopyBoard(initialBoard)

  start1 := time.Now()
  boards1 := SerialSandPile(initialBoard)
  elapse1 := time.Since(start1)
  fmt.Println("Serial sand pile takes ", elapse1)
  finalBoard1 := DrawGameBoard(boards1, cellWidth)
  finalBoard1.SaveToPNG("serial.png")



  start2 := time.Now()
  boards2 := ParallelSandPile(initialBoard2, n)
  elapse2 := time.Since(start2)
  fmt.Println("Parallel sand pile takes ", elapse2)
  finalBoard2 := DrawGameBoard(boards2, cellWidth)
  finalBoard2.SaveToPNG("parallel.png")

}
