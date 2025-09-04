package visual

import (
	"fmt"
	"os"
	"gomoku/internal/game"
	"github.com/fogleman/gg"
)

const (
	CellSize   = 40
	BoardSize  = 15
	Margin     = 18  // Balanced padding for coordinate visibility
	OutputFile = "gomoku.png"
	MinPadding = 1  // Minimum padding around occupied area
	MaxPadding = 1  // Maximum padding around occupied area
)

// DisplayArea represents the area of the board to display
type DisplayArea struct {
	MinRow, MaxRow int
	MinCol, MaxCol int
	Width, Height  int
}

// calculateOptimalDisplayArea determines the best area to display based on stone positions
func calculateOptimalDisplayArea(gs *game.GameState) DisplayArea {
	minRow, maxRow := BoardSize, -1
	minCol, maxCol := BoardSize, -1
	
	// Find the bounds of occupied positions
	for row := 0; row < BoardSize; row++ {
		for col := 0; col < BoardSize; col++ {
			if gs.Board[row][col] != '+' {
				if row < minRow {
					minRow = row
				}
				if row > maxRow {
					maxRow = row
				}
				if col < minCol {
					minCol = col
				}
				if col > maxCol {
					maxCol = col
				}
			}
		}
	}
	
	// If no stones, show center area (7x7 around center)
	if maxRow == -1 {
		center := BoardSize / 2
		return DisplayArea{
			MinRow: center - 3,
			MaxRow: center + 3,
			MinCol: center - 3,
			MaxCol: center + 3,
			Width:  7,
			Height: 7,
		}
	}
	
	// Add padding around occupied area
	padding := MinPadding
	
	// Calculate display bounds with padding
	displayMinRow := max(0, minRow-padding)
	displayMaxRow := min(BoardSize-1, maxRow+padding)
	displayMinCol := max(0, minCol-padding)
	displayMaxCol := min(BoardSize-1, maxCol+padding)
	
	return DisplayArea{
		MinRow: displayMinRow,
		MaxRow: displayMaxRow,
		MinCol: displayMinCol,
		MaxCol: displayMaxCol,
		Width:  displayMaxCol - displayMinCol + 1,
		Height: displayMaxRow - displayMinRow + 1,
	}
}

// max returns the larger of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// min returns the smaller of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// GenerateBoardImage creates a PNG visualization of the current board state
func GenerateBoardImage(gs *game.GameState) error {
	// Calculate optimal display area
	displayArea := calculateOptimalDisplayArea(gs)
	
	// Calculate image dimensions based on display area
	imgWidth := displayArea.Width*CellSize + 2*Margin
	imgHeight := displayArea.Height*CellSize + 2*Margin + 25 // Compact space for move history

	// Create drawing context
	dc := gg.NewContext(imgWidth, imgHeight)

	// Set background color (beige/wood color)
	dc.SetRGB(0.86, 0.78, 0.63)
	dc.Clear()

	// Draw grid lines
	drawGrid(dc, displayArea)

	// Draw coordinate labels
	drawCoordinateLabels(dc, displayArea)

	// Draw stones
	drawStones(dc, gs, displayArea)

	// Draw move history
	drawMoveHistory(dc, gs, displayArea)

	// Save to PNG file
	return dc.SavePNG(OutputFile)
}

// drawGrid draws the board grid lines
func drawGrid(dc *gg.Context, displayArea DisplayArea) {
	dc.SetRGB(0, 0, 0) // Black lines
	dc.SetLineWidth(0.5)

	// Vertical lines
	for i := 0; i < displayArea.Width; i++ {
		x := float64(Margin + i*CellSize + CellSize/2)
		y1 := float64(Margin + CellSize/2)
		y2 := float64(Margin + (displayArea.Height-1)*CellSize + CellSize/2)
		dc.DrawLine(x, y1, x, y2)
		dc.Stroke()
	}

	// Horizontal lines
	for i := 0; i < displayArea.Height; i++ {
		y := float64(Margin + i*CellSize + CellSize/2)
		x1 := float64(Margin + CellSize/2)
		x2 := float64(Margin + (displayArea.Width-1)*CellSize + CellSize/2)
		dc.DrawLine(x1, y, x2, y)
		dc.Stroke()
	}
}

// drawCoordinateLabels draws A-O column labels and 01-15 row labels
func drawCoordinateLabels(dc *gg.Context, displayArea DisplayArea) {
	dc.SetRGB(0, 0, 0) // Black text
	
	// Load font (try system font, fallback to default)
	if err := dc.LoadFontFace("/usr/share/fonts/truetype/liberation/LiberationSans-Regular.ttf", 12); err != nil {
		// Fallback - no font loading, text will use default
	}

	// Column labels (A-O) - only for displayed columns
	for i := 0; i < displayArea.Width; i++ {
		actualCol := displayArea.MinCol + i
		x := float64(Margin + i*CellSize + CellSize/2)
		y := float64(Margin - 8)  // Closer to board
		label := string(rune('A' + actualCol))
		dc.DrawStringAnchored(label, x, y, 0.5, 0.5)
	}

	// Row labels (01-15) - only for displayed rows
	for i := 0; i < displayArea.Height; i++ {
		actualRow := displayArea.MinRow + i
		x := float64(Margin - 8)   // Adjusted for smaller margin
		y := float64(Margin + i*CellSize + CellSize/2)
		label := fmt.Sprintf("%02d", actualRow+1)
		dc.DrawStringAnchored(label, x, y, 0.5, 0.5)
	}
}

