mod structs;
mod loops;
mod pointer_ref;
mod types;
mod print;
mod vars;
fn main() {
    print::run();
    let name = "lxd";
    let mut age = 1;
    println!("I'm {1}, I'm {0} years old",age,name);
    println!("a: {a}, b:{b}",a = "A",b = "B");
    println!("Binary: {:b} Hex: {0:x} Octal: {0:o}",10);
    println!("10 + 10 = {}",10+10);
    vars::run();
    types::run();
    structs::run();
    pointer_ref::run();
    loops::run();
}
