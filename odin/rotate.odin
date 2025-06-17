package jawbreaker


import "core:fmt"

// Rotation direction enum
Rotation_Direction :: enum {
    CLOCKWISE,
    COUNTER_CLOCKWISE,
}

// Method 2: In-place rotation (more memory efficient)
rotate_array :: proc(direction: Rotation_Direction) {
    switch direction {
    case .CLOCKWISE:
        rotate_clockwise_inplace()
    case .COUNTER_CLOCKWISE:
        rotate_counter_clockwise_inplace()
    }
}

// Helper: Clockwise in-place rotation
rotate_clockwise_inplace :: proc() {
// Transpose the matrix
    for i in 0..<NUM_BLOCKS {
        for j in i+1..<NUM_BLOCKS {
            board[i][j], board[j][i] = board[j][i], board[i][j]
        }
    }

    // Reverse each row
    for i in 0..<NUM_BLOCKS {
        for j in 0..<NUM_BLOCKS/2 {
            board[i][j], board[i][NUM_BLOCKS-1-j] = board[i][NUM_BLOCKS-1-j], board[i][j]
        }
    }
}

// Helper: Counter-clockwise in-place rotation
rotate_counter_clockwise_inplace :: proc() {
// Reverse each row first
    for i in 0..<NUM_BLOCKS {
        for j in 0..<NUM_BLOCKS/2 {
            board[i][j], board[i][NUM_BLOCKS-1-j] = board[i][NUM_BLOCKS-1-j], board[i][j]
        }
    }

    // Then transpose the matrix
    for i in 0..<NUM_BLOCKS {
        for j in i+1..<NUM_BLOCKS {
            board[i][j], board[j][i] = board[j][i], board[i][j]
        }
    }
}
