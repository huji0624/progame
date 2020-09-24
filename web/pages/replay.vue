<template>
  <div class="container">
    <el-page-header @back="$router.go(-1)" content="回放详情"> </el-page-header>
    <!-- <div class="btns">
      <el-input
        class="btn"
        v-model="column"
        placeholder="请输入列数"
      ></el-input>

      <el-input class="btn" v-model="rows" placeholder="请输入行数"></el-input>

      <el-button @click="onChange">确定</el-button>
    </div> -->
    <div class="main" :style="{ width: mainW + 'px', height: mainH + 'px' }">
      <div @click="onClick(it)" class="items" v-for="(it, i) in total" :key="i">
        <div class="item"></div>
        <div class="item"></div>
        <div class="item"></div>
        <div class="item"></div>
        <div class="item"></div>
        <div class="item"></div>
        <div class="item"></div>
        <div class="item"></div>
        <div class="item"></div>
      </div>
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
      column: 12, //列
      rows: 6, //行
      mainW: '',
      mainH: '',
      total: 0,
    };
  },
  async asyncData({ app, query }) {
    let data = { gid: query.gid };
    let res = await app.$axios.get('game', { params: data });
    console.log(res);
  },
  methods: {
    /**
     * 初始化
     */
    init() {
      _this.onChange();

      const afterArr = [];
      let y = _this.column;
      let x = _this.rows;

      for (let j = 0; j < x; j++) {
        for (let i = 0; i < y; i++) {
          const item = {
            type: 'init', //设置初始属性
            isCheck: false, //是否点击过
            pos: [i, j], //格子的坐标
            isRepeat: 'not', //是否递归过
            isTip: false, //用户点击数字时的提示
            coin: 10, //金币
          };
          afterArr.push(item);
        }
      }
      // console.log(afterArr);
      _this.total = afterArr;
    },

    onChange() {
      _this.mainW = _this.column * 100 + 2;
      _this.mainH = _this.rows * 100 + 2;
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
    _this.init();
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
        border: 1px solid #aaa;
        margin: 3px;
      }
    }
    :hover {
      background: #eee;
    }
  }
}
</style>
