class JB {
    // Piece constants
    static tokenSpace = 1000;
    static WHITE = 0;
    static PURPLE = 1000;
    static BLUE = 2000;
    static GREEN = 3000;
    static RED = 4000;
    static YELLOW = 5000;
    static GRAY = 6000;

    // standard game pieces
    static GAME_PIECES = [
        JB.PURPLE,
        JB.BLUE,
        JB.GREEN,
        JB.RED,
        JB.YELLOW
    ];

    static COLOR_MAP = new Map([
        [JB.PURPLE, "#8a2be2"],
        [JB.BLUE, "#00a0ff"],
        [JB.GREEN, "#00cc66"],
        [JB.RED, "#ff3333"],
        [JB.YELLOW, "#ffcc00"],
        [JB.GRAY, "#999999"],
        [JB.WHITE, "#ffffff"],
    ]);

    static POWER_X = 1;
    static POWER_PLUS = 2;
    static POWER_CIRCLE = 3;
    static POWER_RECT = 4;
    static POWER_FILL = 5;

    static POWER_PIECES = [
        JB.POWER_X,
        JB.POWER_PLUS,
        JB.POWER_CIRCLE,
        JB.POWER_RECT,
        JB.POWER_FILL,
    ];

    static shapeStrokeColor = "white";
    static shapeBorderColor = "black";
    static shapeLineWidth = 4;
    static shapeBorderWidth = 2;

    static POWER_UP_X_DIRECTIONS = [
        { sr: -1, sc: -1, dr: -1, dc: -1 }, // Top-left
        { sr: -1, sc:  1, dr: -1, dc:  1 }, // Top-right
        { sr:  1, sc: -1, dr:  1, dc: -1 }, // Bottom-left
        { sr:  1, sc:  1, dr:  1, dc:  1 }, // Bottom-right
    ];

    static POWER_UP_PLUS_DIRECTIONS = [
        { sr: -1, sc:  0, dr: -1, dc:  0 }, // Up
        { sr:  1, sc:  0, dr:  1, dc:  0 }, // Down
        { sr:  0, sc: -1, dr:  0, dc: -1 }, // Left
        { sr:  0, sc:  1, dr:  0, dc:  1 }, // Right
    ];

    static POWER_UP_RECT_DIRECTIONS = [
        { sr: -2, sc: -2, dc:  0, dr:  1 },
        { sr:  2, sc: -2, dc:  1, dr:  0 },
        { sr:  2, sc:  2, dc:  0, dr: -1 },
        { sr: -2, sc:  2, dc: -1, dr:  0 },
    ]

    static POWER_UP_CIRCLE_DIRECTIONS = [
        { sr: -1, sc: -3, dc:  0, dr:  1 },
        { sr:  3, sc: -1, dc:  1, dr:  0 },
        { sr:  1, sc:  3, dc:  0, dr: -1 },
        { sr: -3, sc:  1, dc: -1, dr:  0 },
    ]

    /**
     * Create a new Jawbreaker game
     */
    constructor(opts) {
        this.rows = this.clamp(opts.rows, 8, 32);
        this.cols = this.clamp(opts.cols, 8, 32);
        this.gap = opts.gap || 1;
        this.border = opts.border || 6;
        this.blockSize = opts.blockSize || 36;
        this.cookieName = opts.cookieName || "jawbreaker";
        this.score = 0;
        this.lastScore = 0;
        this.bestScore = 0;
        this.board = this.newBoard();
        this.hoverList = new Set();

        this.canvas = document.getElementById("game-canvas");
        this.ctx = this.canvas.getContext("2d");

        const canvasWidth = this.cols * this.blockSize +
            (this.cols - 1) * this.gap + 2 * this.border;
        const canvasHeight = this.rows * this.blockSize +
            (this.rows - 1) * this.gap + 2 * this.border;

        this.canvas.width = canvasWidth;
        this.canvas.height = canvasHeight;

        // Read scores cookie
        const scoresCookie = this.getCookie(this.cookieName);
        if (scoresCookie) {
            const [lastScore, bestScore] = scoresCookie.split('|');
            this.lastScore = parseInt(lastScore, 10) || 0;
            this.bestScore = parseInt(bestScore, 10) || 0;

            // Update UI
            document.getElementById('last-score').innerText = this.lastScore;
            document.getElementById('best-score').innerText = this.bestScore;
        }

        this.registerEvents();
        this.renderBoard();
    }

