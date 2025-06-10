package jawbreaker

import "core:fmt"
import "core:math"
import "core:math/rand"
import "core:mem"
import rl "vendor:raylib"

Vec2i :: [2]int

NUM_BLOCKS :: 8
BLOCK_SIZE :: 25
DISPLAY_PADDING :: 24
BOARD_PADDING :: 4
BLOCK_PADDING :: 2

Block_Color :: enum {
    Empty,
    Purple,
    Blue,
    Green,
    Red,
    Yellow,
}

block_color_values := [Block_Color]rl.Color {
    .Empty  = { 35, 38, 44, 255 },
    .Purple = rl.DARKPURPLE,
    .Blue   = rl.BLUE,
    .Green  = rl.LIME,
    .Red    = rl.RED,
    .Yellow = rl.GOLD,
}

highlight_color_values := [Block_Color]rl.Color {
    .Empty  = { 105, 114, 132, 255 },
    .Purple = rl.PURPLE,
    .Blue   = rl.SKYBLUE,
    .Green  = rl.GREEN,
    .Red    = rl.PINK,
    .Yellow = rl.YELLOW,
}

board: [NUM_BLOCKS][NUM_BLOCKS]Block_Color // columns x rows
game_score := 0
game_over := false
game_bonus := 0

get_board_coords :: proc(mouse: rl.Vector2) -> Vec2i {
// Adjust for screen padding and board border
    grid_x := mouse.x - DISPLAY_PADDING - BOARD_PADDING
    grid_y := mouse.y - DISPLAY_PADDING - BOARD_PADDING

    // Check if within board bounds
    if grid_x < 0 ||
    grid_x >= BLOCK_SIZE * NUM_BLOCKS ||
    grid_y < 0 ||
    grid_y >= BLOCK_SIZE * NUM_BLOCKS {
        return Vec2i{ -1, -1 }
    }

    // Compute board coordinates
    board_x := grid_x / BLOCK_SIZE
    board_y := grid_y / BLOCK_SIZE

    // Validate board coordinates
    if board_x >= NUM_BLOCKS || board_y >= NUM_BLOCKS {
        return Vec2i{ -1, -1 }
    }

    return Vec2i{ int(math.floor(board_x)), int(math.floor(board_y)) }
}

get_connected_pieces :: proc(start: Vec2i, allocator := context.temp_allocator) -> map[Vec2i]bool {
// Validate input
    if start.x < 0 ||
    start.x >= NUM_BLOCKS ||
    start.y < 0 ||
    start.y >= NUM_BLOCKS {
        return nil
    }

    // Get target value to match
    target_value := board[start.x][start.y]
    if target_value == .Empty {
        return nil
    }

    // Initialize result array and visited set
    result := make(map[Vec2i]bool, 0, allocator)
    visited := make(map[Vec2i]bool, allocator)
    defer delete(visited)

    // Stack for DFS
    stack := make([dynamic]Vec2i, 0, allocator)
    defer delete(stack)
    append(&stack, start)

    // Possible directions (up, right, down, left)
    directions := [4]Vec2i{ { -1, 0 }, { 0, 1 }, { 1, 0 }, { 0, -1 } }

    for len(stack) > 0 {
    // Pop current point
        current := pop(&stack)
        if current in visited {
            continue
        }

        // Mark as visited and add to result
        visited[current] = true
        result[current] = true

        // Check all four directions
        for dir in directions {
            next := Vec2i{ current.x + dir.x, current.y + dir.y }

            // Check bounds
            if next.x < 0 ||
            next.x >= NUM_BLOCKS ||
            next.y < 0 ||
            next.y >= NUM_BLOCKS {
                continue
            }

            // Check if value matches and not visited
            if board[next.x][next.y] == target_value && !(next in visited) {
                append(&stack, next)
            }
        }
    }

    return result
}

is_game_over :: proc() -> bool {
    for c in 0 ..< NUM_BLOCKS {
        for r in 0 ..< NUM_BLOCKS {
            if board[c][r] != .Empty {
                pieces := get_connected_pieces(Vec2i{ c, r })
                if len(pieces) > 1 {
                    return false
                }
            }
        }
    }

    return true
}

calculate_bonus :: proc() -> int {
    threshold :: 10 // Bonus if 10 or fewer pieces remain
    pieces := 0
    for c in 0 ..< NUM_BLOCKS {
        for r in 0 ..< NUM_BLOCKS {
            if board[r][c] != .Empty {
                pieces += 1
            }
        }
    }

    if pieces == 0 {
        return 2000
    } else if pieces <= threshold {
        return (threshold - pieces + 1) * 100
    }

    return 0
}

applyGravity :: proc() {
// Apply gravity in two passes: first downward, then rightward

// Pass 1: Move all non-zero values down within each column
    for col in 0 ..< NUM_BLOCKS {
    // Start from the bottom and work upward
        write_pos := NUM_BLOCKS - 1

        // Scan from bottom to top
        for row := NUM_BLOCKS - 1; row >= 0; row -= 1 {
            if board[col][row] != .Empty {
            // If we found a non-zero value, move it to the write position
                if write_pos != row {
                    board[col][write_pos] = board[col][row]
                    board[col][row] = .Empty
                }
                write_pos -= 1
            }
        }
    }

    // Pass 2: Move all non-zero values right within each row
    for row in 0 ..< NUM_BLOCKS {
    // Start from the right and work leftward
        write_pos := NUM_BLOCKS - 1

        // Scan from right to left
        for col := NUM_BLOCKS - 1; col >= 0; col -= 1 {
            if board[col][row] != .Empty {
            // If we found a non-zero value, move it to the write position
                if write_pos != col {
                    board[write_pos][row] = board[col][row]
                    board[col][row] = .Empty
                }
                write_pos -= 1
            }
        }
    }
}

