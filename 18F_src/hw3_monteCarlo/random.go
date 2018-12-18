package main

import (
  "math/rand"
  //"time"
)

// =====================================
// Monte Carlo Relative Prime
// Write a function RelativelyPrimeProbability that takes three int as input(x,
// y, numPairs). Uniformly generate pairs of integers between x and y, inclusive
// (numPairs total), and return an estimate of probability of the pairs being
// relatively prime(type float64).

// Your program needs to panic if any of the input is non-positive.

// numPairs <= 100000. For all inputs, either 1 <= x, y <= 1000, or
// x = 1 and y <= 1000000.

// The PRNG will be seeded by the autograder and you are
// free to use it in whatever way you like. We will grade your implementation
// by checking that your output is within a certain amount of the correct answer.
// If your implementation is correct, it has at least 99.9999%
// probability to pass. If you believe your implementation
// is correct but autograder have it wrong, buy a lottery today as it's your
// (super) lucky day!

// Sample Input: (2, 4, 10)
// Sample Output: 0.5
// Note: You can get anything from 0 to 1. The actual probability for these x and y is 4/9 = .444.

// =====================================
func CalculateGCD(a, b int) int {
// CalculateGCD is to calculate the greatest common number of two nonegative integers
  if a == b {
    return a
  } else if a < b {
    return CalculateGCD(a, b - a)
  } else { // b < a
    return CalculateGCD(b, a - b)
  } //end else
} //CalculateGCD()

func IsRelativePrime(a, b int) bool {
// IsRelativePrime is to test wether two nonegative integers are relative primers
// Two integers are relative primers iff their GCD is 1.
  return CalculateGCD(a, b) == 1
} //IsRelativePrime()

func GenerateRandomNumber(x, y int) (int, int) {
//GenerateRandomNumber generates two random numbers within the range x and y
  a := rand.Intn(y - x + 1) + x
  b := rand.Intn(y - x + 1) + x
  return a, b
} //GenerateRandomNumber()

func RelativelyPrimeProbability(x, y, numPairs int) float64 {
// RelativelyPrimeProbability simulates two integers within the range from x to y for numPairs times,
// and determine the probability of haivng relative prime numbers
  if (x <= 0 || y <= 0 || numPairs <= 0) {
    panic("Input should be positive integers!")
  } // check inputs are positive

  primeCount := 0
  for i := 0; i < numPairs; i++ {
    a, b := GenerateRandomNumber(x, y)
    if IsRelativePrime(a, b) == true {
      primeCount++
    } // end if
  } // end for

  return float64(primeCount)/float64(numPairs)
} // RelativelyPrimeProbability()



// =====================================
// Repeat Sequences
// Write a function HasRepeat that takes a slice of integers as input and
// returns a boolean indicating whether there is a repeated value (a value that
// appears twice or more in the slice).

// Length of slice <= 500. Elements in slice <= 100000000 in absolute value.
// Sample Input: [1, 2, 3, 4, 5, 6, 7, 1, 2, 3]
// Sample Output: true
// =====================================

func HasRepeat(numbers []int) bool {
  hasRepeat := false
  length := len(numbers)
  for i := 0; i < length; i++ {
    for j := i + 1; j < length; j++ {
      if numbers[i] == numbers[j] {
        hasRepeat = true
        break
      } //end if
    } // end for j
    if hasRepeat == true {
      break
    } // end if
  } // end for i
  return hasRepeat
} // HasRepeat()


// =====================================
// Monte Carlo Birthday Paradox
// Write a function BirthdayParadox that takes two integers as input(numPeople, numTrials).
// Simulate numPeople random birthdays and estimate the probability that at least two
// of them are the same (sample numTrials times total). Return the estimate
// as a float64.

// Your program should panic if either of the inputs is non-positive.

// numPeople <= 50, numTrials <= 100000. We have the same guarantee regarding
// testing your answer for correctness as in the Monte Carlo relatively prime problem.

// Sample Input: numPeople = 2, numTrials = 1000
// Sample Output: 0.003
// Note: You can get anything from 0 to ~0.01. The actual probability is
// 1/365 = 0.0027...

// =====================================
func GenerateBirthdays(numPeople int) []int {
  birthdays := make([]int, 0) // store numPeople people's birthdays
  for i := 0; i < numPeople; i++ {
    //rand.Seed(time.Now().UTC().UnixNano())
    birthday := rand.Intn(365) + 1
    birthdays = append(birthdays, birthday)
  } //end for
  return birthdays
} //GenerateBirthdays()

