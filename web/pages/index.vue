<template>
  <div class="container">
    <div class="btns">
      <el-input
        class="btn"
        v-model="column"
        placeholder="请输入列数"
      ></el-input>

      <el-input class="btn" v-model="rows" placeholder="请输入行数"></el-input>

      <el-button @click="onChange">确定</el-button>
    </div>
    <div class="main" :style="{ width: mainW + 'px', height: mainH + 'px' }">
      <div
        @click="onClick(it)"
        class="item"
        v-for="(it, i) in total"
        :key="i"
      ></div>
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
      column: 5, //列
      rows: 5, //行
      mainW: '',
      mainH: '',
      total: 25,
    };
  },
  methods: {
    init() {
      _this.onChange();
      const afterArr = [];
      let y = _this.column;
      let x = _this.rows;
      for (let i = 0; i < y; i++) {
        for (let j = 0; j < x; j++) {
          const item = {
            type: 'init', //设置初始属性
            isCheck: false, //是否点击过
            pos: [j, i], //格子的坐标
            isRepeat: 'not', //是否递归过
            isTip: false, //用户点击数字时的提示
            isFlag: false, //用户是否插了旗子
          };
          afterArr.push(item);
        }
      }
      _this.total = afterArr;
    },
    onChange() {
      _this.mainW = _this.column * 100 + 2;
      _this.mainH = _this.rows * 100 + 2;
      _this.total = _this.column * _this.rows;
    },
    onClick(it) {
      alert(it.pos);
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
    .item {
      width: 100px;
      height: 100px;
      border: 1px solid #eee;
      float: left;
      user-select: none;
      background: #888;
      transition: all 0.2s;
    }
    :hover {
      background: #eee;
    }
  }
}
</style>
