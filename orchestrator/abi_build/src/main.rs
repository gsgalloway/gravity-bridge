use ethers::contract::Abigen;
use std::process;

fn main() {
    let abigen = match Abigen::new("Gravity", "../abi/Gravity.json") {
        Ok(abigen) => abigen,
        Err(e) => {
            println!("Could not open Gravity.json: {}", e);
            process::exit(1);
        }
    };

    let abi = match abigen
        .add_event_derive("serde::Deserialize")
        .add_event_derive("serde::Serialize")
        .generate()
    {
        Ok(abi) => abi,
        Err(e) => {
            println!("Could not generate abi from Gravity.json: {}", e);
            process::exit(1);
        }
    };

    match abi.write_to_file("../gravity_utils/src/gravity.rs") {
        Ok(_) => (),
        Err(e) => println!("Error writing gravity.rs: {}", e),
    }
}