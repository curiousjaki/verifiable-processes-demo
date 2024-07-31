#![no_main]
// If you want to try std support, also update the guest Cargo.toml file
#![no_std]  // std support is experimental

use risc0_zkvm::guest::env;

risc0_zkvm::guest::entry!(main);


fn main() {
    // TODO: Implement your guest code here

    // read the input
    let consumption: f64 = env::read();
    let emission_factor: f64 = env::read();

    // TODO: do something with the input
    let co2_emissions: f64 = consumption * emission_factor;
    // write public output to the journal
    env::commit(&co2_emissions);
    //env::commit(&emission_factor);
}