    registerEvents() {
        this.canvas.addEventListener("click", (e) => {
            this.handleClick(e);
        });

        this.canvas.addEventListener("mousemove", (e) => {
            this.handleMouseMove(e);
        });

        this.canvas.addEventListener("mouseleave", () => {
            this.hoverList.clear();
            this.hoverIndex = -1;
            this.renderBoard();
        });

        const newGameBtn = document.querySelector(".new-game-btn");
        if (newGameBtn) {
            newGameBtn.addEventListener("click", (e) => {
                this.resetGame();
                document
                .getElementById("game-over-overlay")
                .classList.remove("visible");
                e.preventDefault();
            });
        }

        // Help button and modal
        const helpBtn = document.getElementById("help-btn");
        const helpModal = document.getElementById("help-modal");
        const closeBtn = document.querySelector(".close");

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

    resetGame() {
        this.board = this.newBoard();
        this.hoverList.clear();
        this.hoverIndex = -1;
        this.score = 0;
        document.getElementById("current-score").innerText = this.score;
        this.renderBoard();
    }

    /**
     * Create a random board with the given dimensions
     * @returns {Array} - Board as an array of integers
     */
    newBoard() {
        let board = [];
        for (let i = 0; i < this.rows * this.cols; i++) {
            let powerUp = 0;
            let color = JB.GAME_PIECES[Math.floor(Math.random() * JB.GAME_PIECES.length)];

            const rnd = Math.random();
            if (rnd < 0.02) {
                powerUp = JB.POWER_PIECES[Math.floor(Math.random() * JB.POWER_PIECES.length)];
                color = JB.GRAY;
            } else if (rnd < 0.07) {
                powerUp = JB.POWER_PIECES[Math.floor(Math.random() * JB.POWER_PIECES.length)];
            }

            board.push(color+powerUp);
        }
        return board;
    }

    getCanvasCoordinates(e) {
        const rect = this.canvas.getBoundingClientRect();
        const x = e.clientX - rect.left;
        const y = e.clientY - rect.top;
        return {x, y};
    }

    getBoardIndexFromCoordinates(x, y) {
        // Adjust coordinates to account for a border
        const adjustedX = x - this.border;
        const adjustedY = y - this.border;

        // Calculate column and row
        const col = Math.floor(adjustedX / (this.blockSize + this.gap));
        const row = Math.floor(adjustedY / (this.blockSize + this.gap));

        if (col < 0 || col >= this.cols || row < 0 || row >= this.rows) {
            return -1;
        }

        return row * this.cols + col;
    }

    getPointFromIndex(index) {
        const col = index % this.cols;
        const row = Math.floor(index / this.cols);
        return {col: col, row: row};
    }

    handleClick(e) {
        const {x, y} = this.getCanvasCoordinates(e);
        const index = this.getBoardIndexFromCoordinates(x, y);

        if (index >= 0 && index < this.board.length) {
            const status = this.move(index);
            document.getElementById("current-score").innerText = status.score;

            if (status.gameOver) {
                this.lastScore = this.score;
                if (this.score > this.bestScore) {
                    this.bestScore = this.score;
                }

                document.getElementById("game-over-score").innerText = this.score;
                document.getElementById("last-score").innerText = this.lastScore;
                document.getElementById("best-score").innerText = this.bestScore;

                this.setCookie(this.cookieName, `${this.lastScore}|${this.bestScore}`);

                document
                .getElementById("game-over-overlay")
                .classList.add("visible");
                return;
            }

            this.hoverList.clear();
            this.hoverIndex = -1;
            const pieces = this.getConnections(index);
            if (pieces.length > 1) {
                this.hoverIndex = index;
                for (const i of pieces) {
                    this.hoverList.add(i);
                }
            }

            this.renderBoard();
        }
    }

    handleMouseMove(e) {
        const {x, y} = this.getCanvasCoordinates(e);
        const index = this.getBoardIndexFromCoordinates(x, y);

        if (index >= 0 && index < this.board.length && index !== this.hoverIndex) {
            this.hoverList.clear();
            this.hoverIndex = index;
            const pieces = this.getConnections(index);
            if (pieces.length > 1) {
                for (const i of pieces) {
                    this.hoverList.add(i);
                }
            }
            this.renderBoard();
        }
    }

    colorMap(piece) {
        return JB.COLOR_MAP.get(this.getPiece(piece));
    }

    getPiece(piece) {
        return Math.floor(piece/JB.tokenSpace) * JB.tokenSpace;
    }

    getPowerUp(piece) {
        return piece % JB.tokenSpace;
    }

    /**
     * Calculate the position of a piece on the canvas
     */
    getPiecePosition(row, col) {
        const x = this.border + col * (this.blockSize + this.gap);
        const y = this.border + row * (this.blockSize + this.gap);
        return { x, y };
    }

    /**
     * Draw an empty (white) space
     */
    drawEmptySpace(x, y, pieceSize, outlineGap, outlineSize) {
        const ctx = this.ctx;
        ctx.save();
        ctx.strokeStyle = this.darkenColor(JB.COLOR_MAP.get(JB.WHITE), 70);
        ctx.lineWidth = 1;
        ctx.beginPath();
        ctx.rect(x + outlineGap, y + outlineGap, outlineSize, outlineSize);
        ctx.stroke();
        ctx.restore();
    }

    /**
     * Draw a colored piece
     */
    drawColoredPiece(x, y, pieceSize, piece, cornerRadius, glow) {
        const ctx = this.ctx;
        ctx.save();

        const gradient = ctx.createLinearGradient(
            x + this.blockSize * 0.8,
            y + this.blockSize * 0.8,
            x + this.blockSize * 0.1,
            y + this.blockSize * 0.1,
        );

        const color = this.colorMap(piece);
        gradient.addColorStop(0, this.lightenColor(color, 12));
        gradient.addColorStop(1, this.darkenColor(color, 6));
        ctx.fillStyle = gradient;

        ctx.beginPath();
        ctx.roundRect(x, y, pieceSize, pieceSize, cornerRadius);
        ctx.fill();

        // Draw top and left highlights
        ctx.beginPath();
        ctx.rect(x + cornerRadius, y, pieceSize - (2 * cornerRadius), glow);
        ctx.rect(x, y + cornerRadius, glow, pieceSize - (2 * cornerRadius));
        ctx.fillStyle = this.lightenColor(color, 20);
        ctx.fill();

        // Draw bottom and right shadows
        ctx.beginPath();
        ctx.rect(x + pieceSize - glow, y + cornerRadius, glow, pieceSize - (2 * cornerRadius));
        ctx.rect(x + cornerRadius, y + pieceSize - glow, pieceSize - (2 * cornerRadius), glow);
        ctx.fillStyle = this.darkenColor(color, 20);
        ctx.fill();

        ctx.restore();
    }

    /**
     * Draw hover effect for a piece
     */
    drawHoverEffect(x, y, pieceSize, color, cornerRadius) {
        const ctx = this.ctx;
        ctx.save();

        ctx.shadowColor = color;
        ctx.shadowBlur = 15;
        ctx.strokeStyle = this.lightenColor(color, 40);
        ctx.lineWidth = 6;

        // Draw rounded rectangle for hover effect
        ctx.beginPath();
        const hoverSize = pieceSize + 2; // Slightly larger than the piece
        const hoverX = x - 1;
        const hoverY = y - 1;

        // Draw a rounded rectangle path for the hover effect
        ctx.moveTo(hoverX + cornerRadius, hoverY);
        ctx.lineTo(hoverX + hoverSize - cornerRadius, hoverY);
        ctx.arcTo(hoverX + hoverSize, hoverY, hoverX + hoverSize, hoverY + cornerRadius, cornerRadius);
        ctx.lineTo(hoverX + hoverSize, hoverY + hoverSize - cornerRadius);
        ctx.arcTo(hoverX + hoverSize, hoverY + hoverSize, hoverX + hoverSize - cornerRadius, hoverY + hoverSize, cornerRadius);
        ctx.lineTo(hoverX + cornerRadius, hoverY + hoverSize);
        ctx.arcTo(hoverX, hoverY + hoverSize, hoverX, hoverY + hoverSize - cornerRadius, cornerRadius);
        ctx.lineTo(hoverX, hoverY + cornerRadius);
        ctx.arcTo(hoverX, hoverY, hoverX + cornerRadius, hoverY, cornerRadius);

        ctx.stroke();
        ctx.restore();
    }

    /**
     * Render the game board
     */
    renderBoard() {
        this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);
        const cornerRadius = 2;
        const pieceSize = this.blockSize - this.gap;
        const outlineGap = 2;
        const outlineSize = pieceSize - (2 * outlineGap);
        const glow = 1;

        // First pass: Draw all pieces (colored or empty)
        for (let row = 0; row < this.rows; row++) {
            for (let col = 0; col < this.cols; col++) {
                const index = row * this.cols + col;
                const piece = this.board[index];
                const { x, y } = this.getPiecePosition(row, col);

                if (piece === JB.WHITE) {
                    this.drawEmptySpace(x, y, pieceSize, outlineGap, outlineSize);
                } else {
                    this.drawColoredPiece(x, y, pieceSize, piece, cornerRadius, glow);
                }
            }
        }

        // Second pass: Draw power-ups and hover effects
        for (let row = 0; row < this.rows; row++) {
            for (let col = 0; col < this.cols; col++) {
                const index = row * this.cols + col;
                const piece = this.board[index];

                if (piece !== JB.WHITE) {
                    const { x, y } = this.getPiecePosition(row, col);

                    // Draw power-up if present
                    if (this.getPowerUp(piece) > 0) {
                        this.drawPowerUp(x, y, pieceSize, this.getPowerUp(piece));
                    }

                    // Draw hover effect if this piece is in the hover list
                    if (this.hoverList.has(index)) {
                        const color = this.colorMap(piece);
                        this.drawHoverEffect(x, y, pieceSize, color, cornerRadius);
                    }
                }
            }
        }
    }

