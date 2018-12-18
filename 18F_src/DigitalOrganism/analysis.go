// Written by Yinjie Gao
// analysis.go contains functions for statistical analysis
// It has 5 major output
// 1. pctFile: percentage of Prokaryote and Eukaryote on the board
// 2. lengthFile: the genome lengths statistics for Prokaryote and Eukaryote on the current board
// 3. genomeFile: genome content in Prokaryote and Eukaryote, i.e. the percentage of 'N', 'E', 'R' gene element
// 4. topFile: contain the genome sequence for top 10-energy organisms
// 5. LengthHeatMap: heatmap for length, colored blue. Darker color represents longer length.
// 6. EnergyHeatMap: heatmap for energy, coloered red. Darker color represents more energy

package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
)

// By default, statistical analysis is execuated every 50 iterations except for topFile and heatmap
// topFile and heatmap are generated every 500 iterations
func GetStat(boards []GameBoard, interval int) {
	pctFile, _ := os.Create(os.Args[9])
	defer pctFile.Close()
	lengthFile, _ := os.Create(os.Args[10])
	defer lengthFile.Close()
	genomeFile, _ := os.Create(os.Args[11])
	defer genomeFile.Close()
	topFile, _ := os.Create(os.Args[12])
	defer topFile.Close()

	pctHeader := "round\tcountE\tcountP\tcountT\tpctE\tpctP\n"
	pctFile.WriteString(pctHeader)
	lengthHeader := "round\tmeanE\tsdE\tminE\tmaxE\tmeanP\tsdP\tminE\tmaxE\tmeanT\tsdT\n"
	lengthFile.WriteString(lengthHeader)
	genomeHeader := "round\tNpctE\tEpctE\tRpctE\tNpctP\tEpctP\tRpctP\tNpctT\tEpctT\tRpctT\n"
	genomeFile.WriteString(genomeHeader)
  topHeader := "round\ttop\tenergy\tgenome content\n"
	topFile.WriteString(topHeader)

	for i, board := range boards {
		if i%interval != 0 {
			continue
		}

		count, length, genomeE, genomeP, energy := board.GetData()
		GetPct(i, count, pctFile)
		LengthDistribution(i, length, count, lengthFile)   // mean, min, max, sd for E/P/Total
		GenomeComposition(i, genomeE, genomeP, genomeFile) // pct of N, E, R in the genome content in each individual & in the population

		if i % 500 != 0 {
			continue
		}

		//DrawHeatmap
		//Get Length range
		minE, maxE := GetMinMax(length[0])
		minP, maxP := GetMinMax(length[1])
		min := minE
		if minP < min {
			min = minP
		}
		max := maxE
		if maxP > max {
			max = maxP
		}

		pngName := "lengthHeatmapBoard" + strconv.Itoa(i) + "_" + strconv.Itoa(min)+ "-" + strconv.Itoa(max) + ".png"

		LengthHeatMap(board, min, max, pngName)

		minE, maxE = GetMinMax(energy[0])
		minP, maxP = GetMinMax(energy[1])
		min = minE
		if minP < min {
			min = minP
		}
		max = maxE
		if maxP > max {
			max = maxP
		}

		pngName = "energyHeatmapBoard" + strconv.Itoa(i) + "_" + strconv.Itoa(min)+ "-" + strconv.Itoa(max) + ".png"

		EnergyHeatMap(board, min, max, pngName)

		TopEnergySequence(i, energy, board, topFile)
	}

}

