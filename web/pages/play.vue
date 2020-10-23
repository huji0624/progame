<template>
  <div class="container">
    <span class="note"> 第 {{ info.RoundID + 1 }} 轮</span>
    <div class="btns">
      <span class="btn" @click="init()"> 开始游戏 </span>
      <span class="btn" @click="onClose">结束游戏</span>
    </div>

    <div class="main" :style="{ width: mainW + 'px', height: mainH + 'px' }">
      <div
        @mouseover="mouseOver(it)"
        @click="onClick(it)"
        class="items"
        v-for="(it, i) in total"
        :key="i"
      >
        <div class="gold">{{ it.gold }}</div>
        <div v-if="it.players">
          <div v-for="(it, i) in it.players" :key="i">
            <div class="item" :class="{ focus: it.isFocus }" v-if="i < 3">
              <div v-if="(it.Name + it.Gold).length > 8">
                <marquee scrollamount="2">
                  {{ it.Name }} - {{ it.Gold }}
                </marquee>
              </div>
              <div v-else>{{ it.Name }} - {{ it.Gold }}</div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 悬浮盒子 -->
    <div class="wrapbox focusOver">
      <div class="conta">
        <div class="tips">棋盘格信息</div>
        <div v-if="!crtPos" class="nodata">鼠标移过棋盘，展示数据</div>
        <div v-else class="crtPos">
          当前坐标：{{ crtPos }} 金币数量：{{ crtGold }}
        </div>
        <div class="list">
          <div v-for="(it, i) in playersInfo" :key="i">
            {{ i + 1 }}、团队：<span class="tname">{{ it.Name }}</span>
            金币：
            <span class="gold">{{ it.Gold }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
let _this;
let name = 'token';
let token = 'token';
let wsurl = 'ws://localhost:8881/ws';
let socket;

export default {
  head() {
    return {
      title: '游戏回放 - 2020程序员节日游戏',
    };
  },
  data() {
    return {
      isClick: false,
      socket: socket,
      x: 0,
      y: 0,
      mainW: 0,
      mainH: 0,
      loopId: 0,
      playersInfo: [],
      crtPos: '',
      crtGold: '',
      total: [],
      isFirst: true,
      info: { RoundID: 0 },
    };
  },

  methods: {
    start() {
      const { x, y, info } = this;
      const { Tilemap: tilemap } = info,
        afterArr = [];

      for (let i = 0; i < y; i++) {
        for (let j = 0; j < x; j++) {
          const it = tilemap[i][j];
          const maps = it.Players || [];
          let newA = [];
          newA = maps.map((it) => {
            if (it.Name == name) it.isFocus = true;
            return it;
          });

          const item = {
            players: newA || [], //玩家属性
            pos: [j, i], //格子的坐标
            gold: it.Gold, //当前格子金币
          };
          afterArr.push(item);
        }
      }
      this.total = afterArr;
    },

    mouseOver(it) {
      this.playersInfo = it.players || [];
      this.crtPos = it.pos;
      this.crtGold = it.gold;
    },

    init() {
      if (typeof WebSocket !== 'undefined') {
        socket = new WebSocket(wsurl);
        // 监听socket连接
        socket.onopen = this.onOpen;
        // 监听socket错误信息
        socket.onerror = this.onError;
        // 监听socket消息
        socket.onmessage = this.onMessage;
        // 监听socket关闭
        socket.onclose = this.close;
      }
    },

    onOpen() {
      console.log('socket连接成功');
      socket.send(JSON.stringify({ msgtype: 0, token: token }));
    },

    onMessage(evt) {
      let received_msg = evt.data;

      let jmsg = JSON.parse(received_msg);
      // console.log('recv msg:', jmsg);
      if (jmsg.Msgtype == 0) {
        console.log('login ok.');
      } else if (jmsg.Msgtype == 1) {
        socket.send(JSON.stringify({ msgtype: 2, token: token }));
      } else if (jmsg.Msgtype == 3) {
        _this.info = _this.$deepCopy(jmsg);
        if (_this.isFirst) {
          _this.x = jmsg.Wid;
          _this.y = jmsg.Hei;
          _this.x = jmsg.Wid;
          _this.mainW = _this.x * 100 + 14;
          _this.mainH = _this.y * 100 + 14;
          _this.isFirst = false;
        }
        _this.start();
        _this.loopId = setTimeout(() => {
          if (_this.isClick) clearTimeout(_this.loopId);
          else {
            let x = Math.floor(Math.random() * jmsg.Wid);
            let y = Math.floor(Math.random() * jmsg.Hei);
            socket.send(
              JSON.stringify({
                msgtype: 4,
                token: token,
                x: x,
                y: y,
                RoundID: _this.info.RoundID,
              })
            );
          }
        }, 700);
      }
    },

    onSend(_data) {},

    onError() {
      console.log('socket发生了错误');
    },

    onClose() {
      socket.close();
      console.log('socket已经关闭,onClose');
    },
    close() {
      console.log('socket已经关闭,close');
    },

    onClick(it) {
      const [x, y] = it.pos;
      console.log(x, y);
      _this.isClick = true;
      socket.send(
        JSON.stringify({
          msgtype: 4,
          token: token,
          x: x,
          y: y,
          RoundID: _this.info.RoundID,
        })
      );
    },
  },
  created() {
    _this = this;
  },
  destroyed() {
    // 销毁监听
    socket.close();
  },
};
function compare(a, b) {
  if (a.isFocus && !b.isFocus) return -1;
}
function sortA(a, b) {
  if (a.Gold > b.Gold) return -1;
}
</script>

<style lang="less" scoped>
@bodercoler: #1dfefe;
.container {
  margin: 0 auto;
  min-height: 100vh;
  text-align: center;
  padding: 10px;
  color: #fff;
  .note {
    font-weight: bold;
    color: @bodercoler;
  }
  .btns {
    margin: 50px;

    .btn {
      color: #04def0;
      font-weight: bold;
      font-family: '微软雅黑';
      padding: 15px 32px;
      margin: 0 10px;
      user-select: none;
      background-image: url('../assets/images/btn.png');
      background-size: 100% 100%;
      cursor: pointer;
    }
    .btn:hover {
      // background: #04def0;
      color: #fff;
    }
  }
  .main {
    width: 702px;
    background: #330a66;
    box-shadow: 0 0 #333333;
    border-radius: 10px;
    border: 7px @bodercoler solid;
    margin: 10px auto;
    box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
    .items {
      width: 100px;
      height: 100px;
      position: relative;
      border: 1px solid @bodercoler;
      float: left;
      user-select: none;
      background: #330a66;
      transition: all 0.2s;
      font-size: 14px;
      display: flex;
      flex-wrap: wrap;
      justify-content: center;
      align-items: center;
      .item {
        width: 95px;
        height: 25px;
        line-height: 23px;
        border: 1px solid red;
        margin: 3px;
        border-radius: 15px;
        z-index: 99;
        color: #1eeaf0;
        background: #635393;
        border: 1px @bodercoler solid;
      }
      .focus {
        color: #0e025e;
        font-weight: bold;
        background-color: #edcc53;
        border-color: #929a19;
      }
      .gold {
        width: 100px;
        height: 100px;
        position: absolute;
        top: 0;
        left: 0;
        color: #969494;
        padding-top: 10px;
        font-size: 60px;
        font-style: italic;
        z-index: 1;
        opacity: 0.3;
      }
    }
    .items:hover {
      background: #eee;
    }
  }

  .wrapbox {
    position: absolute;
    top: 170px;
    right: 20px;
    width: 300px;
    font-size: 14px;
    padding: 12px;
    color: #1ceaee;
    word-break: break-all;
    border-radius: 10px;
    border: 3px @bodercoler solid;
    background: #1d0957;
    box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
    .tips {
      font-size: 20px;
      font-weight: bold;
      line-height: 30px;
      span {
        font-size: 13px;
        font-weight: 300;
      }
    }
    .list {
      min-width: 150px;
      z-index: 2000;
      width: 300px;
      max-height: 200px;
      line-height: 1.4;
      text-align: justify;
      overflow-y: auto;
    }
  }

  .focusOver {
    top: 170px;
    left: 20px;
    .crtPos {
      text-align: left;
      color: #fcf8a7;
      line-height: 30px;
    }
  }

  .tname {
    display: inline-block;
    font-size: 16px;
    font-weight: bold;
    width: 90px;
    color: #add5fd;
  }
  .gold {
    font-weight: bold;
    color: #ffeb00;
  }
  .nodata {
    line-height: 50px;
  }

  .noGame {
    color: #fff;
    width: 300px;
    margin: 350px auto;
    padding: 50px;
    border-radius: 5px;
    color: #64dbf3;
    border: @bodercoler 2px solid;
  }
}
</style>

<style lang="less">
.el-checkbox__input.is-checked .el-checkbox__inner {
  background-color: #04def0;
  border-color: #04def0;
}
.el-checkbox__input.is-checked + .el-checkbox__label {
  color: #04def0;
}
.el-checkbox {
  color: #04def0;
  line-height: 30px;
}
.el-checkbox__input.is-indeterminate .el-checkbox__inner {
  background-color: #04def0;
  border-color: #04def0;
}
.el-page-header__content {
  font-size: 18px;
  color: #04def0;
}
.el-checkbox__label {
  width: 90px;
}
</style>