    drawPowerUp(x, y, size, powerUp) {
        const f = 8;
        size -= f;
        x += (f/2);
        y += (f/2);
        switch (powerUp) {
            case JB.POWER_X:
                this.drawX(x, y, size);
                break;
            case JB.POWER_PLUS:
                this.drawPlus(x, y, size);
                break;
            case JB.POWER_RECT:
                this.drawRect(x, y, size);
                break;
            default:
                this.drawCircle(x, y, size);
                break;
        }
    }

    /**
     * Common method to prepare context for shape drawing
     */
    prepareShapeContext() {
        const ctx = this.ctx;
        ctx.save();
        ctx.beginPath();
        return ctx;
    }

    drawX(x, y, size) {
        const ctx = this.prepareShapeContext();
        const padding = size * 0.15;

        ctx.moveTo(x + padding, y + padding);
        ctx.lineTo(x + size - padding, y + size - padding);
        ctx.moveTo(x + size - padding, y + padding);
        ctx.lineTo(x + padding, y + size - padding);

        this.drawShape(ctx);
        ctx.restore();
    }

    drawPlus(x, y, size) {
        const ctx = this.prepareShapeContext();
        const padding = size * 0.15;

        ctx.moveTo(x + size / 2, y + padding);
        ctx.lineTo(x + size / 2, y + size - padding);
        ctx.moveTo(x + padding, y + size / 2);
        ctx.lineTo(x + size - padding, y + size / 2);

        this.drawShape(ctx);
        ctx.restore();
    }

