package analysis

import (
	"context"
	"fmt"
	"os"
	"strings"
	"gomoku/internal/game"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
)

// BedrockProvider provides AWS Bedrock LLM functionality
type BedrockProvider struct {
	client  *bedrockruntime.Client
	modelID string
}

// NewBedrockProvider creates a new BedrockProvider instance
func NewBedrockProvider(modelID string) (*BedrockProvider, error) {
	// Check for required AWS environment variables
	if os.Getenv("AWS_BEARER_TOKEN_BEDROCK") == "" {
		return nil, fmt.Errorf("missing required AWS environment variable: AWS_BEARER_TOKEN_BEDROCK")
	}

	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(getAWSRegion()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := bedrockruntime.NewFromConfig(cfg)

	return &BedrockProvider{
		client:  client,
		modelID: modelID,
	}, nil
}

// getAWSRegion returns the AWS region from environment or default
func getAWSRegion() string {
	if region := os.Getenv("AWS_REGION"); region != "" {
		return region
	}
	return "us-west-2"
}

// IsConfigured checks if AWS credentials are properly configured
func IsConfigured() bool {
	// Only use AWS_BEARER_TOKEN_BEDROCK
	return os.Getenv("AWS_BEARER_TOKEN_BEDROCK") != ""
}

// AnalyzeGamePosition performs strategic analysis of the current game position with caching
func (bp *BedrockProvider) AnalyzeGamePosition(gs *game.GameState) (string, error) {
	// Check if cached analysis is valid for current board state
	if gs.Analysis != "" && gs.AnalysisHash == gs.CurrentBoardHash {
		return fmt.Sprintf("[CACHED ANALYSIS]\n%s", gs.Analysis), nil
	}

	// Cache miss or invalid, generate fresh analysis
	boardStr := formatBoardForAnalysis(gs)
	prompt := createAnalysisPrompt(gs, boardStr)
	
	systemPrompt := `You are a professional Gomoku (Five in a Row) game analyst and commentator. 
Provide strategic insights like a sports commentator, analyzing the current position, player advantages/disadvantages, 
and tactical opportunities. Be clear, engaging, and educational in your analysis.`

	analysis, err := bp.Generate(prompt, &systemPrompt)
	if err != nil {
		return "", err
	}

	// Cache the analysis result
	gs.Analysis = analysis
	gs.AnalysisHash = gs.CurrentBoardHash

	return analysis, nil
}

// Generate performs synchronous text generation (same as reference implementation)
func (bp *BedrockProvider) Generate(prompt string, systemPrompt *string) (string, error) {
	messages := []types.Message{
		{
			Role: types.ConversationRoleUser,
			Content: []types.ContentBlock{
				&types.ContentBlockMemberText{
					Value: prompt,
				},
			},
		},
	}

	input := &bedrockruntime.ConverseInput{
		ModelId:  aws.String(bp.modelID),
		Messages: messages,
		InferenceConfig: &types.InferenceConfiguration{
			MaxTokens:   aws.Int32(1000),
			Temperature: aws.Float32(0.6),
		},
	}

	// Add system prompt if provided
	if systemPrompt != nil && *systemPrompt != "" {
		input.System = []types.SystemContentBlock{
			&types.SystemContentBlockMemberText{
				Value: *systemPrompt,
			},
		}
	}

	response, err := bp.client.Converse(context.TODO(), input)
	if err != nil {
		return "", fmt.Errorf("bedrock converse failed: %w", err)
	}

	// Extract text from response
	if msgOutput, ok := response.Output.(*types.ConverseOutputMemberMessage); ok {
		if len(msgOutput.Value.Content) > 0 {
			if textBlock, ok := msgOutput.Value.Content[0].(*types.ContentBlockMemberText); ok {
				return textBlock.Value, nil
			}
		}
	}

	return "", fmt.Errorf("no text content in response")
}

// formatBoardForAnalysis creates a text representation of the board for LLM analysis
func formatBoardForAnalysis(gs *game.GameState) string {
	var sb strings.Builder
	
	// Header with column labels
	sb.WriteString("   ")
	for col := 0; col < game.BoardSize; col++ {
		sb.WriteString(fmt.Sprintf(" %c", 'A'+col))
	}
	sb.WriteString("\n")
	
	// Board with row labels
	for row := 0; row < game.BoardSize; row++ {
		sb.WriteString(fmt.Sprintf("%02d ", row+1))
		for col := 0; col < game.BoardSize; col++ {
			sb.WriteString(fmt.Sprintf(" %c", gs.Board[row][col]))
		}
		sb.WriteString("\n")
	}
	
	return sb.String()
}

// createAnalysisPrompt creates the analysis prompt for the LLM with enhanced accuracy
func createAnalysisPrompt(gs *game.GameState, boardStr string) string {
	var prompt strings.Builder
	
	prompt.WriteString("You are analyzing a Gomoku (Five in a Row) board. ACCURACY IS CRITICAL - please read the board very carefully.\n\n")
	
	// Clear symbol explanation
	prompt.WriteString("**BOARD SYMBOLS:**\n")
	prompt.WriteString("- '+' = Empty intersection\n")
	prompt.WriteString("- 'O' = WHITE stone\n") 
	prompt.WriteString("- 'X' = BLACK stone\n\n")
	
	prompt.WriteString("**COORDINATE SYSTEM:**\n")
	prompt.WriteString("- Columns: A-O (left to right)\n")
	prompt.WriteString("- Rows: 01-15 (top to bottom)\n")
	prompt.WriteString("- Example: H-08 means column H, row 8\n\n")
	
	// Board display
	prompt.WriteString("**CURRENT BOARD POSITION:**\n")
	prompt.WriteString(boardStr)
	prompt.WriteString("\n")
	
	// Move history for cross-verification
	prompt.WriteString("**MOVE HISTORY FOR VERIFICATION:**\n")
	if len(gs.Moves) > 0 {
		for _, move := range gs.Moves {
			prompt.WriteString(fmt.Sprintf("%d. %s = %c (%s player)\n", move.MoveNum, move.Position, move.Piece, move.Player))
		}
	} else {
		prompt.WriteString("No moves played yet - empty board\n")
	}
	prompt.WriteString("\n")
	
	// Game context
	prompt.WriteString(fmt.Sprintf("**GAME STATUS:** %s\n", gs.GetGameStatus()))
	prompt.WriteString(fmt.Sprintf("**TOTAL MOVES:** %d\n\n", len(gs.Moves)))
	
	// Step-by-step verification requirement
	prompt.WriteString("**STEP 1: BOARD VERIFICATION (REQUIRED)**\n")
	prompt.WriteString("Before analysis, please verify the board by listing ALL stone positions:\n")
	prompt.WriteString("- List each WHITE stone position (O) with its coordinates\n")
	prompt.WriteString("- List each BLACK stone position (X) with its coordinates\n") 
	prompt.WriteString("- Cross-check these positions against the move history above\n")
	prompt.WriteString("- If any discrepancies found, note them clearly\n\n")
	
	// Analysis requirements
	prompt.WriteString("**STEP 2: STRATEGIC ANALYSIS**\n")
	prompt.WriteString("After verifying the board, provide comprehensive analysis:\n\n")
	prompt.WriteString("1. **Position Summary**: Describe stone formations and patterns\n")
	prompt.WriteString("2. **BLACK's Position**: Analyze advantages, threats, opportunities\n")
	prompt.WriteString("3. **WHITE's Position**: Analyze advantages, threats, opportunities\n")
	prompt.WriteString("4. **Tactical Assessment**: Immediate threats and key intersections\n")
	prompt.WriteString("5. **Strategic Outlook**: Who has the better position and why\n")
	
	if !gs.GameOver {
		prompt.WriteString(fmt.Sprintf("6. **Recommended Move**: Best next move for %s with detailed reasoning\n", gs.CurrentTurn))
	}
	
	prompt.WriteString("\n**CRITICAL REMINDERS:**\n")
	prompt.WriteString("- Double-check each stone position against coordinates\n")
	prompt.WriteString("- 'O' = WHITE, 'X' = BLACK - do not confuse these\n")
	prompt.WriteString("- Verify your stone counts match the move history\n")
	prompt.WriteString("- Be extremely careful with coordinate mapping (A-O columns, 01-15 rows)\n")
	
	return prompt.String()
}