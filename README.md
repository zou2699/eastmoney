# eastmoney
获取今日的基金估值

## Usage

```bash
git clone https://github.com/zou2699/eastmoney.git

cd eastmoney

go build -o eastmoney  main.go

./eastmoney
```

```bash
# 获取config.yaml下配置的基金代码
curl 127.0.0.1:8080

# 通过天天基金的编号id直接进行查询
curl 127.0.0.1:8080/110022
```

当前字段介绍
```
Name   基金名字
Id     基金ID
Gzzf   估值涨幅
GzDate 估值时间
```