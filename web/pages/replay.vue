<template>
  <div class="container">
    <div v-if="nodata" class="nodata">该局游戏不存在</div>
    <div v-else>
      <el-page-header @back="$router.go(-1)" content="回放详情">
      </el-page-header>
      <span class="note"> 第 {{ roundNo + 1 }} 轮</span>
      <div class="btns">
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
              <!-- <marquee  behavior="alternate" scrollamount="2"> -->
              <div class="item" :class="{ focus: it.isFocus }" v-if="i < 3">
                <div v-if="it.Name.length > 4">
                  <marquee scrollamount="2">
                    {{ it.Name }} - {{ it.Gold }}
                  </marquee>
                </div>
                <div v-else>{{ it.Name }} - {{ it.Gold }}</div>
              </div>
              <!-- {{ it.Name.slice(0, 4) }} - {{ it.Gold }} -->
              <!-- </marquee> -->
            </div>
          </div>
        </div>
      </div>
      <div class="popover">
        <div class="conta">
          <div class="tips">棋盘格信息</div>
          <div v-if="crtPos" class="crtPos">
            当前坐标：{{ crtPos }} 金币数量：{{ crtGold }}
          </div>
          <div v-for="(it, i) in playersInfo" :key="i">
            {{ i + 1 }}、团队：<span class="name">{{ it.Name }}</span>
            金币：
            <span class="gold">{{ it.Gold }}</span>
          </div>
        </div>
      </div>
      <div class="playerlist">
        <div class="list">
          <div class="tips">选择关注的游戏队伍 <span>最多3个</span></div>

          <!-- <el-checkbox
            :indeterminate="isIndeterminate"
            v-model="checkAll"
            @change="onCheckAll"
            >全选</el-checkbox
          > -->
          <div style="margin: 15px 0"></div>
          <el-checkbox-group
            v-model="focusPlayers"
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
  </div>
</template>

<script>
let _this;
const allPlayers = [];
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
      isIndeterminate: false,
      allPlayers: [],
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
      let { x, y, players, roundNo, allRound, focusPlayers } = this;
      if (roundNo > allRound.length - 1) return false;

      const tilemap = allRound[roundNo].Tilemap,
        afterArr = [];

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
            _this.focusPlayers.includes(it.Name) && (it.isFocus = true);
            !_this.focusPlayers.includes(it.Name) && (it.isFocus = false);
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
    onCheckAll(val) {
      this.focusPlayers = val ? this.allPlayers : [];
      this.isIndeterminate = false;
    },
    onChangePlayer(value) {
      let checkedCount = value.length;
      this.checkAll = checkedCount === this.allPlayers.length;
      this.isIndeterminate =
        checkedCount > 0 && checkedCount < this.allPlayers.length;
      this.start();
    },
  },
  created() {
    _this = this;
    if (!_this.nodata) _this.start();
    else
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
  .popover {
    position: absolute;
    top: 170px;
    left: 20px;
    width: 300px;

    .conta {
      padding: 12px;
      z-index: 2000;
      min-height: 150px;
      line-height: 1.4;
      text-align: justify;
      font-size: 14px;
      box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
      word-break: break-all;
      border-radius: 10px;
      border: 3px @bodercoler solid;
      color: #1ceaee;
      background: #1d0957;
      .crtPos {
        color: #fcf8a7;
      }
      .name {
        font-weight: bold;
        width: 150px;
        color: #409eff;
      }
      .gold {
        font-weight: bold;
        color: #929a19;
      }
    }
  }
  .playerlist {
    position: absolute;
    top: 170px;
    right: 20px;
    width: 300px;
    .list {
      position: absolute;
      min-width: 150px;
      padding: 12px;
      z-index: 2000;
      width: 300px;
      min-height: 150px;
      line-height: 1.4;
      text-align: justify;
      font-size: 14px;
      box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
      word-break: break-all;
      border-radius: 10px;
      border: 3px @bodercoler solid;

      color: #1ceaee;
      background: #1d0957;
    }
  }
  .nodata {
    color: #fff;
    width: 300px;
    margin: 350px auto;
    padding: 50px;
    border-radius: 5px;
    color: #64dbf3;
    border: @bodercoler 2px solid;
  }
  .tips {
    font-size: 20px;
    font-weight: bold;
    span {
      font-size: 13px;
      font-weight: 300;
    }
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
