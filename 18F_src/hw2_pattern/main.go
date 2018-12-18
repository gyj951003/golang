package main
import (
  "fmt"
  //"strings"
)

func main() {
  fmt.Println(PatternCount("GCGCG", "GCG"))
  m := make(map[string]int)
  m["a"] = -4
  m["b"] = -5
  m["c"] = -3
  m["d"] = -1
  fmt.Println(MaxDict(m))
  fmt.Println(FrequencyMap("GCGCGGCGCG", 3))
}

func PatternCount(pattern, text string) int {
  count := 0
  for i := 0; i <= len(text) - len(pattern); i++ {
    if text[i : i + len(pattern)] == pattern {
      count++
    }
  }
    return count
}

func MaxDict(dict map[string]int) int {
  i := 1
  max := 0
  for k := range dict {
    if i == 1 {
      max = dict[k]
    } else {
        if dict[k] >= max {
          max = dict[k]
        }
    }
    i++
    }
    return max
}

func FrequencyMap(text string, k int) map[string]int {
  frequencyMap := make(map[string]int)
  for i := 0; i <= len(text) - k; i++ {
    key := text[i : i + k]
    count, ok := frequencyMap[key]
    if ok == true {
      frequencyMap[key] = count + 1
    } else {
      frequencyMap[key] = 1
    }
  }
  return frequencyMap
}

func FrequentWords(text string, k int) []string {
  frequencyMap := FrequencyMap(text, k)
  max := MaxDict(frequencyMap)
  frequentString := make([]string, 0)
  for k := range frequencyMap {
    if frequencyMap[k] == max {
      frequentString = append(frequentString, k)
    }
  }
  return frequentString
}

func Complement(c byte) string {
  switch c {
  case 'A':
    return "T"
  case 'T':
    return "A"
  case 'C':
    return "G"
  case 'G':
    return "C"
  }
  return ""
}

func ReverseComplement(pattern string) string {
  n := len(pattern)
  result := ""
  for i := 0; i < len(pattern); i++ {
      result += Complement(pattern[n - i - 1])
  }
  return result
}

func PatternMatching(pattern, text string) []int {
  array := make([]int, 0)
  for i := 0; i <= len(text) - len(pattern); i++ {
    if text[i : i + len(pattern)] == pattern {
      array = append(array, i)
    }
  }
  return array
}

func ClumpFinding(genome string, k, L, t int) []string {
  clump := make([]string, 0)
  if len(genome) < L {
    L = len(genome)
  }
  for i := 0; i <= len(genome) - L; i++ {
    frequencyMap := FrequencyMap(genome[i : i + L], k)
    for k := range frequencyMap{
      if frequencyMap[k] >= t {
        repeat := false
        if len(clump) != 0 {
          for j := 0; j < len(clump); j++ {
            if clump[j] == k {
              repeat = true
            }
          }
        }
        if repeat == false {
          clump = append(clump, k)
        }
      }
    }
  }
  return clump
}

func SkewArray(genome string) []int {
  skewArray := make([]int, 0)
  skewArray = append(skewArray, 0)
  for i := 0; i < len(genome); i++ {
    if genome[i] == 'C' {
      skewArray = append(skewArray, skewArray[i] - 1)
    } else if genome[i] == 'G' {
        skewArray = append(skewArray, skewArray[i] + 1)
    } else {
        skewArray = append(skewArray, skewArray[i])
    }
  }
  return skewArray
}

func MinimumSkew(genome string) []int {
  minSkewArray := make([]int, 0)
  skewArray := SkewArray(genome)
  min := 0
  for i := 0; i < len(skewArray); i++ {
    if skewArray[i] < min {
      min = skewArray[i]
    }
  }
  for i := 0; i < len(skewArray); i++ {
    if skewArray[i] == min {
      minSkewArray = append(minSkewArray, i)
    }
  }
  return minSkewArray
}
