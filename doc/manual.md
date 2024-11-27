# ISHOCON1マニュアル
## 時間制限
あなたがISHOCON1に興味を持ち続けている間

## インスタンスの作成
AWSのイメージのみ作成しました。
* AMI: ami-073e32684d5ff32c8
* Instance Type: c7i.xlarge
* EBS-optimized: Yes
* Root Volume: 8GB, General Purpose SSD (gp3)

参考画像
* Security Groupの設定で `TCP 22 (SSH)` と `TCP 80 (HTTP)` を `Inbound 0.0.0.0/0` からアクセスできるようにしてください。
![](https://raw.githubusercontent.com/showwin/ISHOCON1/master/doc/images/instance3.png)

## アプリケーションの起動
### インスタンスにログインする
例:
```
$ ssh -i ~/.ssh/your_private_key.pem ubuntu@xx.xx.xx.xx
```

### ishocon ユーザに切り替える
```
$ sudo su - ishocon
```

```
$ ls
benchmark # ベンチマーカー
data      # DBの初期化データ (後で説明)
webapp    # 最適化するアプリケーション
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
$ go build -o webapp *.go
$ ./webapp
```

#### Scala の場合 (メンテナンス外)
```
$ cd ~/webapp/scala/ishocon1
$ sbt
> ~;jetty:stop;jetty:start
```

#### Node.js (TypeScript) の場合 (メンテナンス外)
```
$ cd ~/webapp/nodejs/
$ npm install
$ npm start
```

TypeScript で書かれています。

`~/webapp/nodejs/dist/` 以下に生成された JavaScript ファイルを直接編集する場合は、 `npm start` する前に `package.json` から TypeScript のコンパイル処理を取り除いて下さい。

```diff
diff --git a/webapp/nodejs/package.json b/webapp/nodejs/package.json
index eb5e600..600ba5e 100644
--- a/webapp/nodejs/package.json
+++ b/webapp/nodejs/package.json
@@ -7,7 +7,7 @@
   "scripts": {
     "clean": "rm -rf ./dist/*",
     "test": "echo \"Error: no test specified\" && exit 1",
-    "start": "tsc && node dist/index.js",
+    "start": "node dist/index.js",
     "build": "tsc"
   },
   "keywords": [
```


#### Crystal の場合 (メンテナンス外)
```
$ cd ~/webapp/crystal
$ shards install
$ crystal build app.cr
$ ./app
```

これでブラウザからアプリケーションが見れるようになるので、 `http://<IPアドレス>` にアクセスしてみましょう。

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
* 初期実装でスコアは400点前後になると思います。(workloadが3の場合)
* 並列度が高い場合は1分以上経っても終了しない場合がありますが、スコアには影響ありません。

## MySQL
3306 番ポートで MySQL(8.0) が起動しています。初期状態では以下のユーザが設定されています。
* ユーザ名: ishocon, パスワード: ishocon
* ユーザ名: root, パスワード: ishocon1

別のバージョンのMySQLに変更することも可能です。
その場合、初期データの挿入は
```
$ cd
$ sudo mysql -u root -pishocon1 ishocon1 < ~/data/ishocon1.dump
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