func BirthdayParadox(numPeople, numTrials int) float64 {
  if (numPeople <= 0 || numTrials <= 0) { // check input
    panic("Input should be positive integers!")
  } //end if

  repeatCount := 0
  for i := 0; i < numTrials; i++ {
    birthdays := GenerateBirthdays(numPeople)
    if HasRepeat(birthdays) == true {
      repeatCount++
    } //end if
  } //end for
  return float64(repeatCount)/float64(numTrials)
}
//The smallest value of numPeople is around 23 for giving >0.5 chance of two people sharing the same birthday


// =====================================
// Period Length
// Write a function ComputePeriodLength that takes a slice of int as input and
// check if you can decide the period of the input sequence. If no, return 0;
// If yes, return the period.

// Length of slice <= 500. Elements in slice <= 100000000 in absolute value.
// Sample Input: [1, 2, 3, 4, 5, 6, 7, 1, 2, 3]
// Sample Output: 7
// =====================================
func ComputePeriodLength(numbers []int) int {
  length := len(numbers)
  interval := 0
  for i := 0; i < length; i++ {
    for j := i + 1; j < length; j++ {
      if numbers[i] == numbers[j] {
        interval = j - i
        break
      } //end if
    } // end for j
    if interval != 0 {
      break
    } // end if
  } // end for i
  return interval
} // ComputePeriodLength()


//=====================================
// Number of Digits
// Write a function CountNumDigits that takes an integer x(type int) as input,
// and return number of digits in x. By convention, 0 have 1 digit, and (-x)
// have same number of digits as x, for any positive integer x.

// You should not rely on any functions from imported packages for this problem.

// Input number <= 10^10.
// Sample Input: 12345
// Sample Output: 5
//=====================================
func CountNumDigits(n int) int {
  if n < 0 {
    n = -1 * n
  }
  if n / 10 != 0 {
    return CountNumDigits(n/10) + 1
  } else {
    return 1
  }
} //CountNumDigits()


// =====================================
// SquareMiddle
// Write a function SquareMiddle that takes two ints: x and numDigits, and
// return the number formed by middle numDigits of x squared.
// Additionally, we ask you to perform sanity checks. If numDigits is not even,
// or if x has more than 2*numDigits of digit, or either input is negative,
// your function should panic.

// For all valid data: 0 <= x < 10^8, 0 < numDigits <= 8
// Your function need to panic if the either of the input falls outside
// of valid range, OR if x is no less than 10^numDigits, OR numDigits is not even.

// Sample Input 1: x = 55, numDigits = 2
// Sample Output 1: 2 (55 * 55 = 3025)
// Sample Input 2: x = 55, numDigits = 3
// Sample Output 2: None(You should panic!)
// =====================================
func Pow10(base, power int) int {
  result := 1
  for i := 0; i < power; i++ {
    result = result * base
  }
  return result
} // Pow10()

func SquareMiddle(x, numDigits int) int {
  //Input check
  if (x < 0 || numDigits < 0) {
    panic("Inputs should be non-negative integers!")
  }
  if (numDigits % 2 != 0) {
    panic("The number of digits should be even!")
  }
  if (CountNumDigits(x) > numDigits) {
    panic("The number of digit of x^2 should be no more than the number of digits!")
  }
  //Complete input check
  squareMiddle := x * x
  squareMiddle = squareMiddle / Pow10(10, numDigits/2)
  squareMiddle = squareMiddle % Pow10(10, numDigits)
  return squareMiddle
}

// =====================================
// SquareMiddle Generator
// Now that you have every tool you need, let's implement the actual generator.
// Write a function GenerateMiddleSquareSequence that takes two ints: seed and
// numDigits, and output the generated sequence right after first repeated
// element appears(stop at this point; This is when you know the period of your
// sequence).

// 0 <= seed < 10^numDigits. 2 <= numDigits <= 6. We guarantee the output
// sequence has less than 500 elements.
// Sample Input: seed = 70, numDigits = 2
// Sample Output: [70, 90, 10, 10]
// =====================================

func GenerateMiddleSquareSequence(seed, numDigits int) []int {
  seq := make([]int, 0)
  seq = append(seq, seed)
  for true {
    seed = SquareMiddle(seed, numDigits)
    seq = append(seq, seed)
    if HasRepeat(seq) == true {
      break
    } //end if
  } //end for
  return seq
} //end GenerateMiddleSquareSequence()
