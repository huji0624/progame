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
    uri = "ws://localhost:8888/ws"
    while True:
        async with websockets.connect(uri) as websocket:
            global gws
            gws = websocket
            await send({"msgtype":0,"token":token})
            res = await websocket.recv()
            res = json.loads(res)
            if res['Msgtype']==0:
                print("login ok.")
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
                    await send({"msgtype":4,"token":token,"x":x,"y":y,"RoundID":res['RoundID']})

asyncio.get_event_loop().run_until_complete(hello())