pub fn run() {
    let mut count = 0;
    
    //Infinite Loop
    loop {
        count += 1;
        println!("{}",count);
        if count == 6 {
            break;
        }
    }

    //While Loop
    while count < 7 {
        if count == 4{
            println!("Woo");
        } else {
            println!("Moo");
        }
        count += 1;
    }

    //for range 
    for i in 1..10 {
        println!("{}",i);
    }
}
