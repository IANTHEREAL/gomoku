package game

// CheckWin checks if the current player has won after the last move
func (gs *GameState) CheckWin(pos Position, player Player) bool {
	piece := GetPieceForPlayer(player)
	
	// Check all four directions: horizontal, vertical, diagonal1, diagonal2
	directions := [][]int{
		{0, 1},  // horizontal
		{1, 0},  // vertical
		{1, 1},  // diagonal \
		{1, -1}, // diagonal /
	}
	
	for _, dir := range directions {
		if gs.countInDirection(pos, piece, dir[0], dir[1]) >= 5 {
			return true
		}
	}
	
	return false
}

// countInDirection counts consecutive pieces in a direction (both ways)
func (gs *GameState) countInDirection(pos Position, piece rune, deltaRow, deltaCol int) int {
	count := 1 // Count the current piece
	
	// Count in positive direction
	count += gs.countInOneDirection(pos, piece, deltaRow, deltaCol)
	
	// Count in negative direction
	count += gs.countInOneDirection(pos, piece, -deltaRow, -deltaCol)
	
	return count
}

// countInOneDirection counts consecutive pieces in one direction
func (gs *GameState) countInOneDirection(pos Position, piece rune, deltaRow, deltaCol int) int {
	count := 0
	row, col := pos.Row+deltaRow, pos.Col+deltaCol
	
	for row >= 0 && row < BoardSize && col >= 0 && col < BoardSize {
		if gs.Board[row][col] == piece {
			count++
			row += deltaRow
			col += deltaCol
		} else {
			break
		}
	}
	
	return count
}

// IsBoardFull checks if the board is completely filled
func (gs *GameState) IsBoardFull() bool {
	for i := 0; i < BoardSize; i++ {
		for j := 0; j < BoardSize; j++ {
			if gs.Board[i][j] == '+' {
				return false
			}
		}
	}
	return true
}