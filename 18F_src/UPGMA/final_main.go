package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type DistanceMatrix [][]float64

type Tree []*Node

type Node struct {
	age            float64
	label          string
	child1, child2 *Node
}

func main() {
	// let's take a file name from the user
	filename := os.Args[1]

	// read in from file and parse matrix and species names
	mtx, speciesNames := ReadMatrixFromFile(filename)

	t := UPGMA(mtx, speciesNames)

	t.PrintGraphViz()
}

//UPGMA!
func UPGMA(mtx DistanceMatrix, speciesNames []string) Tree {
	t := InitializeTree(speciesNames)
	clusters := t.InitializeClusters() // pointing clusters at leaves

	numLeaves := len(speciesNames)

	// engine of UPGMA: set ages and children of all internal nodes according to algorithm, updating matrix and clusters as we go
	for k := numLeaves; k < 2*numLeaves-1; k++ {
		row, col, val := FindMinElt(mtx)
		fmt.Println(row, col)
		// big assumption: col > row

		t[k].age = val / 2
		t[k].child1 = clusters[row]
		t[k].child2 = clusters[col]


		// we have set t[k], now we update matrix and clusters
		mtx = AddRowCol(mtx, clusters, row, col)
		mtx = DelRowCol(mtx, row, col)

		clusters = append(clusters, t[k])
		clusters = DelClusters(clusters, row, col)

		// ALERT: DEBUGGING LINE SINCE IT SEEMS ISSUE IS WITH LAST EDGE
		if k == 2*numLeaves-2 {
			fmt.Println("Row:", row)
			fmt.Println("Col:", col)
			fmt.Println("Age:", t[k].age)
			fmt.Println("Child 1", t[k].child1.label)
			fmt.Println("Child 2", t[k].child2.label)
		}

	}

	return t
}

//InitializeTree takes the n names of our present-day species (leaves) and returns a rooted binary tree with 2n-1 total nodes, where the leaves are the first n and have the associated species names.
func InitializeTree(speciesNames []string) Tree {
	numLeaves := len(speciesNames)
	var t Tree = make([]*Node, 2*numLeaves-1)

	//create our 2n-1 nodes, and assign labels (no children yet)
	for i := range t {
		// create a node (default age: 0)
		var vx Node
		if i < numLeaves { // set the species name of leaves
			vx.label = speciesNames[i]
		} else { // set internal node names to their integer label
			vx.label = "Ancestor Species " + strconv.Itoa(i)
		}
		// point t[i] at current Node
		t[i] = &vx
	}
	return t
}

//InitializeClusters is a Tree method that returns a slice of pointers to the leaves of the Tree
func (t Tree) InitializeClusters() []*Node {
	numLeaves := (len(t) + 1) / 2
	clusters := make([]*Node, numLeaves)

	// our pointers are nil.  So point them at leaves!
	for i := 0; i < numLeaves; i++ {
		clusters[i] = t[i]
	}

	return clusters
}

//FindMinElt takes a DistanceMatrix and returns the row, column, and value corresponding to a minimum element.
//Assumption that col > row
func FindMinElt(mtx DistanceMatrix) (int, int, float64) {
	if len(mtx) <= 1 || len(mtx[0]) <= 1 {
		panic("We gave too small a matrix to FindMinElt")
	}
	row := 0
	col := 1
	val := mtx[row][col]

	// range over distance matrix values, seeing if we can do better
	for i := 0; i < len(mtx); i++ {
		for j := i + 1; j < len(mtx[i]); j++ {
			if mtx[i][j] < val {
				val = mtx[i][j]
				row = i
				col = j
			}
		}
	}

	return row, col, val
}

