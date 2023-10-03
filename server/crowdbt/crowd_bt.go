package crowdbt

import (
	"math"

	"gonum.org/v1/gonum/mathext"
)

const GAMMA = 0.1
const KAPPA = 0.0001
const MU_PRIOR = 0.0
const SIGMA_SQ_PRIOR = 1.0
const ALPHA_PRIOR = 10.0
const BETA_PRIOR = 1.0
const EPSILON = 0.25
const MIN_VIEWS = 4

func Update(alpha float64, beta float64, muWinner float64, sigmaSqWinner float64, muLoser float64, sigmaSqLoser float64) (float64, float64, float64, float64, float64, float64) {
	updatedAlpha, updatedBeta, _ := UpdateAnnotator(alpha, beta, muWinner, sigmaSqWinner, muLoser, sigmaSqLoser)
	updatedMuWinner, updatedMuLoser := UpdateMus(alpha, beta, muWinner, sigmaSqWinner, muLoser, sigmaSqLoser)
	updatedSigmaSqWinner, updatedSigmaSqLoser := UpdateSigmaSqs(alpha, beta, muWinner, sigmaSqWinner, muLoser, sigmaSqLoser)
	return updatedAlpha, updatedBeta, updatedMuWinner, updatedSigmaSqWinner, updatedMuLoser, updatedSigmaSqLoser
}

func ExpectedInformationGain(alpha float64, beta float64, muA float64, sigmaSqA float64, muB float64, sigmaSqB float64) float64 {
	alpha1, beta1, c := UpdateAnnotator(alpha, beta, muA, sigmaSqA, muB, sigmaSqB)
	muA1, muB1 := UpdateMus(alpha, beta, muA, sigmaSqA, muB, sigmaSqB)
	sigmaSqA1, sigmaSqB1 := UpdateSigmaSqs(alpha, beta, muA, sigmaSqA, muB, sigmaSqB)
	probARankedAbove := c

	alpha2, beta2, _ := UpdateAnnotator(alpha, beta, muB, sigmaSqB, muA, sigmaSqA)
	muB2, muA2 := UpdateMus(alpha, beta, muB, sigmaSqB, muA, sigmaSqA)
	sigmaSqB2, sigmaSqA2 := UpdateSigmaSqs(alpha, beta, muB, sigmaSqB, muA, sigmaSqA)

	return probARankedAbove*(DivergenceGaussian(muA1, sigmaSqA1, muA, sigmaSqA)+DivergenceGaussian(muB1, sigmaSqB1, muB, sigmaSqB)+GAMMA*DivergenceBeta(alpha1, beta1, alpha, beta)) + (1.0-probARankedAbove)*DivergenceGaussian(muA2, sigmaSqA2, muA, sigmaSqA) + (DivergenceGaussian(muB2, sigmaSqB2, muB, sigmaSqB) + GAMMA*DivergenceBeta(alpha2, beta2, alpha, beta))
}

// https://en.wikipedia.org/wiki/Normal_distribution (See Kullback-Leibler divergence)
func DivergenceGaussian(mu1 float64, sigmaSq1 float64, mu2 float64, sigmaSq2 float64) float64 {
	sigmaRatio := sigmaSq1 / sigmaSq2
	leftTerm := math.Pow(mu1-mu2, 2) / (2.0 * sigmaSq2)
	rightTerm := (sigmaRatio - 1.0 - math.Log(sigmaRatio)) / 2.0
	return leftTerm + rightTerm
}

// https://en.wikipedia.org/wiki/Beta_distribution (Also Kullback-Leibler divergence)
func DivergenceBeta(alpha1 float64, beta1 float64, alpha2 float64, beta2 float64) float64 {
	lnTerm := mathext.Lbeta(alpha2, beta2) - mathext.Lbeta(alpha1, beta1)
	aTerm := (alpha1 - alpha2) * mathext.Digamma(alpha1)
	bTerm := (beta1 - beta2) * mathext.Digamma(beta1)
	abTerm := (alpha2 - alpha1 + beta2 - beta1) * mathext.Digamma(alpha1+beta1)
	return lnTerm + aTerm + bTerm + abTerm
}

func UpdateAnnotator(alpha float64, beta float64, muWinner float64, sigmaSqWinner float64, muLoser float64, sigmaSqLoser float64) (float64, float64, float64) {
	eMuWinner := math.Exp(muWinner)
	eMuLoser := math.Exp(muLoser)
	c1 := eMuWinner/(eMuWinner+eMuLoser) + 0.5*(sigmaSqWinner+sigmaSqLoser)*(eMuWinner*eMuLoser*(eMuLoser-eMuWinner))/math.Pow(eMuWinner+eMuLoser, 3)
	c2 := 1.0 - c1
	c := (c1*alpha + c2*beta) / (alpha + beta)

	expected := (c1*(alpha+1.0)*alpha + c2*alpha*beta) / (c * (alpha + beta + 1.0) * (alpha + beta))
	expectedSq := (c1*(alpha+2.0)*(alpha+1.0)*alpha + c2*(alpha+1.0)*alpha*beta) / (c * (alpha + beta + 2.0) * (alpha + beta + 1.0) * (alpha + beta))
	variance := expectedSq - math.Pow(expected, 2)

	updatedAlpha := (expected - expectedSq) * expected / variance
	updatedBeta := (expected - expectedSq) * (1.0 - expected) / variance

	return updatedAlpha, updatedBeta, c
}

func UpdateMus(alpha float64, beta float64, muWinner float64, sigmaSqWinner float64, muLoser float64, sigmaSqLoser float64) (float64, float64) {
	eMuWinner := math.Exp(muWinner)
	eMuLoser := math.Exp(muLoser)
	mult := (alpha*eMuWinner)/(alpha*eMuWinner+beta*eMuLoser) - eMuWinner/(eMuWinner+eMuLoser)

	updatedMuWinner := muWinner + mult*sigmaSqWinner
	updatedMuLoser := muLoser - mult*sigmaSqLoser

	return updatedMuWinner, updatedMuLoser
}

func UpdateSigmaSqs(alpha float64, beta float64, muWinner float64, sigmaSqWinner float64, muLoser float64, sigmaSqLoser float64) (float64, float64) {
	eMuWinner := math.Exp(muWinner)
	eMuLoser := math.Exp(muLoser)
	mult := (alpha*eMuWinner*beta*eMuLoser)/math.Pow(alpha*eMuWinner+beta*eMuLoser, 2) - eMuWinner*eMuLoser/math.Pow(eMuWinner+eMuLoser, 2)

	updatedSigmaSqWinner := sigmaSqWinner * math.Max(1.0+mult*sigmaSqWinner, KAPPA)
	updatedSigmaSqLoser := sigmaSqLoser * math.Max(1.0+mult*sigmaSqLoser, KAPPA)

	return updatedSigmaSqWinner, updatedSigmaSqLoser
}
