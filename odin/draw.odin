package jawbreaker

import rl "vendor:raylib"

draw_block :: proc() {
    for row in 0 ..< NUM_BLOCKS {
        for col in 0 ..< NUM_BLOCKS {
            rec := rl.Rectangle {
                f32(col * BLOCK_SIZE + WINDOW_PADDING + GRID_OUTER),
                f32(row * BLOCK_SIZE + WINDOW_PADDING + GRID_OUTER),
                BLOCK_SIZE - GRID_INNER,
                BLOCK_SIZE - GRID_INNER,
            }

            lo := block_color_values[board[col][row]]
            hi := highlight_color_values[board[col][row]]

            if board[col][row] == .Empty {
                rl.DrawRectangleRounded(rec, 0.3, 16, hi)
            } else {
                rl.DrawRectangleGradientEx(rec, hi, lo, lo, lo)
            }
            block_type := board[col][row]
            #partial switch(board[col][row]) {
            case .PowerPlus:   draw_plus(rec)
            case .PowerMinus:  draw_minus(rec)
            case .PowerPipe:   draw_pipe(rec)
            case .PowerTimes:  draw_times(rec)
            case .PowerRect:   draw_rectangle(rec)
            case .PowerCircle: draw_circle(rec)
            }
        }
    }
}

draw_highlight :: proc () {
    for row in 0 ..< NUM_BLOCKS {
        for col in 0 ..< NUM_BLOCKS {
            rec := rl.Rectangle {
                f32(col * BLOCK_SIZE + WINDOW_PADDING + GRID_OUTER),
                f32(row * BLOCK_SIZE + WINDOW_PADDING + GRID_OUTER),
                BLOCK_SIZE - GRID_INNER,
                BLOCK_SIZE - GRID_INNER,
            }

            if connected_pieces[col][row] {
                color := highlight_color_values[board[col][row]]
                rl.DrawRectangleRoundedLinesEx(rec, 0.3, 16, 2, color)
            }
        }
    }
}

