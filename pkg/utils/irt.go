package utils

import "math"

// Rasch Model: P(correct) = exp(θ-β) / (1 + exp(θ-β))
func RaschProbability(theta, beta float64) float64 {
	return math.Exp(theta-beta) / (1 + math.Exp(theta-beta))
}

// UpdateTheta memperbarui kemampuan siswa (theta) berdasarkan jawaban.
// Jika benar: theta naik. Jika salah: theta turun.
func UpdateTheta(theta, probability float64, correct bool) float64 {
	u := 0.0
	if correct {
		u = 1.0
	}
	learningRate := 0.3
	return theta + learningRate*(u-probability)
}

// UpdateBeta memperbarui kesulitan soal (beta) berdasarkan jawaban user.
// Kebalikan UpdateTheta: jika user benar → beta turun (soal lebih mudah dari dugaan),
// jika user salah → beta naik (soal lebih sulit dari dugaan).
func UpdateBeta(beta, theta float64, correct bool) float64 {
	u := 0.0
	if correct {
		u = 1.0
	}
	p := RaschProbability(theta, beta)
	learningRate := 0.3
	return beta + learningRate*(p-u)
}
