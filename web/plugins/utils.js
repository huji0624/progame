import Vue from 'vue';
import { util } from '@wnl/util';
import { wnlShare, wxShare } from '@wnl/ui';
import { TxMapKey } from '@/core/config';

/**
 * 导航到某个路由
 */
Vue.prototype.$toURL = function (url) {
  this.$router.push(url);
};

/**
 * 简单深拷贝
 */
Vue.prototype.$deepCopy = function (obj) {
  return JSON.parse(JSON.stringify(obj));
};

/**
 * 检测对象是否为空
 */
Vue.prototype.$checkNullObj = function (obj) {
  return Object.keys(obj).length === 0 && obj.constructor === Object;
};

/**
 * 返回
 */
Vue.prototype.$back = function (url) {
  let length = history.length;
  if (length == 2 || length == 1) return this.$router.push('/');
  else if (url) this.$router.push(url);
  else return this.$router.go(-1);
};

/**
 * 百度埋点自定义事件提交
 * @param String activeName 活动名称 JSJ2020
 * @param String eventName 事件名  start
 * @param String eventType 事件类型 click
 */
Vue.prototype.$hm = function (activeName, eventName, eventType) {
  let clt = '.wx',
    OS = 'az';

  util.isWnl && (clt = '.wnl');
  util.isIOS && (OS = 'ios');

  let evn = ['_trackEvent', activeName + '.' + eventName + clt, eventType, OS];
  window._hmt.push(evn);
};

Vue.prototype.$share = function (share) {
  // share = {
  //   title: '主标题',
  //   text: '副标题',
  //   image: '分享图地址，与imgUrl相同',
  //   imgUrl: '分享图地址，与image相同',
  //   url: '分享页面地址',
  //   callback: () => {},
  // };
  if (util.isWnl) wnlShare.setShareData(share);
  if (util.isWeixin) new wxShare(share);
};

Vue.prototype.$getLocationTxmap = function () {
  let geolocation = new qq.maps.Geolocation(TxMapKey, 'myapp');
  let re = null;
  return new Promise(function (resolve, reject) {
    geolocation.getLocation(
      (position) => {
        re = {
          code: 'success',
          body: position,
        };
        resolve(re);
      },
      (err) => {
        re = {
          code: 'fail',
          body: err,
        };
        resolve(re);
      }
    );
  });
};
