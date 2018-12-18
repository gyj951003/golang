//Written by Haonan Sun & Yinjie Gao

package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var MismatchRate float64
var IndelRate float64
var ProIndelRate float64
var ProMismatchRate float64
var Operations = []byte{'N', 'R', 'E', 'N', 'N', 'N'}
var rewardPatterns = map[string]int{"E": 2, "EE": 4, "EEE": 8, "EEEE": 16}
var StartingGene = []byte{'N', 'R', 'E'}

type Color [3]uint8

type Cell interface {
	TakeAction(geneNum, r, c int, newBoard *GameBoard)
	CopyCell() Cell
	ReportType() string
	ReportColor() Color
	ReportEnergy() int
	ReportGene() []Gene
}

type Prokaryote struct {
	energy   int
	color    Color
	genes    Gene
	rpcGenes Gene
	rpcPos   int
	oprtPos  int //operation Position for gene
}

type Gene []byte

type Eukaryote struct {
	energy   int
	color    Color
	genes    [2]Gene
	rpcGenes [2]Gene
	rpcPos   [2]int
	oprtPos  [2]int //operation Position for gene
	//typeName string // Eukaryote

}

type GameBoard [][]Cell

//Return cell type of Eukaryote as a string
func (cell *Eukaryote) ReportType() string {
	return "Eukaryote"
}

//Return cell type of Prokaryote as a string
func (cell *Prokaryote) ReportType() string {
	return "Prokaryote"
}

//Return the gene of a Eukaryote in the form of a slice of two genes
func (cell *Eukaryote) ReportGene() []Gene {
	genes := make([]Gene, 2)
	genes[0] = cell.genes[0]
	genes[1] = cell.genes[1]
	return genes
}

//Return the gene of a Prokaryote in the form of a slice of one genes
func (cell *Prokaryote) ReportGene() []Gene {
	genes := make([]Gene, 1)
	genes[0] = cell.genes
	return genes
}

//Return the color of a Eukaryote in the form of color
func (cell *Eukaryote) ReportColor() Color {
	return cell.color
}

//Return the color of a Prokaryote in the form of color
func (cell *Prokaryote) ReportColor() Color {
	return cell.color
}

//Return the current energy of a Eukaryote
func (cell *Eukaryote) ReportEnergy() int {
	return cell.energy
}

//Return the current energy of a Prokaryote
func (cell *Prokaryote) ReportEnergy() int {
	return cell.energy
}

/*==============================================================================
InitialBoard
Preconditions:
rowNum and colNum are integers.
Postconditions:
A GameBoard of rowNum*colNum is returned.
 *============================================================================*/
func InitialBoard(rowNum, colNum int) GameBoard {
	board := make(GameBoard, rowNum)
	for r := range board {
		board[r] = make([]Cell, colNum)
	}
	return board
}

/*==============================================================================
CopyCell
Preconditions:
Cell is a Eukaryote.
Postconditions:
A newcell that is an exact copy of the called Eukaryote is returned.
 *============================================================================*/

func (cell *Eukaryote) CopyCell() Cell {
	var newCell Eukaryote
	newCell.genes[0] = make(Gene, len(cell.genes[0]))
	copy(newCell.genes[0], cell.genes[0])
	newCell.genes[1] = make(Gene, len(cell.genes[1]))
	copy(newCell.genes[1], cell.genes[1])
	newCell.rpcGenes[0] = make(Gene, len(cell.rpcGenes[0]))
	copy(newCell.rpcGenes[0], cell.rpcGenes[0])
	newCell.rpcGenes[1] = make(Gene, len(cell.rpcGenes[1]))
	copy(newCell.rpcGenes[1], cell.rpcGenes[1])
	newCell.rpcPos[0] = cell.rpcPos[0]
	newCell.rpcPos[1] = cell.rpcPos[1]
	newCell.oprtPos[0] = cell.oprtPos[0]
	newCell.oprtPos[1] = cell.oprtPos[1]
	newCell.energy = cell.energy
	newCell.color = cell.color
	//newCell.typeName = cell.typeName
	return &newCell
}

