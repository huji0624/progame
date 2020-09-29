<template>
  <div class="container">
    <el-page-header @back="$router.go(-1)" content="回放详情"> </el-page-header>
    <el-button @click="onPre()" type="primary" plain>上一轮</el-button>
    <span> 第 {{ roundNo + 1 }} 轮</span>
    <el-button @click="onNext()" type="primary" plain>下一轮</el-button>
    <el-button @click="onAutoPlay()" type="primary" plain>
      {{ loopId ? '停止播放' : '自动播放' }}
    </el-button>

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
              {{ it.Name }} - {{ it.Gold }}
            </div>
          </div>
        </div>
      </div>
    </div>
    <div class="popover">
      <el-popover
        placement="top-start"
        title="玩家信息"
        width="270"
        v-model="popShow"
      >
        <div v-if="crtPos" class="crtPos">当前坐标：{{ crtPos }}</div>
        <div v-for="(it, i) in playersInfo" :key="i">
          {{ i + 1 }}、团队：<span class="name">{{ it.Name }}</span>
          金币：
          <span class="gold">{{ it.Gold }}</span>
        </div>
      </el-popover>
    </div>
    <div class="playerlist">
      <div class="list">
        选择你关注的游戏队伍
        <br />
        <el-checkbox
          :indeterminate="isIndeterminate"
          v-model="checkAll"
          @change="onCheckAll"
          >全选</el-checkbox
        >
        <!-- <div style="margin: 15px 0"></div> -->
        <el-checkbox-group v-model="focusPlayers" @change="onChangePlayer">
          <el-checkbox v-for="it in allPlayers" :label="it" :key="it">
            {{ it }}
          </el-checkbox>
        </el-checkbox-group>
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
      title: '程序员节日活动',
    };
  },
  data() {
    return {
      roundNo: 0,
      loopId: 0,
      popShow: true,
      playerlistShow: true,
      playersInfo: [],
      crtPos: '',
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

    let allRound = [];
    res.forEach((it) => {
      allRound.push(JSON.parse(it));
    });
    const { Wid: x, Hei: y } = allRound[0];
    const mainW = x * 100 + 2,
      mainH = y * 100 + 2;

    return { allRound, x, y, mainW, mainH };
  },

  methods: {
    start() {
      let { x, y, players, roundNo, allRound, focusPlayers } = this;
      if (roundNo > allRound.length - 1) {
        return false;
      }

      const round = allRound[roundNo],
        tilemap = round.Tilemap,
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
    },
    mouseLeave() {
      this.playersInfo = [];
      this.crtPos = '';
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
      if (this.roundNo > 0) this.roundNo--;
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
      clearInterval(_this.loopId);
      _this.loopId = 0;
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
    },
  },
  created() {
    _this = this;
    _this.start();
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
.container {
  margin: 0 auto;
  min-height: 100vh;
  text-align: center;
  padding: 10px;
  .btns {
    margin: 50px;

    .btn {
      width: 20%;
    }
  }
  .main {
    width: 702px;
    background: cornsilk;
    box-shadow: 0 0 #333333;
    border: 1px solid #888;
    margin: 10px auto;
    box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
    .items {
      width: 100px;
      height: 100px;
      position: relative;
      border: 1px solid #eee;
      float: left;
      user-select: none;
      background: #dedede;
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
        color: #409eff;
        background: #ecf5ff;
        border-color: #b3d8ff;
      }
      .focus {
        color: #fff;
        background-color: #f56c6c;
        border-color: #f56c6c;
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
        opacity: 0.1;
      }
    }
    .items:hover {
      background: #eee;
    }
  }
  .popover {
    position: absolute;
    top: 80px;
    left: 20px;
    .crtPos {
      color: crimson;
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
  .playerlist {
    position: absolute;
    top: 80px;
    right: 20px;
    width: 300px;
    .list {
      position: absolute;
      background: #fff;
      min-width: 150px;
      border: 1px solid #ebeef5;
      padding: 12px;
      z-index: 2000;
      color: #606266;
      line-height: 1.4;
      text-align: justify;
      font-size: 14px;
      box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
      word-break: break-all;
    }
  }
}
</style>
