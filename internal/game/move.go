package game

import (
	"fmt"
	"strings"
)

// ValidateAndParseMove validates and parses a move string like "H-08-X"
func ValidateAndParseMove(moveStr string) (string, Player, error) {
	parts := strings.Split(moveStr, "-")
	if len(parts) != 3 {
		return "", "", fmt.Errorf("invalid move format: expected COL-ROW-PIECE, got %s", moveStr)
	}
	
	position := parts[0] + "-" + parts[1] // "H-08"
	pieceStr := parts[2]
	
	// Validate piece
	if len(pieceStr) != 1 {
		return "", "", fmt.Errorf("invalid piece: %s", pieceStr)
	}
	
	piece := rune(pieceStr[0])
	var player Player
	
	switch piece {
	case 'X':
		player = PlayerBlack
	case 'O':
		player = PlayerWhite
	default:
		return "", "", fmt.Errorf("invalid piece: %c (must be X or O)", piece)
	}
	
	// Validate position format
	if _, err := ParsePosition(position); err != nil {
		return "", "", fmt.Errorf("invalid position: %w", err)
	}
	
	return position, player, nil
}

// MakeMove attempts to make a move on the board
func (gs *GameState) MakeMove(moveStr string) error {
	// Check if game is over
	if gs.GameOver {
		if gs.Winner != nil {
			return fmt.Errorf("game is over, %s has won", *gs.Winner)
		}
		return fmt.Errorf("game is over (draw)")
	}
	
	// Parse and validate move
	position, player, err := ValidateAndParseMove(moveStr)
	if err != nil {
		return err
	}
	
	// Check if it's the correct player's turn
	if player != gs.CurrentTurn {
		return fmt.Errorf("it's %s's turn, not %s's turn", gs.CurrentTurn, player)
	}
	
	// Parse position
	pos, err := ParsePosition(position)
	if err != nil {
		return err
	}
	
	// Check if position is valid and empty
	if !gs.IsValidPosition(pos) {
		return fmt.Errorf("position %s is out of bounds", position)
	}
	
	if !gs.IsEmptyPosition(pos) {
		return fmt.Errorf("position %s is already occupied", position)
	}
	
	// Make the move
	piece := GetPieceForPlayer(player)
	gs.Board[pos.Row][pos.Col] = piece
	
	// Add to move history
	move := Move{
		Player:   player,
		Position: position,
		MoveNum:  len(gs.Moves) + 1,
		Piece:    piece,
	}
	gs.Moves = append(gs.Moves, move)
	
	// Check for win
	if gs.CheckWin(pos, player) {
		gs.GameOver = true
		gs.Winner = &player
		return nil
	}
	
	// Check for draw
	if gs.IsBoardFull() {
		gs.GameOver = true
		gs.Winner = nil // Draw
		return nil
	}
	
	// Switch turn
	gs.CurrentTurn = NextPlayer(gs.CurrentTurn)
	
	// Update current board hash and clear cached analysis
	gs.CurrentBoardHash = gs.GenerateBoardHash()
	gs.Analysis = ""
	gs.AnalysisHash = ""
	
	return nil
}

// GetMoveHistory returns formatted move history
func (gs *GameState) GetMoveHistory() []string {
	history := make([]string, len(gs.Moves))
	for i, move := range gs.Moves {
		history[i] = fmt.Sprintf("%d. %s-%c (%s)", move.MoveNum, move.Position, move.Piece, move.Player)
	}
	return history
}

// GetGameStatus returns current game status
func (gs *GameState) GetGameStatus() string {
	if gs.GameOver {
		if gs.Winner != nil {
			return fmt.Sprintf("Game Over - %s wins!", *gs.Winner)
		}
		return "Game Over - Draw!"
	}
	return fmt.Sprintf("%s to move", gs.CurrentTurn)
}