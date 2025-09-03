# Gomoku Game Simulator

A professional command-line Gomoku (Five in a Row) game simulator with AI-powered strategic analysis, built in Go.

## Features

🎮 **Complete Game Engine**
- Full 15×15 Gomoku implementation with win detection
- Turn-based gameplay with move validation
- Automatic game state persistence (JSON)

🤖 **AI Strategic Analysis** 
- Professional game commentary and strategic insights
- Position evaluation and move recommendations  
- Powered by AWS Claude via Bedrock
- Intelligent caching for instant repeated analysis

📊 **Rich Visualization**
- Auto-generated PNG board visualization
- Real-time updates after each move
- Clean coordinate system (A-O columns, 01-15 rows)

⚡ **Developer Friendly**
- Command-line interface for easy integration
- Comprehensive move history tracking
- Robust error handling and validation

## Quick Start

### Prerequisites

1. **Go 1.21+** installed
2. **AWS Credentials** (optional, for AI analysis):
   ```bash
   export AWS_BEARER_TOKEN_BEDROCK="your-token"
   ```

### Installation

```bash
# build
cd gomoku_game
go build -o gomoku cmd/main.go
```

### Verify Installation
```bash
./gomoku
# Should show usage information
```

## How to Play Gomoku

### Game Setup

**Step 1: Start Fresh Game**
```bash
# Game auto-initializes on first move
./gomoku status
# Output: WHITE to move (Total Moves: 0)
```

### Playing the Game (Two Players)

**Step 2: Player 1 (WHITE) Makes Opening Move**
```bash
# WHITE goes first, play center for good strategy
./gomoku move H-08-O
# Output: Move H-08-O successful! Status: BLACK to move
```

**Step 3: Player 2 (BLACK) Responds**  
```bash
# BLACK player makes their move
./gomoku move I-08-X
# Output: Move I-08-X successful! Status: WHITE to move
```

**Step 4: Continue Alternating Turns**
```bash
# Keep alternating until someone wins or board fills
./gomoku move G-08-O  # WHITE's turn
./gomoku move J-08-X  # BLACK's turn
./gomoku move F-08-O  # WHITE's turn
# ... continue playing
```

### Monitoring the Game

**Check Current Status**
```bash
./gomoku status
# Output: WHITE to move (Total Moves: 6, Last Move: J-08-X)
```

**View Complete History**
```bash
./gomoku history
# Shows numbered list of all moves played
```

**Get AI Strategic Analysis**
```bash
./gomoku analyze
# Professional commentary on current position
# Includes move recommendations and strategic insights
```

**View Board Visualization**
```bash
# Check the auto-generated gomoku.png file
# Updates automatically after each move
ls -la gomoku.png
```

### Game End

**When Someone Wins**
```bash
./gomoku move E-08-O  # Completing 5 in a row
# Output: Move E-08-O successful! Status: Game Over - WHITE wins!

# Further moves are rejected
./gomoku move A-01-X
# Output: Error: game is over, WHITE has won
```

## Command Reference

| Command | Description | Example |
|---------|-------------|---------|
| `./gomoku move <pos>` | Make a move | `./gomoku move H-08-X` |
| `./gomoku status` | Show current turn | Shows who moves next |
| `./gomoku history` | List all moves | Numbered move history |  
| `./gomoku analyze` | AI analysis | Strategic commentary |

### Move Format
- **Pattern**: `{Column}-{Row}-{Piece}`
- **Columns**: A-O (A=leftmost, O=rightmost)  
- **Rows**: 01-15 (01=top, 15=bottom)
- **Pieces**: X=Black, O=White
- **Example**: `H-08-X` = Black stone at center

## Typical Game Flow

### For Two Players Sharing One Computer:

1. **Initialize**: Run `./gomoku status` to confirm empty board
2. **Player 1 (WHITE)**: `./gomoku move H-08-O` (center opening)
3. **Player 2 (BLACK)**: `./gomoku move I-08-X` (adjacent response)
4. **Continue alternating** until someone wins
5. **Optional**: Use `./gomoku analyze` for strategic guidance
6. **Monitor**: Check `./gomoku history` and `gomoku.png` visualization


## AI Analysis Features

The AI analysis provides professional-grade commentary including:

- **Position Verification**: Confirms all stone positions
- **Strategic Assessment**: Evaluates both players' positions  
- **Tactical Threats**: Identifies immediate winning/blocking moves
- **Move Recommendations**: Suggests best next moves with reasoning
- **Game Commentary**: Professional analysis like a sports commentator

## File Outputs

- `gamestate.json` - Complete game state (auto-created)
- `gomoku.png` - Visual board representation (auto-updated)  
- `analysis_cache.json` - AI analysis cache (auto-managed)

## Game Rules

- **Objective**: Get 5 stones in a row (horizontal, vertical, or diagonal)
- **Turn Order**: WHITE moves first, then alternating
- **Board**: 15×15 grid with A-O columns and 01-15 rows
- **Winning**: Game ends immediately when 5-in-a-row is achieved
- **No Moves**: After game ends, further moves are rejected

## Troubleshooting

**"Invalid move" errors:**
- Check if it's your turn (`./gomoku status`)
- Verify position format (e.g., `H-08-X` not `h8x`)  
- Ensure position isn't occupied (`./gomoku history`)

**AI analysis unavailable:**
- Set `AWS_BEARER_TOKEN_BEDROCK` environment variable
- Game works fine without AI - only analysis is affected

**Build issues:**
- Ensure Go 1.21+ is installed
- Run `go mod tidy` to fetch dependencies

---

🎮 **Enjoy your Gomoku games!** The simulator handles all the game logic while you focus on strategy.
