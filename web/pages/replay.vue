<template>
  <div class="container">
    <div v-if="nodata" class="noGame">该局游戏不存在</div>
    <div v-else>
      <el-page-header @back="$router.go(-1)" content="回放详情">
      </el-page-header>

      <div class="btns">
        <span class="note"> 第 {{ roundNo + 1 }} 轮</span>
        <span class="btn" @click="onPre()">上 一 轮</span>
        <span class="btn" @click="onNext()">下 一 轮</span>
        <span class="btn" @click="onAutoPlay()">
          {{ loopId ? '停止播放' : '自动播放' }}
        </span>
      </div>

      <div class="main" :style="{ width: mainW + 'px', height: mainH + 'px' }">
        <div
          @mouseover="mouseOver(it)"
          @mouseleave="mouseLeave()"
          class="items"
          v-for="(it, i) in total"
          :key="i"
        >
          <div class="gold">{{ it.gold }}</div>
          <div v-if="it.players">
            <div v-for="(it, i) in it.players" :key="i">
              <div class="item" :class="{ focus: it.isFocus }" v-if="i < 3">
                <div v-if="(it.Name + it.Gold).length > 6">
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
      <div class="wrapbox rank">
        <div class="conta">
          <div class="tips">本局游戏排名</div>
          <div v-if="allPlayersRank.length == 0" class="nodata">
            暂无游戏队伍
          </div>
          <div v-else class="list">
            <div class="crtPos" v-for="(it, i) in allPlayersRank" :key="i">
              <span class="ranking"> {{ i + 1 }}</span>
              团队：<span class="tname">{{ it.Name }}</span>
              金币：
              <span class="gold">{{ it.Gold }}</span>
            </div>
          </div>
        </div>
      </div>
      <div class="wrapbox playerlist">
        <div class="">
          <div class="tips">关注的游戏队伍 <span>最多3个</span></div>
          <div v-if="allPlayers.length == 0" class="nodata">暂无游戏队伍</div>
          <div v-else class="list">
            <el-checkbox-group
              v-model="focusChecked"
              :max="3"
              text-color="#eee"
              @change="onChangePlayer"
            >
              <el-checkbox v-for="it in allPlayers" :label="it" :key="it">
                {{ it.slice(0, 8) }}
              </el-checkbox>
            </el-checkbox-group>
          </div>
        </div>
      </div>
      <div class="wrapbox focusPlayer">
        <div class="">
          <div class="tips">关注的游戏队伍</div>
          <div v-if="focusPlayers.length == 0" class="nodata">
            请选择游戏队伍
          </div>
          <div v-else class="list">
            <div class="crtPos" v-for="(it, i) in focusPlayers" :key="i">
              <span class="ranking"> {{ i + 1 }}、</span>
              <span class="tname">{{ it.Name }}</span>
              的金币变更为
              <span class="gold">{{ it.Gold }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
let _this;
export default {
  head() {
    return {
      title: '游戏回放 - 2020程序员节日游戏',
    };
  },
  data() {
    return {
      roundNo: 0,
      loopId: 0,
      playerlistShow: true,
      playersInfo: [],
      crtPos: '',
      crtGold: '',
      players: [],
      total: [],
      isFirst: true,
      checkAll: false,
      focusPlayers: [],
      focusChecked: [],
      isIndeterminate: false,
      allPlayers: [],
      allPlayersRank: [],
    };
  },
  watch: {
    roundNo(n) {
      this.start();
    },
  },
  async asyncData({ app, query }) {
    const data = { gid: query.gid };
    const res = await app.$axios.get('game', { params: data });

    if (res instanceof Array) {
      let allRound = [];
      res.forEach((it) => {
        allRound.push(JSON.parse(it));
      });
      const { Wid: x, Hei: y } = allRound[0];
      const mainW = x * 100 + 14,
        mainH = y * 100 + 14;
      return { allRound, x, y, mainW, mainH, nodata: false };
    } else return { nodata: true };
  },

  methods: {
    start() {
      let { x, y, players, roundNo, allRound, focusChecked } = this;
      if (roundNo > allRound.length - 1) return false;

      const tilemap = allRound[roundNo].Tilemap,
        afterArr = [];
      let focusPlayers = [];

      for (let i = 0; i < y; i++) {
        for (let j = 0; j < x; j++) {
          const it = tilemap[i][j];
          const maps = it.Players || [];
          let newA = [];
          //只有第一次进来才遍历全部玩家
          if (this.isFirst && it.Players) {
            maps.map((it) => {
              this.allPlayers.push(it.Name);
            });
          }
          //添加关注玩家
          newA = maps.map((it) => {
            _this.focusChecked.includes(it.Name) &&
              ((it.isFocus = true), focusPlayers.push(it));
            !_this.focusChecked.includes(it.Name) && (it.isFocus = false);
            return it;
          });
          //关注玩家排名靠前
          newA = newA.sort(compare);
          const item = {
            players: newA || [], //玩家属性
            pos: [j, i], //格子的坐标
            gold: it.Gold, //当前格子金币
          };
          afterArr.push(item);
        }
      }
      this.isFirst = false;
      this.total = afterArr;
      _this.focusPlayers = focusPlayers.sort(sortA);
    },

    mouseOver(it) {
      this.playersInfo = it.players || [];
      this.crtPos = it.pos;
      this.crtGold = it.gold;
    },
    mouseLeave() {
      // this.playersInfo = [];
      // this.crtPos = '';
    },
    onNext() {
      if (this.roundNo < this.allRound.length - 1) this.roundNo++;
      else if (_this.loopId) {
        this.$notify({
          message: '播放结束',
          type: 'success',
          center: true,
        });
        this.onStop();
      }
    },
    onPre() {
      this.roundNo > 0 && this.roundNo--;
    },
    onAutoPlay() {
      if (!_this.loopId) {
        this.$notify({
          message: '播放开始',
          type: 'success',
          center: true,
        });
        _this.loopId = setInterval(() => {
          _this.onNext();
        }, 3000);
      } else _this.onStop();
    },
    onStop() {
      if (_this.loopId) {
        clearInterval(_this.loopId);
        _this.loopId = 0;
      }
    },

    onChangePlayer(value) {
      this.start();
    },
    getPlayersGold() {
      const { x, y, allRound } = this;
      const tilemap = allRound[allRound.length - 1].Tilemap,
        afterArr = [];
      for (let i = 0; i < y; i++) {
        for (let j = 0; j < x; j++) {
          const it = tilemap[i][j];
          const maps = it.Players || [];
          if (it.Players) {
            maps.map((it) => {
              afterArr.push(it);
            });
          }
        }
      }
      this.allPlayersRank = afterArr.sort(sortA);
      console.log(this.allPlayersRank);
    },
  },
  created() {
    _this = this;
    if (!_this.nodata) {
      _this.getPlayersGold();
      _this.start();
    } else
      setTimeout((_) => {
        _this.$router.go(-1);
      }, 3000);
  },
  beforeDestroy() {
    _this.$once('hook:beforeDestroy', () => {
      _this.onStop();
    });
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
    margin: 30px auto;

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
    margin: 10px;
    margin-left: 400px;
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
    left: 20px;
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
  .playerlist {
    top: 400px;
  }
  .focusPlayer {
    top: 270px;
    .ranking {
      font-weight: bold;
    }
  }
  .focusOver {
    top: 70px;
    .crtPos {
      text-align: left;
      color: #fcf8a7;
      line-height: 30px;
    }
    .list {
      max-height: 100px;
    }
  }
  .rank {
    top: 675px;
    .ranking {
      font-size: 16px;
      color: #fcf8a7;
      display: inline-block;
      font-weight: bold;
      width: 22px;
    }
  }
  .tname {
    display: inline-block;
    font-size: 14px;
    font-weight: bold;
    width: 94px;
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
  line-height: 24px;
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
