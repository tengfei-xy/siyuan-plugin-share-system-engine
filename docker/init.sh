#!/bin/bash
function main() {
    read_yaml '.services.spss_mysql.environment.MYSQL_ROOT_PASSWORD'
    root_passwd=${RET}

    read_yaml '.services.spss_mysql.environment.MYSQL_USER'
    user=${RET}

    read_yaml '.services.spss_mysql.environment.MYSQL_PASS'
    passwd=${RET}

    read_yaml '.services.spss_mysql.environment.MYSQL_DATABASE'
    database=${RET}

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
    ret=$(yq  "$1" ./docker-compose.yml )
    RET=${ret#*\"}
    RET=${RET%*\"}
    RET=${RET%*\'}
    RET=${RET#*\'}

}

# 当容器不存在时，直接退出
function check_docker(){
    sudo docker ps  | grep "spss_mysql" > /dev/null 2>&1 || { echo "mysql容器不存在" ; exit 255 ; }
}

# 安装yq命令，用于解析yaml文件
function check_yq(){
    os=$(uname)
    case $os in
    # macOS基本命令检测
    Darwin)
        which yq >/dev/null 2>&1 || {
            echo "准备安装jq命令..."
            brew install yq || {
                error "brew install yq 执行失败"
                exit 255
            }
        }
        return
        ;;
    Linux)
        # Centos 基本命令检测
        test -r /etc/redhat-release && grep "CentOS" /etc/redhat-release >/dev/null 2>&1 && {

            which yq >/dev/null 2>&1 || {
                echo "准备安装yq命令..."
                sudo yum -y install yq || {
                    error "sudo yum -y install yq 执行失败"
                    exit 255
                }
            }
            return
        }
        # Ubuntu 基本命令检测
        lsb_release -a 2>/dev/null | grep "Ubuntu" >/dev/null 2>&1 && {
            
            which yq >/dev/null 2>&1 || {
                echo "准备安装yq命令..."
                sudo apt -y install yq || {
                    error "sudo apt -y install yq 执行失败"
                    exit 255
                }
            }
            return
        }
        # Debian 基本命令检测
        lsb_release -a 2>/dev/null | grep "Dibian" >/dev/null 2>&1 && {
            
            which yq >/dev/null 2>&1 || {
                echo "准备安装yq命令..."
                sudo apt -y install yq || {
                    error "sudo apt -y install yq 执行失败"
                    exit 255
                }
            }
            return
        }
        ;;

    esac

}
check_docker
check_yq
main