package main

import "fmt"

// Character represents a game character linked to an account
type Character struct {
	CharID  int
	AccID   int
	ClassID int
}

// NewCharacter creates a new character
func NewCharacter(charID, accID, classID int) *Character {
	return &Character{
		CharID:  charID,
		AccID:   accID,
		ClassID: classID,
	}
}

// Display prints the character details
func (c *Character) Display() {
	fmt.Printf("Character ID: %d, Account ID: %d, Class ID: %d\n", c.CharID, c.AccID, c.ClassID)
}
