#!/bin/bash
function main() {
    root_passwd=$(read_yaml 'MYSQL_ROOT_PASSWORD')

    user=$(read_yaml 'MYSQL_USER')

    passwd=$(read_yaml 'MYSQL_PASS')

    database=$(read_yaml 'MYSQL_DATABASE')

    # 创建用户并赋权
    sudo docker exec -it spss_mysql sh -c  "mysql -u root -p${root_passwd} << EOF
grant ALL PRIVILEGES on ${database}.* to '${user}'@'%' identified by '${passwd}';
flush privileges;
EOF" > /dev/null 2>&1

    sudo docker cp ./mysql/ddl.sql spss_mysql:/tmp/ddl.sql > /dev/null 2>&1
    # 初始化库
    sudo docker exec -it spss_mysql sh -c "mysql -u root -p${root_passwd} < /tmp/ddl.sql" > /dev/null 2>&1

    # 检查是否初始化完成
    sudo docker exec -it spss_mysql sh -c  "mysql -u root -p${root_passwd} << EOF
show databases; 
EOF" | grep "${database}" > /dev/null 2>&1 && echo "初始化完成"
}
function read_yaml() {
    local ret
    ret=$(grep "\ \{0,\}#" -v  ./docker-compose.yml | grep "$1" | awk  -F ':' '{print $2}')
    ret=${ret#*\"}
    ret=${ret%*\"}
    ret=${ret%*\'}
    ret=${ret#*\'}
    echo "$ret"
}

# 当容器不存在时，直接退出
function check_docker(){
    sudo docker ps  | grep "spss_mysql" > /dev/null 2>&1 || { echo "mysql容器不存在" ; exit 255 ; }
}

check_docker
main