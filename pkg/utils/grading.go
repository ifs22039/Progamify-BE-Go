package utils

import (
	"boysitorus/Progamify-Restful-API/internal/config"
	"encoding/json"
	"log"

	"github.com/go-resty/resty/v2"
)

func EssayGrading(expected string, actual string) (float64, error) {
	apiKey := config.LoadConfig().HuggingfaceApiKey
	log.Println("Using API Key:", apiKey)

	client := resty.New()

	log.Println("Starting EssayGrading function...")

	// Kirim permintaan ke API Hugging Face
	response, err := client.R().
		SetHeader("Authorization", "Bearer "+apiKey).
		SetBody(map[string]interface{}{
			"inputs": map[string]interface{}{
				"source_sentence": expected,
				"sentences":       []string{actual},
			},
		}).
		Post("https://api-inference.huggingface.co/models/sentence-transformers/all-MiniLM-L6-v2")

	if err != nil {
		log.Println("Error sending request:", err)
		return 0, err
	}

	// Log status code & response body
	log.Println("Response Status Code:", response.StatusCode())
	log.Println("Response Body:", string(response.Body()))

	// Parsing JSON response
	var data []float64
	if err := json.Unmarshal(response.Body(), &data); err != nil {
		log.Println("Error parsing JSON:", err)
		return 0, err
	}

	// Hitung skor kesamaan
	similarityScore := 0.0
	if len(data) > 0 {
		similarityScore = data[0] * 100
	}

	log.Println("Similarity Score:", similarityScore)

	return similarityScore, nil
}