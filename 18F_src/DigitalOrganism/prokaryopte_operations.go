package main

import (
	//"fmt"

	"math/rand"
)

/*==============================================================================
NonOperation
Preconditions:
cell is a Prokaryote.
Postconditions:
Cell's energy -1. Cell's operation Position +1.
 *============================================================================*/
func (cell *Prokaryote) NonOperation() {
	cell.oprtPos += 1
	cell.energy -= 1
}

/*==============================================================================
EnergyGain
Preconditions:
cell is a Prokaryote. len(rewardPatterns)>0.
Postconditions:
rewardPatterns is a map of string with int value asscoicated with them.
The energy of the cell increase according to the highest value pattern matched by the rewardPatterns.
 *============================================================================*/
func (cell *Prokaryote) EnergyGain(rewardPatterns map[string]int) {

	gene := cell.genes
	pos := cell.oprtPos % len(cell.genes)

	//Check if there is a match in the reward pattern, find the largest reward for the reward Patterns
	largestReward := 0
	for pattern, reward := range rewardPatterns {
		n := len(pattern)
		if pos+n > len(gene) {
			continue
		}
		if string(gene[pos:pos+n]) == pattern && reward > largestReward {
			largestReward = reward
		}
	}

	cell.energy += largestReward
	cell.oprtPos += 1
	cell.energy -= 1
}

/*==============================================================================
Replication
Preconditions:
cell is a Prokaryote. r,c represent the location of the cell.  geneNum is a useless
parameter here because Eukaryote's Replication use it.
Postconditions:
The rpcgenes add one nucleotides of genes if without mutation. Otherwise, it might
gain two nucleotides or zero or more, base on the mutation. After each addiiton,
rpcPos would +1.oprtPos also +1. cell.energy would -1. If the cell is ready
to Reproduce, then the cell would reproduce.
 *============================================================================*/

func (cell *Prokaryote) Replication(geneNum, r, c int, newBoard *GameBoard) {
	if cell.rpcPos < len(cell.genes) {
		nucleotide := cell.genes[cell.rpcPos]
		nucleotides := make(Gene, 0)
		(&nucleotides).ProMutate(nucleotide) //return nucleotides to be added
		cell.rpcGenes = append(cell.rpcGenes, nucleotides...)
		cell.rpcPos += 1
		cell.oprtPos += 1
		cell.energy -= 1
	} else if cell.ReadyToReproduce() {
		//fmt.Println(cell.genes[1], cell.rpcGenes[1])
		cell.Reproduce(r, c, newBoard)
	}
}

// Mutate wil decide what the next copied nucleotide to be
// if no mutation, it will copy the template gene
// Otherwise, it has a chance to have a deletion, insertion, and mismatch.
func (nucleotides *Gene) ProMutate(nucleotide byte) Gene {
	//fmt.Println(ProIndelRate, ProMismatchRate)
	if Deletion(ProIndelRate) == false && nucleotides.Insertion(nucleotide, ProIndelRate) == false && nucleotides.Mismatch(nucleotide, ProMismatchRate) == false {
		*nucleotides = append(*nucleotides, nucleotide)

	}
	return *nucleotides
}

/*==============================================================================
ReadyToReproduce
Preconditions:
cell is a Prokaryote.
Postconditions:
Return true if the rpcpos equal to len(cell.genes). Return false otherwise.
 *============================================================================*/
func (cell *Prokaryote) ReadyToReproduce() bool {
	if cell.rpcPos == len(cell.genes) && len(cell.genes) != 0 {
		return true
	}
	return false
}

/*==============================================================================
Reproduce
Preconditions:
cell is a Prokaryote. The cell is ready to reproduce.
Postconditions:
Generate a new Prokaryote base on the cell and its position, r and c. Reset rpcGenes and rpcPos.
 *============================================================================*/

func (cell *Prokaryote) Reproduce(r, c int, newBoard *GameBoard) {
	newBoard.GenerateNewProkaryote(cell, r, c)
	rpcGeneNew := make(Gene, 0)
	cell.rpcGenes = rpcGeneNew
	cell.rpcPos = 0
}

/*==============================================================================
GenerateNewProkaryote
Preconditions:
cell is a Prokaryote. The cell is ready to reproduce. r,c
Postconditions:
Generate a new Prokaryote base on the cell. Reset rpcGenes and rpcPos.
 *============================================================================*/

