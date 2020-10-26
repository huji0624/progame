import json
import pprint

def getResult(name):
    f = open(name+".json",encoding='utf-8')
    t = f.read()
    f.close()
    return json.loads(t)

def calFinal(fn,res,ratio):
    for v in res["Sorted"]:
        name = v['Name']
        score = v['Gold']
        gc = v['GameCount']
        if gc>=64:
            fn[name] = fn.get(name,0) + score*ratio

first = getResult("first")
second = getResult("second")
third = getResult("third")

give = third['Gives']
move = third["Moves"]

final = {}
calFinal(final,first,0.1)
calFinal(final,second,0.3)
calFinal(final,third,0.6)

final = sorted(final.items(),key=lambda x:x[1],reverse=True)

give = sorted(give.items(),key=lambda x:x[1],reverse=True)
print("Gives Rank:")
pprint.pprint(give)

move = sorted(move.items(),key=lambda x:x[1],reverse=True)
print("Moves Rank:")
pprint.pprint(move)

print("Final Score:")
pprint.pprint(final)

