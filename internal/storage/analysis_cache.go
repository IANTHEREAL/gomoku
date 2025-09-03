package storage

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"os"
	"time"
	"gomoku/internal/game"
)

const AnalysisCacheFile = "analysis_cache.json"

// AnalysisCache represents cached analysis results
type AnalysisCache struct {
	Cache map[string]CachedAnalysis `json:"cache"`
}

// CachedAnalysis represents a single cached analysis result
type CachedAnalysis struct {
	Analysis   string    `json:"analysis"`
	Timestamp  time.Time `json:"timestamp"`
	MoveCount  int       `json:"move_count"`
	BoardHash  string    `json:"board_hash"`
}

// NewAnalysisCache creates a new analysis cache
func NewAnalysisCache() *AnalysisCache {
	return &AnalysisCache{
		Cache: make(map[string]CachedAnalysis),
	}
}

// LoadAnalysisCache loads analysis cache from file
func LoadAnalysisCache() (*AnalysisCache, error) {
	// Check if file exists
	if _, err := os.Stat(AnalysisCacheFile); os.IsNotExist(err) {
		// File doesn't exist, return new cache
		return NewAnalysisCache(), nil
	}

	data, err := os.ReadFile(AnalysisCacheFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read analysis cache file: %w", err)
	}

	var cache AnalysisCache
	err = json.Unmarshal(data, &cache)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal analysis cache: %w", err)
	}

	// Initialize map if nil
	if cache.Cache == nil {
		cache.Cache = make(map[string]CachedAnalysis)
	}

	return &cache, nil
}

// SaveAnalysisCache saves analysis cache to file
func SaveAnalysisCache(cache *AnalysisCache) error {
	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal analysis cache: %w", err)
	}

	err = os.WriteFile(AnalysisCacheFile, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write analysis cache file: %w", err)
	}

	return nil
}

// GenerateBoardHash generates a unique hash for the current board state
func GenerateBoardHash(gs *game.GameState) string {
	// Create a string representation of the board state
	boardData := ""
	
	// Include board state
	for i := 0; i < game.BoardSize; i++ {
		for j := 0; j < game.BoardSize; j++ {
			boardData += string(gs.Board[i][j])
		}
	}
	
	// Include move count and current turn for additional uniqueness
	boardData += fmt.Sprintf("_%d_%s", len(gs.Moves), gs.CurrentTurn)
	
	// Generate MD5 hash
	hash := md5.Sum([]byte(boardData))
	return fmt.Sprintf("%x", hash)
}

// GetCachedAnalysis retrieves cached analysis if available
func (ac *AnalysisCache) GetCachedAnalysis(gs *game.GameState) (string, bool) {
	boardHash := GenerateBoardHash(gs)
	
	cached, exists := ac.Cache[boardHash]
	if !exists {
		return "", false
	}
	
	// Verify the cached analysis matches current state
	if cached.MoveCount != len(gs.Moves) {
		// Move count doesn't match, remove invalid cache
		delete(ac.Cache, boardHash)
		return "", false
	}
	
	return cached.Analysis, true
}

// CacheAnalysis stores analysis result in cache
func (ac *AnalysisCache) CacheAnalysis(gs *game.GameState, analysis string) {
	boardHash := GenerateBoardHash(gs)
	
	ac.Cache[boardHash] = CachedAnalysis{
		Analysis:   analysis,
		Timestamp:  time.Now(),
		MoveCount:  len(gs.Moves),
		BoardHash:  boardHash,
	}
}

// CleanOldCache removes outdated cache entries (optional cleanup)
func (ac *AnalysisCache) CleanOldCache(maxAge time.Duration) {
	cutoff := time.Now().Add(-maxAge)
	
	for hash, cached := range ac.Cache {
		if cached.Timestamp.Before(cutoff) {
			delete(ac.Cache, hash)
		}
	}
}

// GetCacheStats returns cache statistics
func (ac *AnalysisCache) GetCacheStats() map[string]interface{} {
	return map[string]interface{}{
		"total_entries": len(ac.Cache),
		"cache_file":    AnalysisCacheFile,
	}
}