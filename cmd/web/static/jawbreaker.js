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
        [this.POWER_UPS.X, { connections: true, multi: true }],
        [this.POWER_UPS.PLUS, { connections: true, multi: true }],
        [this.POWER_UPS.CIRCLE, { connections: true, multi: true }],
        [this.POWER_UPS.RECT, { connections: true, multi: true }],
        [this.POWER_UPS.DISC, { connections: true, multi: false }],
        [this.POWER_UPS.FILL, { connections: false, multi: false }],
        [this.POWER_UPS.RIGHT, { connections: false, multi: false }],
        [this.POWER_UPS.LEFT, { connections: false, multi: false }],
        [this.POWER_UPS.EXCHANGE, { connections: false, multi: false }],
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

class CookieManager {
    static getCookie(name) {
        const cookies = document.cookie.split(";");
        for (const cookie of cookies) {
            const trimmed = cookie.trim();
            if (trimmed.startsWith(`${name}=`)) {
                return trimmed.substring(name.length + 1);
            }
        }
        return "";
    }

    static setCookie(name, value, days = 365) {
        const date = new Date();
        date.setTime(date.getTime() + days * 24 * 60 * 60 * 1000);
        document.cookie = `${name}=${value};expires=${date.toUTCString()};path=/`;
    }
}

class Board {
    constructor(rows, cols) {
        this.rows = rows;
        this.cols = cols;
        this.board = Array(rows * cols).fill(GameConfig.COLORS.WHITE);
    }

    getIndex(row, col) {
        return row * this.cols + col;
    }

    getPoint(index) {
        return {
            row: Math.floor(index / this.cols),
            col: index % this.cols,
        };
    }

    initialize() {
        const size = this.rows * this.cols;
        for (let i = 0; i < size; i++) {
            this.board[i] = this.randomPiece();
        }
        const indices = this.createShuffledIndices(size);
        GameConfig.POWER_CONFIG.forEach((opts, power) => {
            const colors = opts.multi
                ? [...GameConfig.GAME_PIECES, GameConfig.COLORS.GRAY]
                : [GameConfig.COLORS.GRAY];
            for (const color of colors) {
                this.board[indices.shift()] = color + power;
            }
        });
    }

    randomPiece() {
        return GameConfig.GAME_PIECES[
            Math.floor(Math.random() * GameConfig.GAME_PIECES.length)
            ];
    }

    createShuffledIndices(n) {
        const arr = Array.from({ length: n }, (_, i) => i);
        for (let i = arr.length - 1; i > 0; i--) {
            const j = Math.floor(Math.random() * (i + 1));
            [arr[i], arr[j]] = [arr[j], arr[i]];
        }
        return arr;
    }

    to2DArray() {
        const arr = [];
        for (let r = 0; r < this.rows; r++) {
            arr.push(this.board.slice(r * this.cols, (r + 1) * this.cols));
        }
        return arr;
    }

    from2DArray(arr) {
        this.board = arr.flat();
    }

    rotate(degrees) {
        const n = this.rows;
        if (!Number.isInteger(Math.sqrt(this.board.length))) {
            throw new Error("Board must be square for rotation.");
        }
        const rotation = ((degrees % 360) + 360) % 360;
        const grid = this.to2DArray();
        const rotated = Array.from({ length: n }, () => Array(n));
        if (rotation === 90) {
            for (let r = 0; r < n; r++) {
                for (let c = 0; c < n; c++) {
                    rotated[c][n - 1 - r] = grid[r][c];
                }
            }
        } else if (rotation === 270) {
            for (let r = 0; r < n; r++) {
                for (let c = 0; c < n; c++) {
                    rotated[n - 1 - c][r] = grid[r][c];
                }
            }
        } else {
            throw new Error("Rotation must be 90 or 270 degrees.");
        }
        this.from2DArray(rotated);
    }

    applyGravity() {
        const board = this.to2DArray();
        for (let c = 0; c < this.cols; c++) {
            let writeRow = this.rows - 1;
            for (let r = this.rows - 1; r >= 0; r--) {
                if (board[r][c] !== GameConfig.COLORS.WHITE) {
                    if (writeRow !== r) {
                        board[writeRow][c] = board[r][c];
                        board[r][c] = GameConfig.COLORS.WHITE;
                    }
                    writeRow--;
                }
            }
        }
        let writeCol = this.cols - 1;
        for (let c = this.cols - 1; c >= 0; c--) {
            const isEmpty = board.every((row) => row[c] === GameConfig.COLORS.WHITE);
            if (!isEmpty) {
                if (writeCol !== c) {
                    for (let r = 0; r < this.rows; r++) {
                        board[r][writeCol] = board[r][c];
                        board[r][c] = GameConfig.COLORS.WHITE;
                    }
                }
                writeCol--;
            }
        }
        this.from2DArray(board);
    }

