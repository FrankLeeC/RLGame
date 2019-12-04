import DataBus    from './databus'

let base = "http://localhost:9091"

let databus = new DataBus()

const white = 36
const black = 6
const pair = white + black
const border = 10 * white + 11 * black
let top = parseInt(window.innerHeight / 4)
let left = (window.innerWidth - border) / 2

let runSpeed = 20
let drawSpeed = 5

let canvas = document.createElement('canvas')
canvas.width = border
canvas.height = border
let ctx   = canvas.getContext('2d')

canvas.style.position = 'fixed'
canvas.style.left = left + 'px'
canvas.style.top = top + 'px'
document.body.appendChild(canvas)

var labyrinthCanvas = document.createElement('canvas')
labyrinthCanvas.width = window.innerWidth
labyrinthCanvas.height = window.innerHeight
var lctx = labyrinthCanvas.getContext('2d')

var pointWidth = parseInt(window.innerWidth / 5)
var pointHeight = parseInt(window.innerHeight / 24)
var pointCanvas = document.createElement('canvas')
pointCanvas.width = pointWidth
pointCanvas.height = pointHeight
var pctx = pointCanvas.getContext('2d')
var pleft = pointCanvas.width * 2
var ptop = parseInt(window.innerHeight / 6)
pointCanvas.style.position = 'fixed'
pointCanvas.style.top = ptop + 'px'
pointCanvas.style.left = pleft + 'px'
document.body.appendChild(pointCanvas)



var bufferedPointCanvas = document.createElement('canvas')
var bpctx = bufferedPointCanvas.getContext('2d')
bufferedPointCanvas.width = pointWidth
bufferedPointCanvas.height = pointHeight

var targetCanvas = document.createElement('canvas')
targetCanvas.width = border
targetCanvas.height = border
var tctx = targetCanvas.getContext('2d')

var aiTargetImg = new Image()
aiTargetImg.src = base + '/static/images/bean_ai.png'
var playerTargetImg = new Image()
playerTargetImg.src = base + '/static/images/bean_player.png'
var aiImg0 = new Image()
aiImg0.src = base + '/static/images/eater_ai_up.png'
var aiImg1 = new Image()
aiImg1.src = base + '/static/images/eater_ai_right.png'
var aiImg2 = new Image()
aiImg2.src = base + '/static/images/eater_ai_down.png'
var aiImg3 = new Image()
aiImg3.src = base + '/static/images/eater_ai_left.png'
var playerImg0 = new Image()
playerImg0.src = base + '/static/images/eater_player_up.png'
var playerImg1 = new Image()
playerImg1.src = base + '/static/images/eater_player_right.png'
var playerImg2 = new Image()
playerImg2.src = base + '/static/images/eater_player_down.png'
var playerImg3 = new Image()
playerImg3.src = base + '/static/images/eater_player_left.png'


/**
 * 游戏主函数
 */
export default class Main {
  
  constructor() {
    // 维护当前requestAnimationFrame的id
    this.aniId    = 0
    this.initTouchEvent()
    this.initLabyrinth()
    this.restart()
  }

  restart() {
    canvas.removeEventListener(
      'touchstart',
      this.touchHandler
    )
    
    this.bindLoop = this.loop.bind(this)
    //   // 清除上一局的动画
    window.cancelAnimationFrame(this.aniId);

    this.aniId = window.requestAnimationFrame(
      this.bindLoop,
      canvas
    )
  }

  initTouchEvent() {
    document.onkeydown = function(){
      var e = window.event;
      if (e.keyCode == 87) {
        databus.setPlayerAction(0)
      } else if (e.keyCode == 68) {
        databus.setPlayerAction(1)
      } else if (e.keyCode == 83) {
        databus.setPlayerAction(2)
      } else if (e.keyCode == 65) {
        databus.setPlayerAction(3)
      } else if (e.keyCode == 32) {
        databus.setPlayerAction(-1)
      }
    }
  }

  initLabyrinth() {
    lctx.lineWidth = black / 2

    this.setBackground()
    this.drawBorder()
    this.drawLabyrinth()
  }