func (newBoard *GameBoard) GenerateNewProkaryote(parent *Prokaryote, r, c int) {
	nbrMap := FindEmptyNbr(r, c, newBoard)
	//Check if there are empty space in the neighborhood

	p1Gene := parent.rpcGenes
	gene1 := make(Gene, len(p1Gene))
	copy(gene1, p1Gene)

	//fmt.Println(p1Gene, p2Gene)

	rpcGene := make(Gene, 0)

	color := CreateColorPro(parent, gene1)
	newBorn := Prokaryote{genes: gene1, rpcGenes: rpcGene, rpcPos: 0, oprtPos: 0, energy: len(gene1), color: color}
	newBorn.InitializeEnergy()
	i := 0

	if len(nbrMap) == 0 {
		parent.BattlePro(&newBorn, r, c, newBoard)
		return
	}
	if len(nbrMap) > 1 {
		i = rand.Intn(len(nbrMap))
	}

	pos := nbrMap[i]
	//fmt.Println(newBorn.genes[0], newBorn.genes[1])
	(*newBoard)[pos[0]][pos[1]] = &newBorn
}

/*==============================================================================
BattlePro
Preconditions:
cell is a Prokaryote. The cell is ready to reproduce. r,c represent the current
location. A pointer of the newboard and a pointer of the newborn is available.
Postconditions:
First, all the neighbor that have a different genome than the rpcgene are mapped.
If none exist, it returns. Else, the program picks one opponent randomly from
the map and fight with it. The result of the fight is based on the energy level
the organism currently have.If the fight win, the selected cell is replaced with
the new born prokaryote. Else, the
 *============================================================================*/

func (cell *Prokaryote) BattlePro(newborn *(Prokaryote), r, c int, newBoard *GameBoard) {
	nbrMap := FindDiffNbr(cell, r, c, newBoard)
	if len(nbrMap) == 0 {
		return
	}
	i := 0
	if len(nbrMap) > 1 {
		i = rand.Intn(len(nbrMap))
	}

	pos := nbrMap[i]
	sbj := (*newBoard)[pos[0]][pos[1]]

	energy := cell.energy
	ratio := float64(energy) / float64(sbj.ReportEnergy())

	win := false

	if energy >= 0 && sbj.ReportEnergy() <= 0 {
		win = true
	}

	if energy > 0 && sbj.ReportEnergy() > 0 && IsWin(ratio, 0.1) == 1 {
		win = true
	}
	if energy < 0 && sbj.ReportEnergy() < 0 && IsWin(ratio, 0.1) == 0 {
		win = true
	}

	if win == true {
		(*newBoard)[pos[0]][pos[1]] = cell
	}

	if win == false && sbj.ReportType() == "Eukaryote" && energy > 0 {
		sbj.(*Eukaryote).energy -= energy
	}

	if win == false && sbj.ReportType() == "Prokaryote" && energy > 0 {
		sbj.(*Prokaryote).energy -= energy
	}
}

/*==============================================================================
FindDiffNbr
Preconditions:
cell is a Prokaryote. r,c represent the current location.
A pointer of the newboard is available.
Postconditions:
First, all the neighbor that have a different genome than the rpcgene are mapped,
and returned as a map of two integers, represeting the row and column of the cell.
 *============================================================================*/
func FindDiffNbr(cell *Prokaryote, r, c int, newBoard *GameBoard) map[int][2]int {
	nbrMap := make(map[int][2]int, 0)
	i := 0
	//Range over all neighbor
	for row := r - 1; row <= r+1; row++ {
		for col := c - 1; col <= c+1; col++ {
			if row == r && col == c {
				continue
			} else if WithinBoard(row, col, newBoard) == false {
				continue
			} else if (*newBoard)[row][col] == nil {
				continue
			} else if (*newBoard)[row][col].ReportType() == "Eukaryote" {
				nbrMap[i] = [2]int{row, col}
				i++
			} else {
				nbrCell := (*newBoard)[row][col].(*Prokaryote)
				if len(nbrCell.genes) != len(cell.rpcGenes) {
					nbrMap[i] = [2]int{row, col}
					i++
					break
				}
				check := false
				for j := range nbrCell.genes {
					if nbrCell.genes[j] != cell.rpcGenes[j] && check == false {
						nbrMap[i] = [2]int{row, col}
						i++
						check = true
						break
					}
				}

			}
		}
	}
	return nbrMap
}

/*==============================================================================
CreateColorPro
Preconditions:
parent is a Prokaryote, cGene represent the copied gene of current parent.
Postconditions:
First, the copied gene is compared with its parent gene. Return parent's color if
they are the same. Otherwise, create a new color that is randomly added in Red and
Blue from its parent, and return it instead.
 *============================================================================*/

func CreateColorPro(parent *Prokaryote, cGene Gene) Color {
	if CompareGenes(parent.genes, cGene) == 0 {
		return parent.color
	}
	R := int(parent.color[0]) + rand.Intn(10)
	G := 0
	B := int(parent.color[1]) + rand.Intn(10)
	//G := int(parent.color[2]) + 10*(rand.Intn(2)-1)
	Red := uint8(R % 205 + 50)
	Blue := uint8(0)
	Green := uint8(G % 205 + 50)
	return [3]uint8{Red, Blue, Green}
}
