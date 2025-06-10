package jawbreaker

import "core:fmt"
import rl "vendor:raylib"

draw_game_piece_colors :: proc() {
    for row in 0 ..< NUM_BLOCKS {
        for col in 0 ..< NUM_BLOCKS {
            rec := rl.Rectangle {
                f32(col * BLOCK_SIZE + SCREEN_PADDING + BOARD_PADDING),
                f32(row * BLOCK_SIZE + SCREEN_PADDING + BOARD_PADDING),
                BLOCK_SIZE - BLOCK_PADDING,
                BLOCK_SIZE - BLOCK_PADDING,
            }

            rl.DrawRectangleRounded(rec, 0.3, 16, block_color_values[board[col][row]])
        }
    }
}

draw_game_piece_highlights :: proc () {
    for row in 0 ..< NUM_BLOCKS {
        for col in 0 ..< NUM_BLOCKS {
            rec := rl.Rectangle {
                f32(col * BLOCK_SIZE + SCREEN_PADDING + BOARD_PADDING),
                f32(row * BLOCK_SIZE + SCREEN_PADDING + BOARD_PADDING),
                BLOCK_SIZE - BLOCK_PADDING,
                BLOCK_SIZE - BLOCK_PADDING,
            }

            if connected_pieces[col][row] {
                color := highlight_color_values[board[col][row]]
                rl.DrawRectangleRoundedLinesEx(rec, 0.3, 16, 2, color)
            }
        }
    }
}

draw_score :: proc(ss: i32, font_size: i32 = 11) {
    score_str := fmt.ctprintf("Score: %v", game_score)
    score_width := rl.MeasureText(score_str, font_size)
    rl.DrawText(score_str,
    ss / 2 - score_width / 2,
    ss - 18,
    font_size,
    rl.WHITE,
    )
}

draw_game_over_score :: proc(ss: i32, font_size: i32 = 15) {
    score_over_text: cstring
    if game_bonus == 0 {
        score_over_text = fmt.ctprintf("Score: %v", game_score)
    } else {
        score_over_text = fmt.ctprintf("Score: %v (Bonus: %v)", game_score, game_bonus)
    }

    score_over_width := rl.MeasureText(score_over_text, font_size)
    rl.DrawText(score_over_text,
    ss / 2 - score_over_width / 2,
    ss / 2 - 35,
    font_size,
    rl.WHITE)
}

draw_game_over_text :: proc(ss: i32, font_size: i32 = 15) {
    game_over_text := fmt.ctprint("Restart: SPACE")
    game_over_width := rl.MeasureText(game_over_text, font_size)
    rl.DrawText(game_over_text,
    ss / 2 - game_over_width / 2,
    ss / 2 - 15,
    font_size,
    rl.WHITE)
}