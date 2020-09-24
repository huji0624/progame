<template>
  <section class="main">
    <div class="head">排行榜</div>
    <el-button @click="onClick">排序</el-button
    ><el-button @click="$router.push('/replay?gid=2')">回放</el-button>
    <div class="list">
      <el-tabs class="tabs" v-model="activeName">
        <el-tab-pane v-for="(it, i) in tabs" :key="i" :label="it.label">
          <div>
            <div class="title">
              <el-row class="row">
                <el-col :span="8"> 名次 </el-col>
                <el-col :span="8"> 队名 </el-col>
                <el-col :span="8"> 得分 </el-col>
              </el-row>
            </div>
            <div v-if="!All[it.name]">比赛还未开始</div>
            <div v-else class="rows" v-for="(it, i) in All[it.name]" :key="i">
              <el-row class="row">
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
        </el-tab-pane>
      </el-tabs>
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
      activeName: '',

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
    const All = res;
    const { Gid } = res;

    return { Gid, All };
  },
  mounted() {},
  methods: {
    async init() {
      // let res = await this.$axios.get('rank');
    },
    onClick() {
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
  .head {
    padding: 20px;
  }
  .list {
    padding: 30px 100px;
    .tabs {
      text-align: center;

      .title {
        background: #99a9bf;
        color: #fff;
        line-height: 78px;
        margin: 3px;
      }
      .rows {
        color: #000;
        background: #d3dce6;
        margin: 3px;
        .row {
          line-height: 70px;
          .no {
            background: red;
          }
        }
      }
      .rows:hover {
        background: #e5e9f2;
      }
    }
  }
}
</style>
