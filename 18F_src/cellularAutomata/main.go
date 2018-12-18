package main

import ("strconv"
  "fmt"
  "os"
)

type GameBoard [][]int


func CountRows(board GameBoard) int {
  return len(board)
}

func CountCols(board GameBoard) int {
  return len(board[0])
}


func CellMatchNeumann(board GameBoard, r int, c int, rule string) bool {
  ruleI := make([]int, 0)
  for i := 0; i < len(rule) - 1; i++ {
    num,_ := strconv.Atoi(rule[i:i+1])
    ruleI = append(ruleI, num)
  }
  //board[r][c] // self; board[r - 1][c] // top; board[r][c + 1] // right
  //board[r + 1][c] // bottom;  board[r][c - 1] // left
  if (board[r][c] != ruleI[0] || board[r - 1][c] != ruleI[1] || board[r][c + 1] != ruleI[2] || board[r + 1][c] != ruleI[3] || board[r][c - 1] != ruleI[4]) {
    return false
  }
  return true
}

func CellMatchMoore(board GameBoard, r int, c int, rule string) bool {
  ruleI := make([]int, 0)
  for i := 0; i < len(rule) - 1; i++ {
    num,_ := strconv.Atoi(rule[i:i+1])
    ruleI = append(ruleI, num)
  }

  if board[r][c] != ruleI[0] || board[r - 1][c - 1] != ruleI[1] || board[r - 1][c] != ruleI[2] || board[r - 1][c + 1] != ruleI[3] || board[r][c - 1] != ruleI[4] || board[r][c + 1] != ruleI[5] || board[r + 1][c - 1] != ruleI[6] || board[r + 1][c] != ruleI[7] || board[r + 1][c + 1]!= ruleI[8] {
    return false
  }
  return true
}

func CellMatch(board GameBoard, r int, c int, nbrHood string, rule string) bool {
  if nbrHood == "vonNeumann" {
    return CellMatchNeumann(board, r, c, rule)
  } else if nbrHood == "Moore" {
    return CellMatchMoore(board, r, c, rule)
  }
  panic("Incorrect neiborhood name given to CellMatch()")
}

func UpdateCells(board GameBoard, r int, c int, nbrHood string, ruleStrings []string) int {
  // range over the rule strings and look for a match
    for i := range ruleStrings {
      rule := ruleStrings[i]
      if CellMatch(board, r, c, nbrHood, rule) {
        lastSymbol := rule[len(rule) - 1:]
        result, err := strconv.Atoi(lastSymbol)
        if err != nil {
          panic("Error coverting the last symbol to an integer")
        }
        return result
      }
    } // end for i
    panic("We did not fine a rule string matching a given cell in Updatecell at row " + strconv.Itoa(r) + "and col " + strconv.Itoa(c))
} // UpdateCells()

func UpdataBoard(board GameBoard, nbrHood string, ruleStrings []string) GameBoard {
  numRows := CountRows(board)
  numCols := CountCols(board)
  newBoard := InitializeBoard(numRows, numCols)

  // range through the board
  // avoiding the boundary
  for r := 1; r < numRows - 1; r++ {
    for c := 1; c < numCols - 1; c++ {
      newBoard[r][c] = UpdateCells(board, r, c, nbrHood, ruleStrings)
    } // end for r
  } // end for c

  return newBoard
} // UpdataBoard()

func InitializeBoard(numRows, numCols int) GameBoard {
  board := make(GameBoard, numRows)
  // range over rows and make the int slices
  for r := range board {
    board[r] = make([]int, numCols)
  }
  return board
}

// PlayAutomaton takes an initial game board as well as a number of generations, the neiborhood type and a collection of rule ruleStrings
// And returns numGens + 1 of boards
func PlayAutomaton(initialBoard GameBoard, numGens int,
  nbrHood string, ruleStrings []string) []GameBoard {
  //Play GOL for numGens rounds
  boards := make([]GameBoard, numGens + 1)
  boards[0] = initialBoard
  for i := 0; i < numGens; i++ {
    boards[i] = UpdataBoard(boards[i-1], nbrHood, ruleStrings)
  }
  return boards
}


func main() {
  if len(os.Args) != 6 {
    fmt.Println("Incorrect number of command line arguments!")
    os.Exit(1)
  }

  initialBoardFile := os.Args[1]
  numGens, err := strconv.Atoi(os.Args[2])
  if err != nil {
    fmt.Println("Error converting numGens!")
    os.Exit(1)
  }
  nbrHood := os.Args[3]
  ruleStringsFileName := os.Args[4]
  outFileName := os.Args[5]

  initialBoard := ReadBoardFromFile(initialBoardFile)

  ruleStrings := ReadStringsFromFile(ruleStringsFileName)

  boards := PlayAutomaton(initialBoard, numGens, nbrHood, ruleStrings)

  fmt.Println("Boards generated, yay! Generating images ...")

  cellWidth := 10

  imageList := DrawGameBoards(boards, cellWidth)
  Process(imageList, outFileName)

  fmt.Println("Program finished successfully!")
  os.Exit(0)
}
