// Written by Yinjie Gao
// eukaryote_operations.go contains all actions that a cell type Eukaryote can take
// Major functions include NonOperation(), EnergyGain(). Replication() along with Reproduce()

package main

import (
	"math/rand"
)

// NonOperation() is executed when the cell's currently expressing gene element is 'N'
// NonOperation() takes 1 unit of energy with no beneficial reward.
func (cell *Eukaryote) NonOperation(geneNum int) {
	cell.oprtPos[geneNum] += 1
	cell.energy -= 1
	return
}

// EnergyGain() is executed when the cell's currently expressing gene element is 'E'
// wEnergyGain() ill find a pattern match in defined rewardPatterns map, and add certain amount of energy into the Eukayote's energy storage
func (cell *Eukaryote) EnergyGain(geneNum int, rewardPatterns map[string]int) {
	gene := cell.genes[geneNum]
	pos := cell.oprtPos[geneNum] % len(cell.genes[geneNum]) // current operation position

	//Check if there is a match in the reward pattern, find the largest reward for the reward Patterns
	largestReward := 0

	for pattern, reward := range rewardPatterns {
		n := len(pattern)

		if pos+n > len(gene) { // out of gene range
			continue
		}

		if string(gene[pos:pos+n]) == pattern && reward > largestReward { // pattern match
			largestReward = reward
		}
	}

	cell.energy += largestReward
	cell.oprtPos[geneNum] += 1
	cell.energy -= 1

	return
}

// Replication() is executed when the cell's currently expressing gene element is 'R'
// It will copy a genome element from the cell's gene and store it into the rpcGene
// When a copy of genome complete replication, the Eukayote will try to find a neighbor to reproduce the next generation
func (cell *Eukaryote) Replication(geneNum, r, c int, newBoard *GameBoard) {

	if cell.rpcPos[geneNum] < len(cell.genes[geneNum]) { //if replication has not finished
		nucleotide := cell.genes[geneNum][cell.rpcPos[geneNum]]
		nucleotides := make(Gene, 0)
		(&nucleotides).Mutate(nucleotide) //return nucleotides to be added
		cell.rpcGenes[geneNum] = append(cell.rpcGenes[geneNum], nucleotides...)
		cell.rpcPos[geneNum] += 1
		cell.oprtPos[geneNum] += 1
		cell.energy -= 1
	}

	//Check if a cell is ready to reproduce
	if cell.ReadyToReproduce() {
		cell.Reproduce(r, c, newBoard)
	}

	return
}

// ReadyToReproduce() return true when the replication has completed for both copies of genes
func (cell *Eukaryote) ReadyToReproduce() bool {
	if cell.rpcPos[0] == len(cell.genes[0]) && cell.rpcPos[1] == len(cell.genes[1]) && len(cell.genes[0]) != 0 && len(cell.genes[1]) != 0{
		return true
	}
	return false
}

// Reproduce() will find a Eukaryote neighbor that is also ready to reproduce and Generate a newborn
// The newBorn will inherite the genome from both side of parents
func (cell *Eukaryote) Reproduce(r, c int, newBoard *GameBoard) {
	nbrEukaryoteMap := FindEukayoteNbr(r, c, newBoard)

	// Range through the cell's neighbor to find another Eukaryote that is ready to reproduce
	for _, pos := range nbrEukaryoteMap {
		nbrCell := (*newBoard)[pos[0]][pos[1]].(*Eukaryote)

		if nbrCell.ReadyToReproduce() == true { // if find neighbor that is ready to reproduce
			newBoard.GenerateNewEukaryote(cell, nbrCell, r, c)

			// After reproduce, reset the parents' rpcGenes to be blank
			// and reset the rpcPos to be 0
			rpcGeneNew := make(Gene, 0)
			rpcGenesNew := [2]Gene{rpcGeneNew, rpcGeneNew}
			cell.rpcGenes = rpcGenesNew
			cell.rpcPos = [2]int{0, 0}
			nbrCell.rpcGenes = rpcGenesNew
			nbrCell.rpcPos = [2]int{0, 0}

			break //stop finding the next neibor
		}
	}

	return
}

