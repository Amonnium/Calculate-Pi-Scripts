fn calculate_pi(iterations: u64) -> f64 {
    let mut pi = 0.0;

    for n in 0..iterations {
        let sign = if n % 2 == 0 { 1.0 } else { -1.0 };

        pi += sign / (2.0 * n as f64 + 1.0);
    }

    4.0 * pi
}

fn main() {
    let result = calculate_pi(100_000_000);

    println!("Pi approximation: {:.15}", result);
}
