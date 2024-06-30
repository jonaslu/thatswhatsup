use std::process::{self};

fn main() {
    let mut rl = rustyline::DefaultEditor::new().expect("Could not create rustyline editor");

    loop {
        let input = rl.readline("$ > ");

        match input {
            Ok(line) => {
                let input_line = line.trim().to_string();

                let error_msg = format!("Could not run command: {}", input_line);
                let mut command = process::Command::new(input_line).spawn().expect(&error_msg);
                command.wait().expect("Couldn't wait for command");
            }
            Err(_) => process::exit(0)
        }
    }
}
