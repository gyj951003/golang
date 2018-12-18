package main

import (
	//"os"
	"io/ioutil"
	"strconv"
	"strings"
	"fmt"
	"math/rand"
)

func main() {
	filename := "debates.txt"
	candidate1Percentages, electoralVotes := ReadPollingData(filename)
	var numTrials uint = 1000000
	marginOfError := 0.1
	probability1, probability2, probabilityTie := SimulateElection(candidate1Percentages, electoralVotes, numTrials, marginOfError)
	candidate1, candidate2 := "Clinton", "Trump"
	fmt.Println(candidate1, "percentage chance of winning is", probability1)
	fmt.Println(candidate2, "percentage chance of winning is", probability2)
	fmt.Println("The percentage chance of a tie is", probabilityTie)
}

//ReadPollingData takes in a filename from which to read data.
//It returns state names, poll results for the two candidates, and electoral college votes.
func ReadPollingData(filename string) (map[string]float64, map[string]uint) {
	data, err := ioutil.ReadFile(filename)
	Check(err)

	// each line of the datadata is stored in the format:
	// stateName, x, y, votes where x = Clinton percentage,
	// y = Trump percentage, and votes = number of Electoral College votes.
	lines := strings.Split(string(data), "\n")

	// we assume that the dataset has 51 lines, one for each state + DC
	//we create four arrays, one to hold each collection of data across all 51 lines
	candidate1Percentages := make(map[string]float64)
	electoralVotes := make(map[string]uint)

	//now, parse out data on each line and add to appropriate array
	for j := 0; j < 51; j++ {
		var line = strings.Split(lines[j], ",")
		stateName := line[0]
		percentage1, err := strconv.ParseFloat(line[1], 64)
		// normalize percentages between 0 and 1
		candidate1Percentages[stateName] = percentage1 / 100.0
		votes, err := strconv.Atoi(line[3])
		electoralVotes[stateName] = uint(votes)
		Check(err)
	}

	return candidate1Percentages, electoralVotes
}

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func AddNoise(poll float64, marginOfError float64) float64 {
	noise := rand.NormFloat64() / 2 * marginOfError
	return poll + noise
}

func SimulateOneElection(candidate1Percentages map[string]float64, electoralVotes map[string]uint, marginOfError float64) (uint, uint) {
	var collegeVote1 uint
	var collegeVote2 uint
	var poll, adjustedpoll float64
	var votes uint
	for k := range candidate1Percentages {
		poll = candidate1Percentages[k]
		votes = electoralVotes[k]
		adjustedpoll = AddNoise(poll, marginOfError)
		if adjustedpoll >= 0.5 {
			collegeVote1 = collegeVote1 + votes
		} else {
			collegeVote2 = collegeVote2 + votes
		}
	}
	return collegeVote1, collegeVote2
}

func SimulateElection(candidate1Percentages map[string]float64,
	electoralVotes map[string]uint, numTrials uint, marginOfError float64) (float64, float64, float64) {
	winCount1 := 0
	winCount2 := 0
	tieCount := 0
	var collegeVote1, collegeVote2 uint
	for i := 0; i < int(numTrials); i++ {
		collegeVote1, collegeVote2 = SimulateOneElection(candidate1Percentages, electoralVotes, marginOfError)
		if collegeVote1 > collegeVote2 {
			winCount1++
		} else if collegeVote1 < collegeVote2 {
			winCount2++
		} else {
			tieCount++
		}
	}
	probability1 := float64(winCount1)/float64(numTrials)
	probability2 := float64(winCount2)/float64(numTrials)
	probabilityTie := float64(tieCount)/float64(numTrials)
	return probability1, probability2, probabilityTie
}
