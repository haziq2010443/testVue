package main

import "fmt"

// Scores represents a score record for a character
type Scores struct {
	ScoreID     int
	CharID      int
	RewardScore int
}

// NewScores creates a new score entry
func NewScores(scoreID, charID, rewardScore int) *Scores {
	return &Scores{
		ScoreID:     scoreID,
		CharID:      charID,
		RewardScore: rewardScore,
	}
}

// Display prints the score details
func (s *Scores) Display() {
	fmt.Printf("Score ID: %d, Character ID: %d, Reward Score: %d\n", s.ScoreID, s.CharID, s.RewardScore)
}
