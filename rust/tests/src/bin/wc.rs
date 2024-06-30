use anyhow::{Context, Result};
use std::io::{self, Read};

fn main() -> Result<(),E> {
    let mut res = String::new();

    io::stdin().read_to_string(&mut res).con;

    Ok(())
}
