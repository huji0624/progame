<template>
  <!-- <section class="main" @mousemove="clickFullscreen"> -->
  <section class="main">
    <!-- <a
      href="https://github.com/huji0624/progame"
      target="_blank"
      class="logo"
      title="查看抢钱大作战的玩法和帮助"
      >抢钱大作战</a
    > -->
    <div class="logo" @click="toURL" title="查看抢钱大作战的玩法和帮助"></div>
    <div class="head"></div>
    <div class="isall">
      <el-popover
        placement="top"
        width="300"
        trigger="hover"
        content="筛选具有排名资格的玩家，即最近游戏中参与超过64局游戏的游戏玩家。"
      >
        <el-checkbox v-model="isAll" slot="reference">
          查看资格对局排名
        </el-checkbox>
      </el-popover>
    </div>
    <div class="list">
      <el-row class="rows">
        <el-col :span="6" v-for="(item, i) in tabs" :key="i">
          <div class="name">{{ item.label }}</div>
          <el-row class="listbody">
            <div class="title">
              <el-row>
                <el-col :span="6"> 名次 </el-col>
                <el-col :span="10"> 队名 </el-col>
                <el-col :span="8"> 得分 </el-col>
              </el-row>
            </div>
            <div class="row">
              <div
                v-if="!All[item.name] || All[item.name].length == 0"
                class="notstart"
              >
                敬请期待
              </div>
              <div v-else>
                <el-row class="item" v-for="(it, i) in All[item.name]" :key="i">
                  <el-col :span="6" title="Nihaoyo ">
                    <img
                      class="img"
                      v-if="i == 0"
                      src="../assets/images/1.png"
                    />
                    <img
                      class="img"
                      v-if="i == 1"
                      src="../assets/images/2.png"
                    />
                    <img
                      class="img"
                      v-if="i == 2"
                      src="../assets/images/3.png"
                    />
                    <div v-if="i > 2">{{ i + 1 }}</div>
                  </el-col>
                  <el-col :span="10">
                    {{ it.Name.slice(0, 12) }}
                  </el-col>
                  <el-col :span="8">
                    {{ it.Gold }}
                  </el-col>
                </el-row>
              </div>
            </div>
          </el-row>
        </el-col>
      </el-row>
    </div>
    <div class="replay">
      <div class="title">点击按钮查看当局回放 <span>仅展示最近30局</span></div>
      <div v-if="Gid < 1" class="notstart">敬请期待</div>
      <div v-else class="btns">
        <span v-for="(item, i) in 30" :key="i">
          <el-button
            type="primary"
            plain
            v-if="Gid - i - 1 > 0"
            class="btn"
            @click="$router.push('/replay?gid=' + (Gid - i - 1))"
          >
            第{{ Gid - i - 1 }}局
          </el-button>
        </span>
      </div>
    </div>
  </section>
</template>

<script>
let _this;
import screenfull from 'screenfull';
export default {
  head() {
    return {
      title: '排行榜 - 2020程序员节日游戏',
    };
  },
  data() {
    return {
      loopId: '',
      isAll: false,
      All: {},
      respGames: {},
      Gid: 0,
      isFullscreen: false,
      tabs: [
        { label: '当前得分', name: 'Total' },
        { label: '第一次排名', name: 'First' },
        { label: '第二次排名', name: 'Second' },
        { label: '第三次排名', name: 'Third' },
      ],
    };
  },
  // async asyncData({ app, params, store }) {
  //   let res = await app.$axios.get('rank');
  //   const All = res,
  //     { Gid } = res;
  //   return { Gid, All };
  // },
  watch: {
    isAll(n) {
      if (n === true) this.checkAll(n);
      else this.All = this.respGames;
    },
  },
  mounted() {
    _this.loopId = setInterval(() => {
      _this.init();
    }, 10000);
  },
  methods: {
    async init() {
      let res = await _this.$axios.get('rank');
      _this.All = res;
      if (this.isAll) this.checkAll();

      _this.respGames = _this.$deepCopy(res);
      _this.Gid = res.Gid || 0;
    },
    onSort() {
      //排序，已废弃
      _this.record.sort(compare('coin'));
    },
    toURL() {
      open('https://github.com/huji0624/progame');
    },
    clickFullscreen() {
      if (this.isFullscreen) return;
      screenfull.toggle();
      this.isFullscreen = true;
    },
    checkAll() {
      let { All } = this;

      let changes = {};
      for (const key in All) {
        if (All[key] && All[key] instanceof Array) {
          const element = All[key].filter((it) => it.GameCount > 63);
          changes[key] = element;
        } else changes[key] = All[key];
      }
      _this.All = changes;
      console.log(_this.All);
    },
  },
  created() {
    _this = this;
    _this.init();
  },
  beforeDestroy() {
    this.$once('hook:beforeDestroy', () => {
      clearInterval(_this.loopId);
    });
  },
};