    drawCircle(x, y, size) {
        const ctx = this.prepareShapeContext();
        const radius = size * 0.3;

        ctx.arc(x + size / 2, y + size / 2, radius, 0, Math.PI * 2);

        this.drawShape(ctx);
        ctx.restore();
    }

    drawRect(x, y, size) {
        const ctx = this.prepareShapeContext();
        const padding = size * 0.15;
        const rectSize = size * 0.7;

        ctx.rect(x + padding, y + padding, rectSize, rectSize);

        this.drawShape(ctx);
        ctx.restore();
    }

    drawShape(ctx) {
        ctx.strokeStyle = JB.shapeBorderColor;
        ctx.lineWidth = JB.shapeLineWidth + JB.shapeBorderWidth;
        ctx.stroke();
        ctx.strokeStyle = JB.shapeStrokeColor;
        ctx.lineWidth = JB.shapeLineWidth;
        ctx.stroke();
    }


    lightenColor(color, percent) {
        const num = parseInt(color.replace("#", ""), 16);
        const amt = Math.round(2.55 * percent);
        const R = (num >> 16) + amt;
        const G = (num >> 8 & 0x00FF) + amt;
        const B = (num & 0x0000FF) + amt;

        return this.rgbString(R, G, B);
    }

    darkenColor(color, percent) {
        const num = parseInt(color.replace("#", ""), 16);
        const amt = Math.round(2.55 * percent);
        const R = (num >> 16) - amt;
        const G = (num >> 8 & 0x00FF) - amt;
        const B = (num & 0x0000FF) - amt;

        return this.rgbString(R, G, B);
    }

