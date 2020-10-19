## 玩法介绍

0x0 这是一个棋盘内吃金币的游戏。

0x0 每一局游戏，进行 96 轮结束。每局游戏开始时，你将被随机分配在棋盘的任意位置，同一个棋盘格子可能会有多个玩家队伍。

0x0 每一轮游戏开始时，系统将发送 8\*6 的游戏棋盘(第一个格子坐标为[0,0]，最后一个格子为 [7,5])信息给玩家，发送形式是二维数组。包含每个格子的金币数量和玩家位置，玩家名称，玩家金币数量。

0x0 每轮游戏收到地图信息后，你需要在 1s 之内向服务器提交你下一步的位置[x,y]。

0x0 你需要通过算法来计算要移动的位置，以抢到金币，以保证最终获得金币最多来赢得比赛胜利。

0x0 每场比赛，玩家的最终得分为：本场比赛截止时间前，最近 m 局游戏结束时的平均金币个数，97>m>64，即你需要参与至少 64 局以上，超过 96 局仅计算 96 局平均金币数量。

##### 特别注意，此 96 局含义为系统局数，不是你参与的游戏局数。

## 游戏时间

游戏服务器会在 10 月 23 日中午 1 点，截止报名以后开启。持续 3 天。

第一场：10 月 23 日 23:00 截止 得分权重占总排名的 10%  
第二场：10 月 24 日 23:00 截止 得分权重占总排名的 30%  
第三场：10 月 25 日 23:00 截止 得分权重占总排名的 60%

## 游戏奖励

第一名：4096 元现金大奖  
第二名：2048 元现金大奖  
第三名：1024 元现金大奖

有钱任性奖（给出金币最多）：512 元现金大奖  
长跑冠军奖（移动步数最多）：512 元现金大奖  
最年轻团队奖（排名前 50%，团队平均年龄最小）：512 元现金大奖

## 游戏规则

0x1、每轮游戏开始时，系统向每个玩家发放 1 个金币，且立刻向玩家发放拥有金币的 20% 利息，利息最高为 10 个金币。

0x1、每轮游戏开始时，系统会在棋盘内随机格子生成 N 个金币（N 可能为负数）。

0x2、每轮游戏收到地图信息后，你需要在 1s 之内向服务器提交你下一步的位置[x,y]。如果你不小心移动到棋盘外，位置则保持不变，扣除 1 个金币。

0x3、每移动 n 步，扣除玩家 n\*1.5(取整) 个金币，原地不动，扣除 1 个金币，当移动所需金币不足时，认为玩家原地不动。每移动一次位置，即为完成一轮游戏。

0x4、每轮移动结束后，根据金币个数不同，获得的金币数量有所不同:

```
a.当金币数量为 -4 时，随机选取该格子中1一个玩家 (马上获得 40% 的利息，利息最高为 10 个金币) 或者 (失去 4个金币) 概率为50%。
b.当金币数量大于0且能整除5，该格子中的玩家金币数量全部为 (该格所有玩家金币总数+该金币数量)/玩家数量 并取整。
c.当金币数量为 7或者11 时，如果该格子只有1个玩家，该玩家获得对应数量金币，如果该格子有多个玩家，随机选取一个玩家失去对应数量金币。
d.当金币数量为 8 时，该格子玩家平分（取整）这8个金币。
e.其他情况，随机选取格子中1一个玩家获得该格金币(正数获得，负数扣除)。
f、当有多个玩家在同一个格子中时，金币最高者的玩家将会给出自己 1/4 的金币给每个金币少者的玩家，当给出的金币数量不够每个人 1 个金币时，不再给出。
g、当同一个格子里有多个最高玩家，则所有最高金币玩家都将给出1/4之一金币给非最高金币玩家。
```

0x5、每局游戏进行 96 轮后结束。

## 游戏准备

服务器会提供一个 websocket 供你连接，获得你的游戏队伍名称和 token 以保证游戏队伍的唯一性。  
我们会开源所有代码，供你搭建环境，调试程序。当然我们也可能会对游戏本身进行调优，不要忘记及时拉取代码哦~

## 你要做的事

编写一个程序，语言不限，通过算法来计算要移动的位置以吃到各种数量的金币，以保证获得金币最多来赢得比赛胜利。

#### 开发和调试

0x1、搭建环境。拉取代码，启动后端服务，后端环境为 go.(go 语言环境搭建教程请自行百度)。

0x2、编写代码，连接到刚启动的 websocket 连接，一般为`"ws://http://localhost:8881/ws`。

0x3、连接上服务器后，你会收到服务器返回的类似如下二维数据

```
{
  "GameID": 5,//当前游戏局数
  "Msgtype": 3,//你收到的消息类别
  //Msgtype //0 登陆 -1 登陆失败 1 请准备 2 准备好了 3 游戏信息 4 玩家行动 5 本局游戏结束

  "RoundID": 3,//当前游戏轮数
  "Wid": 3,//棋盘宽
  "Hei": 4,//棋盘高
  "Tilemap": [
    [
      { "Gold": 1, "Players": [{ "Name": "YOU", "Gold": 2 }] },//这是你的位置，包含你的名字、金币数、所在位置[0,0]
      { "Gold": 3 },
      { "Gold": 3, "Players": [{ "Name": "player1", "Gold": 0 }] }
    ],
    [
      { "Gold": 4 },
      {
        "Gold": 0,
        "Players": [
          { "Name": "player2", "Gold": 3 },
          { "Name": "player3", "Gold": 0 }
        ]
      },
      { "Gold": 4 }
    ],
    [{ "Gold": 0 }, { "Gold": 3 }, { "Gold": 2 }],
    [{ "Gold": 2 }, { "Gold": 2 }, { "Gold": 0 }]
  ]
}

```

0x3、让代码跑起来。调试，调优。  
你向服务器发送的消息类别字段为`msgtype`（注意区分大小写）  
msgtype：0 //0 登陆 -1 登陆失败 1 请准备 2 准备好了 3 游戏信息 4 玩家行动 5 本局游戏结束

每轮游戏最终，你需要向服务器提交类似如下具有 x，y 坐标的移动位置。`msgtype:4`

```
ws.send(JSON.stringify({msgtype:4,token:token,x:x,y:y,RoundID:RoundID}));
```

0x4、后端服务启动后，你可以在浏览器打开`http://localhost:8881`查看本地对局排名情况，查看本地进 30 局回放。

##### 我们在 `testclient` 目录下存放了使用 JavaScript 和 python 语言编写的测试案例，用于帮助你理解游戏规则，也帮助你清晰怎么做。

#### 正式环境

你只需要将 websocket 的链接改为我们线上服务器的链接，就可以啦~

不要忘记在游戏开始时，加入游戏。

## 实时赛况

0x1 你可以[点击这里](https://testmobile.51wnl-cq.com/20201024/)观看比赛即时排名情况。  
0x2 你可以点击【回放】按钮查看他们的激烈对决。学习和探讨他们是如何取得胜利的，源代码会在赛后开源。  
0x3 回放过程中，你可以添加你感兴趣的游戏队伍，系统将为你特别关注。
