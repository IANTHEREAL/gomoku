package cli

import (
	"fmt"
	"strings"
	"gomoku/internal/storage"
	"gomoku/internal/visual"
	"gomoku/internal/analysis"
)

// Handler handles different CLI commands
type Handler struct {
	analyzer *analysis.BedrockProvider
}

// NewHandler creates a new CLI handler
func NewHandler() (*Handler, error) {
	// Initialize Bedrock analyzer (optional - will handle error gracefully if not configured)
	analyzer, _ := analysis.NewBedrockProvider("us.anthropic.claude-3-7-sonnet-20250219-v1:0")
	
	return &Handler{
		analyzer: analyzer,
	}, nil
}

// HandleCommand processes the command line arguments and executes the appropriate action
func (h *Handler) HandleCommand(args []string) error {
	if len(args) == 0 {
		return h.showUsage()
	}

	command := strings.ToLower(args[0])

	switch command {
	case "move":
		if len(args) != 2 {
			return fmt.Errorf("usage: gomoku move <position> (e.g., gomoku move H-08-X)")
		}
		return h.handleMove(args[1])

	case "status":
		return h.handleStatus()

	case "history":
		return h.handleHistory()

	case "analyze":
		return h.handleAnalyze()

	default:
		return fmt.Errorf("unknown command: %s", command)
	}
}

// handleMove processes a move command
func (h *Handler) handleMove(moveStr string) error {
	// Load current game state
	gs, err := storage.LoadGameState()
	if err != nil {
		return fmt.Errorf("failed to load game state: %w", err)
	}

	// Make the move
	err = gs.MakeMove(moveStr)
	if err != nil {
		// Enhanced error message with format explanation
		errorMsg := fmt.Sprintf("invalid move: %v\n\n", err)
		errorMsg += "Correct format: <Column>-<Row>-<Piece>\n"
		errorMsg += "Examples:\n"
		errorMsg += "  H-08-X  (Black stone at center)\n"
		errorMsg += "  A-01-O  (White stone at top-left)\n"
		errorMsg += "  O-15-X  (Black stone at bottom-right)\n\n"
		errorMsg += "Rules:\n"
		errorMsg += "  Columns: A-O (A=left, O=right)\n"
		errorMsg += "  Rows: 01-15 (01=top, 15=bottom)\n"
		errorMsg += "  Pieces: X=Black, O=White"
		return fmt.Errorf("%s", errorMsg)
	}

	// Save updated game state
	err = storage.SaveGameState(gs)
	if err != nil {
		return fmt.Errorf("failed to save game state: %w", err)
	}

	// Generate updated board image
	err = visual.GenerateBoardImage(gs)
	if err != nil {
		return fmt.Errorf("failed to generate board image: %w", err)
	}

	// Show result
	fmt.Printf("Move %s successful!\n", moveStr)
	fmt.Printf("Status: %s\n", gs.GetGameStatus())
	fmt.Printf("Board visualization updated: %s\n", visual.OutputFile)

	return nil
}

// handleStatus shows current game status with ASCII board
func (h *Handler) handleStatus() error {
	gs, err := storage.LoadGameState()
	if err != nil {
		return fmt.Errorf("failed to load game state: %w", err)
	}

	fmt.Printf("Game Status: %s\n", gs.GetGameStatus())
	fmt.Printf("Total Moves: %d\n", len(gs.Moves))
	
	if len(gs.Moves) > 0 {
		lastMove := gs.Moves[len(gs.Moves)-1]
		fmt.Printf("Last Move: %d. %s-%c (%s)\n", lastMove.MoveNum, lastMove.Position, lastMove.Piece, lastMove.Player)
	}

	fmt.Println("\nCurrent Board:")
	fmt.Println(gs.String())

	return nil
}

// handleHistory shows move history
func (h *Handler) handleHistory() error {
	gs, err := storage.LoadGameState()
	if err != nil {
		return fmt.Errorf("failed to load game state: %w", err)
	}

	if len(gs.Moves) == 0 {
		fmt.Println("No moves have been made yet.")
		return nil
	}

	fmt.Printf("Move History (%d moves):\n", len(gs.Moves))
	fmt.Println(strings.Repeat("-", 40))
	
	history := gs.GetMoveHistory()
	for _, moveStr := range history {
		fmt.Println(moveStr)
	}

	fmt.Println(strings.Repeat("-", 40))
	fmt.Printf("Current Status: %s\n", gs.GetGameStatus())

	return nil
}

// handleAnalyze performs AI analysis of current position
func (h *Handler) handleAnalyze() error {
	gs, err := storage.LoadGameState()
	if err != nil {
		return fmt.Errorf("failed to load game state: %w", err)
	}

	if h.analyzer == nil {
		return fmt.Errorf("AI analysis not available - AWS credentials not configured")
	}

	// Check if AWS is properly configured
	if !analysis.IsConfigured() {
		return fmt.Errorf("AI analysis not available - missing AWS credentials")
	}

	fmt.Println("Analyzing current position...")
	fmt.Println(strings.Repeat("=", 80))

	analysisResult, err := h.analyzer.AnalyzeGamePosition(gs)
	if err != nil {
		return fmt.Errorf("failed to analyze position: %w", err)
	}

	fmt.Println(analysisResult)
	fmt.Println(strings.Repeat("=", 80))

	return nil
}

// showUsage displays usage information
func (h *Handler) showUsage() error {
	fmt.Println("Gomoku Game Simulator")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  gomoku move <position>   Make a move (e.g., gomoku move H-08-X)")
	fmt.Println("  gomoku status            Show current game status and turn")
	fmt.Println("  gomoku history           Show complete move history")
	fmt.Println("  gomoku analyze           AI strategic analysis (requires AWS credentials)")
	fmt.Println("")
	fmt.Println("Move Format:")
	fmt.Println("  <Column>-<Row>-<Piece>")
	fmt.Println("  Column: A-O (A=leftmost, O=rightmost)")
	fmt.Println("  Row: 01-15 (01=top, 15=bottom)")
	fmt.Println("  Piece: X=Black, O=White")
	fmt.Println("  Example: H-08-X (Black stone at center)")
	fmt.Println("")
	fmt.Println("Notes:")
	fmt.Println("  - Game auto-initializes on first move")
	fmt.Println("  - Board visualization (gomoku.png) updates automatically")
	fmt.Println("  - White moves first")

	return nil
}