    fillSpaces(connections, weighted = true) {
        if (!connections) {
            connections = this.board
            .map((piece, i) => (piece === GameConfig.COLORS.WHITE ? i : null))
            .filter((i) => i !== null);
        }
        if (weighted) {
            const counts = {};
            for (const piece of this.board) {
                const color = this.getColor(piece);
                if (
                    color !== GameConfig.COLORS.WHITE &&
                    color !== GameConfig.COLORS.GRAY
                ) {
                    counts[color] = (counts[color] || 0) + 1;
                }
            }
            const choices = Object.entries(counts).map(([value, count]) => ({
                value: parseInt(value),
                count,
            }));
            const normalized = this.normalizeGaps(choices, 0.33);
            for (const i of connections) {
                this.board[i] = this.weightedRandom(normalized);
            }
        } else {
            for (const i of connections) {
                this.board[i] = this.randomPiece();
            }
        }
    }

    fillPieces(connections) {
        if (!connections) {
            connections = this.board
            .map((piece, i) =>
                GameConfig.GAME_PIECES.includes(piece) ? i : null
            )
            .filter((i) => i !== null);
        }
        this.fillSpaces(connections, true);
    }

    weightedRandom(choices) {
        const total = choices.reduce((sum, obj) => sum + obj.count, 0);
        let r = Math.random() * total;
        for (const choice of choices) {
            if (r < choice.count) {
                return choice.value;
            }
            r -= choice.count;
        }
        return choices[0].value;
    }

    normalizeGaps(choices, maxGapPercent) {
        if (choices.length <= 1 || maxGapPercent < 0 || maxGapPercent > 1) {
            throw new Error("Invalid normalizeGaps input.");
        }
        const normalized = JSON.parse(JSON.stringify(choices));
        normalized.sort((a, b) => b.count - a.count);
        const ratios = [];
        for (let i = 0; i < normalized.length - 1; i++) {
            ratios.push(
                normalized[i].count > 0
                    ? normalized[i + 1].count / normalized[i].count
                    : 1
            );
        }
        for (let i = 0; i < normalized.length - 1; i++) {
            const current = normalized[i];
            const next = normalized[i + 1];
            const gap = current.count - next.count;
            const maxGap = current.count * maxGapPercent;
            if (gap > maxGap) {
                next.count = Math.round(current.count * (1 - maxGapPercent));
            } else {
                const target = Math.round(current.count * ratios[i]);
                if (current.count - target <= gap) {
                    next.count = target;
                }
            }
        }
        return normalized;
    }

    getColor(piece) {
        return (
            Math.floor(piece / GameConfig.TOKEN_SPACE) * GameConfig.TOKEN_SPACE
        );
    }

    getPowerUp(piece) {
        return piece % GameConfig.TOKEN_SPACE;
    }
}

class AnimationManager {
    constructor(renderer, board, speed) {
        this.renderer = renderer;
        this.board = board;
        this.speed = speed;
        this.startBoard = [];
        this.endBoard = [];
        this.indices = [];
        this.position = 0;
        this.piecesPerFrame = 8;
        this.resolve = null;
        this.frameId = null;
    }

    start(startBoard, endBoard, indices, piecesPerFrame) {
        this.startBoard = startBoard.slice();
        this.endBoard = endBoard.slice();
        this.indices = indices.slice();
        this.piecesPerFrame = piecesPerFrame;
        this.position = 0;
        return new Promise((resolve) => {
            this.resolve = resolve;
            this.frameId = requestAnimationFrame(() => this.step());
        });
    }

    step() {
        this.board.board = this.startBoard.slice();
        const newPos = this.position + this.piecesPerFrame;
        const endPos = Math.min(newPos, this.indices.length);
        for (let i = 0; i < endPos; i++) {
            this.board.board[this.indices[i]] = this.endBoard[this.indices[i]];
        }
        this.renderer.render();
        this.startBoard = this.board.board.slice();
        this.position = endPos;
        if (this.position < this.indices.length) {
            this.frameId = requestAnimationFrame(() => this.step());
        } else {
            this.cleanup();
        }
    }

    cleanup() {
        this.startBoard = [];
        this.endBoard = [];
        this.indices = [];
        this.position = 0;
        this.frameId = null;
        if (this.resolve) {
            this.resolve();
            this.resolve = null;
        }
    }

    isAnimating() {
        return this.frameId !== null;
    }
}

class Renderer {
    constructor(ctx, board, blockSize, rows, cols) {
        this.ctx = ctx;
        this.board = board;
        this.blockSize = blockSize;
        this.rows = rows;
        this.cols = cols;
        this.pieceSize = blockSize - GameConfig.GAP;
        this.outlineGap = 2;
        this.outlineSize = this.pieceSize - 2 * this.outlineGap;
        this.cornerRadius = 3;
    }

    getPosition(row, col) {
        return {
            x: GameConfig.BORDER_WIDTH + col * (this.blockSize + GameConfig.GAP),
            y: GameConfig.BORDER_WIDTH + row * (this.blockSize + GameConfig.GAP),
        };
    }

    drawEmptySpace(x, y) {
        this.ctx.save();
        this.ctx.strokeStyle = this.darkenColor(
            GameConfig.COLOR_MAP.get(GameConfig.COLORS.WHITE),
            70
        );
        this.ctx.lineWidth = 1;
        this.ctx.beginPath();
        this.ctx.rect(
            x + this.outlineGap,
            y + this.outlineGap,
            this.outlineSize,
            this.outlineSize
        );
        this.ctx.stroke();
        this.ctx.restore();
    }

