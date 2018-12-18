package main

import (
	"os"
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"image"
)

// The data stored in a single cell of a field
type Cell struct {
	strategy  string //represents "C" or "D" corresponding to the type of prisoner in the cell
	score float64 //represents the score of the cell based on the prisoner's relationship with neighboring cells
}

// The game board is a 2D slice of Cell objects
type GameBoard [][]Cell

// Initialize board
func InitialBoard(numRows, numCols int) GameBoard {
	board := make(GameBoard, numRows)

	for r := range board {
    board[r] = make([]Cell, numCols)
  }

	return board
}

// Read initial state from the file and returns a initial GameBoard
func ReadBoardFromFile(fileName string) GameBoard {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error: reading file.")
		os.Exit(1)
	}
	scanner := bufio.NewScanner(file)
	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	line1 := strings.Split(lines[0], " ")
	numRows, err1 := strconv.Atoi(line1[0])
	if err1 != nil {
		fmt.Println("Error: Converting numRows.")
		os.Exit(1)
	}
	numCols, err2 := strconv.Atoi(line1[1])
	if err2 != nil {
		fmt.Println("Error: Converting numCols.")
		os.Exit(1)
	}

	initialBoard := InitialBoard(numRows, numCols)

	for i := 0; i < numRows; i++ {
		for j:= 0; j < numCols; j++ {
			s := string(lines[i + 1][j])
			initialBoard[i][j] = Cell{strategy: s, score: 0}
		}
	}

	file.Close()
	return initialBoard
}

func WithinBoard(currentBoard GameBoard, r, c int) bool {
	numRows := len(currentBoard)
	numCols := len(currentBoard[0])
	if r < 0 || r >= numRows {
		return false
	}
	if c < 0 || c >= numCols {
		return false
	}
	return true
}

func Compete(self Cell, nbr Cell, b float64) float64 {
	if self.strategy == "C" {
		if nbr.strategy == "C" {
			return 1.0
		} else {
			return 0.0
		}
	} else if self.strategy == "D" {
		if nbr.strategy == "C" {
			return b
		} else {
			return 0.0
		}
	} else {
		panic("Stategy input is wrong!\n")
	}
}

func UpdateScore(currentBoard GameBoard, r int, c int, b float64) float64 {
	score := 0.0
	for i := r-1; i < r + 2; i++ {
		for j := c-1; j < c+2; j++ {
			if (i != r || j != c) && WithinBoard(currentBoard, i, j) {
				score = score + Compete(currentBoard[r][c], currentBoard[i][j], b)
			}
		}
	}
	return score
}

func UpdateCell(currentBoard GameBoard, r int, c int) string {
	bestStrategy := currentBoard[r][c].strategy
	bestScore := currentBoard[r][c].score
	for i := r-1; i < r + 2; i++ {
		for j := c-1; j < c+2; j++ {
			if (i != r || j != c) && WithinBoard(currentBoard, i, j) {
				if currentBoard[i][j].score > bestScore {
					bestStrategy = currentBoard[i][j].strategy
					bestScore = currentBoard[i][j].score
				}
			}
		}
	}
	return bestStrategy
}

func OneRoundPrisonerDilemma(currentBoard GameBoard, b float64) GameBoard {
	numRows := len(currentBoard)
	numCols := len(currentBoard[0])
	newBoard := InitialBoard(numRows, numCols)

	for r := 0; r < numRows; r++ {
		for c := 0; c < numCols; c++ {
			currentBoard[r][c].score = UpdateScore(currentBoard, r, c, b)
		}
	}

	for r := 0; r < numRows; r++  {
		for c := 0; c < numCols; c++ {
			newBoard[r][c].strategy = UpdateCell(currentBoard, r, c)
		}
	}

	return newBoard
}

// Demostrate the Prisoners Dilemma and returns a slice of GameBoard
func DemoPrisonerDilemma(initialBoard GameBoard, numGens int, b float64) []GameBoard {
	boards := make([]GameBoard, numGens + 1)
	boards[0] = initialBoard
	for i := 0; i < numGens; i++ {
		boards[i + 1] = OneRoundPrisonerDilemma(boards[i], b)
	}
	return boards
}

func DrawGameBoards(boards []GameBoard, cellWidth int) []image.Image {
	numGenerations := len(boards)
	imageList := make([]image.Image, numGenerations)
	for i := range boards {
		c := DrawGameBoard(boards[i], cellWidth)
		imageList[i] = c.img
	}
	return imageList
}

//DrawGameBoard takes a single game board along with a cellWidth and produces
//an image.Image object corresponding to the board with each cell having cellWidth x cellWidth pixels
func DrawGameBoard(board GameBoard, cellWidth int) Canvas {
	height := len(board) * cellWidth
	width := len(board[0]) * cellWidth
	c := CreateNewCanvas(width, height)

	// declare colors
	blue := MakeColor(0, 0, 255)
	red := MakeColor(255, 0, 0)

	// fill in colored squares
	for i := range board {
		for j := range board[i] {
			if board[i][j].strategy == "C" {
				c.SetFillColor(blue)
			} else if board[i][j].strategy == "D" {
				c.SetFillColor(red)
			} else {
				panic("Error: Out of range value.\n")
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
	initialBoard := ReadBoardFromFile(os.Args[1])
	b, err1 := strconv.ParseFloat(os.Args[2], 64)
	if err1 != nil {
		fmt.Println("Error: reading b.")
		os.Exit(1)
	}
	numGens, err2 := strconv.Atoi(os.Args[3])
	if err2 != nil {
		fmt.Println("Error: reading numGens.")
		os.Exit(1)
	}
	boards := DemoPrisonerDilemma(initialBoard, numGens, b)

	cellWidth := 10
  //imageList := DrawGameBoards(boards, cellWidth)
	finalCanvas := DrawGameBoard(boards[numGens], cellWidth)
	finalCanvas.SaveToPNG("Prisoners.png")
  //Process(imageList, "Prinsoners")

	os.Exit(0)
}
