// import Vue from 'vue';
// import vant from 'vant';
import 'vant/lib/index.css';
// import VConsole from 'vconsole/dist/vconsole.min.js';
import MtaH5 from 'mta-h5-analysis';
// import { isDebug } from '../core/config';

let isfirst = true;

export default ({ app, store }) => {
  // isDebug && new VConsole();
  // Vue.use(vant);

  app.router.afterEach((to, from) => {
    if (isfirst) {
      mtaInit();
      isfirst = false;
    }
    MtaH5.pgv();
  });
};

/**
 * 数据统计初始化
 * @param {object} store store
 */
function mtaInit() {
  MtaH5.init({
    sid: '500732196', //必填，统计用的appid
    cid: '500732209', //如果开启自定义事件，此项目为必填，否则不填
    autoReport: 1, //是否开启自动上报(1:init完成则上报一次,0:使用pgv方法才上报)
    senseHash: 0, //hash锚点是否进入url统计
    senseQuery: 0, //url参数是否进入url统计
    performanceMonitor: 0, //是否开启性能监控
    ignoreParams: [], //开启url参数上报时，可忽略部分参数拼接上报
  });
}
