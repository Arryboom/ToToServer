![crab](http://www.crab.pub/crabs.png)
## ToTo是什么?
一个GO语言编写的、支持正向和反向的端口转发工具，ToToServer是服务端，ToToClient是客户端，服务端可以单独使用(模式1)，也可以搭配客户端使用（模式2），也可以混合使用（模式1、2）

## ToTo服务端的模式1和模式2是什么意思？

* 模式1（还没开发完成）
    *  模式1是反向转发，在服务端配置文件设置密码和控制端口，ToTo客户端通过正确的端口和密码连接服务端之后，客户端可以把本机端口A映射到服务端X(设定或随机)端口上，当别人访问服务器的X端口，就相当于访问客户机的端口A
    *  应用场景：客户机A是某个局域网的电脑，A开了游戏服务器想让B加入，正常情况下B无法直接访问到A。A可以通过ToTo把游戏服务端的端口映射到具有公网IP的ToTo服务器的端口X上，B只要连接ToTo服务器的X端口，就相当于访问A的游戏服务端
    *  说明：支持TCP和UDP转发，支持加密，支持多客户端。如果想一路畅通，请确保ToTo服务器能被互联网访问和足够的带宽和较低的延迟
* 模式2（已开发完成）
    *  模式2是正向转发，可以脱离客户端使用，在服务端配置好远程地址X与协议与本地端口A，当别人访问服务器的端口A就相当于访问远程地址X
    *  应用场景：电脑A和电脑B属于同一个局域网，互联网上的电脑可以访问到A但不能访问到B，AB之间相互连通。假设B开了Web服务器端口80，又想让外网用户访问到的话，可以在电脑A使用ToTo服务端，并配置电脑A的端口80指向电脑B的端口80，这样当外网用户访问电脑A的80端口就相当于访问电脑B的80端口
    *  说明：支持TCP和UDP转发，支持HTTP头的域名替换。如果想一路畅通，请确保ToTo服务器能被互联网访问和足够的带宽和较低的延迟


## 编译

复制到GOPATH目录编译。对于没接触过GO语言的人来说可能有点难，可以联系我帮你编译，发送邮件给我就行，记得附上目标平台。

## 配置

把下面的内容修改后保存到服务端主程序目录下，文件名为config.json
```json
    {
        "ServerPassword": "123456",
        "ServerPort": "5244",
        "EnableIPV6Support": true,
        "EnableMode1": true,
        "EnableMode2": true,
        "Mode2UdpAliveTime":5,
        "Debug":true,
        "Mode2Config": [
            {
                "RemoteAddr": "45.32.34.191",
                "RemotePort": "80",
                "Protocol":"TCP",
                "LocalPort": "80",
                "HTTPHeaderHostReplace":"crab.pub"
            },
            {
                "RemoteAddr": "45.32.34.191",
                "RemotePort": "22",
                "Protocol":"TCP",
                "LocalPort": "23"
            },
            {
                "RemoteAddr": "114.114.114.114",
                "RemotePort": "53",
                "Protocol":"UDP",
                "LocalPort": "53"
            }
        ]
    }
```
上面示例配置由于模式1还没开发完成，所以请忽略掉ServerPassword、ServerPort、EnableIPV6Support、EnableMode1。

* EnableMode2字段为是否开启模式2
* Mode2UdpAliveTime字段为模式2的UDP转发的心跳时间，被转发的远程机器每发送一个封包，心跳时间就更新一次。超时将断开连接
* Debug字段为输出调试信息的开关
* Mode2Config字段为一个数组，里面放了需要转发的端口的信息，看下面详解
```json
    {
            "RemoteAddr": "45.32.34.191",
            "RemotePort": "80",
            "Protocol":"TCP",
            "LocalPort": "80",
            "HTTPHeaderHostReplace":"crab.pub"
    }
```
* 上面的意思是，转发45.32.34.191的80端口的crab.pub域名的HTTP服务到ToToServer本机的80端口
* RemoteAddr和RemotePort为远程地址和端口
* Protocol为协议，选择有TCP和UDP
* LocalPort为转发到本机哪个端口
* HTTPHeaderHostReplace为可选字段，假设你想转发远程HTTP端口，然后远程地址的80端口绑定了很多个域名，你只想转发某个域名，就可以使用此字段来替换HTTP头的Host属性
* 如果你没用上HTTPHeaderHostReplace，请删掉，否则会出现不可预料的结果！！

## 有问题反馈
在使用中有任何问题，欢迎反馈给我，可以用以下联系方式跟我交流

* 邮件(529493022#qq.com, 把#换成@)
* QQ: 529493022