    drawColoredPiece(x, y, piece) {
        this.ctx.save();
        const size = this.blockSize;
        const angle = (225 * Math.PI) / 180;
        const diagonal = Math.sqrt(size * size + size * size);
        const startX = x + size / 2 + (Math.cos(angle) * diagonal) / 2;
        const startY = y + size / 2 + (Math.sin(angle) * diagonal) / 2;
        const endX = x + size / 2 - (Math.cos(angle) * diagonal) / 2;
        const endY = y + size / 2 - (Math.sin(angle) * diagonal) / 2;
        const gradHighlight = this.ctx.createLinearGradient(x, y, x + size, y + size);
        const colors = GameConfig.GRADIENT_MAP.get(this.board.getColor(piece));
        gradHighlight.addColorStop(0, colors[1]);
        gradHighlight.addColorStop(1, this.lightenColor(colors[0], 5));
        this.ctx.fillStyle = gradHighlight;
        this.ctx.beginPath();
        this.ctx.roundRect(x, y, this.pieceSize, this.pieceSize, this.cornerRadius);
        this.ctx.fill();
        const gradBody = this.ctx.createLinearGradient(startX, startY, endX, endY);
        gradBody.addColorStop(0, this.darkenColor(colors[0], 4));
        gradBody.addColorStop(1, this.lightenColor(colors[1], 8));
        this.ctx.fillStyle = gradBody;
        this.ctx.beginPath();
        this.ctx.roundRect(
            x + 2,
            y + 2,
            this.pieceSize - 4,
            this.pieceSize - 4,
            this.cornerRadius
        );
        this.ctx.fill();
        this.ctx.restore();
    }

    drawHoverEffect(x, y, piece, glows) {
        this.ctx.save();
        const color = GameConfig.COLOR_MAP.get(this.board.getColor(piece));
        this.ctx.shadowColor = color;
        this.ctx.shadowBlur = 15;
        this.ctx.strokeStyle = glows
            ? "rgba(255, 255, 255, 0.7)"
            : this.lightenColor(color, 40);
        this.ctx.lineWidth = glows ? 8 : 4;
        const hoverSize = this.pieceSize + (glows ? 4 : 2);
        const hoverX = x - (glows ? 2 : 1);
        const hoverY = y - (glows ? 2 : 1);
        this.ctx.beginPath();
        this.ctx.roundRect(
            hoverX,
            hoverY,
            hoverSize,
            hoverSize,
            this.cornerRadius + 1
        );
        this.ctx.stroke();
        this.ctx.restore();
    }

    drawPowerUp(x, y, powerUpType) {
        const padding = 8;
        const shapeSize = this.pieceSize - padding;
        const shapeX = x + padding / 2;
        const shapeY = y + padding / 2;
        const drawFn =
            {
                [GameConfig.POWER_UPS.X]: this.drawX.bind(this),
                [GameConfig.POWER_UPS.PLUS]: this.drawPlus.bind(this),
                [GameConfig.POWER_UPS.RECT]: this.drawRect.bind(this),
                [GameConfig.POWER_UPS.FILL]: this.drawPowerFill.bind(this),
                [GameConfig.POWER_UPS.RIGHT]: this.drawRotateRight.bind(this),
                [GameConfig.POWER_UPS.LEFT]: this.drawRotateLeft.bind(this),
                [GameConfig.POWER_UPS.DISC]: this.drawDisc.bind(this),
                [GameConfig.POWER_UPS.EXCHANGE]: this.drawExchange.bind(this),
                [GameConfig.POWER_UPS.CIRCLE]: this.drawCircle.bind(this),
            }[powerUpType] || this.drawCircle.bind(this);
        drawFn(shapeX, shapeY, shapeSize);
    }

    drawX(x, y, size) {
        const ctx = this.prepareShapeContext();
        const padding = size * 0.2;
        ctx.moveTo(x + padding, y + padding);
        ctx.lineTo(x + size - padding, y + size - padding);
        ctx.moveTo(x + size - padding, y + padding);
        ctx.lineTo(x + padding, y + size - padding);
        this.finalizeShapeDraw(ctx);
    }

    drawPlus(x, y, size) {
        const ctx = this.prepareShapeContext();
        const padding = size * 0.2;
        ctx.moveTo(x + size / 2, y + padding);
        ctx.lineTo(x + size / 2, y + size - padding);
        ctx.moveTo(x + padding, y + size / 2);
        ctx.lineTo(x + size - padding, y + size / 2);
        this.finalizeShapeDraw(ctx);
    }

    drawCircle(x, y, size) {
        const ctx = this.prepareShapeContext();
        const radius = size * 0.3;
        ctx.arc(x + size / 2, y + size / 2, radius, 0, Math.PI * 2);
        this.finalizeShapeDraw(ctx);
    }

