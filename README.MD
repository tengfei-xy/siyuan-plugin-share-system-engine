# 思源笔记分享笔记插件的服务器

插件实现思路：从插件中导出html文件并上传到web服务器，web服务器同时生成对应的index.html作为首页。

服务器实现思路：请求特定的链接时，从mysql中查询app_id和doc_id，并重定向到nginx中。

## docker部署

1. 克隆项目

   ```bash
   git clone git@github.com:tengfei-xy/siyuan-plugin-share-system.git
   cd siyuan-plugin-share-system
   ```

2. 构建应用镜像

   ```bash
   sudo docker build -t spss_engine  -f ./docker/Dockerfile .
   ```

3. docker-compose启动所有容器

   

​	



## 启动说明

设置配置文件夹

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

## 相关服务

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

参考：[建库文件](https://github.com/tengfei-xy/siyuan-plugin-share-system-engine/blob/main/sql/ddl.sql)