    rgbString(R, G, B) {
        return "#" + (
            0x1000000 +
            (R < 255 ? R < 1 ? 0 : R : 255) * 0x10000 +
            (G < 255 ? G < 1 ? 0 : G : 255) * 0x100 +
            (B < 255 ? B < 1 ? 0 : B : 255)
        ).toString(16).slice(1);
    }

      getConnections(index, filtered) {
        if (index < 0 || index >= this.board.length) {
            return [];
        }

        if (filtered === undefined) {
            filtered = true;
        }

        const target = this.getPiece(this.board[index]);
        if (target === JB.WHITE) {
            return [];
        }

        const powerUp = this.getPowerUp(this.board[index]);
        if (powerUp > 0) {
            let c = this.getPowerUpConnections(index, target, powerUp);
            if (c.length > 0) {
                c.unshift(index);
            }

            return filtered ? c.filter(value => value >= 0) : c;
        }

        const connectedIndices = [];
        const stack = [index];
        const visited = new Array(this.rows * this.cols).fill(false);

        while (stack.length > 0) {
            const i = stack.pop();

            if (i < 0 || i >= this.board.length || this.getPiece(this.board[i]) !== target || visited[i]) {
                continue;
            }

            visited[i] = true;
            connectedIndices.push(i);

            const row = Math.floor(i / this.cols);
            const col = i % this.cols;

            if (row > 0) stack.push(i - this.cols);
            if (row < this.rows - 1) stack.push(i + this.cols);
            if (col > 0) stack.push(i - 1);
            if (col < this.cols - 1) stack.push(i + 1);
        }

        if (connectedIndices.length < 2) {
            return [];
        }

          return filtered ? connectedIndices.filter(value => value >= 0) : connectedIndices;
    }

    getPowerUpConnections(index, target, powerUp) {
        switch(powerUp) {
            case JB.POWER_X:
                return this.getXConnections(index, target);
            case JB.POWER_PLUS:
                return this.getPlusConnections(index, target)
            case JB.POWER_RECT:
                return this.getRectConnections(index, target);
            default:
                return this.getCircularConnections(index, target);
        }
    }

    getPlusConnections(index, target) {
        const maxIters = Math.max(this.rows, this.cols);
        return this.geConnectionsWithDirections(index, target, JB.POWER_UP_PLUS_DIRECTIONS, maxIters);
    }

    getXConnections(index, target) {
        const maxIters = Math.max(this.rows, this.cols);
        return this.geConnectionsWithDirections(index, target, JB.POWER_UP_X_DIRECTIONS, maxIters);
    }

    getRectConnections(index, target) {
        return this.geConnectionsWithDirections(index, target, JB.POWER_UP_RECT_DIRECTIONS, 4);
    }

    getCircularConnections(index, target) {
        let c1 = this.geConnectionsWithDirections(index, target, JB.POWER_UP_RECT_DIRECTIONS, 1);
        let c2 = this.geConnectionsWithDirections(index, target, JB.POWER_UP_CIRCLE_DIRECTIONS, 3);
        return [...c1, ...c2];
    }

    geConnectionsWithDirections(index, target, directions, max) {
        const point = this.getPointFromIndex(index);
        let startRow = point.row;
        let startCol = point.col;
        let connectedIndices = [];

        for (let i = 0; i < max; i++) {
            for (const dir of directions) {
                const r = startRow + dir.sr + dir.dr * i;
                const c = startCol + dir.sc + dir.dc * i;

                if (r >= 0 && r < this.rows && c >= 0 && c < this.cols) {
                    index = r * this.cols + c;
                    const piece = this.getPiece(this.board[index]);
                    if (piece === target || (target === JB.GRAY && target !== JB.WHITE)) {
                        connectedIndices.push(index);
                    } else {
                        connectedIndices.push(-1);
                    }
                }
            }
        }

        return connectedIndices;
    }

