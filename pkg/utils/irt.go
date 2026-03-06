package utils

import "math"

// RaschProbability menghitung probabilitas jawaban benar berdasarkan Rasch Model.
//
// Formula (1) dari paper:
//
//	P_ni = exp(θ - β) / (1 + exp(θ - β))
//
// Parameters:
//   - theta: kemampuan mahasiswa (logit scale)
//   - beta:  tingkat kesulitan soal (logit scale)
//
// Return:
//   - Probabilitas antara 0 dan 1.
//   - Jika P ≈ 0.5, soal sangat sesuai dengan kemampuan mahasiswa (ideal untuk adaptive testing).
//   - Jika P > 0.5, soal terlalu mudah untuk mahasiswa.
//   - Jika P < 0.5, soal terlalu sulit untuk mahasiswa.
func RaschProbability(theta, beta float64) float64 {
	// use a numerically stable logistic function.  when theta-beta is very
	// large or very negative the raw exp() call can overflow/underflow and the
	// result is not useful.  clamp manually to avoid that.
	diff := theta - beta
	if diff > 700 { // e^700 is already ~1e304, near float64 max
		return 1.0
	}
	if diff < -700 {
		return 0.0
	}
	expv := math.Exp(diff)
	return expv / (1 + expv)
}

// UpdateTheta memperbarui estimasi kemampuan mahasiswa (θ) menggunakan metode logit
// berdasarkan akumulasi seluruh jawaban mahasiswa.
//
// Pendekatan dari paper (formula 2 & 5):
//   - θ⁰ = ln(p / q), di mana p = proporsi benar, q = proporsi salah
//   - Nilai logit ini kemudian di-scale ke interval scale yang sama dengan beta
//
// Parameters:
//   - correctCount: jumlah jawaban benar mahasiswa hingga saat ini
//   - totalCount:   jumlah total soal yang sudah dijawab mahasiswa
//
// Return:
//   - Nilai theta baru dalam logit scale.
//   - Jika mahasiswa menjawab semua benar (p=1) atau semua salah (p=0),
//     dikembalikan nilai ekstrem ±4.595 (praktis mendekati ±∞ dalam logit).
func UpdateTheta(correctCount, totalCount int) float64 {
	if totalCount == 0 {
		return 0.0
	}

	p := float64(correctCount) / float64(totalCount)
	q := 1.0 - p

	// Hindari ln(0) — sesuai catatan paper bahwa nilai β > -6 atau β > 6
	// dianggap tidak informatif (probability mendekati 1 atau 0).
	// Clamp p ke range (0.01, 0.99) agar tetap dalam range logit yang bermakna.
	const epsilon = 0.01
	if p < epsilon {
		p = epsilon
	}
	if p > 1-epsilon {
		p = 1 - epsilon
	}
	q = 1.0 - p

	// Formula (2): θ⁰ = ln(p / q)
	return math.Log(p / q)
}

// UpdateBeta memperbarui estimasi kesulitan soal (β) menggunakan metode logit
// berdasarkan akumulasi seluruh jawaban semua mahasiswa terhadap soal tersebut.
//
// Pendekatan dari paper (formula 3 & 4):
//   - β⁰ = ln(q / p), di mana p = proporsi benar oleh semua peserta, q = proporsi salah
//   - Kebalikan dari theta: semakin banyak yang menjawab benar → soal semakin mudah (β turun)
//
// Parameters:
//   - correctCount: jumlah mahasiswa yang menjawab soal ini dengan benar
//   - totalCount:   jumlah total mahasiswa yang sudah menjawab soal ini
//
// Return:
//   - Nilai beta baru dalam logit scale.
func UpdateBeta(correctCount, totalCount int) float64 {
	if totalCount == 0 {
		return 0.0
	}

	p := float64(correctCount) / float64(totalCount)
	q := 1.0 - p

	// Clamp — sama seperti UpdateTheta
	const epsilon = 0.01
	if p < epsilon {
		p = epsilon
	}
	if p > 1-epsilon {
		p = 1 - epsilon
	}
	q = 1.0 - p

	// Formula (3): β⁰ = ln(q / p)
	// Kebalikan theta: soal mudah (p tinggi) → β negatif (mudah)
	// Soal sulit (p rendah) → β positif (sulit)
	return math.Log(q / p)
}

// ItemInformation menghitung nilai information function untuk sebuah soal
// pada kemampuan theta tertentu.
//
// Formula (6) dari paper:
//
//	I_i(θ) = P_i(θ) × Q_i(θ)
//
// Nilai information function maksimum (0.25) tercapai ketika P = 0.5,
// yaitu saat kesulitan soal tepat sesuai kemampuan mahasiswa.
func ItemInformation(theta, beta float64) float64 {
	p := RaschProbability(theta, beta)
	q := 1.0 - p
	return p * q
}

// MeasurementError menghitung standar error pengukuran berdasarkan
// total information function dari semua soal yang sudah dijawab.
//
// Formula (7) & (8) dari paper:
//
//	I(θ) = D² × Σ I_j(θ)   di mana D = 1.7 (faktor koreksi logistik-normal)
//	SE(θ) = 1 / √I(θ)
//
// Parameters:
//   - theta:  kemampuan mahasiswa saat ini
//   - betas:  slice berisi nilai beta semua soal yang sudah dijawab
//
// Return:
//   - Nilai SE. Semakin kecil SE, semakin akurat estimasi kemampuan mahasiswa.
//   - Adaptive testing berhenti ketika SE < threshold yang ditetapkan (lihat paper).
func MeasurementError(theta float64, betas []float64) float64 {
	const D = 1.7
	totalInfo := 0.0
	for _, beta := range betas {
		totalInfo += ItemInformation(theta, beta)
	}
	totalInfo *= D * D

	if totalInfo == 0 {
		return math.Inf(1)
	}
	return 1.0 / math.Sqrt(totalInfo)
}

// SelectNextItem memilih soal berikutnya dari bank soal berdasarkan prinsip adaptive testing.
//
// Dari paper (bagian 2, tahap 7): soal berikutnya dipilih berdasarkan kondisi:
//
//	|θ_i - β_j| = minimum
//
// Artinya dipilih soal yang tingkat kesulitannya (beta) paling mendekati
// kemampuan mahasiswa (theta) saat ini — sehingga P(correct) ≈ 0.5
// dan information function mendekati nilai maksimumnya (0.25).
//
// Parameters:
//   - theta:        kemampuan mahasiswa saat ini
//   - availableBetas: map dari ID soal ke nilai beta soal tersebut
//
// Return:
//   - ID soal yang paling sesuai dengan kemampuan mahasiswa saat ini.
//   - Mengembalikan -1 jika tidak ada soal tersedia.
//
// Example usage (see repository/exercise_repository.go for a real-world
// interface):
//
//    betas := map[int]float64{101: -1.2, 102: 0.3, 103: 1.5}
//    id := SelectNextItem(userTheta, betas)
//    // id will be the question whose beta is closest to userTheta
//
// In a production system we normally iterate, removing the selected ID from
// the map and calling the function again until the desired number of items
// have been picked.
func SelectNextItem(theta float64, availableBetas map[int]float64) int {
	bestID := -1
	bestDiff := math.Inf(1)

	for id, beta := range availableBetas {
		diff := math.Abs(theta - beta)
		if diff < bestDiff {
			bestDiff = diff
			bestID = id
		}
	}

	return bestID
}