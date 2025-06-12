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
    PowerPlus,
    PowerMinus,
    PowerTimes,
    PowerRect,
    PowerCircle,
}


block_color_values := [Block_Color]rl.Color {
    .Empty       = { 35, 38, 44, 255 },
    .Purple      = rl.DARKPURPLE,
    .Blue        = rl.BLUE,
    .Green       = rl.LIME,
    .Red         = rl.RED,
    .Yellow      = rl.GOLD,
    .PowerPlus   = rl.DARKGRAY,
    .PowerMinus  = rl.DARKGRAY,
    .PowerTimes  = rl.DARKGRAY,
    .PowerRect   = rl.DARKGRAY,
    .PowerCircle = rl.DARKGRAY,
}

highlight_color_values := [Block_Color]rl.Color {
    .Empty       = { 105, 114, 132, 255 },
    .Purple      = rl.PURPLE,
    .Blue        = rl.SKYBLUE,
    .Green       = rl.GREEN,
    .Red         = rl.PINK,
    .Yellow      = rl.YELLOW,
    .PowerPlus   = rl.GRAY,
    .PowerMinus  = rl.GRAY,
    .PowerTimes  = rl.GRAY,
    .PowerRect   = rl.GRAY,
    .PowerCircle = rl.GRAY,
}

Selection :: struct{
    Start: Vec2i,
    Dir: Vec2i,
}

TimesSelection := []Selection {
    {{-1, -1}, {-1, -1}},
    {{ 1, -1}, { 1, -1}},
    {{-1,  1}, {-1,  1}},
    {{ 1,  1}, { 1,  1}},
}

PlusSelection := []Selection{
    {{ 0, -1 }, { 0, -1}},
    {{ 0,  1},  { 0,  1}},
    {{-1,  0},  {-1,  0}},
    {{ 1,  0},  { 1,  0}},
}

MinusSelection := []Selection{
    {{-1,  0},  {-1,  0}},
    {{ 1,  0},  { 1,  0}},
}

RectSelection := []Selection {
    {{-2, -2}, { 0,  1}},
    {{-2,  2}, { 1,  0}},
    {{ 2,  2}, { 0, -1}},
    {{ 2, -2}, {-1,  0}},
}

CircleSelection := []Selection {
    {{-3, -1}, { 0,  1}},
    {{-1,  3}, { 1,  0}},
    {{ 3,  1}, { 0, -1}},
    {{ 1, -3}, {-1,  0}},
}


board: [NUM_BLOCKS][NUM_BLOCKS]Block_Color // columns x rows
undo_board: [NUM_BLOCKS][NUM_BLOCKS]Block_Color // columns x rows
connected_pieces: [NUM_BLOCKS][NUM_BLOCKS]bool

game_score := 0
last_score := 0
best_score := 0
undo_score := 0
can_undo := false
game_over := false
game_bonus := 0


restart :: proc() {
    for row in 0 ..< NUM_BLOCKS {
        for col in 0 ..< NUM_BLOCKS {
            board[col][row] = random_block_color()
        }
    }

    // place power ups
    power_ups := []Block_Color{.PowerPlus, .PowerMinus, .PowerTimes, .PowerRect, .PowerCircle}

    pos := select_random_positions(len(power_ups))
    for p, i in pos {
        board[p.x][p.y] = power_ups[i]
    }

    last_score = game_score
    if game_score > best_score {
        best_score = game_score
    }
    game_score = 0
    game_over = false
    game_bonus = 0
    can_undo = false
}

main :: proc() {
    rl.SetConfigFlags({ .VSYNC_HINT })
    screen_size := BLOCK_SIZE * NUM_BLOCKS + (SCREEN_PADDING + BOARD_PADDING) * 2

    rl.InitWindow(i32(screen_size * 2), i32(screen_size * 2), "Jawbreaker")
    rl.SetTargetFPS(60)
    restart()

    for !rl.WindowShouldClose() {
        rl.BeginDrawing()
        rl.ClearBackground({ 43, 60, 80, 255 })

        if game_over {
            if rl.IsKeyPressed(.SPACE) || rl.IsMouseButtonPressed(.LEFT) {
                restart()
            }
        } else {
            if can_undo && (rl.IsKeyPressed(.Z) || rl.IsKeyPressed(.U)) {
                board = undo_board
                game_score = undo_score
                can_undo = false
            }

            if rl.IsKeyPressed(.R) {
                restart()
            }
        }

        camera_zoom := f32(rl.GetScreenHeight()) / f32(screen_size)
        camera := rl.Camera2D {
            zoom = camera_zoom,
        }

        rl.BeginMode2D(camera)

        board_size := f32(NUM_BLOCKS) * BLOCK_SIZE + BOARD_PADDING * 2

        board_rec := rl.Rectangle{ SCREEN_PADDING, SCREEN_PADDING, board_size, board_size }


        mouse_pos := get_board_coords(rl.GetMousePosition() / camera_zoom)
        num_connected := get_connected_pieces(mouse_pos)

        made_move := false
        if rl.IsMouseButtonPressed(.LEFT) && num_connected > 1 {
            undo_score = game_score
            undo_board = board
            make_move(num_connected)
            made_move = true
            can_undo = true
        }

        // draw game board
        rl.DrawRectangleRounded(board_rec, 0.02, 16, { 20, 24, 33, 255 })

        draw_block()

        if !made_move && num_connected > 1 {
            draw_highlight()
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
