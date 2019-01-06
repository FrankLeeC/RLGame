import axios from 'axios'

let instance

let policy = new Map()
let transition = new Map()

/**
 * 全局状态管理器
 */
export default class DataBus {

  constructor() {
    if ( instance )
      return instance

    instance = this


    this.reset()
    this.prepared = 0
    this.aiOver = true
    this.playerOver = true

    this.initAIState()
    this.randomAITarget()
    this.initPlayerState()
    this.randomPlayerTarget()
    this.initPolicy()
    this.initTransition()
  }

  initAIState() {
    if (this.aiState < 0) {
      this.aiState = parseInt(Math.random() * 100)
      this.aiLastState = this.aiState
    }
  }

  initPlayerState() {
    if (this.playerState < 0) {
      this.playerState = parseInt(Math.random() * 100)
      this.playerLastState = this.playerState
    }
  }

  initPolicy() {
    var that = this
    axios.get('https://eater.liwanyi.me/policy')
      .then(function (res) {
        for (var k in res.data) {
          var v = res.data[k]
          var t = new Map()
          for (var s in v) {
            t.set(parseInt(s), parseInt(v[s]))
          }
          policy.set(parseInt(k), t)
        }
        that.prepared += 1
      })
      .catch(function (error) {
        console.log(error);
      });
  }

  initTransition() {
    var that = this
    axios.get('https://eater.liwanyi.me/transition')
      .then(function (res) {
        for (var k in res.data) {
          var v = res.data[k]
          transition.set(parseInt(k), v)
        }
        that.prepared += 1
      })
      .catch(function (error) {
        console.log(error);
      });
  }

  /**
   * 数据是否准备完毕
   */
  isPrepared() {
    return this.prepared >= 2
  }

  /**
   * 重置
   */
  reset() {
    this.frame      = 0
    this.aiState = -1
    this.aiLastState = -1
    this.aiAction = -1
    this.aiTarget = -1

    this.playerState = -1
    this.playerLastState = -1
    this.playerAction = -1
    this.playerTarget = -1
    this.playerLastAction = this.playerAction
    this.playerLastTouchTime = 0
    this.playerShouldStop = false

    this.playerPoint = 0
    this.aiPoint = 0

    this.gameOver = false
    this.rate = 0.5
    this.lastRate = 0.0
  }

  getPlayerPoint() {
    if (this.playerPoint > 9) {
      return this.playerPoint
    }
    return '0' + this.playerPoint
  }

  getAIPoint() {
    if (this.aiPoint > 9) {
      return this.aiPoint
    }
    return '0' + this.aiPoint
  }

  /**
   * 随机一个AI目标
   */
  randomAITarget() {
    this.aiLastTarget = this.aiTarget
    while (true) {
      var r = parseInt(Math.random() * 100)
      if (r != this.playerTarget && r != this.aiState && r != this.playerState) {
        this.aiTarget = r
        break
      }
    }
  }

  getPointRate() {
    return this.rate
  }

  calculateRate() {
    if (this.aiPoint > this.playerPoint) {
      var d = this.aiPoint - this.playerPoint
      this.rate = (50 + d) / 100
    } else if (this.aiPoint < this.playerPoint) {
      var d = this.playerPoint - this.aiPoint
      this.rate = (50 - d) / 100
    }
  }

  shouldDrawPoint() {
    return this.lastRate != this.rate
  }

  resetLastRate() {
    this.lastRate = this.rate
  }

  checkGameOver() {
    if (this.aiPoint >= 50 || this.playerPoint >= 50) {
      this.gameOver = true
    }
    this.gameOver = false
    return this.gameOver
  }
  
  /**
   * 随机一个玩家目标
   */
  randomPlayerTarget() {
    this.playerLastTarget = this.playerTarget
    while (true) {
      var r = parseInt(Math.random() * 100)
      if (r != this.aiTarget && r != this.playerState && r != this.aiState) {
        this.playerTarget = r
        break
      }
    }
  }

  /**
   * 获取AI下一个状态
   */
  nextAIState() {
    this.aiAction = this.getAction(this.aiState, this.aiTarget)
    this.aiLastState = this.aiState
    this.aiState = this.getState(this.aiState, this.aiAction)
    return this.aiState
  }

  nextPlayerState() {
    if (this.playerAction == -1) {  // 静止
      return this.playerState
    } else {
      this.playerLastState = this.playerState
      this.playerState = this.getState(this.playerState, this.playerAction)
      return this.playerState
    }
  }

  /**
   * 设置玩家的移动方向
   */
  setPlayerAction(a) {
    this.playerLastAction = this.playerAction
    this.playerAction = a
  }

  /**
   * AI是否已经结束当前轮次
   */
  isAIOver() {
    this.aiOver = this.aiState == this.aiTarget
    if (this.aiOver) {
      this.aiPoint += 1
      this.calculateRate()
    }
    return this.aiOver
  }

  /**
   * 玩家是否已经结束当前轮次
   */
  isPlayerOver() {
    this.playerOver = this.playerState == this.playerTarget
    if (this.playerOver) {
      this.playerPoint += 1
      this.calculateRate()
    }
    return this.playerOver
  }

  /**
   * 获取当前状态下要执行的动作
   */
  getAction(state, target) {
    return policy.get(target).get(state)
  }

  /**
   * 获取状态s下执行动作a后，下一个状态
   */
  getState(s, a) {
    return transition.get(s)[a]
  }

}