get_random_block_color :: proc() -> Block_Color {
    valid_colors := []Block_Color{ .Purple, .Blue, .Green, .Red, .Yellow }

    random_index := rand.int_max(len(valid_colors))
    return valid_colors[random_index]
}

reset :: proc() {
    for row in 0 ..< NUM_BLOCKS {
        for col in 0 ..< NUM_BLOCKS {
            board[col][row] = get_random_block_color()
        }
    }

    game_score = 0
    game_over = false
    game_bonus = 0
}

main :: proc() {
    rl.SetConfigFlags({ .VSYNC_HINT })
    screen_size :=
    BLOCK_SIZE * NUM_BLOCKS + (DISPLAY_PADDING + BOARD_PADDING) * 2

    rl.InitWindow(i32(screen_size * 2), i32(screen_size * 2), "Jawbreaker")
    rl.SetTargetFPS(60)
    reset()
    mouse_pressed: bool
    mouse_pressed_pos: Vec2i

    for !rl.WindowShouldClose() {
        rl.BeginDrawing()
        rl.ClearBackground({ 43, 60, 80, 255 })

        if game_over {
            if rl.IsKeyPressed(.SPACE) {
                reset()
            }
        }

        camera_zoom := f32(rl.GetScreenHeight()) / f32(screen_size)
        camera := rl.Camera2D {
            zoom = camera_zoom,
        }

        rl.BeginMode2D(camera)

        board_size := f32(NUM_BLOCKS) * BLOCK_SIZE + BOARD_PADDING * 2
        //rl.DrawRectangle(DISPLAY_PADDING, DISPLAY_PADDING, board_size, board_size, {20, 24, 33, 255})
        board_rec := rl.Rectangle {
            DISPLAY_PADDING,
            DISPLAY_PADDING,
            board_size,
            board_size,
        }
        rl.DrawRectangleRounded(board_rec, 0.02, 16, { 20, 24, 33, 255 })

        selection := get_connected_pieces(get_board_coords(rl.GetMousePosition() / camera_zoom))

        if rl.IsMouseButtonPressed(.LEFT) && len(selection) > 1 {
            n := len(selection)
            game_score += n * n - 1
            for p in selection {
                board[p.x][p.y] = .Empty
            }

            applyGravity()

            selection = nil

            game_over = is_game_over()

            if game_over {
                game_bonus = calculate_bonus()
                game_score += game_bonus
            }
        }

        for row in 0 ..< NUM_BLOCKS {
            for col in 0 ..< NUM_BLOCKS {
                rec := rl.Rectangle {
                    f32(col * BLOCK_SIZE + DISPLAY_PADDING + BOARD_PADDING),
                    f32(row * BLOCK_SIZE + DISPLAY_PADDING + BOARD_PADDING),
                    BLOCK_SIZE - BLOCK_PADDING,
                    BLOCK_SIZE - BLOCK_PADDING,
                }

                rl.DrawRectangleRounded(rec, 0.3, 16, block_color_values[board[col][row]])
            }
        }

        for row in 0 ..< NUM_BLOCKS {
            for col in 0 ..< NUM_BLOCKS {
                rec := rl.Rectangle {
                    f32(col * BLOCK_SIZE + DISPLAY_PADDING + BOARD_PADDING),
                    f32(row * BLOCK_SIZE + DISPLAY_PADDING + BOARD_PADDING),
                    BLOCK_SIZE - BLOCK_PADDING,
                    BLOCK_SIZE - BLOCK_PADDING,
                }

                point := Vec2i{ col, row }
                if point in selection {
                    rl.DrawRectangleRoundedLinesEx(rec, 0.3, 16, 2, highlight_color_values[board[col][row]])
                }
            }
        }

        score_str := fmt.ctprintf("Score: %v", game_score)
        font_size : i32 = 11
        score_width := rl.MeasureText(score_str, font_size)
        ss := i32(screen_size)
        rl.DrawText(
        score_str,
        ss / 2 - score_width / 2,
        ss - 18,
        font_size,
        rl.WHITE,
        )

        if game_over {
            if game_bonus > 0 {

            }

            font_size = 15
            score_over_text: cstring
            if game_bonus == 0 {
                score_over_text = fmt.ctprintf("Score: %v", game_score)
            } else {
                score_over_text = fmt.ctprintf("Score: %v (Bonus: %v)", game_score, game_bonus)
            }

            score_over_width := rl.MeasureText(score_over_text, font_size)
            rl.DrawText(score_over_text, ss / 2 - score_over_width / 2, ss / 2 - 35, font_size, rl.WHITE)


            game_over_text := fmt.ctprint("Restart: SPACE")
            game_over_width := rl.MeasureText(game_over_text, font_size)
            rl.DrawText(game_over_text, ss / 2 - game_over_width / 2, ss / 2 - 15, font_size, rl.WHITE)
        }

        rl.EndMode2D()
        rl.EndDrawing()

        free_all(context.temp_allocator)
    }

    rl.CloseWindow()

}
