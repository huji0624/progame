<template>
  <div class="container">
    <el-page-header @back="$router.go(-1)" content="回放详情"> </el-page-header>
    <el-button @click="onPre()">Pre</el-button>
    <el-button @click="onNext()">Next</el-button>

    <div class="main" :style="{ width: mainW + 'px', height: mainH + 'px' }">
      <div
        @mouseover="mouseOver(it)"
        @mouseleave="mouseLeave()"
        @click="onClick(it)"
        class="items"
        v-for="(it, i) in total"
        :key="i"
      >
        <div class="item">
          {{ it.gold }}
        </div>
        <!-- <div class="item"></div>
        <div class="item"></div>
        <div class="item"></div>
        <div class="item"></div>
        <div class="item"></div>
        <div class="item"></div>
        <div class="item"></div>
        <div class="item"></div> -->
      </div>
    </div>
    <div class="popover">
      <el-popover
        placement="top-start"
        title="玩家信息"
        width="100%"
        v-model="popShow"
      >
        <div v-for="(it, i) in playersInfo" :key="i">
          <span>团队名：{{ it.Name }}</span>
          <span>金币数：{{ it.Gold }}</span>
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
      popShow: true,
      playersInfo: [],
      players: [],
      total: [],
    };
  },
  watch: {
    roundNo(n) {
      console.log(n);
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
        tilemap = round.Tilemap;
      const afterArr = [];
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
      roundNo++;
      console.log(round);
    },

    mouseOver(it) {
      this.playersInfo = it.players || [];
    },
    mouseLeave() {
      this.playersInfo = [];
    },
    onNext() {
      this.roundNo++;
    },
    onPre() {
      if (this.roundNo > 0) this.roundNo--;
    },
    onClick(it) {
      alert(it.pos);
    },

    /**
     * 根据坐标计算索引
     * @param {array}
     * @return {number}
     */
    calcIdx([a, b]) {
      const row = this.boardSize[0];
      return b * row + a;
    },
  },
  created() {
    _this = this;
    this.start();
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
      display: flex;
      flex-wrap: wrap;
      justify-content: center;
      align-items: center;
      .item {
        width: 25px;
        height: 25px;
        border-radius: 50%;
        // border: 1px solid #aaa;
        margin: 3px;
      }
    }
    :hover {
      background: #eee;
    }
  }
}
</style>
