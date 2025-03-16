package main

import "fmt"

type Number interface {
	int64 | float64
}

func main() {
	// Initialize a map for the integer values
	ints := map[string]int64{
		"first":  34,
		"second": 12,
	}

	// Initialize a map for the float values
	floats := map[string]float64{
		"first":  35.98,
		"second": 26.99,
	}

	fmt.Printf("Generic Sums: %v and %v\n",
		SumIntsOrFloats(ints),
		SumIntsOrFloats(floats))
}

func SumInts(m map[string]int64) int64 {
	var result int64

	for _, v := range m {
		result += v
	}

	return result
}

func SumFloats(m map[string]float64) float64 {
	var result float64

	for _, v := range m {
		result += v
	}

	return result
}

func SumIntsOrFloats[K comparable, V Number](m map[K]V) V {
	var result V

	for _, v := range m {
		result += v
	}

	return result
}
