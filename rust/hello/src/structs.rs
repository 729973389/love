struct Person {
    Name: String,
    Age: u8,
}

impl Person{
fn GetPerson(s: &str,a: u8) -> Person{
    Person {
        Name: s.to_string(),
        Age:  a,
    }
    }
}

pub fn run () {
    let m = Person::GetPerson("lxd",24);
    println!("{:?}",(m.Name,m.Age));
}
