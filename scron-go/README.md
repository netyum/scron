scron-go
======
### 编译

glide 安装 已经安装，就不许安装了，并运行 

```
glide install
```

解决依赖

```
./build.sh
```

本程序会使用upx压缩  除upx为linux_64位，其他平台相应处理

### 运行

本程序不包含Web操作，可自行设计Web程序，操作数据库即可。

*默认会查找 /etc/scron.conf文件和/etc/scron_warning.conf*
```
./scron
```

*指定scron.conf文件*
```
./scron -c scron.conf -w scron_warning
```

*指定部署环境*
```
./scron -c scron.conf -e dev
```

### 部署

*部署在类(nix)在机器的crontab中
```
*/1 * * * * ./scron
```

### scron.conf

```
[pro]
db_host = 192.168.0.2
db_name = test
db_port = 3306
db_user = test
db_pass = test
log_path = /opt/logs
cron_table_name = sys_crontab
host_table_name = sys_crontab_host
iskill = 1

[test]
db_host = 10.241.221.106
db_name = test
db_port = 3306
db_user = sp_sys
db_pass = test
log_path = /data/logs
cron_table_name = sys_crontab
host_table_name = sys_crontab_host
iskill = 1

[dev]
db_host = 172.16.10.44
db_name = db_sys
db_port = 3306
db_user = write
db_pass = write
log_path = /data/logs
cron_table_name = sys_crontab
host_table_name = sys_crontab_host
iskill = 1
```

### scron_waring.conf
```
[dev]
phones = 1234567890  123456789
warning_script =

[test]
phones = 1234567890  123456789
warning_script =

[pro]
phones = 1234567890  123456789
warning_script =
```

cron_table_name是任务表名，host_table_name是主机表名, iskill=1表示超时后，会自动杀掉，不杀掉请设置成0


### 报警机制

在scron_warning设置报警脚本
```
warning_script "msg" "手机号"

```

warning_script 需要支持上面的格式