// GenerateNewEukaryote() will takes the 2 parents for the newborn and generate a newborn
// Before reproduce, the parents' rpcGenes will undergo crossover and pass one of the rpcGenes to the newborn
// The newborn will be put into a empty space, if there is no empty space nearby, the parent will fight with the organisms nearby and if the parent win, the newborn can be put into that space.
func (newBoard *GameBoard) GenerateNewEukaryote(parent1, parent2 *Eukaryote, r, c int) {
	//parents rpcGenes undergo CrossOver
	parent1.rpcGenes = CrossOver(parent1.rpcGenes)
	parent2.rpcGenes = CrossOver(parent2.rpcGenes)

	//Select one of the parents' genes to pass to the next generation
	p1Gene := parent1.rpcGenes[rand.Intn(2)]
	gene1 := make(Gene, len(p1Gene))
	copy(gene1, p1Gene)
	p2Gene := parent2.rpcGenes[rand.Intn(2)]
	gene2 := make(Gene, len(p2Gene))
	copy(gene2, p2Gene)

	//Create a newborn Eukaryote
	newGenes := [2]Gene{gene1, gene2}
	rpcGene := make(Gene, 0)
	rpcGenes := [2]Gene{rpcGene, rpcGene}
	newBorn := Eukaryote{genes: newGenes, rpcGenes: rpcGenes, rpcPos: [2]int{0, 0}, oprtPos: [2]int{0, 0}}
	newBorn.CreateColor(parent1, parent2)
	newBorn.InitializeEnergy()

	//Check if there are empty space in the neighborhood
	nbrMap := FindEmptyNbr(r, c, newBoard)

	if len(nbrMap) == 0 { // No empty neighbor, battle with other cells nearby
		(&newBorn).Battle(r, c, parent1, parent2, newBoard)
	} else { // has empty neibor, Put the newborn into an empty space
		i := 0

		if len(nbrMap) > 1 { // more than 1 empty space, randomly choose one
			i = rand.Intn(len(nbrMap))
		}

		pos := nbrMap[i]
		(*newBoard)[pos[0]][pos[1]] = &newBorn
	}

}

//Battle() takes the parent's energy and compare it with the battle subject's energy, the one with more energy has higher chance to win
func (cell *Eukaryote)Battle(r, c int, p1, p2 *Eukaryote, newBoard *GameBoard) {
	// Find all neighbors that is not the same specied
	energy := (p1.energy + p2.energy) / 2

  nbrMap := FindEukaryoteBattleNbr(r, c, newBoard, cell)

  if len(nbrMap) == 0 {
    return
  }

  i := 0

  if len(nbrMap) != 1 { // randomly choose one battle subject in the neighborhood
    i = rand.Intn(len(nbrMap))
  }

  pos := nbrMap[i]
  sbj := (*newBoard)[pos[0]][pos[1]]
  ratio := float64(energy) / float64(sbj.ReportEnergy()) / 2 // the chance of
	win := false

	if energy >= 0 && sbj.ReportEnergy() <= 0 {
		win = true
	}

  if energy > 0 && sbj.ReportEnergy() > 0 && IsWin(ratio, 0.05) == 1 {
    win = true
  }
	if energy < 0 && sbj.ReportEnergy() < 0 && IsWin(ratio, 0.05) == 0 {
		win = true
	}

	if win == true {
		(*newBoard)[pos[0]][pos[1]] = cell
	}

	if win == false && sbj.ReportType() == "Eukaryote" && energy > 0{
		sbj.(*Eukaryote).energy -= energy
	}

	if win == false && sbj.ReportType() == "Prokaryote" && energy > 0{
		sbj.(*Prokaryote).energy -= energy
	}

	return
}

// IsWin() takes the ratio and form a normal distribution with mean = ratio
// If the number chosen >= 0.5 the newborn will win, otherwise, it will lose and will not be produced
func IsWin(ratio float64, variance float64) int {
  chance := rand.NormFloat64() * variance + ratio

  if chance < 0.5 {
    return 0
  } else {
    return 1
  }
}