/*==============================================================================
CopyCell
Preconditions:
Cell is a Prokaryote.
Postconditions:
A newcell that is an exact copy of the called Prokaryote is returned.
 *============================================================================*/
func (cell *Prokaryote) CopyCell() Cell {
	var newCell Prokaryote
	newCell.genes = make(Gene, len(cell.genes))
	copy(newCell.genes, cell.genes)
	newCell.rpcGenes = make(Gene, len(cell.rpcGenes))
	copy(newCell.rpcGenes, cell.rpcGenes)
	newCell.rpcPos = cell.rpcPos

	newCell.oprtPos = cell.oprtPos

	newCell.energy = cell.energy
	newCell.color = cell.color
	//newCell.typeName = cell.typeName
	return &newCell
}

/*==============================================================================
TakeAction
Preconditions:
Cell is a Eukaryote.
Postconditions:
The cell takes a action that is at its operation position according to its genome.
 *============================================================================*/
func (cell *Eukaryote) TakeAction(geneNum, r, c int, newBoard *GameBoard) {
	if len(cell.genes[geneNum]) > 0 {
		pos := cell.oprtPos[geneNum] % len(cell.genes[geneNum])
		oprt := cell.genes[geneNum][pos]

		switch oprt {
		case 'N':
			cell.NonOperation(geneNum)
		case 'R':
			cell.Replication(geneNum, r, c, newBoard)
		case 'E':
			cell.EnergyGain(geneNum, rewardPatterns)
		default:
			panic("gene contain unknown operations!")
		}
	}
}

/*==============================================================================
TakeAction
Preconditions:
Cell is a Prokaryote.
Postconditions:
The cell takes a action that is at its operation position according to its genome.
 *============================================================================*/

func (cell *Prokaryote) TakeAction(geneNum, r, c int, newBoard *GameBoard) {
	if len(cell.genes) > 0 {
		pos := cell.oprtPos % len(cell.genes)
		oprt := cell.genes[pos]

		switch oprt {
		case 'N':
			cell.NonOperation()
		case 'R':
			cell.Replication(geneNum, r, c, newBoard)
		case 'E':
			cell.EnergyGain(rewardPatterns)
		default:
			panic("gene contain unknown operations!")
		}
	}
}

func UpdateCell(newBoard *GameBoard, r, c int) {
	if (*newBoard)[r][c] == nil { //No cell on the box, no action
		return
	} else if (*newBoard)[r][c].ReportType() == "Eukaryote" {
		currentCell := (*newBoard)[r][c]
		currentCell.TakeAction(0, r, c, newBoard)
		currentCell.TakeAction(1, r, c, newBoard)
	} else if (*newBoard)[r][c].ReportType() == "Prokaryote" {
		currentCell := (*newBoard)[r][c]
		currentCell.TakeAction(0, r, c, newBoard)
	}
}

func UpdateBoard(currentBoard GameBoard) GameBoard {
	newBoard := InitialBoard(len(currentBoard), len(currentBoard[0]))

	//copy through the board
	for r := range newBoard {
		for c := range newBoard[r] {
			if currentBoard[r][c] == nil {
				continue
			}
			newBoard[r][c] = currentBoard[r][c].CopyCell()
		}
	}

	//Update throught the board
	for r := range currentBoard {
		for c := range currentBoard[r] {
			UpdateCell(&newBoard, r, c)
		}
	}
	return newBoard
}

func (cell *Eukaryote) InitializeEnergy() {
	cell.energy = (len(cell.genes[0]) + len(cell.genes[1])) * 5
}

func (cell *Prokaryote) InitializeEnergy() {
	cell.energy = (len(cell.genes)) * 5
}

func SimulateDigitalOriganism(initialBoard GameBoard, numGens int) []GameBoard {
	boards := make([]GameBoard, numGens)
	boards[0] = initialBoard

	for i := 0; i < numGens-1; i++ {
		boards[i+1] = UpdateBoard(boards[i])
		if i == numGens / 4 {
			fmt.Println("25% finished")
		}
		if i == numGens / 2 {
			fmt.Println("50% finished")
		}
		if i == 3* numGens / 4 {
			fmt.Println("75% finished")
		}
	}

	return boards
}

