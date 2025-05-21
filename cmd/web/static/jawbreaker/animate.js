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