  drawLine(arr) {
    var start = arr[0]
    var a = start[0] * pair
    if (start[0] == 10) {
      a += black
    }
    var b = start[1] * pair
    if (start[1] == 10) {
      b += black
    }
    lctx.moveTo(a, b)
    for(var i = 1; i < arr.length; i++) {
      var a = arr[i][0] * pair
      if (arr[i][0] == 10) {
        a += black
      }
      var b = arr[i][1] * pair
      if (arr[i][1] == 10) {
        b += black
      }
      lctx.lineTo(a, b)
    }
  }

  drawLabyrinth() {
    this.drawLine([[5, 0], [5, 1], [4, 1], [4, 3]])
    this.drawLine([[0, 1], [2, 1], [2, 2], [3, 2], [3, 1]])
    this.drawLine([[7, 1], [10, 1]])
    this.drawLine([[6, 1], [6, 3]])
    this.drawLine([[5, 2], [9, 2], [9, 6]])
    this.drawLine([[1, 2], [1, 3], [5, 3]])
    this.drawLine([[8, 3], [7, 3], [7, 8]])
    this.drawLine([[3, 3], [3, 5]])
    this.drawLine([[3, 4], [5, 4]])
    this.drawLine([[0, 4], [2, 4], [2, 5], [1, 5]])
    this.drawLine([[8, 4], [6, 4], [6, 5], [4, 5], [4, 7]])
    this.drawLine([[9, 5], [8, 5], [8, 7], [10, 7]])
    this.drawLine([[1, 6], [4, 6]])
    this.drawLine([[2, 6], [2, 9]])
    this.drawLine([[0, 7], [1, 7]])
    this.drawLine([[1, 8], [1, 10]])
    this.drawLine([[3, 7], [3, 10]])
    this.drawLine([[3, 8], [5, 8], [5, 6]])
    this.drawLine([[4, 9], [6, 9], [6, 8], [8, 8], [8, 9], [9, 9], [9, 8]])
    this.drawLine([[5, 9], [5, 10]])
    this.drawLine([[7, 9], [7, 10]])    
    lctx.stroke()

  }

  drawBorder() {
    lctx.strokeStyle = 'black'
    lctx.strokeRect(0, 0, border, border)
  }

  /**
   * 设置背景
   */
  setBackground() {
    lctx.fillStyle = 'white'
    lctx.fillRect(0, 0, labyrinthCanvas.width, labyrinthCanvas.height)
  }

  // 游戏结束后的触摸事件处理逻辑
  touchEventHandler(e) {
     e.preventDefault()
  }



  clearAITarget() {
    var d = 20
    var x = databus.aiTarget % 10
    var y = parseInt(databus.aiTarget / 10)
    var l = (x + 1) * black + x * white
    var t = (y + 1) * black + y * white

    tctx.clearRect(l, t, d, d)
  }

  clearPlayerTarget() {
    var d = 20
    var x = databus.playerTarget % 10
    var y = parseInt(databus.playerTarget / 10)
    var l = (x + 1) * black + x * white
    var t = (y + 1) * black + y * white
    

    tctx.clearRect(l, t, d, d)
  }

  drawAITarget() {
    var aiTarget = databus.aiTarget
    var tx = aiTarget % 10
    var ty = parseInt(aiTarget / 10)
    var tl = (tx + 1) * black + tx * white
    var tt = (ty + 1) * black + ty * white

    tctx.drawImage(aiTargetImg, tl, tt)
    ctx.drawImage(targetCanvas, 0, 0)
  }

  drawPlayerTarget() {
    var playerTarget = databus.playerTarget
    var tx = playerTarget % 10
    var ty = parseInt(playerTarget / 10)
    var tl = (tx + 1) * black + tx * white
    var tt = (ty + 1) * black + ty * white

    tctx.drawImage(playerTargetImg, tl, tt)
    ctx.drawImage(targetCanvas, 0, 0)
  }

