# medical

添加依赖：
```
cd medical && go mod tidy
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

夸克网盘链接：https://pan.quark.cn/s/4f4afd35cd55
提取码：XDLM

登录名/root :  `jack`

## Ubuntu22.04虚拟机环境

网盘链接：https://www.123pan.com/s/q9USVv-8d3l

提取码:1J6E

登录名ExcitedFrog，密码与root：19660813

> 说明：
>
> - 区块链版本为Fabric 2.2.0，搭建步骤见https://blog.csdn.net/weixin_44165950/article/details/124857431
> - golang版本为1.18.8
> - 含有vscode、vim和notepad++开发环境
> - 已更换为国内镜像源，换源步骤和内容（版本选22.04）见https://mirrors.pku.edu.cn/Help/Ubuntu
> - MySQL版本为8.0，用户和密码都为root，端口为3306
> - Navicat在主目录下的navicat文件夹内
> - docker使用的是docker.io，未换源