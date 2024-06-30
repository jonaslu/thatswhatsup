#[derive(Debug)]
enum List {
    Cons(u32, Box<List>),
    Nil,
}

use List::{Cons, Nil};


impl List {
    fn new() -> Box<List> {
        Box::new(Nil)
    }

    fn prepend(self: Box<Self>, val: u32) -> Box<List> {
        Box::new(Cons(val, self))
    } 
}

fn main() {
    let yak = List::new().prepend(1).prepend(2);
    println!("Value: {:?}", yak);
}
