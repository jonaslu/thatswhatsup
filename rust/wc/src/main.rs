use std::{env, fs, io::Read};
// use std::fs::File;

// fn print_and_exit<E: std::fmt::Display>(msg: &str, e: &E) -> ! {
//     eprintln!("Error {}: {}", msg, e);
//     std::process::exit(1);
// }

fn get_or_exit<T, E: std::fmt::Display>(msg: &str, r: Result<T, E>) -> T {
    match r {
        Err(e) => {
            eprintln!("{}, error: {}", msg, e);
            std::process::exit(1);
        }
        Ok(v) => v,
    }
}

enum CountMode {
    Bytes,
    Words,
    Lines,
}

struct Config {
    file_name: String,

    count_mode: CountMode,
}

impl Config {
    fn new(args: &[String]) -> Config {
        let file_name = args[1].to_string();
        let mode_str = args[2].as_str();

        let count_mode: CountMode;

        match mode_str {
            "-l" => count_mode = CountMode::Lines,
            "-w" => count_mode = CountMode::Words,
            "-b" => count_mode = CountMode::Bytes,
            _ => {
                eprintln!("Mode not recognized: {}", mode_str);
                std::process::exit(1);
            }
        }

        Config {
            file_name,
            count_mode,
        }
    }
}

fn main() {
    let args: Vec<String> = env::args().collect();

    let config = Config::new(&args);

    let mut file_fd = get_or_exit(
        &format!("Failed to open file {}", config.file_name),
        fs::File::open(config.file_name),
    );

    let mut file_contents = String::new();
    let _ = get_or_exit(
        "Couldn't read file",
        file_fd.read_to_string(&mut file_contents),
    );

    match config.count_mode {
        CountMode::Bytes => {
            println!("{}", file_contents.len())
        }
        CountMode::Words => {
            println!("{}", file_contents.split_whitespace().count())
        }
        CountMode::Lines => {
            println!("{}", file_contents.split("\n").count())
        }
    }
}
