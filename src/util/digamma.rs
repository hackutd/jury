/**
 * Bernardo, J. M. (1976). Algorithm AS 103: Psi (Digamma) Function. 
 * Journal of the Royal Statistical Society. 
 * Series C (Applied Statistics), 25(3), 315-317. 
 * https://doi.org/10.2307/2347257
 */

pub fn psi(x: f64) -> f64 {
    let euler_mascheroni = 0.57721566490153286060;
    let mut r: f64;
    let mut value: f64;
    let mut x2: f64;

    // Make sure input is greater than 0
    if x <= 0.0 {
        eprintln!("Error: Digamma psi function called with argument <= 0");
        return 0.0;
    }

    // Approximate for small values
    if x <= 0.000001 {
        return - euler_mascheroni - 1.0 / x + 1.6449340668482264365 * x;
    }

    // Reduce to DIGAMMA(X + N)
    value = 0.0;
    x2 = x;
    while x2 < 8.5 {
        value -= 1.0 / x2;
        x2 += 1.0;
    }

    // Use Stirling's (actually de Moivre's) expansion.
    r = 1.0 / x2;
    value += f64::ln(x2) - 0.5 * r;

    r *= r;
    value = value
        - r * (1.0 / 12.0
        - r * (1.0 / 120.0
        - r * (1.0 / 252.0
        - r * (1.0 / 240.0
        - r * (1.0 / 132.0)))));

    return value;
}
