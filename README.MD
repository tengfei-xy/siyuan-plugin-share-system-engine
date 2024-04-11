# 思源笔记分享笔记插件的服务器

插件实现思路：从插件中导出html文件并上传到web服务器，web服务器同时生成对应的index.html作为首页。

服务器实现思路：请求特定的链接时，从mysql中查询app_id和doc_id，并重定向到nginx中。

## 一、docker部署

注意：docker的部署不是特别稳定，请有动手能力的小伙伴小心尝试！当前代码不具有升级的能力，有可能应用升级后，数据无法恢复。除非自行导出！

### 1.1容器说明

1. spss表示siyuan-publig-share-system
2. spss_nginx表示nginx容器，spss_mysql表示mysql容器，spss_engine表示应用容器

### 1.2容器启动说明

1. [下载](https://github.com/tengfei-xy/siyuan-plugin-share-system-engine/releases)并解压适合你的系统版本的压缩包

2. docker-compose启动所有容器

   ```bash
   cd docker/
   sudo docker-compose up -d
   ```

3. 运行以上命令后，spss_mysql将执行初始化，spss_engine由于无法连接到数据库而自动重启，重启直到初始化完成

4. spss_engine启动后，运行如下命令

   ```bash
   ./init.sh
   ```

5. 修改思源笔记的分享笔记插件的服务器地址并保存，每打开设置窗口，本插件就会获取链接，当控制台（切换为所有消息级别）输出`{"err":3,"msg":"此页面没有共享","data":""}`表示服务已正常启动。

### 1.3 容器日志输出

nginx容器日志

```bash
sudo dcoker logs spsse_nginx
```


应用容器日志

```bash
sudo dcoker logs spsse_engine
```

### 1.4 应用容器的镜像删除

当应用容器需要变更时，仅仅运行`sudo docker-compose up -d`是无法让新的应用容器运行，因为已经存在旧的镜像，因此需要先删除旧的容器，再执行`sudo docker-compose up -d`。

其他参考指令

```bash
sudo docker run -d \
  --name spss_engine \
  -e SPSS_STARTUP_ENV=docker \
  -e SHARE_BASE_LINK=http://127.0.0.1 \
  -e SAVE_PATH=/data \
  -e LISTEN_PORT=25934 \
  -v $(pwd):/data \
  --restart always \
  spss_engine
```



### 1.5 容器所需参数和细节

```yaml
# spss表示siyuan-publig-share-system

services:
  # 服务名,不可更改
  spss_nginx:

    # 镜像名,不建议修改
    image: "nginx:1.25.4"
    # 容器名,不建议修改，启停时所需要的参数
    container_name: spss_nginx

    # 端口映射，本地端口:容器端口，本地端口将作为分享链接的中的地址端口，如果是80，分享地址将忽略
    ports:
      - "80:80"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf
      - ./nginx/conf.d:/etc/nginx/conf.d
      # html的资源文件夹，spsse_engine的资源文件将保存与此
      - ./nginx/html:/usr/share/nginx/html
      - ./nginx/logs:/var/log/nginx

    # 停止时，是否重启，不建议修改
    restart: always
    # 网络配置，不可删除
    networks:
      - spss_network

  spss_mysql:
    # 镜像名,不建议修改
    image: "mysql:5.7.43"

    # 容器名,不建议修改，启停时所需要的参数
    container_name: spss_mysql
    volumes:
      - ./mysql/my.cnf:/etc/mysql/my.cnf
    environment:
      MYSQL_ROOT_PASSWORD: "password"
      MYSQL_USER: 'share_system'
      MYSQL_PASS: 'password'
      MYSQL_DATABASE: 'share_system'
    restart: always
    privileged: true

    # 网络配置，不可删除
    networks:
      - spss_network

  spss_engine:
    # 构建的镜像名和tag，不建议修改
    image: "spss_engine:0.1.0"

    # 容器名,不建议修改，启停时所需要的参数
    container_name: spss_engine

    # 构建容器，不可修改
    build:
      context: ../
      dockerfile: ./docker/Dockerfile
      
    # 容器内的变量
    environment:
      # 不可修改
      SPSS_STARTUP_ENV: "docker"

      # 分享链接的基本地址，可以修改，
      # 此参数结合nginx的映射的本地端口，组成完整的分享地址
      SHARE_BASE_LINK: "http://127.0.0.1"
      # 不建议修改，如果要修改，同时修改spss_engine的映射地址
      SAVE_PATH: "/data"
      # 不建议修改，如果要修改，同时修改nginx.conf里的对应的端口
      LISTEN_PORT: 25934
      # 不建议修改，本程序适用的HTTP框架
      GIN_MODE: release
    # 意外停止时的重启动作，不建议修改
    restart: always
    # 挂载的文件夹，不建议修改
    volumes:
      - ./:/data
    # 网络配置，不可删除
    networks:
      - spss_network
# 网络配置，不可删除
networks:
  spss_network:
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