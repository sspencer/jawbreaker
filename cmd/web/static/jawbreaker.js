class JB {
    // Piece constants
    static WHITE = 'w';
    static PURPLE = 'p';
    static BLUE = 'b';
    static GREEN = 'g';
    static RED = 'r';
    static YELLOW = 'y';

    // Color pieces array for random generation
    static COLOR_PIECES = [
        JB.PURPLE,
        JB.BLUE,
        JB.GREEN,
        JB.RED,
        JB.YELLOW
    ];

    // Color map for HTML class names
    static COLOR_MAP = new Map([
        [JB.PURPLE, "#8a2be2"],
        [JB.BLUE, "#00a0ff"],
        [JB.GREEN, "#00cc66"],
        [JB.RED, "#ff3333"],
        [JB.YELLOW, "#ffcc00"],
        [JB.WHITE, "#ffffff"], // was transparent
    ]);

    /**
     * Create a new Jawbreaker game
     */
    constructor(opts) {
        this.rows = opts.rows || 12;
        this.cols = opts.cols || 12;
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

        const canvasWidth = this.rows * this.blockSize +
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
     * @returns {string} - Board as a string
     */
    newBoard() {
        let board = '';
        for (let i = 0; i < this.rows * this.cols; i++) {
            const rnd = Math.floor(Math.random() * JB.COLOR_PIECES.length);
            board += JB.COLOR_PIECES[rnd];
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
            const pieces = this.getConnectedPieces(index);
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
            const pieces = this.getConnectedPieces(index);
            if (pieces.length > 1) {
                for (const i of pieces) {
                    this.hoverList.add(i);
                }
                this.renderBoard();
            } else if (this.hoverList.size > 0) {
                this.hoverList.clear();
                this.hoverIndex = -1;
                this.renderBoard();
            }
        }
    }

    renderBoard() {
        this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);
        const cornerRadius = 2;
        const pieceSize = this.blockSize - this.gap;
        const outlineGap = 2;
        const outlineSize = pieceSize-(2*outlineGap);

        for (let row = 0; row < this.rows; row++) {
            for (let col = 0; col < this.cols; col++) {
                const index = row * this.cols + col;
                const color = this.board[index];
                const x = this.border + col * (this.blockSize + this.gap);
                const y = this.border + row * (this.blockSize + this.gap);

                if (color === JB.WHITE) {
                    this.ctx.save();
                    this.ctx.strokeStyle = this.darkenColor(JB.COLOR_MAP.get(JB.WHITE), 70);
                    this.ctx.lineWidth = 1;
                    this.ctx.beginPath();
                    this.ctx.rect(x+outlineGap, y+outlineGap, outlineSize, outlineSize);
                    this.ctx.stroke();
                    this.ctx.restore();
                } else {
                    this.ctx.save();

                    const gradient = this.ctx.createLinearGradient(
                        x + this.blockSize * 0.8,
                        y + this.blockSize * 0.8,
                        x + this.blockSize * 0.1,
                        y + this.blockSize * 0.1,
                    );

                    const baseColor = JB.COLOR_MAP.get(color);
                    gradient.addColorStop(0, this.lightenColor(baseColor, 12));
                    gradient.addColorStop(1, this.darkenColor(baseColor, 6));
                    this.ctx.fillStyle = gradient;

                    this.ctx.beginPath();
                    const glow = 1;
                    const pieceX = x;
                    const pieceY = y;

                    this.ctx.roundRect(pieceX, pieceY, pieceSize, pieceSize, cornerRadius);
                    this.ctx.fill();

                    this.ctx.beginPath();
                    this.ctx.rect(pieceX + cornerRadius, pieceY, pieceSize - (2*cornerRadius), glow);
                    this.ctx.rect(pieceX, pieceY + cornerRadius, glow, pieceSize - (2*cornerRadius));
                    this.ctx.fillStyle = this.lightenColor(baseColor, 20);
                    this.ctx.fill();
                    this.ctx.beginPath();
                    this.ctx.rect(pieceX + pieceSize - glow, pieceY + cornerRadius, glow, pieceSize - (2*cornerRadius));
                    this.ctx.rect(pieceX + cornerRadius, pieceY + pieceSize - glow , pieceSize - (2*cornerRadius), glow);
                    this.ctx.fillStyle = this.darkenColor(baseColor, 20);
                    this.ctx.fill();

                    this.ctx.restore();
                }
            }
        }


        for (let row = 0; row < this.rows; row++) {
            for (let col = 0; col < this.cols; col++) {
                const index = row * this.cols + col;
                const color = this.board[index];
                const x = this.border + col * (this.blockSize + this.gap);
                const y = this.border + row * (this.blockSize + this.gap);

                if (color !== JB.WHITE) {
                    this.ctx.save();

                    const baseColor = JB.COLOR_MAP.get(color);

                    // highlight
                    if (this.hoverList.has(index)) {
                        this.ctx.shadowColor = baseColor;
                        this.ctx.shadowBlur = 15;
                        this.ctx.strokeStyle = this.lightenColor(baseColor, 40);
                        this.ctx.lineWidth = 6;

                        // Draw rounded rectangle for hover effect
                        this.ctx.beginPath();
                        const hoverSize = pieceSize + 2; // Slightly larger than the piece
                        const hoverX = x - 1;
                        const hoverY = y - 1;

                        // Draw a rounded rectangle path for the hover effect
                        this.ctx.moveTo(hoverX + cornerRadius, hoverY);
                        this.ctx.lineTo(hoverX + hoverSize - cornerRadius, hoverY);
                        this.ctx.arcTo(hoverX + hoverSize, hoverY, hoverX + hoverSize, hoverY + cornerRadius, cornerRadius);
                        this.ctx.lineTo(hoverX + hoverSize, hoverY + hoverSize - cornerRadius);
                        this.ctx.arcTo(hoverX + hoverSize, hoverY + hoverSize, hoverX + hoverSize - cornerRadius, hoverY + hoverSize, cornerRadius);
                        this.ctx.lineTo(hoverX + cornerRadius, hoverY + hoverSize);
                        this.ctx.arcTo(hoverX, hoverY + hoverSize, hoverX, hoverY + hoverSize - cornerRadius, cornerRadius);
                        this.ctx.lineTo(hoverX, hoverY + cornerRadius);
                        this.ctx.arcTo(hoverX, hoverY, hoverX + cornerRadius, hoverY, cornerRadius);

                        this.ctx.stroke();
                    }

                    this.ctx.restore();
                }
            }
        }
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

    getConnectedPieces(index) {
        if (index < 0 || index >= this.board.length) {
            return [];
        }

        const target = this.board[index];
        if (target === JB.WHITE) {
            return [];
        }

        const connectedIndices = [];
        const stack = [index];
        const visited = new Array(this.rows * this.cols).fill(false);

        while (stack.length > 0) {
            const i = stack.pop();

            if (i < 0 || i >= this.board.length || this.board[i] !== target || visited[i]) {
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

        return connectedIndices;
    }

    floodFill(index) {
        const connectedPieces = this.getConnectedPieces(index);
        if (connectedPieces.length < 2) {
            return 0;
        }

        const boardArray = this.board.split("");

        for (const i of connectedPieces) {
            boardArray[i] = JB.WHITE;
        }

        this.board = boardArray.join("");
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

        let newBoard = "";
        for (let row = 0; row < this.rows; row++) {
            for (let col = 0; col < this.cols; col++) {
                newBoard += boardArray[row][col];
            }
        }
        this.board = newBoard;
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
            if (this.getConnectedPieces(i).length > 1) {
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
