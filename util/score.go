package util

import (
	"encoding/json"
	"fmt"
	"os"
)

const scoresFile = "scores.json" // File to store scores

// ScoresData represents the scores that will be saved to disk
type scoresData struct {
	LastScore int `json:"lastScore"`
	BestScore int `json:"bestScore"`
}

// SaveScores saves the last and best scores to disk
func SaveScores(lastScore, bestScore int) error {
	data := scoresData{
		LastScore: lastScore,
		BestScore: bestScore,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshaling scores: %w", err)
	}

	err = os.WriteFile(scoresFile, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error writing scores file: %w", err)
	}

	return nil
}

// LoadScores loads the last and best scores from disk
func LoadScores() (int, int, error) {
	// Check if the file exists
	if _, err := os.Stat(scoresFile); os.IsNotExist(err) {
		return 0, 0, nil
	}

	jsonData, err := os.ReadFile(scoresFile)
	if err != nil {
		return 0, 0, fmt.Errorf("error reading scores file: %w", err)
	}

	var data scoresData
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return 0, 0, fmt.Errorf("error unmarshaling scores: %w", err)
	}

	return data.LastScore, data.BestScore, nil
}
