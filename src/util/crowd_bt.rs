use rgsl::psi::diagamma::psi;
use statrs::function::beta;

pub const GAMMA: f64 = 0.1;
// const LAMBDA: f64 = 1.0; // regularization parameter - not used??
const KAPPA: f64 = 0.0001;
pub const MU_PRIOR: f64 = 0.0;
pub const SIGMA_SQ_PRIOR: f64 = 1.0;
pub const ALPHA_PRIOR: f64 = 10.0;
pub const BETA_PRIOR: f64 = 1.0;
pub const EPSILON: f64 = 0.25;

pub fn update(
    alpha: f64,
    beta: f64,
    mu_winner: f64,
    sigma_sq_winner: f64,
    mu_loser: f64,
    sigma_sq_loser: f64,
) -> (f64, f64, f64, f64, f64, f64) {
    let (updated_alpha, updated_beta, _) = update_annotator(
        alpha,
        beta,
        mu_winner,
        sigma_sq_winner,
        mu_loser,
        sigma_sq_loser,
    );
    let (updated_mu_winner, updated_mu_loser) = update_mus(
        alpha,
        beta,
        mu_winner,
        sigma_sq_winner,
        mu_loser,
        sigma_sq_loser,
    );
    let (updated_sigma_sq_winner, updated_sigma_sq_loser) = update_sigma_sqs(
        alpha,
        beta,
        mu_winner,
        sigma_sq_winner,
        mu_loser,
        sigma_sq_loser,
    );
    return (
        updated_alpha,
        updated_beta,
        updated_mu_winner,
        updated_sigma_sq_winner,
        updated_mu_loser,
        updated_sigma_sq_loser,
    );
}

pub fn expected_information_gain(
    alpha: f64,
    beta: f64,
    mu_a: f64,
    sigma_sq_a: f64,
    mu_b: f64,
    sigma_sq_b: f64,
) -> f64 {
    let (alpha1, beta1, c) = update_annotator(alpha, beta, mu_a, sigma_sq_a, mu_b, sigma_sq_b);
    let (mu_a1, mu_b1) = update_mus(alpha, beta, mu_a, sigma_sq_a, mu_b, sigma_sq_b);
    let (sigma_sq_a1, sigma_sq_b1) =
        update_sigma_sqs(alpha, beta, mu_a, sigma_sq_a, mu_b, sigma_sq_b);
    let prob_a_ranked_above = c;

    let (alpha2, beta2, _) = update_annotator(alpha, beta, mu_b, sigma_sq_b, mu_a, sigma_sq_a);
    let (mu_a2, mu_b2) = update_mus(alpha, beta, mu_b, sigma_sq_b, mu_a, sigma_sq_a);
    let (sigma_sq_a2, sigma_sq_b2) =
        update_sigma_sqs(alpha, beta, mu_b, sigma_sq_b, mu_a, sigma_sq_a);

    return prob_a_ranked_above
        * (divergence_gaussian(mu_a1, sigma_sq_a1, mu_a, sigma_sq_a)
            + divergence_gaussian(mu_b1, sigma_sq_b1, mu_b, sigma_sq_b)
            + GAMMA * divergence_beta(alpha1, beta1, alpha, beta))
        + (1.0 - prob_a_ranked_above)
            * (divergence_gaussian(mu_a2, sigma_sq_a2, mu_a, sigma_sq_a)
                + divergence_gaussian(mu_b2, sigma_sq_b2, mu_b, sigma_sq_b)
                + GAMMA * divergence_beta(alpha2, beta2, alpha, beta));
}

// https://en.wikipedia.org/wiki/Normal_distribution (See Kullback-Leibler divergence)
fn divergence_gaussian(mu1: f64, sigma_sq1: f64, mu2: f64, sigma_sq2: f64) -> f64 {
    let sigma_ratio = sigma_sq1 / sigma_sq2;
    let left_term = f64::powi(mu1 - mu2, 2) / (2.0 * sigma_sq2);
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

fn update_annotator(
    alpha: f64,
    beta: f64,
    mu_winner: f64,
    sigma_sq_winner: f64,
    mu_loser: f64,
    sigma_sq_loser: f64,
) -> (f64, f64, f64) {
    let e_mu_winner = f64::exp(mu_winner);
    let e_mu_loser = f64::exp(mu_loser);
    let c1 = e_mu_winner / (e_mu_winner + e_mu_loser)
        + 0.5
            * (sigma_sq_winner + sigma_sq_loser)
            * (e_mu_winner * e_mu_loser * (e_mu_loser - e_mu_winner))
            / f64::powi(e_mu_winner + e_mu_loser, 3);
    let c2 = 1.0 - c1;
    let c = (c1 * alpha + c2 * beta) / (alpha + beta);

    let expected = (c1 * (alpha + 1.0) * alpha + c2 * alpha * beta)
        / (c * (alpha + beta + 1.0) * (alpha + beta));
    let expected_sq = (c1 * (alpha + 2.0) * (alpha + 1.0) * alpha
        + c2 * (alpha + 1.0) * alpha * beta)
        / (c * (alpha + beta + 2.0) * (alpha + beta + 1.0) * (alpha + beta));

    let variance = expected_sq - f64::powi(expected, 2);

    let updated_alpha = (expected - expected_sq) * expected / variance;
    let updated_beta = (expected - expected_sq) * (1.0 - expected) / variance;

    return (updated_alpha, updated_beta, c);
}

fn update_mus(
    alpha: f64,
    beta: f64,
    mu_winner: f64,
    sigma_sq_winner: f64,
    mu_loser: f64,
    sigma_sq_loser: f64,
) -> (f64, f64) {
    let e_mu_winner = f64::exp(mu_winner);
    let e_mu_loser = f64::exp(mu_loser);
    let mult = (alpha * e_mu_winner) / (alpha * e_mu_winner + beta * e_mu_loser)
        - e_mu_winner / (e_mu_winner + e_mu_loser);

    let updated_mu_winner = mu_winner + sigma_sq_winner * mult;
    let updated_mu_loser = mu_loser - sigma_sq_loser * mult;

    return (updated_mu_winner, updated_mu_loser);
}

fn update_sigma_sqs(
    alpha: f64,
    beta: f64,
    mu_winner: f64,
    sigma_sq_winner: f64,
    mu_loser: f64,
    sigma_sq_loser: f64,
) -> (f64, f64) {
    let e_mu_winner = f64::exp(mu_winner);
    let e_mu_loser = f64::exp(mu_loser);
    let mult = (alpha * e_mu_winner * beta * e_mu_loser)
        / f64::powi(alpha * e_mu_winner + beta * e_mu_loser, 2)
        - (e_mu_winner * e_mu_loser / f64::powi(e_mu_winner + e_mu_loser, 2));

    let updated_sigma_sq_winner = sigma_sq_winner * f64::max(1.0 + sigma_sq_winner * mult, KAPPA);
    let updated_sigma_sq_loser = sigma_sq_loser * f64::max(1.0 + sigma_sq_loser * mult, KAPPA);

    return (updated_sigma_sq_winner, updated_sigma_sq_loser);
}
