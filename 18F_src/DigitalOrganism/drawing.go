// Drawing.go contains the functions to visualize Gameboard and to generate heatmap
// Written by Yinjie Gao & Haonan Sun.

package main

import (
	"image"
)

// Copied from codes in lecture. Credit to Phillip Compeau
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

// DrawGameBoard will generate a size * size cells board with cells filled by the color of individual digital organisms
func DrawGameBoard(board GameBoard, cellWidth int) image.Image {
	height := len(board) * cellWidth
	width := len(board[0]) * cellWidth
	c := CreateNewCanvas(width, height)
  //fmt.Println("Success Creat Canvas")
	// declare colors
	darkGray := MakeColor(50, 50, 50)
	// black := MakeColor(0, 0, 0)

	// fill in colored squares
	for i := range board {
		for j := range board[i] {
			if board[i][j] == nil {
				c.SetFillColor(darkGray)
			} else {
        color := board[i][j].ReportColor()
        //fmt.Println(color)
        cellColor := MakeColor(color[0], color[1], color[2])
        c.SetFillColor(cellColor)
      }
			x := j * cellWidth
			y := i * cellWidth
			c.ClearRect(x, y, x+cellWidth, y+cellWidth)
			c.Fill()
		}
	}
  img := c.img
	return img

}

// DrawLengthBoard() generate a heatmap for genome length. Organisms with longer genome will have darker color
func DrawLengthBoard(board GameBoard, min, max int, cellWidth int) Canvas {
	height := len(board) * cellWidth
	width := len(board[0]) * cellWidth
	c := CreateNewCanvas(width, height)
  black := MakeColor(0, 0, 0)
	// declare colors: Red: B = 225, RG: range through(250, 0)
	for i := range board {
		for j := range board[i] {
			if board[i][j] == nil {
				c.SetFillColor(black)
			} else {
				length := 0
				if board[i][j].ReportType()=="Eukaryote" {
					length = len(board[i][j].(*Eukaryote).genes[0]) + len(board[i][j].(*Eukaryote).genes[0])
				} else if board[i][j].ReportType()=="Prokaryote" {
					length = len(board[i][j].(*Prokaryote).genes)
				}

				rg := 0
				if length != max {
					rg = 250 * (max - length) / (max - min)
				}
        cellColor := MakeColor(uint8(rg), uint8(rg), 255)
        c.SetFillColor(cellColor)
      }
			x := j * cellWidth
			y := i * cellWidth
			c.ClearRect(x, y, x+cellWidth, y+cellWidth)
			c.Fill()
		}
	}
	return c
}

// DrawEnergyBoard() generate a heatmap for energy. Organisms with more energy will have darker color
func DrawEnergyBoard(board GameBoard, min, max int, cellWidth int) Canvas {
	height := len(board) * cellWidth
	width := len(board[0]) * cellWidth
	c := CreateNewCanvas(width, height)
  black := MakeColor(0, 0, 0)
	// declare colors: Red: R = 225, BG range through(250, 0)
	for i := range board {
		for j := range board[i] {
			if board[i][j] == nil {
				c.SetFillColor(black)
			} else {
				energy := board[i][j].ReportEnergy()

				bg := 0
				if energy != max {
					bg = 250 * (max - energy) / (max - min)
				}
        cellColor := MakeColor(255, uint8(bg), uint8(bg))
        c.SetFillColor(cellColor)
      }
			x := j * cellWidth
			y := i * cellWidth
			c.ClearRect(x, y, x+cellWidth, y+cellWidth)
			c.Fill()
		}
	}
	return c
}

// DrawGameBoards() combines images generated form DrawGameBoard() into imageList
func DrawGameBoards(boards []GameBoard, cellWidth, interval int) []image.Image {
	//numGenerations := len(boards)
	imageList := make([]image.Image, 0)
	for i := 0; i <  len(boards); i++ {
    if i % interval == 0 {
      imageList = append(imageList, DrawGameBoard(boards[i], cellWidth))
    }
	}
	return imageList
}

// Written by Yinjie Gao and Haonan Sun
// Dec, 03, 2018
