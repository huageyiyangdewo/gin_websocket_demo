# gin_websocket_demo
gin websocket 即时聊天demo

```bash
docker exec -it mysql bash

mysql -uroot -p123456

update user set host='%' where user='root';
flush privileges;
create database `chat` default character set utf8mb4;

```

https://www.jianshu.com/p/b039fccf37c9

https://www.cnblogs.com/BillyLV/articles/12842922.html


docker compose  TypeError: expected string or bytes-like object

```bash
docker exec -it mongodb bash
mongo -uroot -p123456

# 新建chat库
use chat
db.createUser({user:"root",pwd:"123456",roles:[{role:'root',db:'chat'}]})
```