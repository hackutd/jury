use gavel3::util::crowd_bt::divergence_gaussian;

fn main() {
    println!("Hello, world!");
    let x = divergence_gaussian(24.3, 0.54, 31.2, 0.65);
    println!("Result of gaussian divergence is {x}");
}
