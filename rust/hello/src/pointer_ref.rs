pub fn run(){
    //primitives
    let i: i32 = 1;
    let i_copy = i;
    println!("i_cpoy: {}, i: {}",i_copy,i);
    //no-primitives
    let v = vec![1,2,3,4];
//    let v_copy = v;
//    println!("v_cpoy: {:?} v: {:?}",v_copy,v);
    /*
     * 7 |     let v = vec![1,2,3,4];
  |         - move occurs because `v` has type `std::vec::Vec<i32>`, which does not implement the `Copy` trait
8 |     let v_copy = v;
  |                  - value moved here
9 |     println!("v_cpoy: {:?} v: {:?}",v_copy,v);
  |                                            ^ value borrowed here after move

     */

    let v_copy = &v;
    println!("v_copy: {:?} v: {:?}",v_copy,v);
}
