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
	Margin     = 60
	OutputFile = "gomoku.png"
)

// GenerateBoardImage creates a PNG visualization of the current board state
func GenerateBoardImage(gs *game.GameState) error {
	// Calculate image dimensions
	imgWidth := BoardSize*CellSize + 2*Margin
	imgHeight := BoardSize*CellSize + 2*Margin + 100 // Extra space for move history

	// Create drawing context
	dc := gg.NewContext(imgWidth, imgHeight)

	// Set background color (beige/wood color)
	dc.SetRGB(0.86, 0.78, 0.63)
	dc.Clear()

	// Draw grid lines
	drawGrid(dc)

	// Draw coordinate labels
	drawCoordinateLabels(dc)

	// Draw stones
	drawStones(dc, gs)

	// Draw move history
	drawMoveHistory(dc, gs)

	// Save to PNG file
	return dc.SavePNG(OutputFile)
}

// drawGrid draws the board grid lines
func drawGrid(dc *gg.Context) {
	dc.SetRGB(0, 0, 0) // Black lines
	dc.SetLineWidth(1)

	// Vertical lines
	for i := 0; i < BoardSize; i++ {
		x := float64(Margin + i*CellSize + CellSize/2)
		y1 := float64(Margin + CellSize/2)
		y2 := float64(Margin + (BoardSize-1)*CellSize + CellSize/2)
		dc.DrawLine(x, y1, x, y2)
		dc.Stroke()
	}

	// Horizontal lines
	for i := 0; i < BoardSize; i++ {
		y := float64(Margin + i*CellSize + CellSize/2)
		x1 := float64(Margin + CellSize/2)
		x2 := float64(Margin + (BoardSize-1)*CellSize + CellSize/2)
		dc.DrawLine(x1, y, x2, y)
		dc.Stroke()
	}
}

// drawCoordinateLabels draws A-O column labels and 01-15 row labels
func drawCoordinateLabels(dc *gg.Context) {
	dc.SetRGB(0, 0, 0) // Black text
	
	// Load font (try system font, fallback to default)
	if err := dc.LoadFontFace("/usr/share/fonts/truetype/liberation/LiberationSans-Regular.ttf", 16); err != nil {
		// Fallback - no font loading, text will use default
	}

	// Column labels (A-O)
	for i := 0; i < BoardSize; i++ {
		x := float64(Margin + i*CellSize + CellSize/2)
		y := float64(Margin - 10)
		label := string(rune('A' + i))
		dc.DrawStringAnchored(label, x, y, 0.5, 0.5)
	}

	// Row labels (01-15)
	for i := 0; i < BoardSize; i++ {
		x := float64(Margin - 20)
		y := float64(Margin + i*CellSize + CellSize/2)
		label := fmt.Sprintf("%02d", i+1)
		dc.DrawStringAnchored(label, x, y, 0.5, 0.5)
	}
}

// drawStones draws black and white stones on the board
func drawStones(dc *gg.Context, gs *game.GameState) {
	stoneRadius := 15.0

	for row := 0; row < BoardSize; row++ {
		for col := 0; col < BoardSize; col++ {
			piece := gs.Board[row][col]
			if piece == '+' {
				continue // Empty position
			}

			x := float64(Margin + col*CellSize + CellSize/2)
			y := float64(Margin + row*CellSize + CellSize/2)

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
		}
	}
}

// drawMoveHistory draws the move history at the bottom of the image
func drawMoveHistory(dc *gg.Context, gs *game.GameState) {
	dc.SetRGB(0, 0, 0) // Black text
	
	yOffset := float64(Margin + BoardSize*CellSize + 20)
	
	// Title
	dc.DrawStringAnchored("Move History:", float64(Margin), yOffset, 0, 0.5)
	yOffset += 25

	// Recent moves (show last 5 moves to avoid overcrowding)
	startIdx := 0
	if len(gs.Moves) > 5 {
		startIdx = len(gs.Moves) - 5
	}

	for i := startIdx; i < len(gs.Moves); i++ {
		move := gs.Moves[i]
		moveText := fmt.Sprintf("%d. %s-%c (%s)", move.MoveNum, move.Position, move.Piece, move.Player)
		dc.DrawStringAnchored(moveText, float64(Margin), yOffset, 0, 0.5)
		yOffset += 20
	}

	// Game status
	if gs.GameOver {
		yOffset += 10
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