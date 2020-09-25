<template>
  <section class="main">
    <div class="head">排行榜</div>
    <div class="list">
      <el-row class="rows">
        <el-col :span="6" v-for="(item, i) in tabs" :key="i">
          <span class="name"> {{ item.label }}</span>
          <el-row>
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
      <div class="title">点击按钮查看当局回放 <span>仅展示最近40局</span></div>
      <div class="btns">
        <span v-for="(item, i) in 40" :key="i">
          <el-button
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
      tabs: [
        { label: '第一次排名', name: 'First' },
        { label: '第二次排名', name: 'Second' },
        { label: '第三次排名', name: 'Third' },
        { label: '总排名', name: 'Total' },
      ],
    };
  },
  async asyncData({ app, params, store }) {
    let res = await app.$axios.get('rank');
    const All = res,
      { Gid } = res;
    return { Gid, All };
  },
  mounted() {},
  methods: {
    onSort() {
      //排序，已废弃
      this.record.sort(compare('coin'));
    },
  },
  created() {
    _this = this;
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
.main {
  margin: 0 auto;
  text-align: center;
  padding: 10px;
  font-size: 14px;
  .head {
    padding: 10px;
    font-size: 20px;
  }
  .list {
    padding: 20px 100px;
    font-size: 16px;
    .name {
      font-size: 20px;
      line-height: 48px;
    }

    .title {
      background: #99a9bf;
      color: #fff;
      line-height: 58px;
      margin: 3px;
    }
    .rows {
      color: #000;

      background: #d3dce6;
      margin: 3px;

      .notstart {
        padding-top: 50px;
      }
      .row {
        font-size: 18px;
        height: 500px;
        overflow: hidden auto;
        line-height: 60px;
        .item {
          margin: 3px;
          background: #e5e9f2;
        }
        .no {
          background: red;
        }
      }
    }
  }
  .replay {
    margin: 10px 100px;
    padding: 10px;
    background: #d3dce6;
    .title {
      line-height: 50px;
      font-size: 18px;
      color: #656464;
      text-align: left;
      margin-left: 10px;
      span {
        font-size: 13px;
      }
    }
    .btns {
      text-align: left;

      .btn {
        width: 100px;
        margin: 5px 10px;
      }
    }
  }
}
::-webkit-scrollbar {
  width: 0 !important;
}
</style>
