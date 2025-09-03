package game

import (
	"crypto/md5"
	"fmt"
	"strings"
)

const BoardSize = 15

type Player string

const (
	PlayerBlack Player = "BLACK"
	PlayerWhite Player = "WHITE"
)

type Position struct {
	Row int
	Col int
}

type Move struct {
	Player   Player `json:"player"`
	Position string `json:"position"` // "H-08"
	MoveNum  int    `json:"move_num"`
	Piece    rune   `json:"piece"` // 'X'/'O'
}

type GameState struct {
	Board            [BoardSize][BoardSize]rune `json:"board"`
	Moves            []Move                     `json:"moves"`
	CurrentTurn      Player                     `json:"current_turn"`
	Winner           *Player                    `json:"winner"`
	GameOver         bool                       `json:"game_over"`
	Analysis         string                     `json:"analysis,omitempty"`
	AnalysisHash     string                     `json:"analysis_hash,omitempty"`
	CurrentBoardHash string                     `json:"current_board_hash,omitempty"`
}

func NewGameState() *GameState {
	gs := &GameState{
		CurrentTurn: PlayerWhite, // White moves first
		GameOver:    false,
		Winner:      nil,
		Moves:       make([]Move, 0),
	}
	
	// Initialize empty board
	for i := 0; i < BoardSize; i++ {
		for j := 0; j < BoardSize; j++ {
			gs.Board[i][j] = '+'
		}
	}
	
	return gs
}

// ParsePosition converts "H-08" to Position{Row: 7, Col: 7}
func ParsePosition(pos string) (Position, error) {
	parts := strings.Split(pos, "-")
	if len(parts) != 2 {
		return Position{}, fmt.Errorf("invalid position format: %s", pos)
	}
	
	// Parse column (A-O -> 0-14)
	colStr := parts[0]
	if len(colStr) != 1 || colStr[0] < 'A' || colStr[0] > 'O' {
		return Position{}, fmt.Errorf("invalid column: %s", colStr)
	}
	col := int(colStr[0] - 'A')
	
	// Parse row (01-15 -> 0-14)
	rowStr := parts[1]
	if len(rowStr) != 2 {
		return Position{}, fmt.Errorf("invalid row format: %s", rowStr)
	}
	
	var row int
	if _, err := fmt.Sscanf(rowStr, "%02d", &row); err != nil {
		return Position{}, fmt.Errorf("invalid row: %s", rowStr)
	}
	
	if row < 1 || row > 15 {
		return Position{}, fmt.Errorf("row out of range: %d", row)
	}
	
	return Position{Row: row - 1, Col: col}, nil
}

// FormatPosition converts Position{Row: 7, Col: 7} to "H-08"
func FormatPosition(pos Position) string {
	col := rune('A' + pos.Col)
	row := pos.Row + 1
	return fmt.Sprintf("%c-%02d", col, row)
}

// IsValidPosition checks if position is within board bounds
func (gs *GameState) IsValidPosition(pos Position) bool {
	return pos.Row >= 0 && pos.Row < BoardSize && pos.Col >= 0 && pos.Col < BoardSize
}

// IsEmptyPosition checks if position is empty on the board
func (gs *GameState) IsEmptyPosition(pos Position) bool {
	if !gs.IsValidPosition(pos) {
		return false
	}
	return gs.Board[pos.Row][pos.Col] == '+'
}

// GetPieceForPlayer returns the piece character for a player
func GetPieceForPlayer(player Player) rune {
	switch player {
	case PlayerBlack:
		return 'X'
	case PlayerWhite:
		return 'O'
	default:
		return '+'
	}
}

// GetPlayerForPiece returns the player for a piece character
func GetPlayerForPiece(piece rune) Player {
	switch piece {
	case 'X':
		return PlayerBlack
	case 'O':
		return PlayerWhite
	default:
		return ""
	}
}

// NextPlayer returns the next player
func NextPlayer(current Player) Player {
	if current == PlayerBlack {
		return PlayerWhite
	}
	return PlayerBlack
}

// String returns a string representation of the board
func (gs *GameState) String() string {
	var sb strings.Builder
	
	// Header with column labels
	sb.WriteString("   ")
	for col := 0; col < BoardSize; col++ {
		sb.WriteString(fmt.Sprintf(" %c", 'A'+col))
	}
	sb.WriteString("\n")
	
	// Board with row labels
	for row := 0; row < BoardSize; row++ {
		sb.WriteString(fmt.Sprintf("%02d ", row+1))
		for col := 0; col < BoardSize; col++ {
			sb.WriteString(fmt.Sprintf(" %c", gs.Board[row][col]))
		}
		sb.WriteString("\n")
	}
	
	return sb.String()
}

// GenerateBoardHash generates a hash for the current board state
func (gs *GameState) GenerateBoardHash() string {
	boardData := ""
	for i := 0; i < BoardSize; i++ {
		for j := 0; j < BoardSize; j++ {
			boardData += string(gs.Board[i][j])
		}
	}
	hash := md5.Sum([]byte(boardData))
	return fmt.Sprintf("%x", hash)
}