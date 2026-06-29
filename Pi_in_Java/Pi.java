import java.math.BigDecimal;
import java.math.RoundingMode;

public class Pi {

    // Define the number of decimal places required
    private static final int SCALE = 100;
    
    // Setting extra precision during internal calculations to prevent rounding errors
    private static final int INTERNAL_SCALE = SCALE + 5;

    public static void main(String[] args) {
        System.out.println("Calculating Pi with " + SCALE + " decimal places...");
        
        long startTime = System.currentTimeMillis();
        BigDecimal pi = calculatePi();
        long endTime = System.currentTimeMillis();

        // Rounding to the final requested scale
        BigDecimal finalPi = pi.setScale(SCALE, RoundingMode.HALF_UP);

        System.out.println("\nCalculated Pi:\n" + finalPi);
        System.out.println("\nExecution time: " + (endTime - startTime) + " ms");
    }

    /**
     * Calculates Pi using Machin's Formula: Pi/4 = 4*arctan(1/5) - arctan(1/239)
     */
    public static BigDecimal calculatePi() {
        BigDecimal term1 = BigDecimal.valueOf(1).divide(BigDecimal.valueOf(5), INTERNAL_SCALE, RoundingMode.HALF_UP);
        BigDecimal term2 = BigDecimal.valueOf(1).divide(BigDecimal.valueOf(239), INTERNAL_SCALE, RoundingMode.HALF_UP);

        BigDecimal arctan5 = calculateArcTan(term1);
        BigDecimal arctan239 = calculateArcTan(term2);

        // 4 * (4 * arctan(1/5) - arctan(1/239))
        BigDecimal multipliedFirstTerm = arctan5.multiply(BigDecimal.valueOf(4));
        BigDecimal difference = multipliedFirstTerm.subtract(arctan239);
        
        return difference.multiply(BigDecimal.valueOf(4));
    }

    /**
     * Calculates arctan(x) using Taylor series expansion: x - x^3/3 + x^5/5 - x^7/7...
     */
    private static BigDecimal calculateArcTan(BigDecimal inverseValue) {
        BigDecimal result = inverseValue;
        BigDecimal currentPower = inverseValue;
        BigDecimal squareValue = inverseValue.multiply(inverseValue);
        
        int sign = -1;
        int denominator = 3;

        while (true) {
            // Next power of x: x^(n+2)
            currentPower = currentPower.multiply(squareValue).setScale(INTERNAL_SCALE, RoundingMode.HALF_UP);
            
            // Current term: x^n / n
            BigDecimal currentTerm = currentPower.divide(BigDecimal.valueOf(denominator), INTERNAL_SCALE, RoundingMode.HALF_UP);

            // If the current term does not affect the required precision, stop the loop
            if (currentTerm.compareTo(BigDecimal.ZERO) == 0) {
                break;
            }

            if (sign == 1) {
                result = result.add(currentTerm);
            } else {
                result = result.subtract(currentTerm);
            }

            // Prepare next iteration variables
            sign = -sign;
            denominator += 2;
        }

        return result;
    }
}