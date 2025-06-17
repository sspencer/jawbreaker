package jawbreaker

import "core:fmt"
import rl "vendor:raylib"

draw_score :: proc(ss: i32, font_size: i32 = 11) {
    score_str := fmt.ctprintf("Score: %d  /  Last: %d  / Best: %d", game_score, last_score, best_score)
    score_width := rl.MeasureText(score_str, font_size)
    rl.DrawText(score_str,
    ss / 2 - score_width / 2,
    12,
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