package utils

import "math"

// Rasch Model
func RaschProbability(theta, beta float64) float64 {
	return math.Exp(theta-beta) / (1 + math.Exp(theta-beta))
}

// Update theta (simple gradient)
func UpdateTheta(theta, probability float64, correct bool) float64 {
	u := 0.0
	if correct {
		u = 1.0
	}
	learningRate := 0.3
	return theta + learningRate*(u-probability)
}
