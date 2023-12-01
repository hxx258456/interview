## mysql常见数据类型

TINYINT、SMALLINT、MEDIUMINT、INT、BIGINT

VARCHAR、CHAR、TEXT、BLOB

DATETIME、DATE 和 TIMESTAMP

## CHAR，VARCHAR 和 Text 的区别？

**长度区别**

- Char 范围是 0～255。
- Varchar 最长是 64k（注意这里的 64k 是整个 row 的长度，要考虑到其它的 column，还有如果存在 not null 的时候也会占用一位，对不同的字符集，有效长度还不一样，比如 utf-8 的，最多 21845，还要除去别的column），但 Varchar 在一般情况下存储都够用了。
- 如果遇到了大文本，考虑使用 Text，最大能到 4G（其中 TEXT 长度 65,535 bytes，约 64kb；MEDIUMTEXT 长度 16,777,215 bytes，约 16 Mb；而 LONGTEXT 长度 4,294,967,295 bytes，约 4Gb）。

**默认值区别**

​	Char 和 Varchar 支持设置默认值，而 Text 不能指定默认值

### 什么是三大范式？

1. 字段唯一
2. 确保表中的每列都和主键相关
3. 非主键列只依赖于主键，不依赖于其他非主键

## mysql有关权限的表有哪几个

user,db,table_priv,columns_priv,host

user:全局权限

db:库级权限

## mysql binlog 有几种录入格式?

- statement，每一条修改数据的sql都会记录在binlog中，不需要记录每一行的变化，减少了binlog日志量，节约了io，提升了性能
- row，不记录sql语句上下文信息，仅保存哪条记录被修改，记录单元为每一行的改动，基本是可以全部记录下来的，但是由于很多操作，日志文件会很大
- mixed，一种折中方案，普通操作使用statement记录，无法使用statement时使用row记录

## mysql binlog

二进制日志,记录ddl(data definition language)和dml(data manipulation language)

常用选项--start-datetime,--stop-datetime,--start-position,--stop-position

## mysql 存储引擎myisam与innodb的区别

- 锁粒度方面: innodb 行级锁，myisam 表级锁
- innodb支持事务,myisam不支持事务
- myisam 数据文件一.myd结尾，innodb数据文件以.idb结尾

## myisam 索引与innodb索引的区别

- innodb索引是聚簇索引，myisam索引是非聚簇索引
- innodb的主键索引的叶子节点存储着行数据，因此主键索引非常高效
- myisam索引的叶子节点存储的是行数据的地址，需要在寻址一次才能得到数据
- innodb非主键索引的叶子节点存储的是主键和其他带索引的列数据，因此查询时会非常高效

## 什么是索引

