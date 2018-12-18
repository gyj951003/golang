package main

// import the "fmt" and "math" packages if useful ...
import "math"
//import "fmt"

/*
func main() {

  // You can use func main() to test functions, but comment this out before you submit to Autolab.
  fmt.Println(SumOfFirstNIntegers(10))

  fmt.Println(Permutation(10, 2))
  fmt.Println(Combination(1000, 998))

  fmt.Println(FibonacciArray(1))
  fmt.Println(FactorialArray(4))

  fmt.Println(MinArray([]int{5, 35, 85, 105, 5}))
  fmt.Println(GCDArray([]int{108, 7}))

  fmt.Println(TrivialPrimeFinder(7))
  fmt.Println(SieveOfEratosthenes(7))
  fmt.Println(ListOfPrimes(7))


  fmt.Println(IsPerfect(6))
  fmt.Println(NextPerfectNumber(17))

  fmt.Println(HailstoneSequence(19))

  fmt.Println(IsPrime(1))
  fmt.Println(IsPrime(4))
  fmt.Println(IsPrime(13))
  fmt.Println(NextTwinPrimers(13))
}
*/
/*==============================================================================
 * 1. Summing n integers
 *
 * Implement a function that takes a positive integer n, and returns
 * the sum of the integers from 1 to n. Return type should be int.
 *
 * n <= 20000.
 *
 *============================================================================*/
/*
Pseudocode 1:
SumOfFirstNIntegers1(n)
  result <- 0
  for i from 1 to n
    result <- result + i
  return result

Pseudocode 2:
SumOfFirstNIntegers2(n)
  result <- 0
  while i <= n
    result <- result + i
    i <- i + 1
*/

func SumOfFirstNIntegers(n int) int {
  result := 0
  for i := 1; i <= n; i++ {
    result = result + i
  }
  return result
}

/*
func SumOfFirstNIntegers2(n int) int {
  result := 0
  i := 1
  for i <= n {
    result = result + i
    i++
  }
  return result
}
*/

/*==============================================================================
 * 2. Combination and Permutation
 *
 * Implement two functions Permutation(n,k) and Combination(n,k) that takes
 * two non-negative integer n and k, and calculate the permutation and combination
 * function in int.
 *
 * Additionally, we ask you to write an implementation that is efficient for
 * small values of k and large values of n, as well as small values of (n-k).
 * For example, you need to be able to compute C(100000, 99999) pretty quickly.
 *
 * n <= 100000, min(k, n-k) <= 15. We guarantee that the return value will not
 * exceed int, but you should take care of any potential overflows during calculation.
 *
 *============================================================================*/
/*
Pseudocode:
Permutation(n, p)
  //requirement: p <= n
  result <- 1
  for i from n to n-p+1
    result = result * i
  return result
*/

func Permutation(n int, p int) int {
  result := 1
  for i := n; i > n - p; i-- {
    result = result * i
  }
  return result
}

/*
Pseudocode:
Combination(n, p)
  //requirement: p <= n
  if p > n/2
    return Combination(n, n-p)
  result <- 1
  for i from 1 to p
    result = result * (n + 1- i) / i
  return result
*/

func Combination(n int, p int) int {
  if p > n/2 {
    return Combination(n, n-p)
  }
  result := 1
  for i := 1; i <= p; i++{
    result = result * (n + 1 - i) / i
  }
  return result
}

/*==============================================================================
* 3. Fibonacci Array
*
* Write a function FibonacciArray(n) that takes a positive integer n and returns
* a slice of integers that contains first n Fibonacci numbers.
* Note that for this problem, the first two Fibonacci numbers are 1 and 1, so
* FibonacciArray(4) should return [1, 1, 2, 3].
*
* We will use array and slice interchangably for the problem descriptions.
*
* n <= 40.
*
*===========================================================================*/
func FibonacciArray(n int) []int {
  array := make([]int, n)
  array[0] = 1
  if n == 1 {
    return array
  }
  array[1] = 1
  if n == 2 {
    return array
  }
  for i := 2; i < n; i++ {
    array[i] = array[i-1] + array[i-2]
  }
  return array
}

/*==============================================================================
 * 4. Factorial Array
 *
 * Write a function FactorialArray(n) that takes a positive integer n and returns
 * a slice of int that satisifes ans[i] = i!, up to n. Note that arrays in
 * Go starts at 0, so FactorialArray(4) should return [1, 1, 2, 6, 24].
 *
 * n <= 15.
 *
 *===========================================================================*/
