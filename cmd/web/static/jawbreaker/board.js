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
