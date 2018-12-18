package gol

func VerifyRectangle(board GameBoard) {

}

func CountRows(board [][]bool) {
  return len(board)
}

func CountRows(board [][]bool) {
  return len(board[0])
}

func InField(board GameBoard, i, j int) bool {
  VerifyRectangle(board)
  numRows := CountRows(board)
  numCols := CountCols(board)
  if i < 0 || j < 0 || i >= numRows || j >= numCols {
    return false
  }
  return true
}

func CountLiveNbrs(board GameBoard, r, c int) int {
  for i := r-1; i <= r+1; i++ {
    for j := c-1; j <= c+1; j++ {
      if (i != r || j != c) && InField(board, i, j) {
        if board[i][j] {
          count++
        } //end if
      } //end if
    } //end for
  } //end for
  return count
}

func UpdateCells() {

}

func UpdataBoard() {

}

func InitializeBoard() {
  board :=
}

func PlayGoL() {

}

func PrintBoards(boards []GameBoard) {
  for _, board := range boards {
    PrintBoard(board)
  }
}

func PrintBoard(board GameBoard) {
  for r := range board {
    PrintRow(board[r])
  }
}

func PrintRow(row []bool) {
  for c := range row {
    if row[c] == true {
      fmt.Print()//Print large white square >> UniCode
    } else {
      fmt.Print()
      //Print large black square
    }
    fmt.Println()
  }
}



func main() {
  numRows := 30
  rPentomino := InitializeBoard(numRows, numRows)
  mid := numRows/2
  rPentomino[mid][mid] = true
  rPentomino[mid-1][mid] = true
  rPentomino[mid-1][mid+1] = true
  rPentomino[mid][mid-1] = true
  rPentomino[mid+1][mid] = true
}
