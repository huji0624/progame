<template>
  <section class="main">
    <div class="head">排行榜</div>
    <div class="list">
      <el-row class="rows">
        <el-col :span="6" v-for="(item, i) in tabs" :key="i">
          <div class="name">{{ item.label }}</div>
          <el-row class="listbody">
            <div class="title">
              <el-row>
                <el-col :span="8"> 名次 </el-col>
                <el-col :span="8"> 队名 </el-col>
                <el-col :span="8"> 得分 </el-col>
              </el-row>
            </div>
            <div class="row">
              <div v-if="!All[item.name]" class="notstart">敬请期待</div>
              <div v-else>
                <el-row class="item" v-for="(it, i) in All[item.name]" :key="i">
                  <el-col :span="8">
                    <div :class="i < 3 ? 'no' : ''">{{ i + 1 }}</div>
                  </el-col>
                  <el-col :span="8">
                    {{ it.Name }}
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
      <div class="title">点击按钮查看当局回放 <span>仅展示最近42局</span></div>
      <div class="btns">
        <span v-for="(item, i) in 42" :key="i">
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
export default {
  head() {
    return {
      title: '金币作战排行榜',
    };
  },
  data() {
    return {
      loopId: '',
      All: {},
      Gid: 0,
      tabs: [
        { label: '第一次排名', name: 'First' },
        { label: '第二次排名', name: 'Second' },
        { label: '第三次排名', name: 'Third' },
        { label: '最终排名', name: 'Total' },
      ],
    };
  },
  // async asyncData({ app, params, store }) {
  //   let res = await app.$axios.get('rank');
  //   const All = res,
  //     { Gid } = res;
  //   return { Gid, All };
  // },
  mounted() {
    _this.loopId = setInterval(() => {
      _this.init();
    }, 10000);
  },
  methods: {
    async init() {
      let res = await _this.$axios.get('rank');
      _this.All = res;
      _this.Gid = res.Gid;
    },
    onSort() {
      //排序，已废弃
      _this.record.sort(compare('coin'));
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
@bodercoler: #1dfefe;
.main {
  margin: 0 auto;
  text-align: center;
  padding: 10px;
  font-size: 14px;
  color: #fff;
  .head {
    padding: 10px;
    font-size: 20px;
  }
  .list {
    padding: 20px 160px;
    font-size: 16px;
    .name {
      width: 94%;
      color: #fff;
      font-size: 20px;
      font-weight: bold;
      line-height: 48px;
      background-color: #1d58db;
      padding: 8px;
      margin: 13px;
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
      opacity: 0.8;
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
        height: 400px;
        overflow: hidden auto;
        line-height: 60px;
        margin: 3px;

        .item {
          margin: 3px;
          // box-shadow: 0px 2px 10px 0px rgba(198, 198, 198, 0.5);
        }
        .item:hover {
          background: #320a65;
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
    margin: 10px 60px;
    padding: 10px;
    background: #2f0365;
    opacity: 0.8;
    border-radius: 10px;
    border: 7px #6ae5ee solid;
    box-shadow: 0 0 10px #ee6a92;
    .title {
      line-height: 50px;
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
        width: 100px;
        margin: 5px 10px;
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
  }
}
</style>
