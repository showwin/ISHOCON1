# ISHOCON1
## ISHOCONとは
ISHOCONとは `Iikanjina SHOwwin CONtest` の略で、[ISUCON](http://isucon.net/)と同じように与えられたアプリケーションの高速化を競うコンテスト(?)です。  
ISUCON5でISUCONに初参戦したのですが、予選を通過できなくて悔しかったので自分の勉強のためにも、問題を作成してみました。  

## 問題概要
今回のテーマは「爆買いECサイト」です。  
近年某国の方々が日本に来て、爆買いをしているようなので、ECサイトもそれに負けずより多くのレスポンスを返せるようにチューニングしましょう。  
![](https://raw.githubusercontent.com/showwin/ISHOCON1/master/doc/images/top.png)

## 問題詳細
* マニュアル: [ISHOCON1マニュアル](https://github.com/showwin/ISHOCON1/blob/master/doc/manual.md)
* AMI: `ami-dd727fb3`
* インスタンスタイプ: `c3.xlarge`
* 参考実装言語: Ruby, Go

## ランキング(v0.3)
ここに載りたい人はどんどんプルリク送ってください！大歓迎です！！

|名前|スコア|ソースコードのレポジトリ or AMI|備考|
|:--:|:--:|:--:|:--:|
| [showwin](https://twitter.com/showwin) | 118918 | [showwin/ISHOCON1_me](https://github.com/showwin/ISHOCON1_me/commits/golang) | 2016/03/05 Go で書きなおした！ |
|[showwin](https://twitter.com/showwin)|78699|[showwin/ISHOCON1_me](https://github.com/showwin/ISHOCON1_me/tree/dc7273d6cfdc90edb43d9490e7538ec63f06a99e)|2015/11/12のぼくではこれが限界だ (Ruby)|
|[showwin](https://twitter.com/showwin)|48642|[showwin/ISHOCON1_me](https://github.com/showwin/ISHOCON1_me/tree/c04b20faee30a5aeb315c33bee6ad8b4c7d87ce7)|2015/11/11のぼくではこれが限界だ (Ruby)|
|初期状態|193|-|初期実装(Go)のworkloadが3で193点になりました。|
|初期状態|104|-|初期実装(Ruby)のworkloadが3で104点になりました。|
