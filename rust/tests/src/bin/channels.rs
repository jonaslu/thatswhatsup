use std::sync::mpsc::{Sender, Receiver};
use std::sync::mpsc;

use std::thread;

static NTHREADS: u32 = 3;

fn main() {
    let (tx, rx): (Sender<u32>, Receiver<u32>) = mpsc::channel();

    let mut children = Vec::new();

    for i in 0..NTHREADS {
        let thread_tx = tx.clone();

        let child = thread::spawn(move || {
            thread_tx.send(i).unwrap();

            println!("thread {} finished", i);
        });

        children.push(child);
    }

    let mut ids = Vec::with_capacity(NTHREADS as usize);

    for _ in 0..NTHREADS {
        ids.push(rx.recv());
    }

    for child in children {
        child.join().unwrap();
    }

    println!("{:?}", ids);
}
