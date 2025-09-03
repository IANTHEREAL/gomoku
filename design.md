# Gomoku Game Simulator Design Document

## Overview

A command-line Gomoku (Five in a Row) game simulator implemented in Go, featuring turn-based gameplay, persistent storage, visualization, and AI-powered strategic analysis.

## Requirements

### Core Features
1. **Turn-based Control**: Strict alternating turns, reject invalid moves
2. **Move Acceptance & Storage**: Accept valid moves and update game state
3. **Current Player Display**: Show whose turn it is (BLACK/WHITE)
4. **Win Detection**: Determine game winner (five in a row)
5. **Visualization**: Generate gomoku.png board image, auto-updated after each move
6. **Strategic Analysis**: AI-powered game analysis with strategic insights
7. **Command-line Interface**: Non-interactive command execution
8. **Move History**: List all moves played in the game

### Technical Requirements
- **Language**: Go
- **Storage**: Efficient local file storage (JSON format)
- **AI Integration**: AWS Bedrock for strategic analysis
- **Visualization**: PNG image generation
- **Coordinate System**: A-O (columns) × 01-15 (rows), format: `H-08-X`

## Architecture

### Project Structure
```
gomoku/
├── cmd/
│   └── main.go           # Command-line entry point
├── internal/
│   ├── game/            # Game core logic
│   │   ├── board.go     # Board state management
│   │   ├── rules.go     # Game rules & win detection  
│   │   └── move.go      # Move processing & validation
│   ├── storage/         # Efficient storage
│   │   └── gamestate.go # JSON format persistence
│   ├── visual/          # Built-in visualization
│   │   └── png.go       # PNG generation
│   ├── analysis/        # AI analysis
│   │   └── bedrock.go   # AWS Bedrock integration
│   └── cli/             # Command handling
│       └── handler.go   # Command routing
└── go.mod
```

## Data Structures

### GameState
```go
type GameState struct {
    Board      [15][15]rune `json:"board"`
    Moves      []Move       `json:"moves"`
    CurrentTurn string      `json:"current_turn"` // "BLACK"/"WHITE"
    Winner     *string      `json:"winner"`       // nil if game in progress
    GameOver   bool         `json:"game_over"`
}
```

### Move
```go
type Move struct {
    Player   string `json:"player"`     // "BLACK"/"WHITE"
    Position string `json:"position"`   // "H-08"
    MoveNum  int    `json:"move_num"`
    Piece    rune   `json:"piece"`      // 'X'/'O'
}
```

## Storage Format (JSON)

```json
{
  "board": [
    ["+","+","+","+","+",...], 
    ["+","X","O","+","+",...],
    ...
  ],
  "moves": [
    {"player":"WHITE","position":"H-08","move_num":1,"piece":"O"},
    {"player":"BLACK","position":"G-08","move_num":2,"piece":"X"}
  ],
  "current_turn": "WHITE",
  "winner": null,
  "game_over": false
}
```

## Command-Line Interface

### Commands
```bash
./gomoku move H-08-X    # Make a move + auto-update PNG
./gomoku status         # Show current turn (BLACK/WHITE to move)  
./gomoku history        # List all moves played
./gomoku analyze        # AI strategic analysis
```

### Built-in Features
- **Auto-initialization**: Creates new game when no save file exists
- **Auto-visualization**: Updates gomoku.png after each move
- **Error handling**: Rejects invalid moves, shows game over status

## Game Rules

### Board Representation
- **Size**: 15×15 grid
- **Empty intersection**: `+`
- **Black stone**: `X` 
- **White stone**: `O`
- **Internal indexing**: 0-based (0,0) to (14,14)

### Coordinate System
- **Columns**: A-O (A=0, B=1, ..., O=14)
- **Rows**: 01-15 (01=0, 02=1, ..., 15=14)
- **Move format**: `{Column}-{Row}-{Piece}` (e.g., `H-08-X`)

### Win Condition
- Five stones in a row (horizontal, vertical, or diagonal)
- Game ends immediately when win condition is met

## Workflow

### Move Processing
1. **Input**: `./gomoku move H-08-X`
2. **Validation**: Check turn, position validity, game status
3. **Update**: Modify board state and move history
4. **Storage**: Save to JSON file
5. **Visualization**: Auto-generate gomoku.png
6. **Output**: Confirm move or show error

### Status Display
1. **Input**: `./gomoku status`
2. **Read**: Load current game state
3. **Output**: Current turn, game status, winner (if any)

### History Display
1. **Input**: `./gomoku history`
2. **Read**: Load move history
3. **Output**: Formatted list of all moves

### Strategic Analysis
1. **Input**: `./gomoku analyze`
2. **Read**: Current board state
3. **AI Call**: AWS Bedrock analysis
4. **Output**: Strategic insights, position evaluation, recommended moves

## Implementation Notes

- **Error Handling**: Graceful error messages for invalid moves, game over scenarios
- **Performance**: Efficient JSON serialization for game state persistence
- **Extensibility**: Modular design allows easy feature additions
- **AWS Integration**: Uses existing Bedrock provider pattern for AI analysis
- **Visualization**: PNG generation based on board state, similar to Python reference implementation