let compare = function (prop) {
  return function (obj1, obj2) {
    var val1 = obj1[prop];
    var val2 = obj2[prop];
    if (!isNaN(Number(val1)) && !isNaN(Number(val2))) {
      val1 = Number(val1);
      val2 = Number(val2);
    }
    if (val1 < val2) {
      return 1;
    } else if (val1 > val2) {
      return -1;
    } else {
      return 0;
    }
  };
};
</script>

<style lang="less" scoped>
@bodercoler: #1dffff;
.main {
  margin: 0 auto;
  text-align: center;
  padding: 10px;
  font-size: 14px;
  color: #fff;
  .head {
    position: fixed;
    left: 0;
    right: 0;
    top: 0;
    bottom: 0;
    margin: auto;
    margin-top: 0;
    width: 448px;
    height: 108px;
    background-image: url('../assets/images/list.png');
    background-size: 100% 100%;
    // padding: 10px;
    font-size: 20px;
    z-index: 10;
  }
  .isall {
    position: fixed;
    left: 60px;
    top: 70px;
    margin-top: 0;
    width: 450px;
    height: 108px;
    font-size: 20px;
    z-index: 10;
  }
  .logo {
    position: fixed;
    left: 0px;
    top: 0px;
    margin-top: 0;
    font-size: 30px;
    z-index: 10;
    color: #fff;
    width: 185px;
    height: 60px;
    background: url(../assets/images/logo.png);
    background-size: 100% 100%;
    cursor: pointer;
    // text-shadow: 0 0 4px #fff, 0 -5px 4px #ff3, 2px -10px 6px #fd3,
    //   -2px -15px 10px #f80, 2px -25px 20px #f20;

    // text-shadow: 0 -1px 0 #123; 凹进效果
    // text-shadow: 0 -1px 1px #eee; 凸出效果
    // text-shadow: 0 1px 1px #123; 凸出效果
  }
  .list {
    padding: 40px 20px 20px;
    margin: 40px 175px 5px;
    font-size: 16px;
    background: #2f0365;
    opacity: 0.9;
    border-radius: 10px;
    border: 7px @bodercoler solid;
    box-shadow: 0 0 10px #ee6a92;
    .name {
      width: 94%;
      color: #fff;
      font-size: 20px;
      font-weight: bold;
      line-height: 48px;
      background-color: #1d58db;
      padding: 5px;
      margin: 11px;
      border-radius: 10px;
      border: @bodercoler 2px solid;
    }

    .title {
      line-height: 58px;
      margin: 3px;
      font-weight: bold;
      border-bottom: @bodercoler 2px dashed;
    }
    .rows {
      background: #320a65;
      opacity: 1;
      margin: 3px;

      .listbody {
        border: @bodercoler 2px solid;
        margin: 10px;
        color: #fff;
        border-radius: 10px;
        background: #320a65;
      }
      .notstart {
        padding-top: 150px;
        color: @bodercoler;
        font-weight: bold;
      }
      .row {
        font-size: 16px;
        height: 390px;
        overflow: hidden auto;
        line-height: 50px;
        margin: 3px;

        .item {
          margin: 3px;
          transition: all 0.2s;
          // box-shadow: 0px 2px 10px 0px rgba(198, 198, 198, 0.5);
          .img {
            padding-top: 10px;
          }
        }
        .item:hover {
          background: #587c82;
        }
        .no {
          color: #fff;
          border-radius: 5px;
          background-color: #f56c6c;
          border-color: #f56c6c;
        }
      }
    }
  }
  .replay {
    height: 200px;
    overflow-y: auto;
    margin: 30px 175px 0;
    padding: 15px;
    background: #2f0365;
    opacity: 0.95;
    border-radius: 10px;
    border: 7px @bodercoler solid;
    box-shadow: 0 0 10px #ee6a92;
    .title {
      line-height: 40px;
      font-size: 18px;
      color: #fff;
      text-align: left;
      margin-left: 10px;
      span {
        font-size: 13px;
        color: #64dbf3;
      }
    }
    .btns {
      text-align: left;

      .btn {
        width: 120px;
        margin: 5px 15px;
        background: transparent;
        color: #64dbf3;
        border-color: #409eff;
      }
      .btn:hover {
        color: #00378a;
        background-color: #64dbf3;
        border-color: #64dbf3;
      }
    }
    .notstart {
      padding: 50px;
      color: @bodercoler;
      font-weight: bold;
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
