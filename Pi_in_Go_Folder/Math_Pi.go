package main

import (
	"fmt"
	"math"
)

func main() {
	// Prints Pi with standard formatting
	fmt.Println("Standard Pi:", math.Pi) 
	
	// Prints Pi with high decimal precision
	fmt.Printf("High Precision Pi: %.48f\n", math.Pi) 
}
