fn is_odd(n: u32) -> bool {
    n % 2 == 0
}

fn main() {
    let upper: u32 = 1000;
    let goat: u32 = (0..).map(|n| n * n).take_while(|n| n < upper)
}
