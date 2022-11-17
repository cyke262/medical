# medical

## 启动项目

添加依赖：
```
cd medical && go mod tidy
```
运行项目：
```
./clean_docker.sh
```

打开浏览器，在`localhost:8088`进行访问

账号：`admin`，密码：`123456`

## 关闭项目

```
lsof -i:8088 
kill -9 xxx
```
查询占用8088端口的进程PID，并kill掉

## 报错

+ 若提示无权限运行.sh文件：
```
chmod 777 ./*.sh
```
+ go mod命令无法运行：
```
source /etc/profile
```

## 区块链浏览器

区块链浏览器启动后在`localhost:8080`进行访问

关闭区块链浏览器后输入docker-compose down -v
