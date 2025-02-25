# 思源笔记分享笔记插件的服务器v2

插件实现思路：从插件中导出html文件并上传到web服务器，服务端同时生成对应的index.html作为首页。

服务器实现思路：请求特定的链接时，从数据库中查询`app_id`和`doc_id`，并重定向到服务端的web接口中。

# 一、架构说明

在v2.0中，仅有一个二进制程序，包含了web服务器。数据库采用了sqlite（本地文件的形式）。

# 二、编译

1. 克隆项目
   ```text
   git clone git@github.com:tengfei-xy/siyuan-plugin-share-system-engine.git
   ```
2. 确保存在go 1.21.4或更高的版本，安装[参考文档](https://www.wolai.com/tengfei-xy/bjUcRE7tSsscWqpszDvbxx "参考文档")
3. 进入项目文件夹并开始编译

   务必设置`CGO_ENABLED`=1
   ```纯文本
   cd siyuan-plugin-share-system-engine
   CGO_ENABLED=1 go build -o spss_engine .
   ```

# 三、部署方式（从0部署）

## 3.1 本地二进制运行

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
   sudo apt install musl
   ```
4. 运行启动
   ```bash
   ./spss_engine -c config.yaml
   ```

## 3.2 docker

- 一键部署 / 从docker hub中拉取并运行容器
  
  注： `-p`的容器端口需要和`SPSS_LISTEN`变量的端口保持一直。作为容器的管理者，需要修改的是-p的宿主端口和`SPSS_SHARE_LINK`变量
  注： 对于SPSS_SHARE_LINK变量表示，生成的链接将以此参数作为开头。这个端口不是容器的端口，是最终外部网络来访问容器的端口。示例：如果你在NAS中或其他容器中，假如你需要让外网来访问2233端口，而容器本身使用了25934端口，那么这个参数的端口依然那是2233。

  ```bash
  # IP访问 / http访问 / 关闭API访问密码保护
  docker run -d -p 25934:25934 -v ./:/data \
  -v /etc/localtime:/etc/localtime --restart=always \
  -e SPSS_LISTEN=0.0.0.0:25934 \
  -e SPSS_SAVE_PATH=/data \
  -e SPSS_SHARE_LINK=http://127.0.0.1:25934 \
  -e SPSS_WEB_FILE_MAX=100 \
  --name=spss_engine \
  tengfeixy/spss_engine:2.4.0


  # https访问 / 启动API访问密码保护
  docker run -d -p 25934:25934 -v ./:/data \
  -v /etc/localtime:/etc/localtime --restart=always \
  -e SPSS_LISTEN=0.0.0.0:25934 \
  -e SPSS_SAVE_PATH=/data \
  -e SPSS_SHARE_LINK=https://www.example.com:25934 \
  -e SPSS_WEB_FILE_MAX=100 \
  -e SPSS_WEB_SSL=true \
  -e SPSS_WEB_SSL_CERT=./cert.pem \
  -e SPSS_WEB_SSL_KEY=./cert.key \
  -e SPSS_DB_API=true \
  -e SPSS_DB_AUTH=username:password \
  --name=spss_engine \
  tengfeixy/spss_engine:2.4.0
  ```
- 本地构建镜像
  1. 构建/拉取镜像，[下载Release](https://github.com/tengfei-xy/siyuan-plugin-share-system-engine/releases "下载Release")并解压
      ```bash
      docker build --no-cache=true -t spss_engine .
      ```
  2. 在当前目录创建[config.yaml](https://raw.githubusercontent.com/tengfei-xy/siyuan-plugin-share-system-engine/refs/heads/main/config.yaml "config.yaml")，并修改
  3. 运行容器
     其他环境变量参数参考上文
     ```bash
     docker run -d -p 25934:25934 -v ./:/data -v /etc/localtime:/etc/localtime --restart=always --name=spss_engine spss_engine
     ```

# 四、排错思路

1. 打开思源笔记的分享插件。（注：每次打开此插件的设置页面，都会访问一次服务器，如果访问失败或访问超市则消息弹窗提示）
2. 根据分享插件的服务器地址，测试此端口在本机是否能访问。如果在浏览器中使用思源，要在插件中勾选“从浏览器访问”
3. 如果在插件服务器在NAS中，应该使用NAS的IP和容器映射的外部端口。
4. 在本机的浏览器中直接访问插件设置里的服务器地址，如果返回404 No found，表示通信正常，此时可以使用这个服务器端。
5. 在插件中，执行创建分享后，获得的分享链接的ip、端口务必需要和分享服务期的ip、端口相同。如果不相同，需要修改配置文件里的shareBaseLink或容器变量SPSS_SHARE_LINK。
6. 根据插件设置页面底部的联系方式找作者，使用向日葵远程debug。

# 五、升级（从旧版本1.x.x升级）

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
3. 修改所有的index.html为index.htm

  linux/macos参考命令
  ```bash
  for line in $(find 数据文件路径 -name 'index.html');do
    sudo mv $line $(dirname $line)/index.htm
  done
  ```

# 六、API

说明：API是用来查看数据库的记录，使用API之前，需要确保API接口在配置文件上被启用。如果包含了用户和密码，在请求中加入`Authorization: Basic base64(用户名:密码)`（[参考文档](https://apifox.com/help/auth/basic-auth/)）

1. 查看所有数据
```
GET /api/db/table
```

2. 删除所有数据
```
DELETE /api/db/table
```

3. 提交旧版本的数据库文件(次文件因包含INSERT)
```
PUT /api/db/table
```