    drawDisc(x, y, size) {
        const ctx = this.prepareShapeContext();
        const radius = size * 0.3;
        const cx = x + size / 2;
        const cy = y + size / 2;
        ctx.beginPath();
        ctx.arc(cx, cy, radius, 0, Math.PI * 2);
        ctx.fillStyle = GameConfig.SHAPE_STROKE_COLOR;
        ctx.fill();
        ctx.strokeStyle = GameConfig.SHAPE_BORDER_COLOR;
        ctx.lineWidth = 2;
        ctx.stroke();
        this.finalizeShapeDraw(ctx);
    }

    drawExchange(x, y, size) {
        this.ctx.save();
        this.ctx.beginPath();
        const arrowColor = "white";
        const borderColor = "black";
        const arrowWidth = size * 0.2;
        const headLength = size * 0.2;
        const cx = x + size / 2;
        const topY = y + 6 + arrowWidth / 2 - 1;
        const bottomY = y + size - 6 - arrowWidth / 2 + 1;
        const drawArrow = (startX, endX, y, direction) => {
            this.ctx.beginPath();
            this.ctx.moveTo(startX, y - arrowWidth / 2);
            this.ctx.lineTo(endX, y - arrowWidth / 2);
            this.ctx.lineTo(endX, y - arrowWidth);
            this.ctx.lineTo(endX + direction * headLength, y);
            this.ctx.lineTo(endX, y + arrowWidth);
            this.ctx.lineTo(endX, y + arrowWidth / 2);
            this.ctx.lineTo(startX, y + arrowWidth / 2);
            this.ctx.closePath();
            this.ctx.fillStyle = arrowColor;
            this.ctx.fill();
            this.ctx.strokeStyle = borderColor;
            this.ctx.lineWidth = 1;
            this.ctx.stroke();
        };
        drawArrow(cx - size * 0.2, cx + size * 0.2 + 3, topY, 1);
        drawArrow(cx + size * 0.2, cx - size * 0.2 - 3, bottomY, -1);
        this.ctx.restore();
    }

    drawRect(x, y, size, opts) {
        const ctx = this.prepareShapeContext();
        const padding = size * 0.2;
        const rectSize = size - 2 * padding;
        ctx.rect(x + padding, y + padding, rectSize, rectSize);
        this.finalizeShapeDraw(ctx, opts);
    }

    drawPowerFill(x, y, size) {
        this.ctx.save();
        const halfSize = size / 2;
        const centerSize = size * 0.3;
        const centerOffset = (size - centerSize) / 2;
        const quadColors = [
            GameConfig.COLORS.BLUE,
            GameConfig.COLORS.GREEN,
            GameConfig.COLORS.RED,
            GameConfig.COLORS.YELLOW,
        ];
        const positions = [
            { qx: x, qy: y },
            { qx: x + halfSize, qy: y },
            { qx: x, qy: y + halfSize },
            { qx: x + halfSize, qy: y + halfSize },
        ];
        for (let i = 0; i < 4; i++) {
            this.ctx.beginPath();
            this.ctx.rect(positions[i].qx, positions[i].qy, halfSize, halfSize);
            this.ctx.fillStyle = GameConfig.COLOR_MAP.get(quadColors[i]);
            this.ctx.fill();
        }
        this.ctx.beginPath();
        this.ctx.rect(x + centerOffset, y + centerOffset, centerSize, centerSize);
        this.ctx.fillStyle = GameConfig.COLOR_MAP.get(GameConfig.COLORS.PURPLE);
        this.ctx.fill();
        this.ctx.restore();
        this.drawRect(x - 8, y - 8, size + 16, "reverse");
    }

    drawRotateRight(x, y, size) {
        const ctx = this.prepareShapeContext();
        const cx = x + size / 2;
        const cy = y + size / 2;
        const r = size * 0.3;
        const arrowLength = size * 0.25;
        const arrowWidth = size * 0.15;
        ctx.arc(cx, cy, r, Math.PI, 2 * Math.PI, false);
        const tipX = cx + r + 2;
        const tipY = cy + 3;
        ctx.moveTo(tipX - arrowLength, tipY - arrowWidth);
        ctx.lineTo(tipX, tipY);
        ctx.lineTo(tipX - arrowLength, tipY + arrowWidth);
        ctx.closePath();
        this.finalizeShapeDraw(ctx);
    }

    drawRotateLeft(x, y, size) {
        const ctx = this.prepareShapeContext();
        const cx = x + size / 2;
        const cy = y + size / 2;
        const r = size * 0.3;
        const arrowLength = size * 0.25;
        const arrowWidth = size * 0.15;
        ctx.arc(cx, cy, r, 0, Math.PI, true);
        const tipX = cx - r - 2;
        const tipY = cy + 3;
        ctx.moveTo(tipX + arrowLength, tipY - arrowWidth);
        ctx.lineTo(tipX, tipY);
        ctx.lineTo(tipX + arrowLength, tipY + arrowWidth);
        ctx.closePath();
        this.finalizeShapeDraw(ctx);
    }

    prepareShapeContext() {
        this.ctx.save();
        this.ctx.beginPath();
        return this.ctx;
    }

