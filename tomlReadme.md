# 使用toml配置文件连接数据库

## toml配置文件的内容

```toml
[dbservers.test]
host = "127.0.0.1"
port = 5432
dbname = "test"
user = "loginuser"
password = "123456"

[dbservers.dborm]
host = "127.0.0.1"
port = 5432
dbname = "test"
user = "loginuser"
password = "123456"

[dbservers.dbsqlx]
host = "127.0.0.1"
port = 5432
dbname = "test"
user = "loginuser"
password = "123456"

[redisservers.redis]
addr = "127.0.0.1:6379"
password = ""
db = 10
```

举个例子 dbservers.test。对于dbservers这个是下面配置信息中代表Config结构体中的DBServers部分。而.后面代表的则是map中的key。



## config数据库的配置信息

首先toml对应的是一个Config的结构体、以下是结构体的代码

```go
type Config struct {
   DBServers    map[string]DBServer    `toml:"dbservers"`
   RedisServers map[string]RedisServer `toml:"redisservers"`
}


// DBServer 表示DB服务器配置
type DBServer struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	DBName   string `toml:"dbname"`
	User     string `toml:"user"`
	Password string `toml:"password"`
}

// RedisServer 表示 redis 服务器配置
type RedisServer struct {
	Addr     string `toml:"addr"`
	Password string `toml:"password"`
	DB       int    `toml:"db"`
}
```



## 解析配置文件

```go
// New 解析toml配置
func New(tomlFile string) (*Config, error) {
   c := &Config{}
   if _, err := toml.DecodeFile(tomlFile, c); err != nil {
      return c, err
   }
   return c, nil
}
```

