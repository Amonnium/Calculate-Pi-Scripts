import kotlin.math.PI

fun main() {
    // 1. Using the built-in constant (Recommended)
    val radius = 5.0
    val circumference = 2 * PI * radius
    println("Circumference: $circumference")

    // 2. Computing Pi using Leibniz formula (Approximation)
    val computedPi = calculatePiLeibniz(1000000)
    println("Computed Pi: $computedPi")
}

fun calculatePiLeibniz(iterations: Int): Double {
    var pi = 0.0
    var denominator = 1
    for (index in 0 until iterations) {
        if (index % 2 == 0) {
            pi += 4.0 / denominator
        } else {
            pi -= 4.0 / denominator
        }
        denominator += 2
    }
    return pi
}