    finalizeShapeDraw(ctx, opts) {
        ctx.strokeStyle = GameConfig.SHAPE_BORDER_COLOR;
        ctx.lineWidth =
            GameConfig.SHAPE_LINE_WIDTH + GameConfig.SHAPE_BORDER_LINE_WIDTH + 1;
        ctx.stroke();
        ctx.strokeStyle = opts
            ? GameConfig.SHAPE_BORDER_COLOR
            : GameConfig.SHAPE_STROKE_COLOR;
        ctx.lineWidth = GameConfig.SHAPE_LINE_WIDTH;
        ctx.stroke();
        ctx.restore();
    }

    lightenColor(color, percent) {
        const num = parseInt(color.replace("#", ""), 16);
        const amt = Math.round(2.55 * percent);
        const R = Math.min((num >> 16) + amt, 255);
        const G = Math.min((num >> 8 & 0x00ff) + amt, 255);
        const B = Math.min((num & 0x0000ff) + amt, 255);
        return `#${((1 << 24) + (R << 16) + (G << 8) + B).toString(16).slice(1)}`;
    }

    darkenColor(color, percent) {
        const num = parseInt(color.replace("#", ""), 16);
        const amt = Math.round(2.55 * percent);
        const R = Math.max((num >> 16) - amt, 0);
        const G = Math.max((num >> 8 & 0x00ff) - amt, 0);
        const B = Math.max((num & 0x0000ff) - amt, 0);
        return `#${((1 << 24) + (R << 16) + (G << 8) + B).toString(16).slice(1)}`;
    }

    render(hoverList = new Set(), hoverIndex = -1) {
        this.ctx.clearRect(0, 0, this.ctx.canvas.width, this.ctx.canvas.height);
        for (let r = 0; r < this.rows; r++) {
            for (let c = 0; c < this.cols; c++) {
                const i = this.board.getIndex(r, c);
                const piece = this.board.board[i];
                const { x, y } = this.getPosition(r, c);
                if (piece === GameConfig.COLORS.WHITE) {
                    this.drawEmptySpace(x, y);
                } else {
                    this.drawColoredPiece(x, y, piece);
                }
            }
        }
        for (let r = 0; r < this.rows; r++) {
            for (let c = 0; c < this.cols; c++) {
                const i = this.board.getIndex(r, c);
                const piece = this.board.board[i];
                if (piece !== GameConfig.COLORS.WHITE) {
                    const { x, y } = this.getPosition(r, c);
                    const powerUp = this.board.getPowerUp(piece);
                    const glows = powerUp > 0 && !GameConfig.POWER_CONFIG.get(powerUp).connections;
                    if (powerUp > 0) {
                        this.drawPowerUp(x, y, powerUp);
                    }
                    if (hoverList.has(i)) {
                        this.drawHoverEffect(x, y, piece, glows);
                    }
                }
            }
        }
    }
}

class JawbreakerGame {
    constructor() {
        this.board = null;
        this.renderer = null;
        this.animator = null;
        this.rows = 0;
        this.cols = 0;
        this.blockSize = 0;
        this.cookieName = "";
        this.score = 0;
        this.bonus = 0;
        this.remainingPieces = 0;
        this.lastScore = 0;
        this.bestScore = 0;
        this.canvas = null;
        this.hoverList = new Set();
        this.hoverIndex = -1;
        this.mobile = false;
        this.undo = null;
        this.weightedFills = true;
    }

    clamp(value, min, max) {
        return Math.max(min, Math.min(max, Number.isInteger(value) ? value : min));
    }

    async init({
                   rows = 8,
                   cols = 8,
                   blockSize = 36,
                   cookieName = "jawbreaker_scores",
                   mobile = false,
               }) {
        this.rows = this.clamp(rows, 8, 20);
        this.cols = this.clamp(cols, 8, 20);
        this.blockSize = blockSize;
        this.cookieName = cookieName;
        this.mobile = mobile;
        this.score = 0;
        this.canvas = document.getElementById("game-canvas");
        if (!this.canvas) {
            throw new Error("Canvas element 'game-canvas' not found.");
        }
        this.canvas.width =
            this.cols * this.blockSize + (this.cols - 1) * GameConfig.GAP + 2 * GameConfig.BORDER_WIDTH;
        this.canvas.height =
            this.rows * this.blockSize + (this.rows - 1) * GameConfig.GAP + 2 * GameConfig.BORDER_WIDTH;
        this.board = new Board(this.rows, this.cols);
        this.renderer = new Renderer(
            this.canvas.getContext("2d"),
            this.board,
            this.blockSize,
            this.rows,
            this.cols
        );
        this.animator = new AnimationManager(
            this.renderer,
            this.board,
            GameConfig.ANIMATE ? GameConfig.ANIMATE_PIECE_SPEED : 0
        );
        this.board.initialize();
        const scores = CookieManager.getCookie(this.cookieName);
        if (scores) {
            const [last, best] = scores.split("|").map((s) => parseInt(s, 10));
            this.lastScore = isNaN(last) ? 0 : last;
            this.bestScore = isNaN(best) ? 0 : best;
        }
        document.getElementById("current-score").innerText = this.score;
        document.getElementById("last-score").innerText = this.lastScore;
        document.getElementById("best-score").innerText = this.bestScore;
        this.registerEvents();
        if (GameConfig.ANIMATE) {
            await this.animateNewBoard();
        } else {
            this.renderer.render();
        }
    }

