<template>
  <div class="container">
    <el-page-header @back="$router.go(-1)" content="回放详情"> </el-page-header>
    <el-button @click="onPre()">上一轮</el-button>
    <span> 第 {{ roundNo + 1 }} 轮</span>
    <el-button @click="onNext()">下一轮</el-button>
    <el-button @click="onAutoPlay()">
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
          <div class="item" v-for="(it, i) in it.players" :key="i" v-if="i < 3">
            {{ it.Name }} - {{ it.Gold }}
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
        <div v-for="(it, i) in playersInfo" :key="i">
          {{ i + 1 }}、团队名：<span class="name">{{ it.Name }}</span>
          金币数：<span class="gold">{{ it.Gold }}</span>
        </div>
      </el-popover>
    </div>
  </div>
</template>

<script>
let _this;
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
      playersInfo: [],
      players: [],
      total: [],
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
      let { x, y, players, roundNo, allRound } = this;
      if (roundNo > allRound.length - 1) {
        alert('播放已结束');
        return false;
      }

      const round = allRound[roundNo],
        tilemap = round.Tilemap,
        afterArr = [];

      for (let i = 0; i < y; i++) {
        for (let j = 0; j < x; j++) {
          const it = tilemap[i][j];
          const item = {
            players: it.Players || [], //玩家属性
            pos: [j, i], //格子的坐标
            gold: it.Gold, //金币
          };
          afterArr.push(item);
        }
      }
      this.total = afterArr;

      console.log(round);
    },

    mouseOver(it) {
      this.playersInfo = it.players || [];
    },
    mouseLeave() {
      // this.playersInfo = [];
    },
    onNext() {
      if (this.roundNo < this.allRound.length - 1) this.roundNo++;
    },
    onPre() {
      if (this.roundNo > 0) this.roundNo--;
    },
    onAutoPlay() {
      if (!_this.loopId)
        _this.loopId = setInterval(() => {
          _this.onNext();
        }, 3000);
      else _this.onStop();
    },
    onStop() {
      clearInterval(_this.loopId);
      _this.loopId = 0;
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
        width: 85px;
        height: 25px;
        line-height: 23px;
        border: 1px solid red;
        margin: 3px;
        border-radius: 15px;
        z-index: 99;
      }
      .gold {
        width: 100px;
        height: 100px;
        position: absolute;
        top: 0;
        left: 0;
        color: #ccc;
        padding-top: 10px;
        font-size: 60px;
        font-style: italic;
        z-index: 1;
        opacity: 0.8;
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
    .name {
      font-weight: bold;
      width: 150px;
      color: #303133;
    }
    .gold {
      font-weight: bold;
      color: #929a19;
    }
  }
}
</style>