    floodFill(index) {
        const connectedPieces = this.getConnections(index);
        if (connectedPieces.length < 2) {
            return 0;
        }

        for (const i of connectedPieces) {
            this.board[i] = JB.WHITE;
        }

        return connectedPieces.length;
    }

    applyGravityAndShiftRight() {
        let boardArray = [];
        for (let row = 0; row < this.rows; row++) {
            const rowArray = [];
            for (let col = 0; col < this.cols; col++) {
                rowArray.push(this.board[row * this.cols + col]);
            }
            boardArray.push(rowArray);
        }

        for (let col = 0; col < this.cols; col++) {
            let writeRow = this.rows - 1;
            for (let row = this.rows - 1; row >= 0; row--) {
                if (boardArray[row][col] !== JB.WHITE) {
                    boardArray[writeRow][col] = boardArray[row][col];
                    if (writeRow !== row) {
                        boardArray[row][col] = JB.WHITE;
                    }
                    writeRow--;
                }
            }
        }

        let writeCol = this.cols - 1;
        for (let col = this.cols - 1; col >= 0; col--) {
            let isEmpty = true;
            for (let row = 0; row < this.rows; row++) {
                if (boardArray[row][col] !== JB.WHITE) {
                    isEmpty = false;
                    break;
                }
            }

            if (!isEmpty) {
                if (writeCol !== col) {
                    for (let row = 0; row < this.rows; row++) {
                        boardArray[row][writeCol] = boardArray[row][col];
                        boardArray[row][col] = JB.WHITE;
                    }
                }
                writeCol--;
            }
        }

        let b = [];
        for (let row = 0; row < this.rows; row++) {
            for (let col = 0; col < this.cols; col++) {
                b.push(boardArray[row][col]);
            }
        }
        this.board = b;
    }

    move(index) {
        const n = this.floodFill(index);
        if (n < 2) {
            return {
                board: this.board,
                score: this.score,
                gameOver: this.isGameOver(),
            };
        }

        this.applyGravityAndShiftRight();
        this.score += this.calculateMoveScore(n);
        const gameOver = this.isGameOver();

        if (gameOver) {
            let remainingPieces = 0;
            for (let i = 0; i < this.board.length; i++) {
                if (this.board[i] !== JB.WHITE) {
                    remainingPieces++;
                }
            }
            this.score += this.calculateRemainingPiecesScore(remainingPieces);
        }

        return {
            board: this.board,
            score: this.score,
            gameOver: gameOver,
        };
    }

    isGameOver() {
        for (let i = 0; i < this.board.length; i++) {
            if (this.board[i] === JB.WHITE) {
                continue;
            }
            if (this.getConnections(i).length > 1) {
                return false;
            }
        }
        return true;
    }

    calculateMoveScore(piecesRemoved) {
        if (piecesRemoved < 2) {
            return 0;
        }
        return piecesRemoved * (piecesRemoved - 1);
    }

    calculateRemainingPiecesScore(remainingPieces) {
        const threshold = 20;
        if (remainingPieces <= threshold) {
            const p = threshold - remainingPieces;
            return p * p;
        }
        return 0;
    }

    clamp(value, min, max) {
        if (value === undefined || !Number.isInteger(value)) {
            return min;
        }

        return Math.min(Math.max(value, min), max);
    }

    getCookie(name) {
        const cookies = document.cookie.split(';');

        for (let i = 0; i < cookies.length; i++) {
            const cookie = cookies[i].trim();
            if (cookie.startsWith(name + '=')) {
                return cookie.substring(name.length + 1);
            }
        }
        return '';
    }

    setCookie(name, value, days = 365) {
        const date = new Date();
        date.setTime(date.getTime() + (days * 24 * 60 * 60 * 1000));
        const expires = "expires=" + date.toUTCString();
        document.cookie = name + "=" + value + ";" + expires + ";path=/";
    }
}
