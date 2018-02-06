scron
======

## scron定时任务程序

### 功能

1.一次部署，通过web添加任务
2.日志输出，
3.超时处理(报警，杀掉)

### 依赖

mysql

### 版本
scron-go为golang版
scron-py为python版

### 注意
程序会能过查看本机hostname获取本机ip，如果本机hostname没有对应ip，程序可能会出错或在数据库中找不到对应的任务