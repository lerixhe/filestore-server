# filestore-server
慕课网：实现企业级分布式云存储系统

## 安装mysql与redis

采用容器方式安装mysql与redis

### 容器安装并运行mysql

```shel
dock pull mysql:5.6
docker run -p 3306:3306 --name mysql -v /usr/local/docker/mysql/conf:/etc/mysql -v /usr/local/docker/mysql/logs:/var/log/mysql -v /usr/local/docker/mysql/data:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=123456 -d mysql:5.6
```

## 修改容器中mysql密码

```shell

#进入容器
docker exec -it [容器名] /bin/bash

mysql -uroot -p

mysql> use mysql;
#更改root密码
mysql> UPDATE user SET Password = password ( 'a123456' ) WHERE User = 'root';
# 刷新权限
mysql> flush privileges;
# 退出

```

### 容器安装并运行redis

```shell
docker pull redis:latest
docker run -itd --name redis-test -p 6379:6379 redis
```

### 修改容器中redis密码

```shell

#进入容器
docker exec -it [容器名] /bin/bash

redis-cli

redis 127.0.0.1:6379> config set requirepass testupload


```

### 其他

window下需要使用虚拟机安装docker,虚拟机使用NAT方式与宿主机通信，这导致访问mysql与redis无法使用localhost,而是虚拟机IP
为保证代码无配置运行，需要进行端口映射。