    async animateNewBoard() {
        const size = this.rows * this.cols;
        const start = Array(size).fill(GameConfig.COLORS.WHITE);
        const end = this.board.board.slice();
        const indices = this.board.createShuffledIndices(size);
        await this.animator.start(start, end, indices, GameConfig.ANIMATE_BOARD_SPEED);
    }

    getCanvasCoordinates(e) {
        const rect = this.canvas.getBoundingClientRect();
        return {
            x: e.clientX - rect.left,
            y: e.clientY - rect.top,
        };
    }

    getBoardIndex(x, y) {
        const adjustedX = x - GameConfig.BORDER_WIDTH;
        const adjustedY = y - GameConfig.BORDER_WIDTH;
        const col = Math.floor(adjustedX / (this.blockSize + GameConfig.GAP));
        const row = Math.floor(adjustedY / (this.blockSize + GameConfig.GAP));
        if (col < 0 || col >= this.cols || row < 0 || row >= this.rows) {
            return -1;
        }
        return this.board.getIndex(row, col);
    }

    getConnectedPieces(index) {
        if (
            index < 0 ||
            index >= this.board.board.length ||
            this.board.board[index] === GameConfig.COLORS.WHITE
        ) {
            return [];
        }
        const piece = this.board.board[index];
        const color = this.board.getColor(piece);
        const powerUp = this.board.getPowerUp(piece);
        if (
            GameConfig.POWER_CONFIG.has(powerUp) &&
            !GameConfig.POWER_CONFIG.get(powerUp).connections
        ) {
            return [index];
        }
        let powerUpIndices = [];
        if (
            GameConfig.POWER_CONFIG.has(powerUp) &&
            GameConfig.POWER_CONFIG.get(powerUp).connections
        ) {
            powerUpIndices = this.getConnectedPowerUps(index, color, powerUp);
            if (powerUpIndices.length > 0) {
                powerUpIndices.unshift(index);
            }
            if (color === GameConfig.COLORS.GRAY) {
                return powerUpIndices;
            }
        }
        const connected = [];
        const stack = [index];
        const visited = new Array(this.board.board.length).fill(false);
        visited[index] = true;
        while (stack.length > 0) {
            const current = stack.pop();
            connected.push(current);
            const { row, col } = this.board.getPoint(current);
            const neighbors = [
                row > 0 ? current - this.cols : -1,
                row < this.rows - 1 ? current + this.cols : -1,
                col > 0 ? current - 1 : -1,
                col < this.cols - 1 ? current + 1 : -1,
            ];
            for (const n of neighbors) {
                if (
                    n !== -1 &&
                    !visited[n] &&
                    this.board.board[n] !== GameConfig.COLORS.WHITE &&
                    this.board.getColor(this.board.board[n]) === color
                ) {
                    visited[n] = true;
                    stack.push(n);
                }
            }
        }
        const combined = [...new Set([...powerUpIndices, ...connected])];
        return combined.length >= 2 ? combined : [];
    }

    getConnectedPowerUps(index, color, powerUp) {
        const handlers = {
            [GameConfig.POWER_UPS.X]: () =>
                this.getXConnections(index, color),
            [GameConfig.POWER_UPS.PLUS]: () =>
                this.getPlusConnections(index, color),
            [GameConfig.POWER_UPS.RECT]: () =>
                this.getRectConnections(index, color),
            [GameConfig.POWER_UPS.CIRCLE]: () =>
                this.getCircularConnections(index, color),
            [GameConfig.POWER_UPS.DISC]: () =>
                this.getDiscConnections(index, color),
        };
        return handlers[powerUp]?.() || [];
    }

    getConnectedDirections(index, color, directions, maxIterations) {
        const { row, col } = this.board.getPoint(index);
        const indices = [];
        for (let iter = 0; iter < maxIterations; iter++) {
            for (const dir of directions) {
                const r = row + dir.sr + dir.dr * iter;
                const c = col + dir.sc + dir.dc * iter;
                if (r >= 0 && r < this.rows && c >= 0 && c < this.cols) {
                    const i = this.board.getIndex(r, c);
                    const pieceColor = this.board.getColor(this.board.board[i]);
                    if (
                        (color === GameConfig.COLORS.GRAY &&
                            pieceColor !== GameConfig.COLORS.WHITE &&
                            pieceColor !== GameConfig.COLORS.GRAY) ||
                        (color !== GameConfig.COLORS.GRAY && color === pieceColor)
                    ) {
                        if (!indices.includes(i)) {
                            indices.push(i);
                        }
                    }
                }
            }
        }
        return indices;
    }

