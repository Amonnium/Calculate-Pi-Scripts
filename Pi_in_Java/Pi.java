public class Pi {
    
    /**
     * Calculate Pi using the Machin formula: pi/4 = 4*arctan(1/5) - arctan(1/239)
     * This formula converges much faster than the Leibniz series
     * @param iterations Number of iterations for the series
     * @return Approximation of Pi
     */
    public static double calculatePiMachin(int iterations) {
        double pi = 0.0;
        
        for (int n = 0; n < iterations; n++) {
            // First term: 4 * arctan(1/5)
            double term1 = Math.pow(1.0 / 5.0, 2 * n + 1) / (2 * n + 1);
            
            // Second term: arctan(1/239)
            double term2 = Math.pow(1.0 / 239.0, 2 * n + 1) / (2 * n + 1);
            
            // Combine with appropriate signs (alternating series)
            double sign = (n % 2 == 0) ? 1.0 : -1.0;
            
            pi += sign * (4 * term1 - term2);
        }
        
        return 4 * pi;
    }
    
    /**
     * Calculate Pi using the Leibniz series (Gregory-Leibniz series)
     * pi/4 = 1 - 1/3 + 1/5 - 1/7 + ...
     * Simpler but converges slower than Machin formula
     * @param iterations Number of iterations for the series
     * @return Approximation of Pi
     */
    public static double calculatePiLeibniz(int iterations) {
        double pi = 0.0;
        
        for (int n = 0; n < iterations; n++) {
            double sign = (n % 2 == 0) ? 1.0 : -1.0;
            pi += sign / (2.0 * n + 1.0);
        }
        
        return 4 * pi;
    }
    
    public static void main(String[] args) {
        // Display built-in constant
        System.out.printf("Built-in Math.PI: %.16f%n", Math.PI);
        System.out.println();
        
        // Calculate using Machin formula (fast convergence)
        double piMachin = calculatePiMachin(100);
        System.out.printf("Machin formula (100 iterations): %.16f%n", piMachin);
        
        // Calculate using Leibniz series (slower, educational)
        double piLeibniz = calculatePiLeibniz(1000000);
        System.out.printf("Leibniz series (1,000,000 iterations): %.16f%n", piLeibniz);
    }
}
