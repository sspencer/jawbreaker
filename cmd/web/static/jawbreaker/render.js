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
		this.ctx.closePath();
	}

	drawColoredPiece(x, y, piece) {
		const size = this.blockSize;
		const angle = (225 * Math.PI) / 180;
		const diagonal = Math.sqrt(size * size + size * size);
		const startX = x + size / 2 + (Math.cos(angle) * diagonal) / 2;
		const startY = y + size / 2 + (Math.sin(angle) * diagonal) / 2;
		const endX = x + size / 2 - (Math.cos(angle) * diagonal) / 2;
		const endY = y + size / 2 - (Math.sin(angle) * diagonal) / 2;

		// Draw highlight gradient
		const gradHighlight = this.ctx.createLinearGradient(
			x,
			y,
			x + size,
			y + size
		);
		const colors = GameConfig.GRADIENT_MAP.get(this.board.getColor(piece));
		gradHighlight.addColorStop(0, colors[1]);
		gradHighlight.addColorStop(1, this.lightenColor(colors[0], 5));
		this.ctx.fillStyle = gradHighlight;
		this.ctx.beginPath();
		this.ctx.roundRect(x, y, this.pieceSize, this.pieceSize, this.cornerRadius);
		this.ctx.fill();
		this.ctx.closePath();

		// Draw body gradient
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
		this.ctx.closePath();
	}

	drawHoverEffect(x, y, piece, glows) {
		this.ctx.save(); // Keep save/restore for shadow effects
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
		this.ctx.closePath();
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
		this.ctx.beginPath();
		const padding = size * 0.2;
		this.ctx.moveTo(x + padding, y + padding);
		this.ctx.lineTo(x + size - padding, y + size - padding);
		this.ctx.moveTo(x + size - padding, y + padding);
		this.ctx.lineTo(x + padding, y + size - padding);
		this.finalizeShapeDraw();
	}

	drawPlus(x, y, size) {
		this.ctx.beginPath();
		const padding = size * 0.2;
		this.ctx.moveTo(x + size / 2, y + padding);
		this.ctx.lineTo(x + size / 2, y + size - padding);
		this.ctx.moveTo(x + padding, y + size / 2);
		this.ctx.lineTo(x + size - padding, y + size / 2);
		this.finalizeShapeDraw();
	}

	drawCircle(x, y, size) {
		this.ctx.beginPath();
		const radius = size * 0.3;
		this.ctx.arc(x + size / 2, y + size / 2, radius, 0, Math.PI * 2);
		this.finalizeShapeDraw();
	}

	drawDisc(x, y, size) {
		this.ctx.beginPath();
		const radius = size * 0.3;
		const cx = x + size / 2;
		const cy = y + size / 2;
		this.ctx.arc(cx, cy, radius, 0, Math.PI * 2);
		this.ctx.fillStyle = GameConfig.SHAPE_STROKE_COLOR;
		this.ctx.fill();
		this.ctx.strokeStyle = GameConfig.SHAPE_BORDER_COLOR;
		this.ctx.lineWidth = 2;
		this.ctx.stroke();
		this.ctx.closePath();
	}

	drawExchange(x, y, size) {
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
		this.ctx.closePath();
	}

	drawRect(x, y, size, opts) {
		this.ctx.beginPath();
		const padding = size * 0.2;
		const rectSize = size - 2 * padding;
		this.ctx.rect(x + padding, y + padding, rectSize, rectSize);
		this.finalizeShapeDraw(opts);
	}

	drawPowerFill(x, y, size) {
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
			this.ctx.closePath();
		}
		this.ctx.beginPath();
		this.ctx.rect(x + centerOffset, y + centerOffset, centerSize, centerSize);
		this.ctx.fillStyle = GameConfig.COLOR_MAP.get(GameConfig.COLORS.PURPLE);
		this.ctx.fill();
		this.ctx.closePath();
		this.drawRect(x - 8, y - 8, size + 16, "reverse");
	}

	drawRotateRight(x, y, size) {
		this.ctx.beginPath();
		const cx = x + size / 2;
		const cy = y + size / 2;
		const r = size * 0.3;
		const arrowLength = size * 0.25;
		const arrowWidth = size * 0.15;
		this.ctx.arc(cx, cy, r, Math.PI, 2 * Math.PI, false);
		const tipX = cx + r + 2;
		const tipY = cy + 3;
		this.ctx.moveTo(tipX - arrowLength, tipY - arrowWidth);
		this.ctx.lineTo(tipX, tipY);
		this.ctx.lineTo(tipX - arrowLength, tipY + arrowWidth);
		this.ctx.closePath();
		this.finalizeShapeDraw();
	}

	drawRotateLeft(x, y, size) {
		this.ctx.beginPath();
		const cx = x + size / 2;
		const cy = y + size / 2;
		const r = size * 0.3;
		const arrowLength = size * 0.25;
		const arrowWidth = size * 0.15;
		this.ctx.arc(cx, cy, r, 0, Math.PI, true);
		const tipX = cx - r - 2;
		const tipY = cy + 3;
		this.ctx.moveTo(tipX + arrowLength, tipY - arrowWidth);
		this.ctx.lineTo(tipX, tipY);
		this.ctx.lineTo(tipX + arrowLength, tipY + arrowWidth);
		this.ctx.closePath();
		this.finalizeShapeDraw();
	}

	finalizeShapeDraw(opts) {
		this.ctx.strokeStyle = GameConfig.SHAPE_BORDER_COLOR;
		this.ctx.lineWidth =
			GameConfig.SHAPE_LINE_WIDTH + GameConfig.SHAPE_BORDER_LINE_WIDTH + 1;
		this.ctx.stroke();
		this.ctx.strokeStyle = opts
			? GameConfig.SHAPE_BORDER_COLOR
			: GameConfig.SHAPE_STROKE_COLOR;
		this.ctx.lineWidth = GameConfig.SHAPE_LINE_WIDTH;
		this.ctx.stroke();
		this.ctx.closePath();
	}

	lightenColor(color, percent) {
		const num = parseInt(color.replace("#", ""), 16);
		const amt = Math.round(2.55 * percent);
		const R = Math.min((num >> 16) + amt, 255);
		const G = Math.min(((num >> 8) & 0x00ff) + amt, 255);
		const B = Math.min((num & 0x0000ff) + amt, 255);
		return `#${((1 << 24) + (R << 16) + (G << 8) + B).toString(16).slice(1)}`;
	}

	darkenColor(color, percent) {
		const num = parseInt(color.replace("#", ""), 16);
		const amt = Math.round(2.55 * percent);
		const R = Math.max((num >> 16) - amt, 0);
		const G = Math.max(((num >> 8) & 0x00ff) - amt, 0);
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
					const glows =
						powerUp > 0 && !GameConfig.POWER_CONFIG.get(powerUp).connections;
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
