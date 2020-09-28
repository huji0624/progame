## 玩法介绍

0x0 这是一个棋盘内吃金币的游戏。

0x0 游戏将被分为 x 局，每局有 10 轮对决。  
0x0 每一轮游戏开始时，系统将生成一个 m\*n 的游戏棋盘，同时，系统将在棋盘内随机分配金币数量。  
0x0 你也将被随机分配在棋盘的任意位置，同一个棋盘格子可能会有多个玩家队伍。  
0x0 你将通过算法来移动位置来吃到各种数量的金币，以保证获得金币最多来赢得比赛胜利。

0x0 游戏进行中，系统会同步统计每局游戏中每个队伍的平均金币数量，金币数量多者排名靠前。你可以在[这里](http:xxxx)进行观看。

0x0 第三天游戏结束后，我们将对游戏队伍进行排名，获得丰厚奖励。

## 游戏规则

0x1、每移动一步，你将失去 1 个金币，同理，移动 n 步，将失去 n 个金币，原地不动，失去一个金币。

0x2、不同玩家移动到同一个棋盘格子，系统将这个格子的金币随机分配给这个格子中的任一队伍，同时，金币最高者的队伍将会分配一个金币给每个金币少者的队伍。

0x3、每轮游戏，你需要在 1s 之内向服务器提交你下一步的位置。[x,y]

0x4、至少 50 局以上游戏玩家才能获得排名的资格。

0x5、

## 游戏准备

服务器会提供一个 websocket 供你连接，获得你的游戏队伍名称和 token 以保证游戏队伍的唯一性。  
我们会开源所有代码，供你搭建环境，调试程序。当然我们也可能会对游戏本身进行调优，不要忘记及时拉取哦~

## 你要做的事

编写一个程序，语言不限。

#### 开发和调试

0x1、搭建环境。拉取代码，后端环境为 go.(go 语言环境搭建教程请自行百度)

0x2、连接到本机的 websocket 连接，一般为`http://localhost:8888`；

0x3、编写代码，让代码跑起来。调试，调优。  
msgtype：0 //0 登陆 -1 登陆失败 1 请准备 2 准备好了 3 游戏信息 4 玩家行动

#### 正式环境

你只需要将 websocket 的链接改为我们线上服务器的链接，就可以啦~

不要忘记在游戏开始时，加入游戏。

## 游戏时间

游戏将分为三个时间段进行。  
第一天：10 月 24 日 19:00-22:00 金币 20%参与排名  
第二天：10 月 25 日 19:00-22:00 金币 30%参与排名  
第三天：10 月 26 日 19:00-22:00 金币 50%参与排名

最后一天，我们将为比较胜利者进行颁奖，不用担心失去前三甲，特别奖可能有你！

## 游戏奖励

第一名：3000 元现金大奖  
第二名：2000 元现金大奖  
第三名：1000 元现金大奖  
特别奖：10 名，每个队伍 500 元

## 另外

0x1 你可以[点击这里](https://gxxx)观看游戏即时排名情况。  
0x2 你可以点击【回放】按钮查看他们的激烈对决。学习和探讨他们是如何取得胜利的，当然，源代码会在赛后开源。