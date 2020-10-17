## 玩法介绍

0x0 这是一个棋盘内吃金币的游戏。

0x0 每一局游戏，进行 120 轮结束。每局游戏开始时，你将被随机分配在棋盘的任意位置，同一个棋盘格子可能会有多个玩家队伍。

0x0 每一轮游戏开始时，系统将发送 8\*6 的游戏棋盘(第一个格子坐标为[0,0]，最后一个格子为 [7,5])信息给玩家，发送形式是二维数组。包含每个格子的金币数量和玩家位置，玩家名称，玩家金币数量。

0x0 每轮游戏收到地图信息后，你需要在 1s 之内向服务器提交你下一步的位置[x,y]。

0x0 你需要通过算法来计算要移动的位置，以抢到金币，以保证最终获得金币最多来赢得比赛胜利。

0x8 每场比赛，玩家的最终得分为：本场比赛截止时间前，最近 64 局游戏结束时的平均金币个数，最近 64 局的游戏中，你需要参与至少 32 局以上。

## 游戏时间

游戏服务器会在 10 月 23 日中午 1 点，截止报名以后开启。持续 3 天。

第一场：10 月 23 日 21:00 截止 得分权重占总排名的 10%  
第二场：10 月 24 日 21:00 截止 得分权重占总排名的 30%  
第三场：10 月 25 日 21:00 截止 得分权重占总排名的 60%

## 游戏奖励

第一名：4096 元现金大奖  
第二名：2048 元现金大奖  
第三名：1024 元现金大奖

有钱任性奖（给出金币最多）：512 元现金大奖  
长跑冠军奖（移动步数最多）：512 元现金大奖  
最年轻团队奖（排名前 50%，团队平均年龄最小）：512 元现金大奖

## 游戏规则

0x1、每轮游戏开始时，系统向每个玩家发放 1 个金币

0x1、每轮游戏开始时，系统会在棋盘内随机格子生成 N 个金币（N可能为负数）

0x2、每轮游戏收到地图信息后，你需要在 1s 之内向服务器提交你下一步的位置[x,y]。如果你不小心移动到棋盘外，位置则保持不变，扣除 1 个金币。

0x3、每移动 n 步，扣除玩家 n*1.5(取整) 个金币，原地不动，扣除 1 个金币，当移动所需金币不足时，认为玩家原地不动。每移动一次位置，即为完成一轮游戏。

0x4、每轮移动结束后，玩家将获得这个所在格子中的金币，如果多个玩家在同一个格子，会随机选取其中一位玩家获得金币。

0x5、每轮移动结束后，当有多个玩家在同一个格子中时，金币最高者的玩家将会给出自己 1/4 的金币给每个金币少者的队伍，当给出的金币数量不够每个人 1 个金币时，则不给了.

0x6、每轮游戏结束后，系统将根据玩家的当前金币数量发放利息，利率为 10%，不足整数个的金币不计入。

0x7、每局游戏进行 120 轮后结束。

## 游戏准备

服务器会提供一个 websocket 供你连接，获得你的游戏队伍名称和 token 以保证游戏队伍的唯一性。  
我们会开源所有代码，供你搭建环境，调试程序。当然我们也可能会对游戏本身进行调优，不要忘记及时拉取代码哦~

## 你要做的事

编写一个程序，语言不限，通过算法来计算要移动的位置以吃到各种数量的金币，以保证获得金币最多来赢得比赛胜利。

#### 开发和调试

0x1、搭建环境。拉取代码，启动后端服务，后端环境为 go.(go 语言环境搭建教程请自行百度)。

0x2、编写代码，连接到刚启动的 websocket 连接，一般为`"ws://http://localhost:8881/ws`。

0x3、让代码跑起来。调试，调优。  
msgtype：0 //0 登陆 -1 登陆失败 1 请准备 2 准备好了 3 游戏信息 4 玩家行动 5 本局游戏结束

每轮游戏最终，你需要向服务器提交类似如下具有x，y坐标的移动位置。  
```
ws.send(JSON.stringify({msgtype:4,token:value,x:x,y:y,RoundID:jmsg.RoundID}));
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