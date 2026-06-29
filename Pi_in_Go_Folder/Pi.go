package main

import "fmt"

func calculatePiLeibniz(iterations int) float64 {
	var sum float64
	numerator := 1.0

	for i := 0; i < iterations; i++ {
		denominator := float64(2*i + 1)
		sum += numerator / denominator
		numerator = -numerator // Alternates the sign (+ and -)
	}

	return sum * 4
}

func main() {
	// A higher number of iterations yields greater accuracy
	iterations := 10000000 
	pi := calculatePiLeibniz(iterations)
	fmt.Printf("Leibniz Pi (%d iterations): %f\n", iterations, pi)
}
