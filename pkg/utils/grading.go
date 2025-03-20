package utils

import (
	"boysitorus/Progamify-Restful-API/internal/config"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
)

func EssayGrading(expected string, actual string) (float64, error) {
	apiKey := config.LoadConfig().HuggingfaceApiKey
	client := resty.New()

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
		return 0, err
	}

	var data []float64
	fmt.Println(response.StatusCode())
	if err := json.Unmarshal(response.Body(), &data); err != nil {
		return 0, err
	}

	similarityScore := 0.0
	if len(data) > 0 {
		similarityScore = data[0] * 100
	}

	return similarityScore, nil
}
