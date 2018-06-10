# ishocon1 #

## Build & Run ##

```sh
$ cd ishocon1
$ sbt
> jetty:start
> browse
```

If `browse` doesn't launch your browser, manually open [http://localhost:8080/](http://localhost:8080/) in your browser.

### Memo

```
$ sbt
> ~;jetty:stop;jetty:start
```

ファイル変更を監視して、更新されればサーバーが再起動する設定。
