package jawbreaker

import rl "vendor:raylib"

Vec2i :: [2]int

NUM_BLOCKS :: 14
BLOCK_SIZE :: 23
SCREEN_PADDING :: 24
BOARD_PADDING :: 3
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
connected_pieces: [NUM_BLOCKS][NUM_BLOCKS]bool

game_score := 0
game_over := false
game_bonus := 0

main :: proc() {
    rl.SetConfigFlags({ .VSYNC_HINT })
    screen_size := BLOCK_SIZE * NUM_BLOCKS + (SCREEN_PADDING + BOARD_PADDING) * 2

    rl.InitWindow(i32(screen_size * 2), i32(screen_size * 2), "Jawbreaker")
    rl.SetTargetFPS(60)
    reset()

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

        board_rec := rl.Rectangle {
            SCREEN_PADDING,
            SCREEN_PADDING,
            board_size,
            board_size,
        }

        // draw game board
        rl.DrawRectangleRounded(board_rec, 0.02, 16, { 20, 24, 33, 255 })


        mouse_pos := get_board_coords(rl.GetMousePosition() / camera_zoom)
        num_connected := get_connected_pieces(mouse_pos)

        made_move := false
        if rl.IsMouseButtonPressed(.LEFT) && num_connected > 1 {
            make_move(num_connected)
            made_move = true
        }

        draw_game_piece_colors()

        if !made_move && num_connected > 1 {
            draw_game_piece_highlights()
        }

        draw_score(i32(screen_size))

        if game_over {
            draw_game_over_score(i32(screen_size))
            draw_game_over_text(i32(screen_size))
        }

        rl.EndMode2D()
        rl.EndDrawing()

        free_all(context.temp_allocator)
    }

    rl.CloseWindow()
}
