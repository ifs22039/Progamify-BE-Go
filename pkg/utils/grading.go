package utils

import (
	"boysitorus/Progamify-Restful-API/internal/config"
	"encoding/json"
	"log"

	"github.com/go-resty/resty/v2"
)

type EssayRequest struct {
	Reference string `json:"reference"`
	Essay     string `json:"essay"`
}

type EssayResponse struct {
	SimilarityScore float64 `json:"similarity"`
	GeminiScore    float64 `json:"gemini_score"`
	FinalScore     float64 `json:"final_score"`
}

// EssayGradeResult is the raw data returned from the external grading API.
// FinalScore and SimilarityScore are in the 0–1 range; callers can convert to
// whichever scale they need.
//
// The previous version of this function only returned a single float (the
// score already converted to percentage).  In order to make the API’s
// `final_score` value available for storage and debugging we now return the
// full result struct.
type EssayGradeResult struct {
	FinalScore      float64
	SimilarityScore float64
}

func EssayGrading(expected string, actual string) (EssayGradeResult, error) {
	url := config.LoadConfig().GradingApiUrl

	client := resty.New()

	log.Println("Starting EssayGrading function...")

	response, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(EssayRequest{
			Reference: expected,
			Essay:     actual,
		}).
		Post(url)

	if err != nil {
		log.Println("Error sending request:", err)
		return EssayGradeResult{}, err
	}

	log.Println("Response Status Code:", response.StatusCode())
	log.Println("Response Body:", string(response.Body()))

	var data EssayResponse
	if err := json.Unmarshal(response.Body(), &data); err != nil {
		log.Println("Error parsing JSON:", err)
		return EssayGradeResult{}, err
	}

	result := EssayGradeResult{
		FinalScore:      data.FinalScore,
		SimilarityScore: data.SimilarityScore,
	}

	// convenience log showing the percentage equivalent used by the repo code
	log.Printf("Calculated essay score: %.2f (final=%.4f, similarity=%.4f)",
		result.FinalScore*100, result.FinalScore, result.SimilarityScore)

	return result, nil
}