// Get length, percentage, genome content data from boards
func (board *GameBoard) GetData() ([]int, [2][]int, [][]byte, [][]byte, [2][]int) {
	count := make([]int, 2) //{E count, P count}
	lengthE := make([]int, 0)
	lengthP := make([]int, 0)
	genomeE := make([][]byte, 0)
	genomeP := make([][]byte, 0)
	energyE := make([]int, 0)
	energyP := make([]int, 0)

	for r := range *board {
		for c := range (*board)[r] {
			if (*board)[r][c] == nil {
				continue
			}

			if (*board)[r][c].ReportType() == "Eukaryote" {
				count[0]++
				len1 := len((*board)[r][c].(*Eukaryote).genes[0])
				len2 := len((*board)[r][c].(*Eukaryote).genes[1])
				lenT := len1 + len2
				lengthE = append(lengthE, lenT)
				energyE = append(energyE, (*board)[r][c].(*Eukaryote).energy)
				genome := make([]byte, lenT)
				copy(genome[0:len1], (*board)[r][c].(*Eukaryote).genes[0])
				copy(genome[len1:], (*board)[r][c].(*Eukaryote).genes[1])
				genomeE = append(genomeE, genome)
			}

			if (*board)[r][c].ReportType() == "Prokaryote" {
				count[1]++
				lenT := len((*board)[r][c].(*Prokaryote).genes)
				lengthP = append(lengthP, lenT)
				energyP = append(energyE, (*board)[r][c].(*Prokaryote).energy)
				genome := make([]byte, lenT)
				copy(genome, (*board)[r][c].(*Prokaryote).genes)
				genomeP = append(genomeP, genome)
			}
		}
	}

	length := [2][]int{lengthE, lengthP} // {E lengths, P lengths}
	energy := [2][]int{energyE, energyP}

	return count, length, genomeE, genomeP, energy
}

// calculate percentage
func GetPct(numInteration int, count []int, pctFile *os.File) {
	countT := count[0] + count[1]
	pctE := float64(count[0]) / float64(countT)
	pctP := float64(count[1]) / float64(countT)
	fmt.Fprintf(pctFile, "%d\t%d\t%d\t%d\t%.*e\t%.*e\n", numInteration, count[0], count[1], countT, 4, pctE, 4, pctP)
}

// Calculate length mean & standard deviation
func LengthDistribution(numInteration int, length [2][]int, count []int, lengthFile *os.File) {
	minE, maxE := GetMinMax(length[0])
	meanE := float64(SumSlice(length[0])) / float64(count[0])
	sdE := CalculateSD(length[0], meanE)

	minP, maxP := GetMinMax(length[1])
	meanP := float64(SumSlice(length[1])) / float64(count[1])
	sdP := CalculateSD(length[1], meanE)

	total := append(length[0], length[1]...)
	meanT := float64(SumSlice(total)) / float64(count[0]+count[1])
	sdT := CalculateSD(total, meanT)

	fmt.Fprintf(lengthFile, "%d\t%.*e\t%.*e\t%.*e\t%.*e\t%.*e\t%.*e\t%.*e\t%.*e\t%.*e\t%.*e\n", numInteration, 4, meanE, 4, sdE, 4, float64(minE), 4, float64(maxE), 4, meanP, 4, sdP, 4, float64(minP), 4, float64(maxP), 4, meanT, 4, sdT)
	return
}

// Report maximum amd minimum number in an integer string
func GetMinMax(s []int) (int, int){
	if len(s) == 0 {
		return 0,0
	}
	min := s[0]
	max := s[0]

	for i := range s {
		if s[i] < min {
			min = s[i]
		}
		if s[i] > max {
			max = s[i]
		}
	}

	return min, max
}

// Calculte genome content percentage and write genome composition file
func GenomeComposition(numInteration int, genomeE, genomeP [][]byte, genomeFile *os.File) {
	numNE, numEE, numRE, numE := CountGenePool(genomeE)
	numNP, numEP, numRP, numP := CountGenePool(genomeP)

	pctNE := float64(numNE) / float64(numE)
	pctEE := float64(numEE) / float64(numE)
	pctRE := float64(numRE) / float64(numE)

	pctNP := float64(numNP) / float64(numP)
	pctEP := float64(numEP) / float64(numP)
	pctRP := float64(numRP) / float64(numP)

	pctNT := float64(numNE+numNP) / float64(numE+numP)
	pctET := float64(numEE+numEP) / float64(numE+numP)
	pctRT := float64(numRE+numRP) / float64(numE+numP)

	fmt.Fprintf(genomeFile, "%d\t%.*e\t%.*e\t%.*e\t%.*e\t%.*e\t%.*e\t%.*e\t%.*e\t%.*e\n", numInteration, 4, pctNE, 4, pctEE, 4, pctRE, 4, pctNP, 4, pctEP, 4, pctRP, 4, pctNT, 4, pctET, 4, pctRT)
}

