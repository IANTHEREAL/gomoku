## **Gomoku Game Rules**

This document outlines the rules for our Gomoku game.

**1. Board Representation**
*   **Structure**: The game is played on a 15x15 2D character grid.
*   **Internal Indexing**: The underlying logic uses a 0-based index `(row, column)`, ranging from `(0, 0)` at the top-left to `(14, 14)` at the bottom-right.

**2. Character Mapping**
*   `+` : Represents an empty intersection.
*   `X` : Represents a Black stone.
*   `O` : Represents a White stone.

**3. Coordinate System**
This defines the mapping between the user-facing representation and the internal indices.
*   **Columns (Horizontal)**:
    *   Represented by letters `A` through `O`.
    *   Mapping: `A` -> column `0`, `B` -> column `1`, ..., `O` -> column `14`.
*   **Rows (Vertical)**:
    *   Represented by two-digit numbers `01` through `15`.
    *   Mapping: `01` -> row `0`, `02` -> row `1`, ..., `15` -> row `14`.

**4. Action Input Format**
*   **Format**: A string concatenating the column letter, row number, and player piece, separated by hyphens.
*   **Examples**:
    *   To play a Black stone at the center `(row=7, column=7)`, the command is: `H-08-X`.
    *   To play a White stone at the top-left corner `(row=0, column=0)`, the command is: `A-01-O`.

**5. Game State and History (`gomoku.md`)**
*   The `gomoku.md` file will persist the game state, including a move log and the current board.
*   A `snap id` system is used to track the history of board states. The initial empty board is the `root snap id`, and the `current snap id` always points to the latest board state.

## Game flow

You are the Black player, you move second.

On your turn, use the Pantheon system to simulate a lookahead. The process is as follows:

1.  **Perform a first-level exploration (playing as the Black player)** - Use the `current snap id` from the previous round as the `source snap id`. Call `parallel explore` to launch 3 concurrent explorations. Set the `shared prompt` to: ["Update the board in gomoku.md with the White player's move ({White player's move command}).", "Then, as the Black player, reason based on the latest board state, decide your move, and update the board in gomoku.md"].

2.  **Perform a second-level exploration (playing as the White player)** - Use `preview` to check the exploration status. to wait the first-level explorations (where Black makes a move) are complete, you can poll at 20-second intervals. For each completed exploration from the first level, use its `latest snap id` to once again call `parallel explore` and launch 3 new concurrent explorations. Set the `shared prompt` to: ["As the White player, reason based on the latest board state, decide the position for your White stone, and update the board in gomoku.md"].

After the two-level exploration described above is complete, you should read the board from each of the second-level explorations, select the board that is most advantageous to you, and then trace back to read (using `read_snap_file`) the `gomoku.md` from the corresponding first-level exploration to get the position of your next move. Print the board after making the move.

Then, wait for the White player to make a decision and continue to the next round.
