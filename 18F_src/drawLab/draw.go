package main
import "gol"

func DrawRectangle() {
    // Draw an 800 pixel by 100 pixel rectangle
  red := MakeColor(255, 0, 0)
  pic := CreateNewCanvas(1000, 300)
  pic.SetStrokeColor(red)
  pic.MoveTo(100, 100)
  pic.LineTo(100, 200)
  pic.LineTo(900, 200)
  pic.LineTo(900, 100)
  pic.LineTo(100, 100)
  pic.Stroke()
  pic.SaveToPNG("myRectangle.png")
}

func DrawGoL(boards gol.GameBoard, cellWidth int) Canvas {
  w := gol.CountCols(boards) * cellWidth
  h := gol.CountRows(boards) * cellWidth

  myCanvas := CreateNewCanvas(w, h)
  black := MakeColor(0, 0, 0)
  // black: alive cells
  myCanvas.SetFillColor(black)

  // draw grid lines
  myCanvas = DrawGridLines(myCanvas, cellWidth)

  // Range over board and fill the alive
  for row := range boards {
    for col := range boards {
      if boards[row][col] ==  1 {
        DrawSquare(myCanvas, row, col, cellWidth)
      }
    }
  }

  return myCanvas
}

func DrawGridLines(myCanvas Canvas, cellWidth int) Canvas{
  w, h := myCanvas.width, myCanvas.height
  // Draw vertical lines
  for i := 1; i < w/cellWidth; i++ {
    x := i * cellWidth
    myCanvas.MoveTo(float64(x), 0.0)
    myCanvas.LineTo(float64(x), float64(h))
  }

  // Draw horizontal lines
  for j := 1; j < h/cellWidth; j++ {
    y := j * cellWidth
    myCanvas.MoveTo(0.0, float64(y))
    myCanvas.LineTo(float64(w), float64(y))
  }

  myCanvas.Stroke()

  return myCanvas
}

func DrawSquare(myCanvas Canvas, row, col, cellWidth int) Canvas {
  x1, y1 := row * cellWidth, col * cellWidth
  x2, y2 := x1 + cellWidth, y1 + cellWidth
  myCanvas.ClearRect(x1, y1, x2, y2)
  myCanvas.Fill()
  return myCanvas
}


func main() {
  DrawRectangle()

  numRows := 30
  rPentomino := gol.InitializeBoard(numRows, numRows)
  mid := numRows/2
  rPentomino[mid][mid] = 1
  rPentomino[mid-1][mid] = 1
  rPentomino[mid-1][mid+1] = 1
  rPentomino[mid][mid-1] = 1
  rPentomino[mid+1][mid] = 1

  cellWidth := 10
  myCanvas := DrawGoL(rPentomino, cellWidth)
  myCanvas.SaveToPNG("myRPentomino.png")

}
