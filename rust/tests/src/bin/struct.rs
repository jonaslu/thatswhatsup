#[derive(Debug)]
struct Rectangle {
    width: f64,
    height: f64
}

impl Rectangle {
    fn consume(&self) {
        println!("I now own this: {}, {}", self.width, self.height)
    }
}

fn main() {
    let rect = Rectangle{width: 6.0, height: 8.1};
    rect.consume();

    println!("I am the god of {}", rect.height);
}
