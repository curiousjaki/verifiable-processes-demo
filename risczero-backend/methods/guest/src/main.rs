#![no_main]
// If you want to try std support, also update the guest Cargo.toml file
#![no_std]  // std support is experimental

use risc0_zkvm::guest::env;
risc0_zkvm::guest::entry!(main);


fn main() {
    // TODO: Implement your guest code here

    // read the input
    let variable_a: f64 = env::read();
    let variable_b: f64 = env::read();
    let operation: String = env::read();
    let result = match operation.as_str() {
        "add" => {
            variable_a + variable_b;
        }
        "sub" => {
            variable_a - variable_b;
        }
        "mul" => {
            variable_a * variable_b;
        }
        "div" => {
            variable_a / variable_b;
        }
        _ => {
            0.0;
        }
    };
    // write public output to the journal
    env::commit(&result);
    //env::commit(&emission_factor);
}
