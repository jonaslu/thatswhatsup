use std::{io::Read, net::TcpListener};

macro_rules! unwrap_or_exit {
    ($e: expr, $s: expr) => {
        match $e {
            Ok(value) => value,
            Err(error) => {
                eprintln!("{}: {error:?}. exiting", $s);
                std::process::exit(1);
            }
        }
    };
}

fn main() {
    let conn = unwrap_or_exit!(
        TcpListener::bind("0.0.0.0:3000"),
        "Error binding on port 3000"
    );

    let (stream, _) = unwrap_or_exit!(conn.accept(), "Could not accept socket");

    for byte in stream.bytes() {
        let bytez = unwrap_or_exit!(byte, "Could not unwrap byte");
        print!("{}", bytez as char);
    }

    println!("Other side hung up");
}
