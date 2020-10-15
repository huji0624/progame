#!/usr/bin/env python

import asyncio
import websockets
import json
import sys


if len(sys.argv)==1:
    token = "hj"
else:
    token = sys.argv[1]

gws = None

async def send(obj):
    global gws
    s = json.dumps(obj)
    print("will send:")
    print(s)
    await gws.send(s)

async def hello():
    #online
    # uri = "ws://pgame.51wnl-cq.com:8881/ws"
    #local debug
    uri = "ws://localhost:8881/ws"
    while True:
        await asyncio.sleep(1)
        try:
            async with websockets.connect(uri) as websocket:
                global gws
                gws = websocket
                await send({"msgtype":0,"token":token})
                res = await websocket.recv()
                res = json.loads(res)
                if res['Msgtype']==0:
                    print("login ok:",token)
                else:
                    await asyncio.sleep(1)
                    continue
                while True:
                    res = await websocket.recv()
                    print("recv:")
                    print(res)
                    res = json.loads(res)
                    if res['Msgtype']==1:
                        await send({"msgtype":2,"token":token})
                    elif res['Msgtype']==3:
                        import random
                        x = random.randint(0,res['Wid']-1)
                        y = random.randint(0,res['Hei']-1)
                        
                        #这就是你最终需要上传的具有x、y轴坐标的位置信息
                        #发送移动位置时，必须带上发送给你的RoundID
                        await send({"msgtype":4,"token":token,"x":x,"y":y,"RoundID":res['RoundID']})
        except IOError as e:
            print(e)
        except websockets.exceptions.ConnectionClosedError as e:
            print(e)
        else:
            continue

asyncio.get_event_loop().run_until_complete(hello())