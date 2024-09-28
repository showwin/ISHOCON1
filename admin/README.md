# 初期データ挿入(AMIから起動した場合は不要)
```
$ mysql -u root ishocon1 < init.sql
$ ruby insert.rb
```

# 初期データ
* users: 5000件
* products: 10000件
* comments: 200000件
* histories: 500000件
