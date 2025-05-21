class GameConfig {
    static TOKEN_SPACE = 1000;
    static COLORS = {
        WHITE: 0,
        PURPLE: 1000,
        BLUE: 2000,
        GREEN: 3000,
        RED: 4000,
        YELLOW: 5000,
        GRAY: 6000,
        BLACK: 7000,
    };
    static GAME_PIECES = [
        this.COLORS.PURPLE,
        this.COLORS.BLUE,
        this.COLORS.GREEN,
        this.COLORS.RED,
        this.COLORS.YELLOW,
    ];
    static COLOR_MAP = new Map([
        [this.COLORS.PURPLE, "#8a2be2"],
        [this.COLORS.BLUE, "#00a0ff"],
        [this.COLORS.GREEN, "#00cc66"],
        [this.COLORS.RED, "#ff3333"],
        [this.COLORS.YELLOW, "#ffcc00"],
        [this.COLORS.GRAY, "#888888"],
        [this.COLORS.BLACK, "#333333"],
        [this.COLORS.WHITE, "#ffffff"],
    ]);
    static GRADIENT_MAP = new Map([
        [this.COLORS.PURPLE, ["#d442ff", "#8a2be2"]],
        [this.COLORS.BLUE, ["#00f0ff", "#0070ff"]],
        [this.COLORS.GREEN, ["#00ff99", "#00cc66"]],
        [this.COLORS.RED, ["#ff5500", "#dd0000"]],
        [this.COLORS.YELLOW, ["#ffee00", "#ff8800"]],
        [this.COLORS.GRAY, ["#777777", "#cccccc"]],
    ]);
    static POWER_UPS = {
        X: 1,
        PLUS: 2,
        CIRCLE: 3,
        RECT: 4,
        FILL: 5,
        RIGHT: 6,
        LEFT: 7,
        DISC: 8,
        EXCHANGE: 9,
    };
    static POWER_CONFIG = new Map([
        [this.POWER_UPS.X, { connections: true, multi: true, rect: true }],
        [this.POWER_UPS.PLUS, { connections: true, multi: true, rect: true }],
        [this.POWER_UPS.CIRCLE, { connections: true, multi: true, rect: true }],
        [this.POWER_UPS.RECT, { connections: true, multi: true, rect: true }],
        [this.POWER_UPS.DISC, { connections: true, multi: false, rect: true }],
        [this.POWER_UPS.FILL, { connections: false, multi: false, rect: true }],
        [this.POWER_UPS.RIGHT, { connections: false, multi: false, rect: false }],
        [this.POWER_UPS.LEFT, { connections: false, multi: false, rect: false }],
        [this.POWER_UPS.EXCHANGE, { connections: false, multi: false, rect: true }],
    ]);
    static POWER_UP_DIRECTIONS = {
        X: [
            { sr: -1, sc: -1, dr: -1, dc: -1 },
            { sr: -1, sc: 1, dr: -1, dc: 1 },
            { sr: 1, sc: -1, dr: 1, dc: -1 },
            { sr: 1, sc: 1, dr: 1, dc: 1 },
        ],
        PLUS: [
            { sr: -1, sc: 0, dr: -1, dc: 0 },
            { sr: 1, sc: 0, dr: 1, dc: 0 },
            { sr: 0, sc: -1, dr: 0, dc: -1 },
            { sr: 0, sc: 1, dr: 0, dc: 1 },
        ],
        RECT: [
            { sr: -2, sc: -2, dc: 0, dr: 1 },
            { sr: 2, sc: -2, dc: 1, dr: 0 },
            { sr: 2, sc: 2, dc: 0, dr: -1 },
            { sr: -2, sc: 2, dc: -1, dr: 0 },
        ],
        CIRCLE: [
            { sr: -1, sc: -3, dc: 0, dr: 1 },
            { sr: 3, sc: -1, dc: 1, dr: 0 },
            { sr: 1, sc: 3, dc: 0, dr: -1 },
            { sr: -3, sc: 1, dc: -1, dr: 0 },
        ],
        DISC1: [
            { sr: -1, sc: -2, dc: 0, dr: 1 },
            { sr: 2, sc: -1, dc: 1, dr: 0 },
            { sr: 1, sc: 2, dc: 0, dr: -1 },
            { sr: -2, sc: 1, dc: -1, dr: 0 },
        ],
        DISC2: [
            { sr: -1, sc: -1, dc: 0, dr: 1 },
            { sr: 1, sc: -1, dc: 1, dr: 0 },
            { sr: 1, sc: 1, dc: 0, dr: -1 },
            { sr: -1, sc: 1, dc: -1, dr: 0 },
        ],
    };
    static GAP = 1;
    static BORDER_WIDTH = 4;
    static SHAPE_STROKE_COLOR = "white";
    static SHAPE_BORDER_COLOR = "black";
    static SHAPE_LINE_WIDTH = 3;
    static SHAPE_BORDER_LINE_WIDTH = 1;
    static ANIMATE_BOARD_SPEED = 6;
    static ANIMATE_PIECE_SPEED = 3;
    static ANIMATE = true;
}
