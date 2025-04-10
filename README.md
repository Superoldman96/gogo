# gogo

blog posts:

- https://chainreactors.github.io/wiki/blog/2022/11/15/gogo-introduce/

![](https://socialify.git.ci/chainreactors/gogo/image?description=1&font=Inter&forks=1&issues=1&language=1&name=1&owner=1&pattern=Circuit%20Board&pulls=1&stargazers=1&theme=Light)


## Features


* 自由的端口配置
* 支持主动/被动指纹识别
* 关键信息提取, 如title, cert 以及自定义提取信息的正则
* 支持nuclei poc, 引擎: https://github.com/chainreactors/neutron
* 无害的扫描, 每个添加的poc都经过人工审核
* 可控的启发式扫描
* 超强的性能, 最快的速度, 尽可能小的内存与CPU占用.
* 最小发包原则, 尽可能少地发包获取最多的信息
* 支持DSL, 可以通过修改的配置文件自定义自己的gogo
* 完善的输出与输出设计
* 几乎不依赖第三方库, 纯原生go编写, 在windows 2003上也可以使用完整的漏洞/指纹识别功能
 
## QuickStart

完整的文档与教程位于wiki: https://chainreactors.github.io/wiki/gogo/

指纹与poc仓库: https://github.com/chainreactors/templates

### 最简使用

指定网段进行默认扫描, 并在命令行输出

`gogo -i 192.168.1.1/24 -p win,db,top2 `

### 端口配置

一些常用的端口配置:

* `-p -`  等于`-p 1-65535`
* `-p 1-1000` 端口范围
* `-p common` tag: common 表示内网常用端口
* `-p top2,top3` 可以同时选择多个tag. 外网常见web端口
* `-p all` 表示所有预设的tag的合集.

通过逗号分割多个配置, 可根据场景进行各种各样的组合配置. 例如:

`gogo -i 1.1.1.1/24 -p 1-1000,common,http,db`

**查看全部端口配置**

`gogo -P port`

可查看所有的tag对应的端口. 

```
当前已有端口配置: (根据端口类型分类)
         top1 :  80,443,8080
         top2 :  70,80,81,82,83,84,85,86,87,88,89,90,443,1080,2000,2001,3000,3001,1443,4443,4430,5000,5001,5601,6000,6001,6002,6003,7000,7001,7002,7003,9000,9001,9002,9003,8080,8081,8082,8083,8084,8085,8086,8087,8088,8089,8090,8091,8000,8001,8002,8003,8004,8005,8006,8007,8008,8009,8010,8011,8012,8013,8014,8015,8016,8017,8018,8019,8020,8820,6443,8443,9443,8787,7080,8070,7070,7443,9080,9081,9082,9083,5555,6666,7777,7788,9999,6868,8888,8878,8889,7890,5678,6789,9090,9091,9092,9093,9094,9095,9096,9097,9098,9099,9100,9988,9876,8765,8099,8763,8848,8161,8060,8899,800,801,888,10000,10001,10002,10003,10004,10005,10006,10007,10008,10009,10010,1081,1082,10080,10443,18080,18000,18088,18090,19090,19091,50070
         top3 :  444,9443,6080,6443,9070,9092,9093,7003,7004,7005,7006,7007,7008,7009,7010,7011,9003,9004,9005,9006,9007,9008,9009,9010,9011,8100,8101,8102,8103,8104,8105,8106,8107,8108,8109,8110,8111,8161,8021,8022,8023,8024,8025,8026,8027,8028,8029,8030,8880,8881,8882,8883,8884,8885,8886,8887,8888,8889,8890,8010,8011,8012,8013,8014,8015,8016,8017,8018,8019,8020,8090,8091,8092,8093,8094,8095,8096,8097,8098,8099,8180,8181,8983,1311,8363,8800,8761,8873,8866,8900,8282,8999,8989,8066,8200,8040,8060,10800,18081
         docker :  2375,2376,2377,2378,2379,2380
         lotus :  1352
         dubbo :  18086,20880,20881,20882
         oracle :  1158,1521,11521,210
         ...
         ...
         ...
```

### 启发式扫描

当目标范围的子网掩码小于24时, 建议启用 smart模式扫描([原理见doc](https://chainreactors.github.io/wiki/gogo/detail/#_7)), 例如子网掩码为16时(输出结果较多, 建议开启--af输出到文件, 命令行只输出日志)

`gogo -i 172.16.1.1/12 -m ss --ping -p top2,win,db --af`

`--af` 表示自动指定文件生成的文件名.

`-m ss` 表示使用supersmart模式进行扫描. 还有ss,sc模式适用于不同场景

`--ping` 表示在指纹识别/信息获取前判断ip是否能被ping通, 减少无效发包. **需要注意的是, 不能被ping通不代表目标一定不存活, 使用时请注意到这一点**

### workflow

启发式扫描的命令有些复杂, 但可以使用workflow将复杂的命令写成配置文件, 快捷调用([内置的workflow细节见doc](https://chainreactors.github.io/wiki/gogo/start/#workflow)).

 `gogo -w 172` 

即可实现与`gogo -i 172.16.1.1/12 -m ss --ping -p top2,win,db --af` 完全相同的配置

**查看所有workflow**

`gogo -P workflow` 

常用的配置已经被集成到workflow中, 例如使用supersmart mod 扫描10段内网, `gogo -w 10`即可. 

还有一些预留配置(即填写了其他配置, 但没有填写目标, 需要-i手动指定目标), 例如:

`gogo -w ss -i 11.0.0.0/8`

workflow中的预设参数优先级低于命令行输入, 因此可以通过命令行覆盖workflow中的参数. 例如:

`gogo -w 10 -i 11.0.0.0/8`

### 示例 

**一个简单的任务**

`gogo -i 81.68.175.32/28 -p top2`

```
gogo -i 81.68.175.32/28 -p top2
[*] Current goroutines: 1000, Version Level: 0,Exploit Target: none, PortSpray: false ,2022-07-07 07:07.07
[*] Start task 81.68.175.32/28 ,total ports: 100 , mod: default ,2022-07-07 07:07.07
[*] ports: 80,81,82,83,84,85,86,87,88,89,90,443,1080,2000,2001,3000,3001,4443,4430,5000,5001,5601,6000,6001,6002,6003,7000,7001,7002,7003,9000,9001,9002,9003,8080,8081,8082,8083,8084,8085,8086,8087,8088,8089,8090,8000,8001,8002,8003,8004,8005,8006,8007,8008,8009,8010,8011,8012,8013,8014,8015,8016,8017,8018,8019,8020,6443,8443,9443,8787,7080,8070,7070,7443,9080,9081,9082,9083,5555,6666,7777,9999,6868,8888,8889,9090,9091,8091,8099,8763,8848,8161,8060,8899,800,801,888,10000,10001,10080 ,2022-07-07 07:07.07
[*] Scan task time is about 8 seconds ,2022-07-07 07:07.07
[+] http://81.68.175.33:80      nginx/1.16.0            nginx                   bd37 [200] HTTP/1.1 200
[+] http://81.68.175.32:80      nginx/1.18.0 (Ubuntu)           nginx                   8849 [200] Welcome to nginx!
[+] http://81.68.175.34:80      nginx           宝塔||nginx                     f0fa [200] 没有找到站点
[+] http://81.68.175.34:8888    nginx           nginx                   d41d [403] HTTP/1.1 403
[+] http://81.68.175.34:3001    nginx           webpack||nginx                  4a9b [200] shop_mall
[+] http://81.68.175.37:80      Microsoft-IIS/10.0              iis10                   c80f [200] HTTP/1.1 200             c0f6 [200] 安全入口校验失败
[*] Alive sum: 5, Target sum : 1594 ,2022-07-07 07:07.07
[*] Totally run: 4.0441884s ,2022-07-07 07:07.07
```

如果要联动其他工具, 可以指定`-q/--quiet`关闭日志信息, 只保留输出结果.

### 输出与再处理

关于输入输出以及各种高级用法请见[output的wiki](https://chainreactors.github.io/wiki/gogo/start/#output)

如果执行`gogo -i 81.68.175.1 --af`

扫描完成后, 可以看到在gogo二进制文件同目录下, 生成了`.81.68.175.1_28_all_default_json.dat1`, 该文件是deflate压缩的json文件.

通过gogo格式化该文件, 获得human-like的结果

```
 gogo  -F .\.81.68.175.1_28_all_default_json.dat1
Scan Target: 81.68.175.1/28, Ports: all, Mod: default
Exploit: none, Version level: 0

[+] 81.68.175.32
        http://81.68.175.32:80  nginx/1.18.0 (Ubuntu)           nginx                   8849 [200] Welcome to nginx!
        tcp://81.68.175.32:22                   *ssh                     [tcp]
        tcp://81.68.175.32:389                                           [tcp]
[+] 81.68.175.33
        tcp://81.68.175.33:3306                 *mysql                   [tcp]
        tcp://81.68.175.33:22                   *ssh                     [tcp]
        http://81.68.175.33:80  nginx/1.16.0            nginx                   bd37 [200] HTTP/1.1 200
[+] 81.68.175.34
        tcp://81.68.175.34:3306                 mysql 5.6.50-log                         [tcp]
        tcp://81.68.175.34:21                   ftp                      [tcp]
        tcp://81.68.175.34:22                   *ssh                     [tcp]
        http://81.68.175.34:80  nginx           宝塔||nginx                     f0fa [200] 没有找到站点
        http://81.68.175.34:8888        nginx           nginx                   d41d [403] HTTP/1.1 403
        http://81.68.175.34:3001        nginx           webpack||nginx                  4a9b [200] shop_mall
[+] 81.68.175.35
        http://81.68.175.35:47001       Microsoft-HTTPAPI/2.0           microsoft-httpapi                       e702 [404] Not Found
[+] 81.68.175.36
        http://81.68.175.36:80  nginx   PHP     nginx                   babe [200] 风闻客栈24小时发卡中心 - 风闻客栈24小时发卡中心
        tcp://81.68.175.36:22                   *ssh                     [tcp]
...
...
```

**导出到其他工具**

一些常用的输出格式.

* `-o full` 默认输出格式, 即上面示例所示.
* `-o color` 带颜色的full输出. 在v2.11.0版本之后, -F 输出到命令行时为默认开启状态. 如果需要关闭, 手动指定`-o full`即可
* `-o jl`  一行一个json, 可以通过管道传给jq实时处理
* `-o json` 一个大的json文件
* `-o url` 只输出url, 通常在`-F`时使用

所有的输出格式见: https://chainreactors.github.io/wiki/gogo/start/#_4

**输出过滤器**

`--filter` 参数可以从dat文件中过滤出指定的数据并输出.

例如过滤指定字段的值: `gogo -F 1.dat --filter framework::redis -o target` 表示从1.dat中过滤出redis的目标, 并输出为target字段.

其中`::` 表示模糊匹配, 还有其他三种语法,如 `==` 为精准匹配, `!=` 为不等于, `!:` 为不包含

`-F 1.json -f file` 重新输出到文件, 也可以`-F 1.dat --af` 自动生成格式化后的文件名. 

## 注意事项

* **(重要)**因为并发过高,可能对路由交换设备造成伤害, 例如某些家用路由设备面对高并发可能会死机, 重启, 过热等后果. 因此在外网扫描的场景下**建议在阿里云,华为云等vps上使用**,如果扫描国外资产,建议在国外vps上使用.本地使用如果网络设备性能不佳会带来大量丢包. 如果在内网扫描需要根据实际情况调整并发数.
* 如果使用中发现疯狂报错,大概率是io问题(例如多次扫描后io没有被正确释放,或者配合proxifier以及类似代理工具使用报错),可以通过重启电脑,或者虚拟机中使用,关闭代理工具解决.如果依旧无法解决请联系我们.
* 还需要注意,upx压缩后的版本虽然体积小,但是有可能被杀软杀,也有可能在部分机器上无法运行.
* 一般情况下无法在代理环境中使用,除非使用-t参数指定较低的速率(默认并发为4000).
* gogo本身并不具备任何攻击性, 也无法对任何漏洞进行利用.
* **使用gogo需先确保获得了授权, gogo反对一切非法的黑客行为**

### 使用场景并发推荐

默认的并发linux为4000, windows为1000, 为企业级网络环境下可用的并发. 不然弱网络环境(家庭, 基站等)可能会导致网络dos

建议根据不同环境,手动使用-t参数指定并发数.

* 家用路由器(例如钓鱼, 物理, 本机扫描)时, 建议并发 100-500
* linux 生产网网络环境(例如外网突破web获取的点), 默认并发4000, 不需要手动修改
* windows 生产网网络环境, 默认并发1000, 不需要手动修改
* 高并发下udp协议漏报较多, 例如获取netbois信息时, 建议单独对udp协议以较低并发重新探测
* web的正向代理(例如regeorg),建议并发 10-30
* 反向代理(例如frp), 建议并发10-100

如果如果发生大量漏报的情况, 大概率是网络环境发生的阻塞, 倒是网络延迟上升超过上限.

因此也可以通过指定 `-d 5 `(tcp默认为2s, tls默认为两倍tcp超时时间,即4s)来提高超时时间, 减少漏报.

未来也许会实现auto-tune, 自动调整并发速率

**这些用法大概只覆盖了一小半的使用场景, 请[阅读文档](https://chainreactors.github.io/wiki/gogo/)**

## Make

### 手动编译

```bash
# download
git clone --recurse-submodules https://github.com/chainreactors/gogo
cd gogo/v2

# sync dependency
go mod tidy   

# generate template.go
# 注意: 如果需要使用go1.10编译windows03可用版本， 也需要先使用高版本的go generate生成相关依赖
go generate  

# build 
go build .

# windows server 2003 compile
GOOS=windows GOARCH=386 go1.10 build .

# 因为go1.10 还没有go mod, 可能会导致依赖报错. 如果发生了依赖报错, 可以使用go1.11 编译. 
# go1.11 官方声明不支持windows server 2003 , 实测可以稳定运行(需要调低并发).
GOOS=windows GOARCH=386 go1.11 build .
```

如果需要编译windows xp/2003的版本, 请先使用高版本的go生成templates. 再使用go 1.11编译即可.

## Similar or related works

* [ServerScan](https://github.com/Adminisme/ServerScan) 早期的简易扫描器, 功能简单但开拓了思路
* [fscan](https://github.com/shadow1ng/fscan) 简单粗暴的扫描器, 细节上有不少问题, 但胜在简单. 参考其直白的命令行，设计了workflow相关功能.
* [kscan](https://github.com/lcvvvv/kscan) 功能全面的扫描器, 从中选取合并了部分指纹
* [ladongo](https://github.com/k8gege/LadonGo) 集成了各种常用功能, 从中学习了多个特殊端口的信息收集
* [cube](https://github.com/JKme/cube) 与fscan类似, 从中学习了NTLM相关协议的信息收集

gogo从这些相似的工作中改进自身. 感谢前人的工作. 

细节上的对比请看[文档](https://chainreactors.github.io/wiki/gogo/design/)

## THANKS

* https://github.com/projectdiscovery/nuclei-templates
* https://github.com/projectdiscovery/nuclei
