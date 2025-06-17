package jawbreaker

import rl "vendor:raylib"

draw_plus :: proc(rect: rl.Rectangle) {
    center_x := rect.x + rect.width / 2
    center_y := rect.y + rect.height / 2

    // Calculate dimensions for the plus shape
    thickness := rect.width / 8  // Thickness of the lines
    length := rect.width * 3 / 5  // Length of each arm

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

    rl.DrawRectangleRec(vertical_rect, rl.WHITE)
    rl.DrawRectangleRec(horizontal_rect, rl.WHITE)
}

draw_pipe :: proc(rect: rl.Rectangle) {
    center_x := rect.x + rect.width / 2
    center_y := rect.y + rect.height / 2

    // Calculate dimensions for the plus shape
    thickness := rect.width / 8  // Thickness of the lines
    length := rect.width * 3 / 5  // Length of each arm

    // Vertical bar
    vertical_rect := rl.Rectangle{
        x = center_x - thickness / 2,
        y = center_y - length / 2,
        width = thickness,
        height = length,
    }

    rl.DrawRectangleRec(vertical_rect, rl.WHITE)
}

draw_minus :: proc(rect: rl.Rectangle) {
    center_x := rect.x + rect.width / 2
    center_y := rect.y + rect.height / 2

    // Calculate dimensions for the minus shape
    thickness := rect.width / 8  // Thickness of the line
    length := rect.width * 3 / 5  // Length of the line

    // Horizontal bar
    minus_rect := rl.Rectangle{
        x = center_x - length / 2,
        y = center_y - thickness / 2,
        width = length,
        height = thickness,
    }

    // Draw white fill
    rl.DrawRectangleRec(minus_rect, rl.WHITE)
}

draw_times :: proc(rect: rl.Rectangle) {
    center_x := rect.x + rect.width / 2
    center_y := rect.y + rect.height / 2

    // Calculate dimensions for the X shape
    thickness := rect.width / 6  // Thickness of the lines
    arm_length := rect.width * 0.55
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

draw_rectangle :: proc(rect: rl.Rectangle) {
    thickness := rect.width / 8

    x := rect.x + thickness + 1
    y := rect.y + thickness + 1

    width := rect.width - 2*thickness -2
    height := rect.height - 2*thickness - 2

    rl.DrawRectangleLinesEx(rl.Rectangle{f32(x), f32(y), f32(width), f32(height)}, thickness, rl.WHITE)
}

draw_circle :: proc(rect: rl.Rectangle) {
    center_x := rect.x + rect.width / 2
    center_y := rect.y + rect.height / 2
    radius := min(rect.width, rect.height) * 0.3

    rl.DrawRing({center_x, center_y}, radius - 1, radius + 1.5, 0, 360, 24, rl.WHITE)
}