func FactorialArray(n int) []int {
  array := make([]int, n + 1)
  array[0] = 1
  if n == 0 {
    return array
  }
  for i := 1; i <= n; i++ {
    array[i] = array[i - 1] * i
  }
  return array
}

/*==============================================================================
 * 5. Min Array
 *
 * Write a function MinArray that takes a slice of ints(at least one element) and
 * returns the minimum value among all elements in the slice.
 *
 * Length of input slice <= 10. Elements in slice <= 10000 (in absolute value).
 *
 *============================================================================*/
/*
Pseudocode:
MinArray(array)
  min <- array[0]
  for each item in array
    if item < min
      update min to the item
  return min
*/

func MinArray(array []int) int {
  min := array[0]
  for i := 1; i < len(array); i++ {
    if array[i] < min {
      min = array[i]
    }
  }
  return min
}
/*==============================================================================
 * 6. GCD (Greatest Common Divisor) Array
 *
 * Write a function GCDArray that takes a slice of positive ints(at least one element)
 * and returns the GCD of all elements in the slice.
 * There are more than one acceptable approaches for this problem. You can try to
 * utilize your MinArray function, and you can also use some nice properties of
 * GCD function.
 *
 * Length of input slice <= 10. Elements in slice <= 10000.
 *
 *===========================================================================*/
/*
Pseudocode:
GCDArray(array)
  min <- MinArray(array)
  for i from min to 1
    if i is the common divisor of all other elements in array
      return i
*/

func GCDArray(array []int) int {
  min := MinArray(array)
  gcd := 1
  for i := min; i >= 1; i-- {
    isGCD := 1
    for j := 0; j < len(array); j++ {
      if array[j] % i != 0 {
        isGCD = 0
        break
      }
    }
    if isGCD == 1 {
      gcd = i
      break
    }
  }
  return gcd
}

/*=======================================================================
 * 7. The Trivial Prime Finder
 *
 * Write a function TrivialPrimeFinder that takes a positive integer n,
 *  and return an array of length (n+1) such that a[i] is true if and only
 * if i is a prime. Note that array indexes in Go starts at 0, so by
 * convention, we have a[0] = a[1] = false.
 *
 * As the name suggests, we ask you to implement the naive algorithm, that
 * is, check if each number is prime(by iterating over its possible factors).
 * We know there are algorithms better suited for this problem, but for now,
 * please implement this naive algorithm first.
 *
 * n <= 100.
 *
 *=======================================================================*/
func TrivialPrimeFinder(n int) []bool {
  isPrime := make([]bool, n+1)
  isPrime[0] = false
  isPrime[1] = false
  if n == 1 {
    return isPrime
  }
  for i := 2; i < n + 1; i++ {
    isPrime[i] = true
    for j := 2; j < i; j++ {
      if isPrime[j] == true {
        if i % j == 0 {
          isPrime[i] = false
          break
        }
      }
    }
  }
  return isPrime
}

/*==============================================================================
 * 8. Sieve Of Eratosthenes
 *
 * Write a function SieveOfEratosthenes that have exactly the same input and output
 * as last problem(TrivialPrimeFinder).
 *
 * As the name suggests, we ask you to implement the Sieve algorithm, not to
 * decide if each number is prime independently. Your program should be able
 * to finish running within seconds for n=5*10^7(but we won't test on it).
 *
 * Think about how you would validate your solutions without calculating some
 * small outputs by hand, given you already implemented TrivialPrimeFinder before.
 *
 * n <= 10000.
 *
 *============================================================================*/
func SieveOfEratosthenes(n int) []bool {
  isPrime := make([]bool, n+1)
  for m:= 0; m <= n; m++ {
    isPrime[m] = true
    }//initialize n+1 elements to be true
  isPrime[0] = false
  isPrime[1] = false
  if n == 1 {
    return isPrime
  }
  for i := 2; i <= n; i++ {
    if isPrime[i] == true {
      for j := 2; j*i <= n; j ++ {
        isPrime[i*j] = false
      }
    }
  }
  return isPrime
}

/*==============================================================================
 * 9. List of Primes
 *
 * Write a function ListOfPrimes that takes a positive integer n, and return a slice
 * containing all primes that are no larger than n, in ascending order.
 * Try to utilize one of the subroutines you wrote.
 *
 * n <= 50.
 *
 *============================================================================*/
/*Psudocode:
ListOfPrimes(n)
  primes //array that stores prime numbers
  isPrime <- SieveOfEratosthenes(n)
  for i from 1 to n
    if isPrime[i] == true
      append i to primes
  return primes
*/

