# 动创云订货管理平台后端[PC端,超管|大B]API
## 超管
admin
dynamic!@#eos

## APP说明
### Company
负责大B业务管理API
### Shop
负责小B业务管理API
### Tools
负责工具类API
### weapp
负责微信小程序API
## RedisDb规划
主要存放客户端的信息
前6个DB存放一些配置类的信息
* db0:存放手机号短信码验证
* db1:存放客户端菜单
* db2:存放客户端首页
* db3:存放客户端商品分类
* db4:存放小B购物车内容
* db5:存放我的菜单配置信息
后3个DB存放订单
* db6:存放订单
* db7:存放订单
* db8:存放订单
## Api使用
```shell
1.项目初始化,数据库初始化
./go-admin migrate -c config/settings.yml

2.项目启动s
./go-admin server -c config/settings.yml

3.生成migrate模板
./go-admin migrate -c config/settings.yml -g false -a true

4.创建app
./go-admin app -n test

```
## 初始化sql问题
1. 有时候MySQL不支持datetime为0的情况,设置db允许时间为空
```shell
[mysqld]
sql_mode=ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION

```
## 构建上传
sudo CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dynamic  && scp -P 26622 -i ~/.ssh/id_rsa -r dynamic root@152.136.36.253:/home/chaoqun

### 项目关联
1. 获取本机公钥地址:/Users/zhaichaoqun/.ssh/id_rsa.pub
2. 在github settings中SSH keys / Add new
3. 在本机ssh-add /Users/zhaichaoqun/.ssh/id_rsa 添加私钥尝试访问github
4.  ssh git@github.com 测试是否可访问
### Systemd 方式启动:
```shell
cat > /etc/systemd/system/api.service << "END"
[Unit]
Description=DyApi
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/usr/local/dynamic/
## 注:根据可执行文件路径修改
ExecStart=/usr/local/dynamic/api server -c config/settings.yml

# auto restart
StartLimitIntervalSec=0
Restart=always
RestartSec=1

[Install]
WantedBy=multi-user.target
##################################
END


systemctl daemon-reload

systemctl start api.service
systemctl status api.service
systemctl enable api.service

```