//AddRowCol takes a DistanceMatrix, a slice of current clusters, and a row/col index (col > row).
//It returns the matrix corresponding to "gluing" clusters[row] and clusters[col] together and forming a new row/col of the matrix (no deletions yet).
func AddRowCol(mtx DistanceMatrix, clusters []*Node, row, col int) DistanceMatrix {
	n := len(mtx)
	newRow := make([]float64, n+1)
	for r := 0; r < n; r++ {
		size1 := clusters[row].CountLeaves()
		size2 := clusters[col].CountLeaves()
		// now we compute the weighted average distance from r to the updated merged cluster.
		newRow[r] = (float64(size1)*mtx[row][r] + float64(size2)*mtx[col][r]) / (float64(size1) + float64(size2))
	}

	//append the new row to our matrix
	mtx = append(mtx, newRow)
	
	// append each value in the new row to end of appropriate column for symmetry
	for c := 0; c < n; c++ {
		mtx[c] = append(mtx[c], newRow[c])
	}

	return mtx
}

//DelRowCol takes a distance matrix and a row/col index and deletes the row and column indicated, returning the resulting matrix
func DelRowCol(mtx DistanceMatrix, row, col int) DistanceMatrix {
	// IMPORTANT! Delete row with index col first
	mtx = append(mtx[:col], mtx[col+1:]...)
	mtx = append(mtx[:row], mtx[row+1:]...)

	// delete columns with indices col (first) and then row : range over the rows and delete given elements in each row
	for i := range mtx {
		mtx[i] = append(mtx[i][:col], mtx[i][col+1:]...)
		mtx[i] = append(mtx[i][:row], mtx[i][row+1:]...)
	}

	return mtx
}

//DelClusters takes a slice of Node pointers along with a row/col index and deletes the clusters in the slice corresponding to these indices.
//Assume col > row
func DelClusters(clusters []*Node, row, col int) []*Node {
	clusters = append(clusters[:col], clusters[col+1:]...)
	clusters = append(clusters[:row], clusters[row+1:]...)
	return clusters
}


//CountLeaves is a Node method that counts the number of leaves in the tree rooted at the node. It returns 1 at a leaf.
func (vx *Node) CountLeaves() int {
	// base case: leaf
	if vx.child1 == nil || vx.child2 == nil {
		return 1
	}
	// otherwise, apply recursive step
	return vx.child1.CountLeaves() + vx.child2.CountLeaves()
}


//ReadMatrixFromFile takes a file name and reads the information in this file to produce
//a distance matrix and a slice of strings holding the species names.  The first line of the
//file should contain the number of species.  Each other line contains a species name
//and its distance to each other species.
func ReadMatrixFromFile(fileName string) (DistanceMatrix, []string) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error: couldn't open the file")
		os.Exit(1)
	}
	var lines []string = make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if scanner.Err() != nil {
		fmt.Println("Sorry: there was some kind of error during the file reading")
		os.Exit(1)
	}
	file.Close()

	mtx := make(DistanceMatrix, 0)
	speciesNames := make([]string, 0)

	for idx, _ := range lines {
		if idx >= 1 {
			row := make([]float64, 0)
			nums := strings.Split(lines[idx], "\t")
			for i, num := range nums {
				if i == 0 {
					speciesNames = append(speciesNames, num)
				} else {
					n, err := strconv.ParseFloat(num, 64)
					if err != nil {
						fmt.Println("Error: Wrong format of matrix!")
						os.Exit(1)
					}
					row = append(row, n)
				}
			}
			mtx = append(mtx, row)
		}
	}
	fmt.Print(mtx)
	return mtx, speciesNames
}

// PrintGraphViz prints the tree in GraphViz format, where directed = true
// if we desire to print a directed graph and directed = false for an
// undirected graph.
func (t Tree) PrintGraphViz() {
	fmt.Println("strict digraph {")
	for i := range t {
		if t[i].child1 != nil && t[i].child2 != nil {
			//print first edge
			fmt.Print("\"", t[i].label, "\"")
			fmt.Print("->")
			fmt.Print("\"", t[i].child1.label, "\"")
			fmt.Print("[label = \"")
			fmt.Printf("%.2f", t[i].age-t[i].child1.age)
			fmt.Print("\"]")
			fmt.Println()

			//print second edge
			fmt.Print("\"", t[i].label, "\"")
			fmt.Print("->")
			fmt.Print("\"", t[i].child2.label, "\"")
			fmt.Print("[label = \"")
			fmt.Printf("%.2f", t[i].age-t[i].child2.age)
			fmt.Print("\"]")
			fmt.Println()
		}
	}
	fmt.Println("}")
}