    getXConnections(index, color) {
        const maxIters = Math.ceil(
            Math.sqrt(this.rows ** 2 + this.cols ** 2)
        );
        return this.getConnectedDirections(
            index,
            color,
            GameConfig.POWER_UP_DIRECTIONS.X,
            maxIters
        );
    }

    getPlusConnections(index, color) {
        const maxIters = Math.max(this.rows, this.cols);
        return this.getConnectedDirections(
            index,
            color,
            GameConfig.POWER_UP_DIRECTIONS.PLUS,
            maxIters
        );
    }

    getRectConnections(index, color) {
        return this.getConnectedDirections(
            index,
            color,
            GameConfig.POWER_UP_DIRECTIONS.RECT,
            4
        );
    }

    getCircularConnections(index, color) {
        const c1 = this.getConnectedDirections(
            index,
            color,
            GameConfig.POWER_UP_DIRECTIONS.RECT,
            1
        );
        const c2 = this.getConnectedDirections(
            index,
            color,
            GameConfig.POWER_UP_DIRECTIONS.CIRCLE,
            3
        );
        return [...c1, ...c2];
    }

    getDiscConnections(index, color) {
        const c1 = this.getConnectedDirections(
            index,
            color,
            GameConfig.POWER_UP_DIRECTIONS.DISC1,
            3
        );
        const c2 = this.getConnectedDirections(
            index,
            color,
            GameConfig.POWER_UP_DIRECTIONS.DISC2,
            2
        );
        return [...c1, ...c2];
    }

    async removePieces(index) {
        const piece = this.board.board[index];
        const powerUp = this.board.getPowerUp(piece);
        let count = 0;
        let isSpecial = false;
        if (powerUp === GameConfig.POWER_UPS.FILL) {
            this.board.board[index] = GameConfig.COLORS.WHITE;
            this.board.applyGravity();
            await this.animateTransformation(() => this.board.fillSpaces());
            count = 1;
            isSpecial = true;
        } else if (powerUp === GameConfig.POWER_UPS.EXCHANGE) {
            this.board.board[index] = GameConfig.COLORS.WHITE;
            this.board.applyGravity();
            await this.animateTransformation(() => this.board.fillPieces());
            count = 1;
            isSpecial = true;
        } else if (powerUp === GameConfig.POWER_UPS.RIGHT) {
            this.board.board[index] = GameConfig.COLORS.WHITE;
            this.board.rotate(90);
            this.board.applyGravity();
            count = 1;
            isSpecial = true;
        } else if (powerUp === GameConfig.POWER_UPS.LEFT) {
            this.board.board[index] = GameConfig.COLORS.WHITE;
            this.board.rotate(270);
            this.board.applyGravity();
            count = 1;
            isSpecial = true;
        } else {
            const pieces = this.getConnectedPieces(index);
            for (const i of pieces) {
                if (this.board.board[i] !== GameConfig.COLORS.WHITE) {
                    this.board.board[i] = GameConfig.COLORS.WHITE;
                    count++;
                }
            }
        }
        return { count, isSpecial };
    }

    async animateTransformation(transformFn) {
        const start = this.board.board.slice();
        transformFn();
        const end = this.board.board.slice();
        const indices = this.board.createShuffledIndices(this.board.board.length).filter(
            (i) => start[i] !== end[i]
        );
        await this.animator.start(start, end, indices, GameConfig.ANIMATE_PIECE_SPEED);
    }

    calculateMoveScore(count) {
        return count < 2 ? 0 : count * (count - 1);
    }

    calculateBonus(remaining) {
        if (remaining === 0) return 2000;
        if (remaining <= 10) return (11 - remaining) * 100;
        return 0;
    }

    isGameOver() {
        for (let i = 0; i < this.board.board.length; i++) {
            if (this.board.board[i] === GameConfig.COLORS.WHITE) continue;
            const powerUp = this.board.getPowerUp(this.board.board[i]);
            if (
                GameConfig.POWER_CONFIG.has(powerUp) &&
                !GameConfig.POWER_CONFIG.get(powerUp).connections
            ) {
                return false;
            }
            if (this.getConnectedPieces(i).length > 0) {
                return false;
            }
        }
        return true;
    }

