function calculatePiBBP(iterations) {
    let pi = 0;
    for (let n = 0; n < iterations; n++) {
        let scale = 1 / Math.pow(16, n);
        let inner = (4 / (8 * n + 1)) - 
                    (2 / (8 * n + 4)) - 
                    (1 / (8 * n + 5)) - 
                    (1 / (8 * n + 6));
        pi += scale * inner;
    }
    return pi;
}

// Just 15 iterations provide maximum JavaScript float accuracy
console.log(`Calculated Pi: ${calculatePiBBP(15)}`);
console.log(`Math.PI Match: ${Math.PI}`);

