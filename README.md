进入prometheus_server_tools文件夹
## 代码结构
main包位于prometheus_server当中，主要流程为初始化controller以后定时器触发聚合，多协程并发采集数据后同步聚合，写入influxdb数据库。

主要的工具包括读取、提取prometheus数据、生成同步时间戳、读写influxdb等请参考prometheus_tools包


对于多指标聚合controller位于prometheus_controller当中。

## 后续完善
下一版本将仿照当前结构添加node相关指标聚合。

## 环境要求
golang 1.14及以上版本下进行开发