// compare gene returns 1 if two genes are not the same, returns 0 if they are the same
func CompareGenes(g1, g2 Gene) int {
	diff := 0

	if len(g1) != len(g2) { // compare length
		return 1
	}

	// range through the gene
	for i := range g1 {
		if g1[i] != g2[i] {
			//fmt.Println("has diff!!!")
			diff = 1
			break
		}
	}
	return diff
}

// IsSameSpecies() returns 1 if the two eukaryote are of the same species, returns 0 if they are not
func IsSameSpecies(cell1, cell2 *Eukaryote) int {
  if CompareGenes(cell1.genes[0],cell2.genes[0]) == 0 && CompareGenes(cell1.genes[1],cell2.genes[1]) == 0 {
    return 1
  }

  if CompareGenes(cell1.genes[0],cell2.genes[1]) == 0 && CompareGenes(cell1.genes[1],cell2.genes[0]) == 0 {
    return 1
  }

  return 0
}

// CreateColor() will create a new color for the newborn if its gene is different from its parents
// If the same, it will inherite parents' color
// Otherwise, it will generate another color within a certain range
func (cell *Eukaryote)CreateColor(p1, p2 *Eukaryote) {
	if IsSameSpecies(p1, cell) == 1{
		cell.color = p1.color
		return
	}

	if IsSameSpecies(p2, cell) == 1{
		cell.color = p2.color
		return
	}

	R := 0
	B := int(p1.color[1]+p2.color[1])/2 + rand.Intn(10)
	G := int(p1.color[2]+p2.color[2])/2 + rand.Intn(10)
	Red := uint8(R)
	Blue := uint8(B % 205 + 50)
	Green := uint8(G % 205 + 50)
	cell.color = [3]uint8{Red, Blue, Green}
	return
}

// Mutate will decide what the next copied nucleotide to be
// if no mutation, it will copy the template gene
// Otherwise, it has a chance to have a deletion, insertion, and mismatch.
func (nucleotides *Gene) Mutate(nucleotide byte) Gene {
	if Deletion(IndelRate) == false && nucleotides.Insertion(nucleotide, IndelRate) == false && nucleotides.Mismatch(nucleotide, MismatchRate) == false {
		*nucleotides = append(*nucleotides, nucleotide)
	}
	return *nucleotides
}

// Create Deletion at a chance of rate
func Deletion(rate float64) bool {
	if rate -0.0 < 0.0001 {
		return false
	}

	num := int(1 / rate)
	if rand.Int()%num == 0 {
		return true
	} else {
		return false
	}
}

// Insertion happens at the chance of pre-defined IndelRate, a genome element will randomly selected
func (nucleotides *Gene) Insertion(nucleotide byte, rate float64) bool {
	if rate -0.0 < 0.0001 {
		return false
	}

	num := int(1 / rate)
	if rand.Int()%num == 0 {
		*nucleotides = append(*nucleotides, nucleotide)
		*nucleotides = append(*nucleotides, Operations[rand.Intn(len(Operations))])
		return true
	} else {
		return false
	}
}

// Insertion happens at the chance of pre-defined MismatchRate, a genome element will randomly selected to replace the current element
func (nucleotides *Gene) Mismatch(nucleotide byte, rate float64) bool {
	if rate -0.0 < 0.0001 {
		return false
	}

	num := int(1 / rate)
	if rand.Int()%num == 0 {
		*nucleotides = append(*nucleotides, Operations[rand.Intn(len(Operations))])
		return true
	} else {
		return false
	}
}

// WithinBoard() check if a given row and column value is within the current GameBoard
func WithinBoard(r, c int, currentBoard *GameBoard) bool {
	if r < 0 || r >= len(*currentBoard) { // row out of range
		return false
	} else if c < 0 || c >= len((*currentBoard)[0]) { // col out of range
		return false
	} else {
		return true
	}
}