func CountGenePool(genome [][]byte) (int, int, int, int) {
	numN := 0
	numE := 0
	numR := 0
	for i := range genome {
		for j := range genome[i] {
			if genome[i][j] == 'N' {
				numN++
			} else if genome[i][j] == 'E' {
				numE++
			} else if genome[i][j] == 'R' {
				numR++
			}
		}
	}
	return numN, numE, numR, numN + numE + numR
}

// Calculate the sum over a slice of integer
func SumSlice(intSlice []int) int {
	total := 0
	for i := range intSlice {
		total += intSlice[i]
	}
	return total
}

// Calculate the standard deviation
func CalculateSD(intSlice []int, mean float64) float64 {
	sum := 0.0
	for i := range intSlice {
		sum += (float64(intSlice[i]) - mean) * (float64(intSlice[i]) - mean)
	}
	sd := math.Sqrt(sum / float64(len(intSlice)-1))
	return sd
}

// Generate LengthHeatMap named as pngName
// Darkest color represents the longest length in the current board
func LengthHeatMap(board GameBoard, min, max int, pngName string) {
	lengthCanvas := DrawLengthBoard(board, min, max, 5)
	lengthCanvas.SaveToPNG(pngName)
}

// Generate EnergyHeatMap named as pngName
// Darkest color represents the maximum energy in the current board
func EnergyHeatMap(board GameBoard, min, max int, pngName string) {
	lengthCanvas := DrawEnergyBoard(board, min, max, 5)
	lengthCanvas.SaveToPNG(pngName)
}

// Write TopEnergyFile
func TopEnergySequence(numInteration int, energy [2][]int, board GameBoard, topFile *os.File) {
	top := FindTop(append(energy[0], energy[1]...), 10)

	for i := 0; i < 10; i++ {
		line := strconv.Itoa(numInteration) + "\t" + strconv.Itoa(i) + "\t"
		PrintMaxGene(i, board, top[i],line, topFile)
	}

	return
}

// Get top-energy organisms' genome sequences
func PrintMaxGene(i int, board GameBoard, max int, line string, topFile *os.File) {
	for r := range board {
		for c := range board[r] {
			if board[r][c] != nil && board[r][c].ReportEnergy() == max {
				if board[r][c].ReportType() == "Eukaryote" {
					line = line + strconv.Itoa(max) + "\t" + string(board[r][c].(*Eukaryote).genes[0]) + ", " + string(board[r][c].(*Eukaryote).genes[1]) + "\t\n"
				} else {
					line = line + strconv.Itoa(max) + "\t" + string(board[r][c].(*Prokaryote).genes) + "\t\n"
				}
				fmt.Fprintf(topFile, line)
				return
			}
		}
	}
}

// Find the 10 largest energy on the current board
func FindTop(intSlice []int, n int) []int {
	top := make([]int, n)
	max := 0

	for i := 0; i < n; i++ {
		max, intSlice = FindMaxAndDelete(intSlice)
	  top[i] = max
	}

	return top
}

// Find maximum number in a slice and delete it
func FindMaxAndDelete(intSlice []int) (int, []int) {
	max := 0

	if len(intSlice) == 0 {
		return max, intSlice
	}

	max = intSlice[0]
	maxIndex := 0

	for i := range intSlice {
		if max < intSlice[i] {
			max = intSlice[i]
			maxIndex = i
		}
	}

	intSlice = append(intSlice[0:maxIndex], intSlice[maxIndex+1:]...)
	return max, intSlice
}

// End of analysis.go
// Written by Yinjie Gao
// Dec,03, 2018
