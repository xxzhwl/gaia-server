configs文件夹下主要包含三个文件夹
1. 通用型配置文件 /configs/common/operate;/configs/common/search
2. 本地服务启动配置文件 /configs/local/config.json
3. 远程配置中心配置缓存

---
1. 通用型配置文件分为两类
    - 通用操作
    - 通用查询
2. 本地服务启动配置文件包含了服务启动的关键性配置
```json
    {
  "SystemEnName": "GoBack",//服务英文名称
  "SystemCnName": "GoBack",//服务中文名称
  "SuperAdmin": "root",//超级管理员设置(通常情况下不需要设置)
  "CronJob": {
    "Mysql": "DSN：charset=utf8mb4&parseTime=True&loc=Local"
  },//不同的服务可以作为schema在外层，内部书写相关配置，比如mysql等
  "Framework": {
    "JaegerTracePoint": "xxxxxxx:4318",
    "Mysql": "DSN:charset=utf8mb4&parseTime=True&loc=Local",
    "ES": {
      "Address": "http://xxxxxx:7901/",
      "UserName": "es",
      "Password": "es"
    },
    "Redis": {
      "Address": "xxxxxxx:6379",
      "UserName": null,
      "Password": "xxxxx"
    }
  },//服务默认的schema是Framework
  "JwtConf": {
    "DurationMinute": 7200
  },//服务中Jwt相关配置，也是可以使用schema配置使用
  "Server": {
    "Port": "8008",
    "Cors": {
      "Enable": false,
      "AllowOrigins": ["http://localhost:3002"],
      "AllowMethods": ["POST", "GET","PUT","DELETE"],
      "AllowHeaders": ["Content-Type", "Authorization"],
      "MaxAge": 12
    },
    "Logger": {
      "PrintConsole": true,
      "DetailMode": false,
      "EnablePushLog": true
    },
    "CommonHandler": {
      "Enable": true
    }
  },//当前服务的启动端口、跨域配置、日志配置、通用查询和操作是否启动配置
  "AsyncTask": {
    "Mysql": "DSN:charset=utf8mb4&parseTime=True&loc=Local",
    "Port": "8080",
    "Cors": {
      "Enable": true,
      "AllowOrigins": ["http://localhost:3002"],
      "AllowMethods": ["POST", "GET","PUT","DELETE"],
      "AllowHeaders": ["Content-Type", "Authorization"],
      "MaxAge": 12
    },
    "Logger": {
      "PrintConsole": false,
      "DetailMode": false,
      "EnablePushLog": true
    },
    "CommonHandler": {
      "Enable": false
    }
  },
  "Gorm": {
    "LogLevel": "Warn"
  },//GORM日志登记配置
  "Message": {
    "FeiShuRobot": "xxx"
  }//消息告警配置
}
```