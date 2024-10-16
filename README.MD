# 思源笔记分享笔记插件的服务器v2

插件实现思路：从插件中导出html文件并上传到web服务器，服务端同时生成对应的index.html作为首页。

服务器实现思路：请求特定的链接时，从数据库中查询app\_id和doc\_id，并重定向到服务端的web接口中。

# 一、架构说明

在v2.0中，仅有一个二进制程序，包含了web服务器。数据库采用了sqlite（本地文件的形式）。

# 二、编译

1. 克隆项目
   ```text
   git clone git@github.com:tengfei-xy/siyuan-plugin-share-system-engine.git
   ```
2. 确保存在go 1.21.4版本，安装[参考文档](https://www.wolai.com/tengfei-xy/bjUcRE7tSsscWqpszDvbxx "参考文档")
3. 进入项目文件夹并开始编译

   务必设置CGO\_ENABLED=1
   ```纯文本
   cd siyuan-plugin-share-system-engine
   CGO_ENABLED=1 go build -o spss_engine .
   ```

## 二、部署方式（从0部署）

## 2.1 本地二进制运行

注：不提供32位系统的二进制文件和windows arm。

1. [下载Release](https://github.com/tengfei-xy/siyuan-plugin-share-system-engine/releases "下载Release")并解压
2. 修改配置文件
3. 安装musl

   centos 7
   ```yaml
   sudo yum install -y musl-libc-static
   ```
   debian / ubuntu
   ```yaml
   sudo apt instal musl
   ```
4. 运行启动
   ```bash
   ./spss_engine -c config.yaml
   ```

## 2.2 docker

- 一键部署
  ```bash
  # IP访问 / http访问 / 关闭API访问密码保护
  docker run -d -p 25934:25934 -v ./:/data \
  -v /etc/localtime:/etc/localtime --restart=always \
  -e SPSS_LISTEN=127.0.0.1:25934 \
  -e SPSS_SAVE_PATH=/data \
  -e SPSS_SHARE_LINK="http://127.0.0.1:25934" \
  -e SPSS_WEB_FILE_MAX=100 \
  --name=spss_engine \
  tengfeixy/spss_engine:2.0.0


  # https访问 / 启动API访问密码保护
  docker run -d -p 25934:25934 -v ./:/data \
  -v /etc/localtime:/etc/localtime --restart=always \
  -e SPSS_LISTEN=0.0.0.0:25934 \
  -e SPSS_SAVE_PATH=/data \
  -e SPSS_SHARE_LINK="https://www.example.com:25934" \
  -e SPSS_WEB_FILE_MAX=100 \
  -e SPSS_WEB_SSL=true \
  -e SPSS_WEB_SSL_CERT=./cert.pem \
  -e SPSS_WEB_SSL_KEY=./cert.key \
  -e SPSS_DB_API=true \
  -e SPSS_DB_AUTH=username:password \
  --name=spss_engine \
  tengfeixy/spss_engine:2.0.0
  ```
- 使用配置文件
  1. 构建/拉取镜像
     - 本地构建镜像,[下载Release](https://github.com/tengfei-xy/siyuan-plugin-share-system-engine/releases "下载Release")并解压
       ```bash
       docker build --no-cache=true -t spss_engine .
       ```
     - DockerHub拉取镜像
       ```bash
       docker build -t spss_engine .
       ```
  2. 在当前目录创建[config.yaml](https://raw.githubusercontent.com/tengfei-xy/siyuan-plugin-share-system-engine/refs/heads/main/config.yaml "config.yaml")，并修改
  3. 运行容器
     ```bash
     docker run -d -p 25934:25934 -v ./:/data -v /etc/localtime:/etc/localtime --restart=always --name=spss_engine spss_engine
     ```

# 三、升级（从旧版本1.x.x升级）

1. 从mysql导出sql，参考语句
   ```bash
   bin/mysqldump -B share_system -u root -p -S /tmp/mysql.sock > ~/share_system.sql
   ```
2. 导入方式

   注：程序只搜索导入的数据文件的开头为INSERT的语句。表初始化和创建又程序启动时自动完成。导入前，将清空数据表。
   - 命令行
     - 服务器在容器中
       1. 将sql文件放入容器中
          ```bash
          docker cp ./share_system.sql  spss_engine:/
          ```
       2. 通过程序导入sql文件
          ```bash
          docker exec -it spss_engine /usr/sbin/spss_engine -d /share_system.sql
          ```
     - 服务器无docker环境
       ```bash
       spss_engine -d ./share_system.sql
       ```
   - 网络

     在服务器已启动后，通过curl或其他工具上传文件。

     注：此方法需要设置配置文件`apiEnable: true`或变量`SPSS_DB_API=true`，并设置用户名、密码。运行正常时，将返回被插入的数据。
     ```bash
     curl -X PUT -u username:password http://<your_address>/api/db/table -F "file=@share_system.sql"
     ```