// drawStones draws black and white stones on the board
func drawStones(dc *gg.Context, gs *game.GameState, displayArea DisplayArea) {
	stoneRadius := 15.0

	// Find the position of the last move for highlighting
	var lastMoveRow, lastMoveCol int = -1, -1
	if len(gs.Moves) > 0 {
		lastMove := gs.Moves[len(gs.Moves)-1]
		if pos, err := game.ParsePosition(lastMove.Position); err == nil {
			lastMoveRow, lastMoveCol = pos.Row, pos.Col
		}
	}

	// Only draw stones in the display area
	for row := displayArea.MinRow; row <= displayArea.MaxRow; row++ {
		for col := displayArea.MinCol; col <= displayArea.MaxCol; col++ {
			piece := gs.Board[row][col]
			if piece == '+' {
				continue // Empty position
			}

			// Convert absolute board coordinates to display coordinates
			displayRow := row - displayArea.MinRow
			displayCol := col - displayArea.MinCol
			
			x := float64(Margin + displayCol*CellSize + CellSize/2)
			y := float64(Margin + displayRow*CellSize + CellSize/2)

			// Check if this is the last move for special highlighting
			isLastMove := (row == lastMoveRow && col == lastMoveCol)

			// Draw the stone
			dc.DrawCircle(x, y, stoneRadius)

			if piece == 'X' { // Black stone
				dc.SetRGB(0, 0, 0)
				dc.FillPreserve()
				dc.SetRGB(0, 0, 0)
				dc.Stroke()
			} else if piece == 'O' { // White stone
				dc.SetRGB(1, 1, 1)
				dc.FillPreserve()
				dc.SetRGB(0, 0, 0)
				dc.Stroke()
			}

			// Add special effect for the last move
			if isLastMove {
				// Draw a colored highlight ring around the last move
				highlightRadius := stoneRadius + 5
				dc.SetLineWidth(3)
				
				// Use different colors for different pieces
				if piece == 'X' { // Black stone - use bright red highlight
					dc.SetRGB(0.9, 0.2, 0.2)
				} else { // White stone - use bright blue highlight  
					dc.SetRGB(0.2, 0.4, 0.9)
				}
				
				dc.DrawCircle(x, y, highlightRadius)
				dc.Stroke()
				
				// Draw a small dot in the center for extra emphasis
				dc.SetRGB(0.8, 0.8, 0.2) // Yellow center dot
				dc.DrawCircle(x, y, 3)
				dc.Fill()
				
				// Reset line width
				dc.SetLineWidth(1)
			}
		}
	}
}

// drawMoveHistory draws the move history at the bottom in a compact single line
func drawMoveHistory(dc *gg.Context, gs *game.GameState, displayArea DisplayArea) {
	dc.SetRGB(0, 0, 0) // Black text
	
	yOffset := float64(Margin + displayArea.Height*CellSize + 10)
	
	// Load smaller font for moves
	if err := dc.LoadFontFace("/usr/share/fonts/truetype/liberation/LiberationSans-Regular.ttf", 10); err != nil {
		// Fallback - no font loading, text will use default
	}
	
	// Show recent moves in single line (last 3 moves)
	startIdx := 0
	if len(gs.Moves) > 3 {
		startIdx = len(gs.Moves) - 3
	}
	
	movesToShow := gs.Moves[startIdx:]
	if len(movesToShow) > 0 {
		// Build move history string with new format
		moveStrs := make([]string, len(movesToShow))
		for i, move := range movesToShow {
			moveStrs[i] = fmt.Sprintf("step(%s-%c)", move.Position, move.Piece)
		}
		
		moveHistory := moveStrs[0]
		for i := 1; i < len(moveStrs); i++ {
			moveHistory += " -> " + moveStrs[i]
		}
		// Remove ultrathink suffix
		
		dc.DrawStringAnchored(moveHistory, float64(Margin), yOffset, 0, 0.5)
		yOffset += 15
	}

	// Game status in same line if space allows, otherwise next line
	if gs.GameOver {
		status := gs.GetGameStatus()
		dc.SetRGB(0.8, 0, 0) // Red text for game over
		dc.DrawStringAnchored(status, float64(Margin), yOffset, 0, 0.5)
	}
}

// ImageExists checks if gomoku.png exists
func ImageExists() bool {
	_, err := os.Stat(OutputFile)
	return !os.IsNotExist(err)
}