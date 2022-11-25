// https://en.wikipedia.org/wiki/Normal_distribution (See Kullback-Leibler divergence)
// TODO: Remove pub; just here for testing rn
pub fn divergence_gaussian(mu_1: f64, sigma_sq_1: f64, mu_2: f64, sigma_sq_2: f64) -> f64 {
    let ratio = sigma_sq_1 / sigma_sq_2;
    return f64::powi(mu_1 - mu_2, 2) / (2.0 * sigma_sq_2) + (ratio - 1.0 - f64::ln(ratio) / 2.0);
}

// TODO: Finish