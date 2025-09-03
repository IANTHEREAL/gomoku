package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"gomoku/internal/game"
)

const GameStateFile = "gamestate.json"

// SaveGameState saves the current game state to JSON file
func SaveGameState(gs *game.GameState) error {
	data, err := json.MarshalIndent(gs, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal game state: %w", err)
	}
	
	err = os.WriteFile(GameStateFile, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write game state file: %w", err)
	}
	
	return nil
}

// LoadGameState loads game state from JSON file
func LoadGameState() (*game.GameState, error) {
	// Check if file exists
	if _, err := os.Stat(GameStateFile); os.IsNotExist(err) {
		// File doesn't exist, return new game state
		return game.NewGameState(), nil
	}
	
	data, err := os.ReadFile(GameStateFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read game state file: %w", err)
	}
	
	var gs game.GameState
	err = json.Unmarshal(data, &gs)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal game state: %w", err)
	}
	
	return &gs, nil
}

// GameStateExists checks if a game state file exists
func GameStateExists() bool {
	_, err := os.Stat(GameStateFile)
	return !os.IsNotExist(err)
}

// DeleteGameState removes the game state file
func DeleteGameState() error {
	if !GameStateExists() {
		return nil // Nothing to delete
	}
	
	err := os.Remove(GameStateFile)
	if err != nil {
		return fmt.Errorf("failed to delete game state file: %w", err)
	}
	
	return nil
}