    async processMove(index) {
        this.undo = { board: this.board.board.slice(), score: this.score };
        document.getElementById("undo-btn").disabled = false;
        const start = this.board.board.slice();
        const { count, isSpecial } = await this.removePieces(index);
        if (count > 0 && !isSpecial) {
            const changed = this.board.board
            .map((p, i) => (p !== start[i] ? i : null))
            .filter((i) => i !== null);
            await this.animator.start(start, this.board.board, changed, 1);
            this.board.applyGravity();
        }
        if (!isSpecial) {
            this.score += this.calculateMoveScore(count);
        }
        const gameOver = this.isGameOver();
        if (gameOver) {
            this.remainingPieces = this.board.board.filter(
                (p) => p !== GameConfig.COLORS.WHITE
            ).length;
            this.bonus = this.calculateBonus(this.remainingPieces);
            this.score += this.bonus;
            this.undo = null;
            this.lastScore = this.score;
            if (this.score > this.bestScore) {
                this.bestScore = this.score;
            }
            document.getElementById("game-over-score").innerText = this.score;
            document.getElementById("last-score").innerText = this.lastScore;
            document.getElementById("best-score").innerText = this.bestScore;
            document.getElementById("bonus-points").innerText = this.bonus;
            document.getElementById("remaining-pieces").innerText = this.remainingPieces;
            document.getElementById("pieces-text").innerText =
                this.remainingPieces === 1 ? "piece" : "pieces";
            document.getElementById("game-over-overlay").classList.add("visible");
            CookieManager.setCookie(
                this.cookieName,
                `${this.lastScore}|${this.bestScore}`
            );
        }
        document.getElementById("current-score").innerText = this.score;
        this.hoverList.clear();
        this.hoverIndex = -1;
        if (this.mobile) {
            const pieces = this.getConnectedPieces(index);
            if (pieces.length > 1) {
                this.hoverIndex = index;
                pieces.forEach((i) => this.hoverList.add(i));
            }
        }
        this.renderer.render(this.hoverList, this.hoverIndex);
        return { score: this.score, gameOver };
    }

    undoMove() {
        if (!this.undo) return;
        this.board.board = this.undo.board.slice();
        this.score = this.undo.score;
        this.undo = null;
        document.getElementById("undo-btn").disabled = true;
        document.getElementById("current-score").innerText = this.score;
        this.renderer.render();
    }

    async handleClick(e) {
        if (this.animator.isAnimating()) return;
        const { x, y } = this.getCanvasCoordinates(e);
        const index = this.getBoardIndex(x, y);
        if (
            index >= 0 &&
            index < this.board.board.length &&
            this.board.board[index] !== GameConfig.COLORS.WHITE
        ) {
            await this.processMove(index);
        }
    }

    handleMouseMove(e) {
        if (this.animator.isAnimating()) return;
        const { x, y } = this.getCanvasCoordinates(e);
        const index = this.getBoardIndex(x, y);
        if (index === this.hoverIndex) return;
        this.hoverList.clear();
        this.hoverIndex = index;
        if (
            index >= 0 &&
            index < this.board.board.length &&
            this.board.board[index] !== GameConfig.COLORS.WHITE
        ) {
            this.getConnectedPieces(index).forEach((i) => this.hoverList.add(i));
        }
        this.renderer.render(this.hoverList, this.hoverIndex);
    }

    handleMouseLeave() {
        if (this.animator.isAnimating() || this.hoverList.size === 0) return;
        this.hoverList.clear();
        this.hoverIndex = -1;
        this.renderer.render();
    }

    handleTouch(e) {
        if (this.animator.isAnimating()) return;
        const { x, y } = this.getCanvasCoordinates(e);
        const index = this.getBoardIndex(x, y);
        if (index === this.hoverIndex) {
            this.handleClick(e);
        } else {
            this.handleMouseMove(e);
        }
    }

    async resetGame() {
        this.board.initialize();
        this.hoverList.clear();
        this.hoverIndex = -1;
        this.score = 0;
        this.bonus = 0;
        this.remainingPieces = 0;
        this.undo = null;
        document.getElementById("undo-btn").disabled = true;
        document.getElementById("current-score").innerText = this.score;
        document.getElementById("game-over-overlay").classList.remove("visible");
        if (GameConfig.ANIMATE) {
            await this.animateNewBoard();
        } else {
            this.renderer.render();
        }
    }

    registerEvents() {
        if (this.mobile) {
            this.canvas.addEventListener("click", (e) => this.handleTouch(e));
        } else {
            this.canvas.addEventListener("click", (e) => this.handleClick(e));
            this.canvas.addEventListener("mousemove", (e) => this.handleMouseMove(e));
            this.canvas.addEventListener("mouseleave", () => this.handleMouseLeave());
        }
        const newGameBtn = document.getElementById("new-game-btn");
        if (newGameBtn) {
            newGameBtn.addEventListener("click", (e) => {
                e.preventDefault();
                this.resetGame();
            });
        }
        const restartBtn = document.getElementById("restart-btn");
        if (restartBtn) {
            restartBtn.addEventListener("click", (e) => {
                e.preventDefault();
                this.resetGame();
            });
        }
        const undoBtn = document.getElementById("undo-btn");
        if (undoBtn) {
            undoBtn.disabled = true;
            undoBtn.addEventListener("click", (e) => {
                e.preventDefault();
                this.undoMove();
            });
        }
        const helpBtn = document.getElementById("help-btn");
        const helpModal = document.getElementById("help-modal");
        const closeBtn = helpModal?.querySelector(".close");
        if (helpBtn && helpModal && closeBtn) {
            helpBtn.addEventListener("click", () => {
                helpModal.style.display = "block";
            });
            closeBtn.addEventListener("click", () => {
                helpModal.style.display = "none";
            });
            window.addEventListener("click", (e) => {
                if (e.target === helpModal) {
                    helpModal.style.display = "none";
                }
            });
        }
    }
}

async function initGame(opts) {
    const game = new JawbreakerGame();
    await game.init(opts);
}