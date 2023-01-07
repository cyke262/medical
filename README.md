# medical

添加依赖：
```
cd medical_testdemo && go mod tidy
```
运行项目：
```
./clean_docker.sh
```
在`localhost:8088`进行访问

```
lsof -i:8088 
kill -9 xxx
```
查询占用8088端口的进程PID，并kill掉

区块链浏览器启动后在`localhost:8080`进行访问

关闭区块链浏览器后输入docker-compose down -v

## 虚拟机环境

### 获取链接

夸克网盘链接：https://pan.quark.cn/s/4f4afd35cd55
提取码：XDLM

登录名/root :  `jack`

### 开发工具使用：

单独开 **两个终端**：

`Navicat`:  localhost:3306 root:root

```shell
cd /home/jack/Software
./ForNavicat的激活与无限试用.sh
```

`GoLand`: 2021.3 Pro
```shell
cd /opt/GoLand-2021.3.3/bin/
./goland.sh
```