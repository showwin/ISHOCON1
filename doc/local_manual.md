# Docker を使ってローカルで環境を整える

```
$ git clone git@github.com:showwin/ISHOCON1.git
$ cd ISHOCON1
```

言語を変更するための make command で dockerfile のパスなどを使いたい言語に合わせる。

```
$ make change-lang ISHOCON_APP_LANG=<your_language>
```

```
$ make build
$ make up
# app_1 と bench_1 のログに 'setup completed.' と出たら起動完了
```

## アプリケーション

```
$ docker exec -it ishocon1_app_1 /bin/bash
```

アプリケーションの起動は [マニュアル](https://github.com/showwin/ISHOCON2/blob/master/doc/manual.md) 参照


## ベンチマーカー

```
$ docker exec -it ishocon2-bench-1 /bin/bash
$ ./benchmark --ip app:80
```

上と同じものを実行できる benchmark 用の make command も用意されています。

```
$ make bench
```
