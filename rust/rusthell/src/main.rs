// Read stuff an execute command

use std::{io::{self, Write}, process};

fn main() {
    let mut input = String::new();
    print!("$ > ");
    io::stdout().flush().expect("Could not flush stdout");
    io::stdin().read_line(&mut input).expect("Could not read stdin");
    
    input = input.trim().to_string();

    let error_msg = format!("Could not run command: {}", input);
    let mut command = process::Command::new(input).spawn().expect(&error_msg);
    command.wait().expect("Couldn't wait for command");
}
