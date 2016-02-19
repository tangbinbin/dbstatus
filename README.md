dbstatus

------------

### 功能介绍
- 输出一个或多个mysql的qps，conn，innodb rows reads delete 等信心，每秒输出一次
- golang1.5

### 使用说明
- git clone https://github.com/tangbinbin/dbstatus.git
- cd dbstatus;make
- ./bin/dbstatus 为可执行文件，可拷贝到可连接数据库的地方

### 使用举例
    ./dbstatus -h
    -h string
        hosts,多个地址之间,分割 (default "127.0.0.1:3306")
    -p string
        password (default "test")
    -u string
        user (default "test")


    单个地址使用示例
    ./dbstatus -u=monitor -p=w123456 -h=127.0.0.1:3306
    - 预期结果
        _____________________________________________________________________________________________________________
                                  |            --QPS--            | --Innodb Rows Status-- |     --Thead--| --bytes--
              addr          time  |  ins   upd   del    sel    qps|  ins   upd   del   read| run  con  cre| recv send
           127.0.0.1:3306 18:47:00|    0     0     0     20    386|    0     0     0      0|   1    1    2|   0k   ok
           127.0.0.1:3306 18:47:01|    0     0     0      0      1|    0     0     0      0|   1    1    0|   0k   ok
           127.0.0.1:3306 18:47:02|    0     0     0      0      1|    0     0     0      0|   1    1    0|   0k   ok
           127.0.0.1:3306 18:47:03|    0     0     0      0      1|    0     0     0      0|   1    1    0|   0k   ok
           127.0.0.1:3306 18:47:04|    0     0     0      0      1|    0     0     0      0|   1    1    0|   0k   ok

    多个地址使用示例
    ./dbstatus -u=monitor -p=w123456 -h=127.0.0.1:3306,localhost:3306
    - 预期结果
        _____________________________________________________________________________________________________________
                                  |            --QPS--            | --Innodb Rows Status-- |     --Thead--| --bytes--
              addr          time  |  ins   upd   del    sel    qps|  ins   upd   del   read| run  con  cre| recv send
           localhost:3306 18:49:28|    0     0     0     22    397|    0     0     0      0|   1    2    2|   0k   0k
           127.0.0.1:3306 18:49:28|    0     0     0     22    398|    0     0     0      0|   1    2    2|   0k   0k
    
           localhost:3306 18:49:29|    0     0     0      0      2|    0     0     0      0|   1    2    0|   0k   0k
           127.0.0.1:3306 18:49:29|    0     0     0      0      2|    0     0     0      0|   1    2    0|   0k   0k
    
           127.0.0.1:3306 18:49:30|    0     0     0      0      1|    0     0     0      0|   1    2    0|   0k   0k
           localhost:3306 18:49:30|    0     0     0      0      3|    0     0     0      0|   1    2    0|   0k   0k
    
           127.0.0.1:3306 18:49:31|    0     0     0      0      2|    0     0     0      0|   1    2    0|   0k   0k
           localhost:3306 18:49:31|    0     0     0      0      2|    0     0     0      0|   1    2    0|   0k   0k
    
           127.0.0.1:3306 18:49:32|    0     0     0      0      2|    0     0     0      0|   1    2    0|   0k   0k
           localhost:3306 18:49:32|    0     0     0      0      2|    0     0     0      0|   1    2    0|   0k   0k