func ListOfPrimes(n int) []int {
  primes := make([]int, n)
  isPrime := SieveOfEratosthenes(n)
  i := 0
  for  j:= 1; j <= n; j++ {
    if isPrime[j] == true {
      primes[i] = j
      i++
    }
  }
  return primes[:i]
}

/*==============================================================================
 * 10. Perfect Numbers I
 *
 * Write a function IsPerfect that takes a positive integer n and return a boolean
 * value indicating if n is a perfect number.
 *
 * A positive integer is a perfect number if and only if sum of its proper positive
 * divisors sum up to itself. For example, proper positive divisors of 6 are 1, 2
 * and 3, and since 6=1+2+3, 6 is a perfect number.
 *
 * There are very few perfect numbers within range of int, but note that your code
 * will be checked by TAs for coding style, so don't cheat!
 *
 * n <= 10000.
 *
 *============================================================================*/
/*
Pseudocode:
IsPerfect(n)
  sum <- 1
  for i from 2 to sqrt(n)
    if n % i == 0
    sum = 1 + i + n/i
  //consider if n = k^2
  if sum == n:
    return true
  else
    return false
*/
func IsPerfect(n int) bool {
  sum := 1
  middle := int(math.Sqrt(float64(n)))
  for i := 2; float64(i) < math.Sqrt(float64(n)); i++ {
    if n % i == 0 {
      sum = sum + i + n/i
    }
  }
  if n % middle == 0 && n / middle == middle {
    sum = sum + middle
  }
  if sum == n {
    return true
  }
  return false
}
/*==============================================================================
 * 11. Perfect Numbers II
 *
 * Write a function NextPerfectNumber that takes a positive integer n and return
 * a int that is the smallest perfect number larger than n.
 *
 * We guarantee the correct answer can be stored in an int. Please use the function
 * you wrote for Problem 8 as a subroutine.
 *
 * n <= 8000. Note that there is a perfect number between 8000 and 10000.
 *
 *============================================================================*/
/* Pseudocode:
NextPerfectNumber(n)
  i <- n + 1
  while IsPerfect(i) != true
    i++

  return i
*/

func NextPerfectNumber(n int) int {
  next := n
  for i := n + 1; IsPerfect(i) != true; i++ {
    next++
  }
  next++
  return next
}
// The next perfect number is 33,550,336 = 2^12(2^13 - 1)
//func NextTwinPrimers(n int) int int {}

/*==============================================================================
* 12. Hailstone Sequence
*
* The Hailstone function is defined as follows. h(n) equals n/2 if n is an even
* number, and n*3+1 if n is odd. The Hailstone sequence of an positive integer
* n is defined as the infinite sequence containing:
*
*    n, h(n), h(h(n)), h(h(h(n))), ...
*
* The famous Collatz Conjecture says for any n, the sequence reaches 1.
* Write a function HailstoneSequence that takes n as input and output the sequence
* up until 1 appears in the sequence.
* For example, HailstoneSequence(7) = [7, 22, 11, 34, 17, 52, 26, 13, 40, 20, 10, 5,
										16, 8, 4, 2, 1]
*
* n <= 1000. Note that the values in the sequence can go beyond n; We assure
* those are all under 10^6.
*
*============================================================================*/
/* Pseudocode:
NextHailstone (n)  // calculate the next value
  if n % 2 == 0
    return n%2
  else
    return 3n + 1

HailstoneSequence (n)
  append n to the array sequence
  while i = NextHailstone(i) not equal to 1
    append i to sequence
  append 1 to sequence
return Sequence
*/
func NextHailstone(n int) int {
  if n % 2 == 0 {
    return n / 2
  }
  return 3 * n + 1
}

func HailstoneSequence(n int) []int {
  var sequence []int
  for i := n; i != 1; i = NextHailstone(i) {
    sequence = append(sequence, i)
  }
  sequence = append(sequence, 1)
  return sequence
}

/* Pseudocode:
IsPrime(n)
  for i from 2 to sqrt(n)
    if n % i == 0
      return false
  return true

NextTwinPrimers(n)
  i <- n + 1
  while IsPrime(i) == false or IsPrime(i + 2) == false:
    i ++

  return i, i+2
*/

func IsPrime(n int) bool {
  if n == 1 {
    return false
  }
  if n == 2 {
    return true
  }
  for i := 2; i <= int(math.Sqrt(float64(n))); i++ {
    if n % i == 0 {
      return false
    }
  }
  return true
}


func NextTwinPrimers(n int) (int, int) {
  twinprime := n + 1
  for i := n + 1; IsPrime(i) == false || IsPrime(i + 2) == false; i++ {
    twinprime = i + 1
  }
  return twinprime , twinprime + 2
}
