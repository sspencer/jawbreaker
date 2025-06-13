package jawbreaker

import "core:math"
import "core:math/rand"
import "core:mem"
import rl "vendor:raylib"
import fmt "core:fmt"

get_board_coords :: proc(mouse: rl.Vector2) -> Vec2i {
    grid_x := mouse.x - SCREEN_PADDING - BOARD_PADDING
    grid_y := mouse.y - SCREEN_PADDING - BOARD_PADDING

    // Check if within board bounds
    if grid_x < 0 ||
    grid_x >= BLOCK_SIZE * NUM_BLOCKS ||
    grid_y < 0 ||
    grid_y >= BLOCK_SIZE * NUM_BLOCKS {
        return Vec2i{ -1, -1 }
    }

    board_x := grid_x / BLOCK_SIZE
    board_y := grid_y / BLOCK_SIZE

    if board_x >= NUM_BLOCKS || board_y >= NUM_BLOCKS {
        return Vec2i{ -1, -1 }
    }

    return { int(math.floor(board_x)), int(math.floor(board_y)) }
}

selected_blocks :: proc(start: Vec2i, allocator := context.temp_allocator) -> int {
    // reset connected_pieces to false
    mem.zero(&connected_pieces, size_of(connected_pieces));

    if start.x < 0 || start.x >= NUM_BLOCKS || start.y < 0 || start.y >= NUM_BLOCKS {
        return 0
    }

    target := board[start.x][start.y]
    if target == .Empty {
        return 0
    }

    #partial switch (target) {
        case .PowerPlus, .PowerMinus, .PowerTimes, .PowerRect, .PowerCircle, .PowerPipe:
            return powerup_selection(start)
        case:
            return color_selection(start, allocator)
    }
}

color_selection :: proc(start: Vec2i, allocator := context.temp_allocator) -> int {
    target := board[start.x][start.y]
    visited := make(map[Vec2i]bool, allocator)
    defer delete(visited)

    stack := make([dynamic]Vec2i, 0, allocator)
    defer delete(stack)
    append(&stack, start)

    directions := [4]Vec2i{ { -1, 0 }, { 0, 1 }, { 1, 0 }, { 0, -1 } }
    count := 0

    for len(stack) > 0 {
        current := pop(&stack)
        if current in visited {
            continue
        }

        visited[current] = true
        connected_pieces[current.x][current.y] = true
        count += 1

        for dir in directions {
            next := Vec2i{ current.x + dir.x, current.y + dir.y }

            // Check bounds
            if next.x < 0 ||
            next.x >= NUM_BLOCKS ||
            next.y < 0 ||
            next.y >= NUM_BLOCKS {
                continue
            }

            if board[next.x][next.y] == target && !(next in visited) {
                append(&stack, next)
            }
        }
    }

    return count
}

powerup_selection :: proc(start: Vec2i) -> int {
    target := board[start.x][start.y]
    sel : []Selection
    #partial switch(target) {
    case .PowerPlus:
        return do_selection(start, PlusSelection, NUM_BLOCKS)
    case .PowerMinus:
        return do_selection(start, MinusSelection, NUM_BLOCKS)
    case .PowerPipe:
        return do_selection(start, PipeSelection, NUM_BLOCKS)
    case .PowerTimes:
        max := int(math.ceil(NUM_BLOCKS * math.sqrt_f32(2)))
        return do_selection(start, TimesSelection, max)
    case .PowerRect:
        return do_selection(start, RectSelection, 4)
    case .PowerCircle:
        return do_selection(start, RectSelection, 1) +
            do_selection(start, CircleSelection, 3)
    }

    return 0
}

do_selection :: proc(start: Vec2i, sel: []Selection, max_iters: int) -> int {
    num := 0

    startRow := start.y
    startCol := start.x

    for iter in 0..<max_iters {
        for s in sel {
            r := startRow + s.Start.y + s.Dir.y * iter
            c := startCol + s.Start.x + s.Dir.x * iter
            if r < 0 || c < 0 || r >= NUM_BLOCKS || c >= NUM_BLOCKS {
                continue
            }

            cur := board[c][r]
            if cur == .Empty || block_color_values[cur] == rl.DARKGRAY {
                continue
            }

            connected_pieces[c][r] = true
            num += 1
        }
    }

    if num > 0 {
        connected_pieces[start.x][start.y] = true
        num += 1
    }

    return num
}

applyGravity :: proc() {
// Pass 1: Move all non-zero values down within each column
    for col in 0 ..< NUM_BLOCKS {
        write_pos := NUM_BLOCKS - 1

        for row := NUM_BLOCKS - 1; row >= 0; row -= 1 {
            if board[col][row] != .Empty {
                if write_pos != row {
                    board[col][write_pos] = board[col][row]
                    board[col][row] = .Empty
                }
                write_pos -= 1
            }
        }
    }

    // Pass 2: Move entire columns right to fill gaps from empty columns
    write_col := NUM_BLOCKS - 1

    for col := NUM_BLOCKS - 1; col >= 0; col -= 1 {
    // Check if this column has any non-empty pieces
        column_has_pieces := false
        for row in 0 ..< NUM_BLOCKS {
            if board[col][row] != .Empty {
                column_has_pieces = true
                break
            }
        }

        // If column has pieces, move it to the write position
        if column_has_pieces {
            if write_col != col {
            // Move entire column
                for row in 0 ..< NUM_BLOCKS {
                    board[write_col][row] = board[col][row]
                    board[col][row] = .Empty
                }
            }
            write_col -= 1
        }
    }
}

is_game_over :: proc() -> bool {
    for c in 0 ..< NUM_BLOCKS {
        for r in 0 ..< NUM_BLOCKS {
            if board[c][r] != .Empty {
                n := selected_blocks(Vec2i{ c, r })
                if n > 1 {
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

random_block_color :: proc() -> Block_Color {
    valid_colors := []Block_Color{ .Purple, .Blue, .Green, .Red, .Yellow }

    random_index := rand.int_max(len(valid_colors))
    return valid_colors[random_index]
}

select_random_positions :: proc(count: int) -> []Vec2i {
    positions := make([]Vec2i, count, context.temp_allocator)
    used := make(map[Vec2i]bool)
    defer delete(used)

    for i in 0..<count {
        for {
            x := rand.int_max(NUM_BLOCKS)
            y := rand.int_max(NUM_BLOCKS)
            pos := Vec2i{x, y}

            if pos not_in used {
                used[pos] = true
                positions[i] = pos
                break
            }
        }
    }

    return positions
}

make_move :: proc(num_connected: int) {
    game_score += num_connected * (num_connected - 1)
    for c in 0 ..< NUM_BLOCKS {
        for r in 0 ..< NUM_BLOCKS {
            if connected_pieces[c][r] {
                board[c][r] = .Empty
            }
        }
    }

    applyGravity()

    game_over = is_game_over()

    if game_over {
        game_bonus = calculate_bonus()
        game_score += game_bonus

        last_score = game_score
        if game_score > best_score {
            best_score = game_score
        }

        state := Game_State{
            best_score = best_score,
            last_score = last_score
        }

        save_high_score(state)
    }

}