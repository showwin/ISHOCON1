# ISHOCON1マニュアル
## 時間制限
あなたがISHOCON1に興味を持ち続けている間

## インスタンスの作成
AWSのイメージのみ作成しました。
* AMI: ami-03134665
* Instance Type: c4.xlarge
* EBS Optimization: なし
* Root Volume: 8GB, Magnetic

要望があればGCPでもイメージを作成するかもしれません。

参考画像  
* 8GB, Magnetic を選択してください。
![](https://raw.githubusercontent.com/showwin/ISHOCON1/master/doc/images/instance1.png)

* 途中で以下のように聞かれることがありますが、Magnetic を選択してください。
![](https://raw.githubusercontent.com/showwin/ISHOCON1/master/doc/images/instance2.png)

* Security Groupの設定で `TCP 22 (SSH)` と `TCP 80 (HTTP)` を `Inbound 0.0.0.0/0` からアクセスできるようにしてください。
![](https://raw.githubusercontent.com/showwin/ISHOCON1/master/doc/images/instance3.png)

## アプリケーションの起動
### インスタンスにログインする
例:
```
$ ssh -i ~/.ssh/your_private_key.pem ec2-user@xx.xx.xx.xx
```

### ishocon ユーザに切り替える
```
$ sudo su - ishocon
```

```
$ ls
 benchmark    #ベンチマーカー
 db_init.sh   #DBの初期化スクリプト(後で説明)
 init         #DB初期化に必要なスクリプト(中身は変更しないでください)
 webapp       #最適化するアプリケーション
```

### Web サーバーを立ち上げる
#### Ruby の場合
```
$ cd ~/webapp/ruby
$ unicorn -c unicorn_config.rb
```

#### Python の場合
```
$ cd ~/webapp/python
$ gunicorn -c gunicorn_config.py app:app
```

#### Go の場合
```
$ cd ~/webapp/go
$ go get -t -d -v ./...
$ go build -o webapp *.go
$ ./webapp
```

#### Crystal の場合
```
$ cd ~/webapp/crystal
$ shards install
$ crystal build app.cr
$ ./app
```

#### Scala の場合
```
$ cd ~/webapp/scala
$ sbt
> ~;jetty:stop;jetty:start
```

これでブラウザからアプリケーションが見れるようになるので、IPアドレスにアクセスしてみましょう。  

**トップページ**
![トップページ](https://raw.githubusercontent.com/showwin/ISHOCON1/master/doc/images/top.png)

`/login` からログインが可能です。
* email: ishocon@isho.con
* password: ishoconpass

**ログイン画面**
![ログイン画面](https://raw.githubusercontent.com/showwin/ISHOCON1/master/doc/images/login.png)


## ベンチマーク
スコアを計測するためのベンチマーカーの使い方の説明です。
```
$ cd ~/
$ ./benchmark --workload 3
```
* ベンチマーカーは並列実行可能で、負荷量を指定することができます。
* 何も指定しない場合は3で実行されます。
* 初期実装でスコアは100点前後になると思います。(workloadが3の場合)
* 並列度が高い場合は1分以上経っても終了しない場合がありますが、スコアには影響ありません。

## MySQL
3306 番ポートで MySQL(5.5) が起動しています。初期状態では以下のユーザが設定されています。
* ユーザ名: ishocon, パスワード: ishocon
* ユーザ名: root, パスワードなし

別のバージョンのMySQLに変更することも可能です。  
その場合、初期データの挿入は
```
$ cd
$ ./db_init.sh
```
で行うことができます。  
既存のMySQLを使う限りはこれを実行する必要はありません。

## スコアについて
* スコアはベンチマーカーが1分間の負荷走行を行っている間にレスポンスが返された
`(status code 200 * 1点) - (status code 4xx * 20) - (status code 5xx * 50)`
により算出されます。
* ただし、200の場合でもベンチマーカーが期待するレスポンスを返す必要があります。
期待しないレスポンスが返ってきた場合にはその場でベンチマーカーが停止し、スコアは表示されません。
* 1分間の負荷走行の前にベンチマーカーが `/initialize` にアクセスをして、データの初期化を行います。
初期化で1分以内にレスポンスが返らない場合には無効となりスコアは表示されません。

## 許されないこと
* インスタンスを複数台用いることや、規定のインスタンスと別のタイプを使用すること。
* ブラウザからアクセスして目視した場合に、初期実装と異なること。
  * 目視で違いが分からなければOKです。
* ベンチマーカーを改変すること。

## 許されること
* DOMを変更する
  * ベンチマーカーにバレなければDOMを変更してもOKです。
* 再起動に耐えられない
  * インスタンスを再起動して、再起動前の状態を復元できる必要はありません。

## 疑問点
[@showwin](https://twitter.com/showwin) にメンションを飛ばすか、 [issues](https://github.com/showwin/ISHOCON1/issues) に書き込んでください。
