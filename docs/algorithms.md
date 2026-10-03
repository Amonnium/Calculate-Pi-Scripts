# How the Pi Examples Work

Pi is the ratio of a circle's circumference to its diameter. These programs demonstrate different ways to approximate it; they do not all calculate the same number of digits or use the same kind of arithmetic.

## Series and formulas

### Leibniz series

The C, C++, Go, Rust, Lua, and Kotlin examples add alternating fractions:

`4 * (1 - 1/3 + 1/5 - 1/7 + ...)`

It is easy to understand, but converges very slowly. The C/C++/Rust examples use 100 million terms and Go uses 10 million; a larger term count does not make this method efficient for many correct digits.

### Chudnovsky algorithm

The Python and Julia examples use the Chudnovsky series, which adds roughly 14 decimal digits of accuracy per term. Python uses the standard `decimal` module. Julia uses arbitrary-precision integers and binary splitting to combine terms efficiently. The Python version asks for a requested digit count; the Julia example prints a 1,000-digit result.

### Machin's formula

Java uses Machin's identity:

`pi = 16 * atan(1/5) - 4 * atan(1/239)`

It evaluates the arctangent series using `BigDecimal`, then rounds to 100 decimal places. The root-level C# example uses a scaled `BigInteger` representation of the same identity.

### Monte Carlo

Python's `Monte_Carlo.py` samples random points in a unit square. The fraction that falls inside a quarter-circle estimates `pi / 4`. The estimate is probabilistic: running it again can produce a different result, and accuracy improves only gradually as the number of samples grows.

### Nilakantha-style and related series

The Ruby example starts at 3 and adds/subtracts terms of the form `4 / (n * (n + 1) * (n + 2))`. The PowerShell example also alternates reciprocal triple-products. These converge faster than Leibniz, but still use ordinary floating-point numbers in these implementations.

### Bailey-Borwein-Plouffe (BBP)

The Raku and JavaScript examples use the BBP series, whose terms contain powers of 1/16 and rational factors. The Raku version uses exact rational values; JavaScript uses floating-point arithmetic and only a small number of terms because its floating-point precision is limited.

## Built-in constants and standard libraries

`Pi_in_Python/Math_Pi.py`, `Pi_in_JavaScript/Math_Pi.js`, `Pi_in_Go_Folder/Math_Pi.go`, `Pi_in_Rust_Folder/Math_Pi.rs`, `Pi_in_Nim_Folder/Math_Pi.nim`, and `Pi.zig` demonstrate a language or library's built-in Pi constant (or a direct `4 * atan(1)` identity). Perl's `Pi.pl` asks its `bignum` support library for 100 digits.

## Compare carefully

- A built-in constant usually gives about 15-16 significant decimal digits because it is a binary floating-point value.
- More iterations do not mean more correct digits once floating-point rounding dominates.
- Monte Carlo output varies from run to run.
- Runtime speed is not a fair accuracy comparison unless algorithms, precision, and workloads are comparable.
- The source comments and constants are part of the experiment; inspect each file to see its chosen iteration count or precision.
