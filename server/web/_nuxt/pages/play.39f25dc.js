(window.webpackJsonp=window.webpackJsonp||[]).push([[3],{335:function(o,e,t){var content=t(343);"string"==typeof content&&(content=[[o.i,content,""]]),content.locals&&(o.exports=content.locals);(0,t(25).default)("fe6b6ee8",content,!0,{sourceMap:!1})},336:function(o,e,t){o.exports=t.p+"img/btn.df4f6ae.png"},337:function(o,e,t){var content=t(345);"string"==typeof content&&(content=[[o.i,content,""]]),content.locals&&(o.exports=content.locals);(0,t(25).default)("811a1ac8",content,!0,{sourceMap:!1})},342:function(o,e,t){"use strict";var n=t(335);t.n(n).a},343:function(o,e,t){var n=t(24),r=t(134),c=t(336);e=n(!1);var l=r(c);e.push([o.i,'.container[data-v-2b9af680]{margin:0 auto;min-height:100vh;text-align:center;padding:10px;color:#fff}.container .note[data-v-2b9af680]{font-weight:700;color:#1dfefe}.container .btns[data-v-2b9af680]{margin:50px}.container .btns .btn[data-v-2b9af680]{color:#04def0;font-weight:700;font-family:"微软雅黑";padding:15px 32px;margin:0 10px;-webkit-user-select:none;-moz-user-select:none;-ms-user-select:none;user-select:none;background-image:url('+l+");background-size:100% 100%;cursor:pointer}.container .btns .btn[data-v-2b9af680]:hover{color:#fff}.container .main[data-v-2b9af680]{width:702px;background:#330a66;box-shadow:0 0 #333;border-radius:10px;border:7px solid #1dfefe;margin:10px auto;box-shadow:0 2px 12px 0 rgba(0,0,0,.1)}.container .main .items[data-v-2b9af680]{width:100px;height:100px;position:relative;border:1px solid #1dfefe;float:left;-webkit-user-select:none;-moz-user-select:none;-ms-user-select:none;user-select:none;background:#330a66;transition:all .2s;font-size:14px;display:flex;flex-wrap:wrap;justify-content:center;align-items:center}.container .main .items .item[data-v-2b9af680]{width:95px;height:25px;line-height:23px;margin:3px;border-radius:15px;z-index:99;color:#1eeaf0;background:#635393;border:1px solid #1dfefe}.container .main .items .focus[data-v-2b9af680]{color:#0e025e;font-weight:700;background-color:#edcc53;border-color:#929a19}.container .main .items .gold[data-v-2b9af680]{width:100px;height:100px;position:absolute;top:0;left:0;color:#969494;padding-top:10px;font-size:60px;font-style:italic;z-index:1;opacity:.3}.container .main .items[data-v-2b9af680]:hover{background:#eee}.container .wrapbox[data-v-2b9af680]{position:absolute;top:170px;right:20px;width:300px;font-size:14px;padding:12px;color:#1ceaee;word-break:break-all;border-radius:10px;border:3px solid #1dfefe;background:#1d0957;box-shadow:0 2px 12px 0 rgba(0,0,0,.1)}.container .wrapbox .tips[data-v-2b9af680]{font-size:20px;font-weight:700;line-height:30px}.container .wrapbox .tips span[data-v-2b9af680]{font-size:13px;font-weight:300}.container .wrapbox .list[data-v-2b9af680]{min-width:150px;z-index:2000;width:300px;max-height:200px;line-height:1.4;text-align:justify;overflow-y:auto}.container .focusOver[data-v-2b9af680]{top:170px;left:20px}.container .focusOver .crtPos[data-v-2b9af680]{text-align:left;color:#fcf8a7;line-height:30px}.container .tname[data-v-2b9af680]{display:inline-block;font-size:16px;font-weight:700;width:90px;color:#add5fd}.container .gold[data-v-2b9af680]{font-weight:700;color:#ffeb00}.container .nodata[data-v-2b9af680]{line-height:50px}.container .noGame[data-v-2b9af680]{color:#fff;width:300px;margin:350px auto;padding:50px;border-radius:5px;color:#64dbf3;border:2px solid #1dfefe}",""]),o.exports=e},344:function(o,e,t){"use strict";var n=t(337);t.n(n).a},345:function(o,e,t){(e=t(24)(!1)).push([o.i,".el-checkbox__input.is-checked .el-checkbox__inner{background-color:#04def0;border-color:#04def0}.el-checkbox,.el-checkbox__input.is-checked+.el-checkbox__label{color:#04def0}.el-checkbox{line-height:30px}.el-checkbox__input.is-indeterminate .el-checkbox__inner{background-color:#04def0;border-color:#04def0}.el-page-header__content{font-size:18px;color:#04def0}.el-checkbox__label{width:90px}",""]),o.exports=e},360:function(o,e,t){"use strict";t.r(e);var n,r,c=t(63);var l={head:function(){return{title:"游戏回放 - 2020程序员节日游戏"}},data:function(){return{isClick:0,socket:r,x:0,y:0,mainW:0,mainH:0,loopId:0,playersInfo:[],crtPos:"",crtGold:"",total:[],isFirst:!0,info:{RoundID:0}}},methods:{start:function(){for(var o=this.x,e=this.y,t=this.info.Tilemap,n=[],i=0;i<e;i++)for(var r=0;r<o;r++){var c=t[i][r],l={players:(c.Players||[]).map((function(o){return"token"==o.Name&&(o.isFocus=!0),o}))||[],pos:[r,i],gold:c.Gold};n.push(l)}this.total=n},mouseOver:function(o){this.playersInfo=o.players||[],this.crtPos=o.pos,this.crtGold=o.gold},init:function(){"undefined"!=typeof WebSocket&&((r=new WebSocket("ws://localhost:8881/ws")).onopen=this.onOpen,r.onerror=this.onError,r.onmessage=this.onMessage)},onOpen:function(){console.log("socket连接成功")},onMessage:function(o){var e=o.data;console.log("recv msg:");var t=JSON.parse(e);0==t.Msgtype?console.log("login ok."):1==t.Msgtype?ws.send(JSON.stringify({msgtype:2,token:"token"})):3==t.Msgtype&&(n.info=n.$deepcopy(t),n.isFirst&&(n.x=t.Wid,n.y=t.Hei,n.x=t.Wid,n.mainW=100*n.x+14,n.mainH=100*n.y+14,n.isFirst=!1),n.loopId=setTimeout((function(){if(n.isClick)clearTimeout(n.loopId);else{var o=Math.floor(Math.random()*t.Wid),e=Math.floor(Math.random()*t.Hei);r.send(JSON.stringify({msgtype:4,token:"token",x:o,y:e,RoundID:n.info.RoundID}))}}),700),onSend())},onSend:function(o){},onError:function(){console.log("socket发生了错误")},onClose:function(){console.log("socket已经关闭")},onClick:function(o){var e=Object(c.a)(o.pos,2),t=e[0],l=e[1];n.isClick=!0,r.send(JSON.stringify({msgtype:4,token:"token",x:t,y:l,RoundID:n.info.RoundID}))}},created:function(){n=this},destroyed:function(){r.close()}},d=(t(342),t(344),t(41)),component=Object(d.a)(l,(function(){var o=this,e=o.$createElement,t=o._self._c||e;return t("div",{staticClass:"container"},[t("span",{staticClass:"note"},[o._v(" 第 "+o._s(o.info.RoundID+1)+" 轮")]),o._v(" "),t("div",{staticClass:"btns"},[t("span",{staticClass:"btn",on:{click:function(e){return o.init()}}},[o._v(" 开始游戏 ")]),o._v(" "),t("span",{staticClass:"btn",on:{click:function(e){return o.socket.close()}}},[o._v("结束游戏")])]),o._v(" "),t("div",{staticClass:"main",style:{width:o.mainW+"px",height:o.mainH+"px"}},o._l(o.total,(function(e,i){return t("div",{key:i,staticClass:"items",on:{mouseover:function(t){return o.mouseOver(e)},click:function(t){return o.onClick(e)}}},[t("div",{staticClass:"gold"},[o._v(o._s(e.gold))]),o._v(" "),e.players?t("div",o._l(e.players,(function(e,i){return t("div",{key:i},[i<3?t("div",{staticClass:"item",class:{focus:e.isFocus}},[(e.Name+e.Gold).length>8?t("div",[t("marquee",{attrs:{scrollamount:"2"}},[o._v("\n                "+o._s(e.Name)+" - "+o._s(e.Gold)+"\n              ")])],1):t("div",[o._v(o._s(e.Name)+" - "+o._s(e.Gold))])]):o._e()])})),0):o._e()])})),0),o._v(" "),t("div",{staticClass:"wrapbox focusOver"},[t("div",{staticClass:"conta"},[t("div",{staticClass:"tips"},[o._v("棋盘格信息")]),o._v(" "),o.crtPos?t("div",{staticClass:"crtPos"},[o._v("\n        当前坐标："+o._s(o.crtPos)+" 金币数量："+o._s(o.crtGold)+"\n      ")]):t("div",{staticClass:"nodata"},[o._v("鼠标移过棋盘，展示数据")]),o._v(" "),t("div",{staticClass:"list"},o._l(o.playersInfo,(function(e,i){return t("div",{key:i},[o._v("\n          "+o._s(i+1)+"、团队："),t("span",{staticClass:"tname"},[o._v(o._s(e.Name))]),o._v("\n          金币：\n          "),t("span",{staticClass:"gold"},[o._v(o._s(e.Gold))])])})),0)])])])}),[],!1,null,"2b9af680",null);e.default=component.exports}}]);