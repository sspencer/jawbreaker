function drawPowerBubble(x, y, size) {
    const ctx = state.ctx;
    ctx.save();

    // Main circle for clipping and outline
    const centerX = x + size / 2;
    const centerY = y + size / 2;
    const radius = size * 0.3;

    // Clip to circle, so color doesn't bleed out
    ctx.beginPath();
    ctx.arc(centerX, centerY, radius, 0, Math.PI * 2);
    ctx.clip();

    // Draw colored pie slices inside clip region
    const quadColors = [BLUE, GREEN, RED, YELLOW];
    for (let i = 0; i < 4; i++) {
        ctx.beginPath();
        ctx.moveTo(centerX, centerY);
        ctx.arc(
            centerX,
            centerY,
            radius,
            (Math.PI / 2) * i,
            (Math.PI / 2) * (i + 1)
        );
        ctx.closePath();
        ctx.fillStyle = COLOR_MAP.get(quadColors[i]);
        ctx.fill();
    }

    // Optional highlight
    ctx.beginPath();
    ctx.arc(centerX, centerY - radius / 2, radius / 3, 0, Math.PI * 2);
    ctx.fillStyle = "rgba(255,255,255,0.28)";
    ctx.fill();

    ctx.restore();

    const f = state.blockSize / 9;
    drawCircle(x-f, y-f, size+f*2);

}


