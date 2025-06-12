package jawbreaker

import "core:fmt"
import rl "vendor:raylib"

draw_block :: proc() {
    for row in 0 ..< NUM_BLOCKS {
        for col in 0 ..< NUM_BLOCKS {
            rec := rl.Rectangle {
                f32(col * BLOCK_SIZE + SCREEN_PADDING + BOARD_PADDING),
                f32(row * BLOCK_SIZE + SCREEN_PADDING + BOARD_PADDING),
                BLOCK_SIZE - BLOCK_PADDING,
                BLOCK_SIZE - BLOCK_PADDING,
            }

            rl.DrawRectangleRounded(rec, 0.3, 16, block_color_values[board[col][row]])

            block_type := board[col][row]
            #partial switch(board[col][row]) {
            case .PowerPlus: draw_plus(rec)
            case .PowerMinus: draw_minus(rec)
            case .PowerTimes: draw_times(rec)
            }
        }
    }
}

draw_highlight :: proc () {
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

draw_plus :: proc(rect: rl.Rectangle) {
    center_x := rect.x + rect.width / 2
    center_y := rect.y + rect.height / 2

    // Calculate dimensions for the plus shape
    thickness := rect.width / 6  // Thickness of the lines
    length := rect.width * 3 / 4  // Length of each arm

    // Vertical bar
    vertical_rect := rl.Rectangle{
        x = center_x - thickness / 2,
        y = center_y - length / 2,
        width = thickness,
        height = length,
    }

    // Horizontal bar
    horizontal_rect := rl.Rectangle{
        x = center_x - length / 2,
        y = center_y - thickness / 2,
        width = length,
        height = thickness,
    }

    // Draw white fill
    rl.DrawRectangleRec(vertical_rect, rl.WHITE)
    rl.DrawRectangleRec(horizontal_rect, rl.WHITE)
//
//    // Draw black borders
//    rl.DrawRectangleLinesEx(vertical_rect, 2, rl.BLACK)
//    rl.DrawRectangleLinesEx(horizontal_rect, 2, rl.BLACK)
}

draw_minus :: proc(rect: rl.Rectangle) {
    center_x := rect.x + rect.width / 2
    center_y := rect.y + rect.height / 2

    // Calculate dimensions for the minus shape
    thickness := rect.width / 6  // Thickness of the line
    length := rect.width * 3 / 4  // Length of the line

    // Horizontal bar
    minus_rect := rl.Rectangle{
        x = center_x - length / 2,
        y = center_y - thickness / 2,
        width = length,
        height = thickness,
    }

    // Draw white fill
    rl.DrawRectangleRec(minus_rect, rl.WHITE)

//    // Draw black border
//    rl.DrawRectangleLinesEx(minus_rect, 2, rl.BLACK)
}

draw_times :: proc(rect: rl.Rectangle) {
    center_x := rect.x + rect.width / 2
    center_y := rect.y + rect.height / 2

    // Calculate dimensions for the X shape
    thickness := rect.width / 8  // Thickness of the lines
    arm_length := rect.width * 3 / 4  // Length of each diagonal arm
    half_arm := arm_length / 2

    // Calculate diagonal offset for thickness
    diag_offset := thickness / 2

    // Draw the X as two diagonal rectangles
    // First diagonal (top-left to bottom-right)
    diagonal1_points := [4]rl.Vector2{
        {center_x - half_arm, center_y - half_arm - diag_offset},
        {center_x - half_arm + diag_offset, center_y - half_arm},
        {center_x + half_arm, center_y + half_arm + diag_offset},
        {center_x + half_arm - diag_offset, center_y + half_arm},
    }

    // Second diagonal (top-right to bottom-left)
    diagonal2_points := [4]rl.Vector2{
        {center_x + half_arm, center_y - half_arm - diag_offset},
        {center_x + half_arm - diag_offset, center_y - half_arm},
        {center_x - half_arm, center_y + half_arm + diag_offset},
        {center_x - half_arm + diag_offset, center_y + half_arm},
    }

    // Draw black borders
    for i in 0..<4 {
        next_i := (i + 1) % 4
        rl.DrawLineEx(diagonal1_points[i], diagonal1_points[next_i], 2, rl.WHITE)
        rl.DrawLineEx(diagonal2_points[i], diagonal2_points[next_i], 2, rl.WHITE)
    }
}

draw_score :: proc(ss: i32, font_size: i32 = 11) {
    score_str := fmt.ctprintf("Score: %d  /  Last: %d  / Best: %d", game_score, last_score, best_score)
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