GOOS=linux GOARCH=amd64 go build
ssh root@10.29.85.185 'kill `cat pgame/pid.log`'
scp ./progame root@10.29.85.185:~/pgame/.
ssh root@10.29.85.185 'cd pgame;nohup ./progame > ./pggame.log 2>&1 &'
echo "Deploy OK!"