// FindEukaryoteBattleNbr() find the neighbor of different species as the battle subject
// It takes the current position r, c, the current board, and the cell.
// Then compare the cell's type and genome content with neighbors, identify those are different from itself.
// return a map of positions of the format: [index]:pos, pos is a 2 int array
func FindEukaryoteBattleNbr(r, c int, newBoard *GameBoard, cell *Eukaryote) map[int][2]int {
  nbrMap := make(map[int][2]int, 0)
  i := 0

  //Range over all neighbor
  for row := r-1; row <= r + 1; row++ {
    for col := c-1; col <= c+1; col++ {
      if row == r && col == c {
        continue
      } else if WithinBoard(row, col, newBoard) == false {
        continue
      } else if (*newBoard)[row][col] == nil {
        continue
      } else {
        sbj := (*newBoard)[row][col]

        if sbj.ReportType() == "Eukaryote" && IsSameSpecies(cell, sbj.(*Eukaryote)) == 1 {
          continue
        }
        nbrMap[i] = [2]int{row, col}
        i++
      }
    }
  }
  return nbrMap
}

// FindEukayoteNbr() returns all Eukayote neighborhood position as a map:
// with their index as key, and there position, i.e. row and col as value
func FindEukayoteNbr(r, c int, newBoard *GameBoard) map[int][2]int {
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
			}
		}
	}
	return nbrMap
}

// FindEukayoteNbr() returns all empty neighborhood position as a map:
// with their index as key, and there position, i.e. row and col as value
func FindEmptyNbr(r, c int, currentBoard *GameBoard) map[int][2]int {
	nbrMap := make(map[int][2]int, 0)
	i := 0
	for row := r - 1; row <= r+1; row++ {
		for col := c - 1; col <= c+1; col++ {
			if row == r && col == c {
				continue
			} else if WithinBoard(row, col, currentBoard) == false {
				continue
			} else if (*currentBoard)[row][col] == nil {
				nbrMap[i] = [2]int{row, col}
				i++
			}
		}
	}
	return nbrMap
}

// CrossOver() takes parent genome and return random cross-over that passes on to the next generation
// The chance of a single crossover happens is 50%
func CrossOver(pairGenes [2]Gene) [2]Gene {
  i := 0
  for rand.Intn(3) != 0 {
    pairGenes = CrossOverOnce(pairGenes)
    i ++
    if i >= 2 {
      break
    }
  }
  return pairGenes
}

// CrossOverOnce() create a single crossover. It randomly choose the beginning and ending along a genome and exchange the segment with another copy of genome
func CrossOverOnce(pairGenes [2]Gene) [2]Gene {
  // Create a distribution that favors picking a center at the ends of a sequence
  // Create such distribution by adding up two normal distribution
  n := len(pairGenes[0])
  stdDev := float64(n) / 12.0
  stdMean1 := float64(n) / 6.0
  stdMean2 := 5.0 * float64(n) / 6.0
  center := (rand.NormFloat64() * stdDev + stdMean1) + (rand.NormFloat64() * stdDev + stdMean2)
  halflength := rand.NormFloat64() * stdDev + stdMean1
  beg := int(center - halflength)
  end := int(center + halflength)

  //Check validity of the beg and end positions
  if int(halflength) <= 0 || beg >= n-1 || end <= 1 {
    return pairGenes
  }
  if beg < 0 {
    beg = 0
  }
  if end > n-1 {
    end = n-1
  }

  // Now, select crossover area on another copy of genome
  n2 := len(pairGenes[1])
  beg2 := -1
  end2 := -1
  j := 0
  for beg2 < 0 || beg2 >= n2 -1 || end2 <= 1 || end2 > n2-1 || beg2 >= end2 {
    beg2 = int(rand.NormFloat64() * 2) + beg
    end2 = int(rand.NormFloat64() * 2) + end
    j ++
    if j > 3 {
      return pairGenes
    }
  }

  gene1 := append(pairGenes[0][0:beg], pairGenes[1][beg2:end2]...)
  gene1 = append(gene1, pairGenes[0][end:]...)

  gene2 := append(pairGenes[1][0:beg2], pairGenes[0][beg:end]...)
  gene2 = append(gene2, pairGenes[1][end2:]...)

  return [2]Gene{gene1, gene2}
}

// End of eukaryote_operations.go
// Written by Yinjie Gao
// Dec,03, 2018
