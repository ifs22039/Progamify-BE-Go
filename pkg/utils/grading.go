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
	SimilarityScore float64 `json:"similarity_score"`
}

func EssayGrading(expected string, actual string) (float64, error) {
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
		return 0, err
	}

	log.Println("Response Status Code:", response.StatusCode())
	log.Println("Response Body:", string(response.Body()))

	var data EssayResponse
	if err := json.Unmarshal(response.Body(), &data); err != nil {
		log.Println("Error parsing JSON:", err)
		return 0, err
	}

	similarityScore := data.SimilarityScore * 100

	log.Println("Similarity Score:", similarityScore)

	return similarityScore, nil
}