  drawTarget() {
    if (databus.aiLastTarget != databus.aiTarget) {
      this.clearAITarget()
    }
    this.drawAITarget()

    if (databus.playerLastTarget != databus.playerTarget) {
      this.clearPlayerTarget()
    }
    this.drawPlayerTarget()
  }

  drawAI() {
    var nextState = databus.aiState
    var x = nextState % 10
    var y = parseInt(nextState / 10)
    var l = (x + 1) * black + x * white
    var t = (y + 1) * black + y * white

    if (databus.aiAction === 0) {
      tctx.drawImage(aiImg0, l, t)
    } else if (databus.aiAction === 1) {
      tctx.drawImage(aiImg1, l, t)
    } else if (databus.aiAction === 2) {
      tctx.drawImage(aiImg2, l, t)
    } else {
      tctx.drawImage(aiImg3, l, t)
    }
    ctx.drawImage(targetCanvas, 0, 0)
  }

  drawPlayer() {
    var action = -1
    if (databus.playerAction == -1) {
      action = databus.playerLastAction
    } else {
      action = databus.playerAction
    }
    var nextState = databus.playerState
    var x = nextState % 10
    var y = parseInt(nextState / 10)
    var l = (x + 1) * black + x * white
    var t = (y + 1) * black + y * white

    if (action === 0) {
      tctx.drawImage(playerImg0, l, t)
    } else if (action === 1) {
      tctx.drawImage(playerImg1, l, t)
    } else if (action === 2) {
      tctx.drawImage(playerImg2, l, t)
    } else {
      tctx.drawImage(playerImg3, l, t)
    }
    ctx.drawImage(targetCanvas, 0, 0)
  }

  drawPoint() {
    pctx.clearRect(0, 0, pointCanvas.width, pointCanvas.height)
    if (databus.shouldDrawPoint()) {
      
      bpctx.clearRect(0, 0, bufferedPointCanvas.width, bufferedPointCanvas.height)

      var rate = databus.getPointRate()
      var start = pointWidth * (rate - 0.05)
      var end = pointWidth * (rate + 0.05)
      var pgrd = pctx.createLinearGradient(start, 0, end, 0)

      pgrd.addColorStop(0, "red");
      pgrd.addColorStop(1, "blue");

      bpctx.fillStyle = pgrd;
      bpctx.fillRect(0, 0, pointWidth, pointHeight);

      databus.resetLastRate()
    }
    pctx.drawImage(bufferedPointCanvas, 0, 0)
  }

  clearAll() {
    ctx.clearRect(0, 0, canvas.width, canvas.height)
    tctx.clearRect(0, 0, targetCanvas.width, targetCanvas.height)
  }

  doRender() {
    this.clearAll()
    ctx.drawImage(labyrinthCanvas, 0, 0)
    this.drawAI()  
    this.drawPlayer()
    this.drawTarget()
    this.drawPoint()
  }

  /**
   * canvas重绘函数
   * 每一帧重新绘制所有的需要展示的元素
   */
  render() {
    this.doRender()
    

    // 游戏结束停止帧循环
    if (databus.gameOver) {
      // this.gameinfo.renderGameOver(ctx, databus.score)

      if ( !this.hasEventBind ) {
        this.hasEventBind = true
        this.touchHandler = this.touchEventHandler.bind(this)
        canvas.addEventListener('touchstart', this.touchHandler)
      }
    }
  }

  // 游戏逻辑更新主函数
  update() {
    databus.checkGameOver()
    if (databus.gameOver)
      return;

    databus.nextAIState()
    databus.nextPlayerState()

    if (databus.isAIOver()) {
      databus.randomAITarget()
    }

    if (databus.isPlayerOver()) {
      databus.randomPlayerTarget()
    } 
  }

  // 实现游戏帧循环
  loop() {
    databus.frame += 1
    if (databus.isPrepared()) {
      if (databus.frame % runSpeed === 0) {
        this.update()
      }
      if (databus.frame % drawSpeed === 0) {
        this.render()
      }
    }


    this.aniId = window.requestAnimationFrame(
      this.bindLoop,
      canvas
    )
  }
}

new Main()

