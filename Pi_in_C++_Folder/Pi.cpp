#include <iostream>
#include <iomanip>

double calculate_pi(long iterations) {
    double pi = 0.0;

    for (long n = 0; n < iterations; n++) {
        double sign = (n % 2 == 0) ? 1.0 : -1.0;
        pi += sign / (2.0 * n + 1.0);
    }

    return 4.0 * pi;
}

int main() {
    double result = calculate_pi(100000000);

    std::cout << std::setprecision(15)
              << "Pi approximation: "
              << result
              << std::endl;

    return 0;
}
