macro_rules! say_hello {
    () => {
        println!("hello")
    };
}

macro_rules! print_debug {
    (debug, $msg:expr) => {
        println!("DEBUG: {}", $msg)
    };
    ($msg:expr) => {
        println!($msg);
    }
}

fn main() {
    say_hello!();
    print_debug!(debug, "frakt");
    print_debug!("fecal");

    let fruits = vec!("abc", "def");
    for fruit in fruits {
        println!("{}", fruit)
    }

    let vet3 = vec![1,2,3];
    for vet in vet3 {
        println!("{}", vet);
    }
}
