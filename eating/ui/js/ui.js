var cvs = document.getElementById('labyrinth');
var ctx = cvs.getContext('2d');

/**
 * fill rectagle
 * @param {string} color color
 * @param {int} x axis-x
 * @param {int} y axis-y
 * @param {int} w width
 * @param {int} h height
 */
function fillArea(color, x, y, w, h) {
    ctx.fillStype = color
    ctx.fillRect(x, y, w, h)
}

function drawFrame(x, y, w, h) {
    fillArea('black', x, y, w, h)
}

function drawLine(x0, y0, x1, y1) {
    ctx.lineWidth = 10
    ctx.beginPath()
    ctx.moveTo(x0, y0)
    ctx.lineTo(x1, y1)
    ctx.stroke()
}

function stroke() {
    ctx.stroke()
}

/**
 * init labyrinth
 * white 50 * 50
 * black 10 * 10
 * 
 * width/height: 50 * 10 + 10 * 11 = 610
 * (0, 0, 610, 0)
 * (0, 0, 0, 60), (300, 10, 300, 60), (600, 10, 10, 50)
 * (0, 60, 130, 10)
 */
function initLabyrinth() {
    drawLine(0, 0, 610, 0)
    drawLine(0, 0, 0, 60)
    drawLine(300, 0, 300, 60)
    drawFrame(610, 0, 600, 60)
    // drawFrame(0, 60, 130, 10)
}

initLabyrinth()