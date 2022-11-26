use statrs::function::beta;
use rgsl::psi::diagamma::psi;

const GAMMA: f64 = 0.1;
const LAMBDA: f64 = 1.0;
const KAPPA: f64 = 0.0001;
const MU_PRIOR: f64 = 0.0;
const SIGMA_SQ_PRIOR: f64 = 1.0;
const ALPHA_PRIOR: f64 = 10.0;
const BETA_PRIOR: f64 = 1.0;
const EPSILON: f64 = 0.25;

// https://en.wikipedia.org/wiki/Normal_distribution (See Kullback-Leibler divergence)
// TODO: Remove pub; just here for testing rn
pub fn divergence_gaussian(mu_1: f64, sigma_sq_1: f64, mu_2: f64, sigma_sq_2: f64) -> f64 {
    let sigma_ratio = sigma_sq_1 / sigma_sq_2;
    let left_term = f64::powi(mu_1 - mu_2, 2) / (2.0 * sigma_sq_2);
    let right_term = (sigma_ratio - 1.0 - f64::ln(sigma_ratio)) / 2.0;
    return left_term + right_term;
}

// https://en.wikipedia.org/wiki/Beta_distribution (Also Kullback-Leibler divergence)
fn divergence_beta(alpha_1: f64, beta_1: f64, alpha_2: f64, beta_2: f64) -> f64 {
    let ln_term = beta::ln_beta(alpha_2, beta_2) - beta::ln_beta(alpha_1, beta_1);
    let a_term = (alpha_1 - alpha_2) * psi(alpha_1);
    let b_term = (beta_1 - beta_2) * psi(beta_1);
    let ab_term = (alpha_2 - alpha_1 + beta_2 - beta_1) * psi(alpha_1 + beta_1);
    return ln_term + a_term + b_term + ab_term;
}

// TODO: Finish