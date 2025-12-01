## 类型映射自定义​
为了兼容 postgres 的 uuid 类型，我们自定义了 uuid 类型。

首先需要生成goctl model的工程化配置文件，在项目目录下执行： `goctl config init goctl.yaml generated in ./goctl.yaml`

在goctl.yaml中添加如下配置：
```yaml
    uuid:
      null_type: sql.NullString
      type: string
```

## 执行生成数据库模板文件
在项目目录下执行： `goctl model pg datasource -url="postgresql://admin:xledger123@localhost:15432/xledger" -table="user" -dir="./service/user/model"`

执行生成带缓存的model文件： `goctl model pg datasource -url="postgresql://admin:xledger123@localhost:15432/xledger" -table="user" -dir="./service/user/model" -cache=true`


简单可以执行 `make model-table TABLE=user` 来生成 user 表的 model