func main() {

	fmt.Println("Parsing CLS)")
	if len(os.Args) != 13 {
		fmt.Println("Incorrect num of command line args!")
		os.Exit(1)
	}
	size, err1 := strconv.Atoi(os.Args[1])
	if err1 != nil {
		fmt.Println("Error converting size")
		os.Exit(1)
	}
	numGens, err2 := strconv.Atoi(os.Args[2])
	if err2 != nil {
		fmt.Println("Error converting numGens")
		os.Exit(1)
	}
	interval, err3 := strconv.Atoi(os.Args[3])
	if err3 != nil {
		fmt.Println("Error converting interval")
		os.Exit(1)
	}
	mutationRate, err4 := strconv.ParseFloat(os.Args[4], 64)
	if err4 != nil {
		fmt.Println("Error converting mutationRate")
		os.Exit(1)
	}
	promutationRate, err5 := strconv.ParseFloat(os.Args[5], 64)
	if err5 != nil {
		fmt.Println("Error converting mutationRate")
		os.Exit(1)
	}
	mismatch, err6 := strconv.ParseFloat(os.Args[6], 64)
	if err6 != nil {
		fmt.Println("Error converting MismatchRate")
		os.Exit(1)
	}
	proMismatch, err6 := strconv.ParseFloat(os.Args[7], 64)
	if err6 != nil {
		fmt.Println("Error converting MismatchRate")
		os.Exit(1)
	}
	genelength, err7 := strconv.Atoi(os.Args[8])
	if err7 != nil {
		fmt.Println("Error converting genelength")
		os.Exit(1)
	}

	IndelRate = mutationRate
	ProIndelRate = promutationRate
	MismatchRate = mismatch
	ProMismatchRate = proMismatch
	initialGene := make(Gene, 0)
	for i := 0; i < genelength; i++ {
		initialGene = append(initialGene, StartingGene...)
	}

	rand.Seed(time.Now().UTC().UnixNano())
	rowNum := size
	colNum := size

	//Initialize board
	initialBoard := InitialBoard(rowNum, colNum)
	//Set ancesters
	//initialGene := Gene{'R', 'N', 'E', 'R', 'N', 'E'}
	initialGenes := [2]Gene{initialGene, initialGene}
	rpcGene := make(Gene, 0)
	rpcGenes := [2]Gene{rpcGene, rpcGene}
	initialColor := Color{0, 255, 255}
	ancester := Eukaryote{color: initialColor, genes: initialGenes, rpcGenes: rpcGenes, rpcPos: [2]int{0, 0}, oprtPos: [2]int{0, 0}}
	ancester.InitializeEnergy()

	ancester1 := ancester.CopyCell()
	ancester2 := ancester.CopyCell()
	ancester3 := ancester.CopyCell()
	//ancester4 := ancester.CopyCell()
	//ancester5 := ancester.CopyCell()
	initialBoard[size-1][size-1] = &ancester
	initialBoard[size-2][size-1] = ancester1
	initialBoard[size-1][size-2] = ancester2
	initialBoard[size-2][size-2] = ancester3
	//initialBoard[151][153] = ancester4
	//initialBoard[153][151] = ancester5

	initialColor2 := Color{255, 255, 0}
	ancesterp := Prokaryote{color: initialColor2, genes: initialGene, rpcGenes: rpcGene, rpcPos: 0, oprtPos: 0}
	ancesterp.InitializeEnergy()
	ancesterpp := ancesterp
	initialBoard[0][0] = &ancesterp
	initialBoard[1][1] = &ancesterpp
	initialBoard[1][0] = &ancesterp
	initialBoard[0][1] = &ancesterp

	fmt.Println("Now, calculate boards..")
	boards := SimulateDigitalOriganism(initialBoard, numGens)

	fmt.Println("Now, draw gif...")

	//Now, visualize the board

	cellWidth := 2
	imageList := DrawGameBoards(boards, cellWidth, interval)
	Process(imageList, "pic")

	fmt.Println("Now, output statistical analysis...")
	interval2 := 50
	GetStat(boards, interval2)
}
