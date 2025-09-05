package visual

import (
	"fmt"
	"os"
	"gomoku/internal/game"
	"github.com/fogleman/gg"
)

const (
	CellSize     = 40
	BoardSize    = 15
	LeftMargin   = 15  // Space for row labels (01-15)
	TopMargin    = 15  // Space for column labels (A-O) 
	RightMargin  = 8   // Minimal right padding
	BottomMargin = 15  // Space for move history
	OutputFile   = "gomoku.png"
	MinPadding   = 2   // Reduced padding around occupied area (1-2 requirement)
	MaxPadding   = 2   // Maximum padding around occupied area
)

// DisplayArea represents the area of the board to display
type DisplayArea struct {
	MinRow, MaxRow int
	MinCol, MaxCol int
	Width, Height  int
}

// calculateOptimalDisplayArea determines the best area to display based on stone positions
// Always returns a square area to meet the requirement
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
	
	// Calculate occupied area size
	occupiedWidth := maxCol - minCol + 1
	occupiedHeight := maxRow - minRow + 1
	
	// Dynamic padding: if occupied area < 4x4, use 3 grid padding; otherwise use 1 grid padding
	var padding int
	if occupiedWidth < 4 && occupiedHeight < 4 {
		padding = 3  // Use 3 circles of extra board for smaller areas
	} else {
		padding = 1  // Use only 1 circle of extra board for areas >= 4x4
	}
	
	// Find optimal square position using enumeration for perfect symmetry
	occupiedRows := maxRow - minRow + 1
	occupiedCols := maxCol - minCol + 1
	requiredSize := max(occupiedRows, occupiedCols) + 2*padding
	
	// Ensure size doesn't exceed board
	requiredSize = min(requiredSize, BoardSize)
	
	bestScore := -1000
	displayMinRow, displayMinCol := -1, -1
	
	// Enumerate all possible square positions to find the most symmetric one
	for startRow := 0; startRow <= BoardSize-requiredSize; startRow++ {
		for startCol := 0; startCol <= BoardSize-requiredSize; startCol++ {
			endRow := startRow + requiredSize - 1
			endCol := startCol + requiredSize - 1
			
			// Check if square contains all stones
			if startRow <= minRow && endRow >= maxRow && 
			   startCol <= minCol && endCol >= maxCol {
				
				// Calculate actual padding distribution
				topPad := minRow - startRow
				bottomPad := endRow - maxRow
				leftPad := minCol - startCol
				rightPad := endCol - maxCol
				
				// Calculate symmetry score (smaller differences = better)
				rowDiff := topPad - bottomPad
				if rowDiff < 0 { rowDiff = -rowDiff }
				colDiff := leftPad - rightPad  
				if colDiff < 0 { colDiff = -colDiff }
				symmetryScore := -(rowDiff + colDiff)
				
				// Prefer positions closer to target padding
				avgPadding := (topPad + bottomPad + leftPad + rightPad) / 4
				paddingDiff := avgPadding - padding
				if paddingDiff < 0 { paddingDiff = -paddingDiff }
				paddingScore := -paddingDiff
				
				totalScore := symmetryScore*10 + paddingScore  // Weight symmetry higher
				
				if totalScore > bestScore {
					bestScore = totalScore
					displayMinRow, displayMinCol = startRow, startCol
				}
			}
		}
	}
	
	displayMaxRow := displayMinRow + requiredSize - 1
	displayMaxCol := displayMinCol + requiredSize - 1
	finalWidth := displayMaxCol - displayMinCol + 1
	finalHeight := displayMaxRow - displayMinRow + 1
	
	return DisplayArea{
		MinRow: displayMinRow,
		MaxRow: displayMaxRow,
		MinCol: displayMinCol,
		MaxCol: displayMaxCol,
		Width:  finalWidth,
		Height: finalHeight,
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
	
	// Calculate image dimensions with optimized margins
	imgWidth := displayArea.Width*CellSize + LeftMargin + RightMargin
	
	// Dynamic height calculation: only add space for text when game is over
	extraSpace := 0
	if gs.GameOver {
		extraSpace = 25  // Space for game over status text
	}
	imgHeight := displayArea.Height*CellSize + TopMargin + BottomMargin + extraSpace

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
		x := float64(LeftMargin + i*CellSize + CellSize/2)
		y1 := float64(TopMargin + CellSize/2)
		y2 := float64(TopMargin + (displayArea.Height-1)*CellSize + CellSize/2)
		dc.DrawLine(x, y1, x, y2)
		dc.Stroke()
	}

	// Horizontal lines
	for i := 0; i < displayArea.Height; i++ {
		y := float64(TopMargin + i*CellSize + CellSize/2)
		x1 := float64(LeftMargin + CellSize/2)
		x2 := float64(LeftMargin + (displayArea.Width-1)*CellSize + CellSize/2)
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
		x := float64(LeftMargin + i*CellSize + CellSize/2)
		y := float64(TopMargin - 8)  // Closer to board
		label := string(rune('A' + actualCol))
		dc.DrawStringAnchored(label, x, y, 0.5, 0.5)
	}

	// Row labels (01-15) - only for displayed rows
	for i := 0; i < displayArea.Height; i++ {
		actualRow := displayArea.MinRow + i
		x := float64(LeftMargin - 8)   // Adjusted for smaller margin
		y := float64(TopMargin + i*CellSize + CellSize/2)
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
			
			x := float64(LeftMargin + displayCol*CellSize + CellSize/2)
			y := float64(TopMargin + displayRow*CellSize + CellSize/2)

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
	
	yOffset := float64(TopMargin + displayArea.Height*CellSize + 10)
	
	// Load smaller font for moves
	if err := dc.LoadFontFace("/usr/share/fonts/truetype/liberation/LiberationSans-Regular.ttf", 10); err != nil {
		// Fallback - no font loading, text will use default
	}
	
	// TODO: Temporarily disabled move history display - keep logic for future use
	// Show recent moves in single line (last 3 moves)
	// startIdx := 0
	// if len(gs.Moves) > 3 {
	// 	startIdx = len(gs.Moves) - 3
	// }
	// movesToShow := gs.Moves[startIdx:]
	// if len(movesToShow) > 0 {
	// 	// Build move history string with new format
	// 	moveStrs := make([]string, len(movesToShow))
	// 	for i, move := range movesToShow {
	// 		moveStrs[i] = fmt.Sprintf("(%s-%c)", move.Position, move.Piece)
	// 	}
	// 	
	// 	moveHistory := moveStrs[0]
	// 	for i := 1; i < len(moveStrs); i++ {
	// 		moveHistory += " -> " + moveStrs[i]
	// 	}
	// 	// Remove ultrathink suffix
	// 	
	// 	dc.DrawStringAnchored(moveHistory, float64(LeftMargin), yOffset, 0, 0.5)
	// 	yOffset += 15
	// }

	// Game status in same line if space allows, otherwise next line
	if gs.GameOver {
		status := gs.GetGameStatus()
		dc.SetRGB(0.8, 0, 0) // Red text for game over
		dc.DrawStringAnchored(status, float64(LeftMargin), yOffset, 0, 0.5)
	}
}

// ImageExists checks if gomoku.png exists
func ImageExists() bool {
	_, err := os.Stat(OutputFile)
	return !os.IsNotExist(err)
}