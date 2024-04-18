# 思源笔记分享笔记插件的服务器

插件实现思路：从插件中导出html文件并上传到web服务器，web服务器同时生成对应的index.html作为首页。

服务器实现思路：请求特定的链接时，从mysql中查询app_id和doc_id，并重定向到nginx中。

## 一、docker部署

注意：docker的部署目前尚且还比较稳定，请有动手能力的小伙伴依然小心尝试！但不具有升级的能力，有可能应用升级后，出现未知问题！不适用于生产环境哦~

### 1.1容器说明

1. spss表示siyuan-publig-share-system
2. spss_nginx表示nginx容器，spss_mysql表示mysql容器，spss_engine表示应用容器
3. 容器所需参数和细节说明参考[docker-compose.yml](https://github.com/tengfei-xy/siyuan-plugin-share-system-engine/blob/main/docker/docker-compose.yml)

### 1.2容器启动说明

1. [下载](https://github.com/tengfei-xy/siyuan-plugin-share-system-engine/releases)并解压适合你的系统版本的压缩包

2. 修改docker-compose.yml，其他参数可以默认，`SHARE_BASE_LINK`不能默认，

3. docker-compose启动所有容器

   ```bash
   cd docker/
   sudo docker-compose up -d
   ```

4. 运行以上命令后，spss_mysql将执行初始化，spss_engine由于无法连接到数据库而自动重启，重启直到数据库初始化完成

5. spss_engine启动后，运行如下命令，仅支持linux和mac，WSL未测试

   ```bash
   ./init.sh
   ```

6. 修改思源笔记的分享笔记插件的服务器地址并保存，每打开设置窗口，本插件就会获取链接，当控制台（切换为所有消息级别）输出`{"err":3,"msg":"此页面没有共享","data":""}`表示服务已正常启动。

### 1.3 容器日志输出

nginx容器日志

```bash

sudo docker exec -it spss_nginx cat /var/log/nginx/access.log
sudo docker exec -it spss_nginx cat /var/log/nginx/error.log

```


应用容器日志

```bash
sudo dcoker logs spsse_engine
```

### 1.4 应用容器的镜像删除

当应用容器需要变更时，仅仅运行`sudo docker-compose up -d`是无法让新的应用容器运行，因为已经存在旧的镜像，因此需要先删除旧的容器，再执行`sudo docker-compose up -d`。

其他指令参考

```bash
sudo docker run -d \
  --name spss_engine \
  -e SPSS_STARTUP_ENV=docker \
  -e SHARE_BASE_LINK=http://127.0.0.1 \
  -e SAVE_PATH=/data \
  -e LISTEN_PORT=25934 \
  -v $(pwd):/data \
  --restart always \
  spss_engin
```

### 1.5 应用容器升级过程（不涉及数据库变更的操作）

首先，从最新版的压缩包中提取`spss_engine`二进制程序文件到`docker/`文件夹下，并修改`docker-compose.yaml`的`spss_engine`的镜像tag名，

注：新版的`docker-compose.yaml`可能其他配置也有变化，而不仅仅是应用容器的tag名

其次，停止容器

```bash
sudo docker-compose down spss_engine
```

可选的步骤，删除所有现存的应用镜像

```bash
sudo docker images | grep spss_engine | awk '{print $3}' | xargs -n1 sudo docker rmi
```

最后，构建镜像并后台运行

```bash
sudo docker-compose up -d spss_engine
```





## 二、无docker环境部署（不推荐）

### 2.1 应用所需的配置文件

设置配置文件，文件名config.yaml

```yaml
mysql:
  ip: "127.0.0.1"
  port: "3306"
  # 写死于程序代码中
  username: "share_system"
  # 可自定义
  password: "password"
  # 写死于程序代码中
  database: "share_system"
  
basic:
  # 监听的这个端口需要和nginx转发的地址相互吻合
  listen: "127.0.0.1:25934"
  
  # 插件上传的保存位置
  savePath: "/usr/local/services/nginx-1.25.4/html"
  
  # 生成的短链接，将以此参数作为开头
  shareBaseLink: "http://124.223.15.220"
```

**权限说明**

由插件上传的文件，均以755权限模式放置于`basic.savePath`路径。但要注意nginx用户的启动程序的用户的两者权限是否冲突。

### 2.2 相关服务

**nginx**

作为主要的web服务器，版本无特别要求，作者使用1.25.4

**参考配置**

```nginx
location / {
      proxy_set_header X-Real-IP $remote_addr;
	    if (!-f $request_filename) {
            	proxy_pass http://127.0.0.1:25934/api/url$uri;
            	break;
            }
            root   html;
        }
location /api {
	            proxy_pass http://127.0.0.1:25934;
                    proxy_redirect off;
                    proxy_set_header X-Real-IP $remote_addr;
		          client_max_body_size 100m;
}
```



**mysql**

作为链接和数据的存储，版本无特别要求，作者使用5.7.43

参考：[建库文件](https://github.com/tengfei-xy/siyuan-plugin-share-system-engine/blob/main/docker/sql/ddl.sql)