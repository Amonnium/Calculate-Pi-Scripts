#include <stdio.h>

double calculate_pi(int iterations) {
    double pi = 0.0;

    for (int n = 0; n < iterations; n++) {
        double sign = (n % 2 == 0) ? 1.0 : -1.0;
        pi += sign / (2.0 * n + 1.0);
    }

    return 4.0 * pi;
}

int main() {
    double result = calculate_pi(100000000);

    printf("Pi approximation: %.15f\n", result);

    return 0;
}
