/*! For license information please see ../LICENSES */
(window.webpackJsonp=window.webpackJsonp||[]).push([[2],{338:function(e,t,n){var content=n(349);"string"==typeof content&&(content=[[e.i,content,""]]),content.locals&&(e.exports=content.locals);(0,n(25).default)("40cef03f",content,!0,{sourceMap:!1})},344:function(e,t,n){e.exports=n.p+"img/1.d95fc63.png"},345:function(e,t,n){e.exports=n.p+"img/2.36dde7e.png"},346:function(e,t,n){e.exports=n.p+"img/3.85838e5.png"},347:function(e,t,n){!function(){"use strict";var t="undefined"!=typeof window&&void 0!==window.document?window.document:{},n=e.exports,r=function(){for(var e,n=[["requestFullscreen","exitFullscreen","fullscreenElement","fullscreenEnabled","fullscreenchange","fullscreenerror"],["webkitRequestFullscreen","webkitExitFullscreen","webkitFullscreenElement","webkitFullscreenEnabled","webkitfullscreenchange","webkitfullscreenerror"],["webkitRequestFullScreen","webkitCancelFullScreen","webkitCurrentFullScreenElement","webkitCancelFullScreen","webkitfullscreenchange","webkitfullscreenerror"],["mozRequestFullScreen","mozCancelFullScreen","mozFullScreenElement","mozFullScreenEnabled","mozfullscreenchange","mozfullscreenerror"],["msRequestFullscreen","msExitFullscreen","msFullscreenElement","msFullscreenEnabled","MSFullscreenChange","MSFullscreenError"]],i=0,r=n.length,o={};i<r;i++)if((e=n[i])&&e[1]in t){for(i=0;i<e.length;i++)o[n[0][i]]=e[i];return o}return!1}(),o={change:r.fullscreenchange,error:r.fullscreenerror},l={request:function(element){return new Promise(function(e,n){var o=function(){this.off("change",o),e()}.bind(this);this.on("change",o);var l=(element=element||t.documentElement)[r.requestFullscreen]();l instanceof Promise&&l.then(o).catch(n)}.bind(this))},exit:function(){return new Promise(function(e,n){if(this.isFullscreen){var o=function(){this.off("change",o),e()}.bind(this);this.on("change",o);var l=t[r.exitFullscreen]();l instanceof Promise&&l.then(o).catch(n)}else e()}.bind(this))},toggle:function(element){return this.isFullscreen?this.exit():this.request(element)},onchange:function(e){this.on("change",e)},onerror:function(e){this.on("error",e)},on:function(e,n){var r=o[e];r&&t.addEventListener(r,n,!1)},off:function(e,n){var r=o[e];r&&t.removeEventListener(r,n,!1)},raw:r};r?(Object.defineProperties(l,{isFullscreen:{get:function(){return Boolean(t[r.fullscreenElement])}},element:{enumerable:!0,get:function(){return t[r.fullscreenElement]}},isEnabled:{enumerable:!0,get:function(){return Boolean(t[r.fullscreenEnabled])}}}),n?e.exports=l:window.screenfull=l):n?e.exports={isEnabled:!1}:window.screenfull={isEnabled:!1}}()},348:function(e,t,n){"use strict";var r=n(338);n.n(r).a},349:function(e,t,n){var r=n(24),o=n(94),l=n(350),c=n(351);t=r(!1);var d=o(l),f=o(c);t.push([e.i,".main[data-v-11ec5d19]{margin:0 auto;text-align:center;padding:10px;font-size:14px;color:#fff}.main .head[data-v-11ec5d19]{right:0;bottom:0;margin:0 auto auto;width:448px;height:108px;background-image:url("+d+');background-size:100% 100%;font-size:20px}.main .head[data-v-11ec5d19],.main .logo[data-v-11ec5d19]{position:fixed;left:0;top:0;z-index:10}.main .logo[data-v-11ec5d19]{font-family:"FZJZJT";margin-top:0;font-size:30px;color:#fff;width:185px;height:60px;background:url('+f+");background-size:100% 100%;cursor:pointer}.main .list[data-v-11ec5d19]{padding:40px 20px 20px;margin:40px 175px 5px;font-size:16px;background:#2f0365;opacity:.9;border-radius:10px;border:7px solid #1dffff;box-shadow:0 0 10px #ee6a92}.main .list .name[data-v-11ec5d19]{width:94%;color:#fff;font-size:20px;font-weight:700;line-height:48px;background-color:#1d58db;padding:5px;margin:11px;border-radius:10px;border:2px solid #1dffff}.main .list .title[data-v-11ec5d19]{line-height:58px;margin:3px;font-weight:700;border-bottom:2px dashed #1dffff}.main .list .rows[data-v-11ec5d19]{background:#320a65;opacity:1;margin:3px}.main .list .rows .listbody[data-v-11ec5d19]{border:2px solid #1dffff;margin:10px;color:#fff;border-radius:10px;background:#320a65}.main .list .rows .notstart[data-v-11ec5d19]{padding-top:150px;color:#1dffff;font-weight:700}.main .list .rows .row[data-v-11ec5d19]{font-size:16px;height:390px;overflow-x:hidden;overflow-y:auto;overflow:hidden auto;line-height:50px;margin:3px}.main .list .rows .row .item[data-v-11ec5d19]{margin:3px;transition:all .2s}.main .list .rows .row .item .img[data-v-11ec5d19]{padding-top:10px}.main .list .rows .row .item[data-v-11ec5d19]:hover{background:#587c82}.main .list .rows .row .no[data-v-11ec5d19]{color:#fff;border-radius:5px;background-color:#f56c6c;border-color:#f56c6c}.main .replay[data-v-11ec5d19]{height:200px;overflow-y:auto;margin:30px 175px 0;padding:15px;background:#2f0365;opacity:.95;border-radius:10px;border:7px solid #1dffff;box-shadow:0 0 10px #ee6a92}.main .replay .title[data-v-11ec5d19]{line-height:40px;font-size:18px;color:#fff;text-align:left;margin-left:10px}.main .replay .title span[data-v-11ec5d19]{font-size:13px;color:#64dbf3}.main .replay .btns[data-v-11ec5d19]{text-align:left}.main .replay .btns .btn[data-v-11ec5d19]{width:120px;margin:5px 15px;background:transparent;color:#64dbf3;border-color:#409eff}.main .replay .btns .btn[data-v-11ec5d19]:hover{color:#00378a;background-color:#64dbf3;border-color:#64dbf3}.main .replay .notstart[data-v-11ec5d19]{padding:50px;color:#1dffff;font-weight:700}",""]),e.exports=t},350:function(e,t,n){e.exports=n.p+"img/list.132173c.png"},351:function(e,t,n){e.exports=n.p+"img/logo.2515a87.png"},353:function(e,t,n){"use strict";n.r(t);n(335),n(49);var r,o=n(4),l=n(347),c=n.n(l),d={head:function(){return{title:"排行榜 - 2020程序员节日游戏"}},data:function(){return{loopId:"",All:{},Gid:0,isFullscreen:!1,tabs:[{label:"当前得分",name:"Total"},{label:"第一次排名",name:"First"},{label:"第二次排名",name:"Second"},{label:"第三次排名",name:"Third"}]}},mounted:function(){r.loopId=setInterval((function(){r.init()}),1e4)},methods:{init:function(){return Object(o.a)(regeneratorRuntime.mark((function e(){var t;return regeneratorRuntime.wrap((function(e){for(;;)switch(e.prev=e.next){case 0:return e.next=2,r.$axios.get("rank");case 2:t=e.sent,r.All=t,r.Gid=t.Gid||0;case 5:case"end":return e.stop()}}),e)})))()},onSort:function(){r.record.sort(f("coin"))},toURL:function(){open("https://github.com/huji0624/progame")},clickFullscreen:function(){this.isFullscreen||(c.a.toggle(),this.isFullscreen=!0)}},created:function(){(r=this).init()},beforeDestroy:function(){this.$once("hook:beforeDestroy",(function(){clearInterval(r.loopId)}))}},f=function(e){return function(t,n){var r=t[e],o=n[e];return isNaN(Number(r))||isNaN(Number(o))||(r=Number(r),o=Number(o)),r<o?1:r>o?-1:0}},m=d,v=(n(348),n(45)),component=Object(v.a)(m,(function(){var e=this,t=e.$createElement,r=e._self._c||t;return r("section",{staticClass:"main",on:{mousemove:e.clickFullscreen}},[r("div",{staticClass:"logo",attrs:{title:"查看抢钱大作战的玩法和帮助"},on:{click:e.toURL}}),e._v(" "),r("div",{staticClass:"head"}),e._v(" "),r("div",{staticClass:"list"},[r("el-row",{staticClass:"rows"},e._l(e.tabs,(function(t,i){return r("el-col",{key:i,attrs:{span:6}},[r("div",{staticClass:"name"},[e._v(e._s(t.label))]),e._v(" "),r("el-row",{staticClass:"listbody"},[r("div",{staticClass:"title"},[r("el-row",[r("el-col",{attrs:{span:6}},[e._v(" 名次 ")]),e._v(" "),r("el-col",{attrs:{span:10}},[e._v(" 队名 ")]),e._v(" "),r("el-col",{attrs:{span:8}},[e._v(" 得分 ")])],1)],1),e._v(" "),r("div",{staticClass:"row"},[e.All[t.name]?r("div",e._l(e.All[t.name],(function(t,i){return r("el-row",{key:i,staticClass:"item"},[r("el-col",{attrs:{span:6,title:"Nihaoyo "}},[0==i?r("img",{staticClass:"img",attrs:{src:n(344)}}):e._e(),e._v(" "),1==i?r("img",{staticClass:"img",attrs:{src:n(345)}}):e._e(),e._v(" "),2==i?r("img",{staticClass:"img",attrs:{src:n(346)}}):e._e(),e._v(" "),i>2?r("div",[e._v(e._s(i+1))]):e._e()]),e._v(" "),r("el-col",{attrs:{span:10}},[e._v("\n                  "+e._s(t.Name.slice(0,12))+"\n                ")]),e._v(" "),r("el-col",{attrs:{span:8}},[e._v("\n                  "+e._s(t.Gold)+"\n                ")])],1)})),1):r("div",{staticClass:"notstart"},[e._v("敬请期待")])])])],1)})),1)],1),e._v(" "),r("div",{staticClass:"replay"},[e._m(0),e._v(" "),e.Gid<1?r("div",{staticClass:"notstart"},[e._v("敬请期待")]):r("div",{staticClass:"btns"},e._l(30,(function(t,i){return r("span",{key:i},[e.Gid-i-1>0?r("el-button",{staticClass:"btn",attrs:{type:"primary",plain:""},on:{click:function(t){e.$router.push("/replay?gid="+(e.Gid-i-1))}}},[e._v("\n          第"+e._s(e.Gid-i-1)+"局\n        ")]):e._e()],1)})),0)])])}),[function(){var e=this.$createElement,t=this._self._c||e;return t("div",{staticClass:"title"},[this._v("点击按钮查看当局回放 "),t("span",[this._v("仅展示最近30局")])])}],!1,null,"11ec5d19",null);t.default=component.exports}}]);