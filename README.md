# ISHOCON1
<img src="https://user-images.githubusercontent.com/1732016/41643273-b4994c02-74a5-11e8-950d-3a1c1e54f44f.png" width="250px">

© [Chie Hayashi](https://www.facebook.com/hayashichie)

## ISHOCONとは
ISHOCONとは `Iikanjina SHOwwin CONtest` の略で、[ISUCON](http://isucon.net/)と同じように与えられたアプリケーションの高速化を競うコンテスト(?)です。
ISUCON5でISUCONに初参戦したのですが、予選を通過できなくて悔しかったので自分の勉強のためにも、問題を作成してみました。

ISUCONは3人チームで取り組むことを基準に課題が作られていますが、ISHOCONは1人で8時間かけて解くことを基準に難易度を設定しています。

## 問題概要
今回のテーマは「爆買いECサイト」です。

近年某国の方々が日本に来て、爆買いをしているようなので、ECサイトもそれに負けずより多くのレスポンスを返せるようにチューニングしましょう。
![](https://raw.githubusercontent.com/showwin/ISHOCON1/master/doc/images/top.png)

## 問題詳細
* マニュアル: [ISHOCON1マニュアル](https://github.com/showwin/ISHOCON1/blob/master/doc/manual.md)
* AMI: `ami-06cda439fc5c0da1b`
* インスタンスタイプ: `c5.xlarge`
* 参考実装言語: Ruby, Go, Python
  * メンテナンス外: Node.js(TypeScript), Crystal(by [@Goryudyuma](https://github.com/Goryudyuma)), Scala(by [@Goryudyuma](https://github.com/Goryudyuma))
* 推奨実施時間: 1人で8時間

## 社内ISUCON等のイベントで使用したい方
自由に使って頂いて構いません。

イベント実施後にブログを書いて [@showwin](https://twitter.com/showwin) まで連絡頂けたら嬉しいです！下の関連リンクに掲載いたします。

サーバーの準備には terraform を使うと便利です。詳しくは [terraform の README](https://github.com/showwin/ISHOCON1/blob/master/terraform/README.md) を参照してください。

## 関連リンク

* [ISUCON5予選と自作ISUCON](http://blog.mmmcorp.co.jp/blog/2015/10/06/isucon5_and_ishocon/) (by [@showwin](https://twitter.com/showwin))
  * ISHOCON作るキッカケになった話。
* [社内ISUCONを開催しました](http://blog.mmmcorp.co.jp/blog/2016/09/01/ishocon_2016/) (by [@showwin](https://twitter.com/showwin))
  * 株式会社MMMでISHOCON1をやった話。
* [社内ISUCONに参加した。](http://yasun.hatenablog.jp/entry/2016/08/31/211927) (by [@yasun](https://twitter.com/_Yasuun_))
  * 上の社内ISHOCONの参加レビュー。
* [ISHOCON1 反省会](https://speakerdeck.com/showwin/ishocon1-fan-sheng-hui) (by [@showwin](https://twitter.com/showwin))
  * Rubyで6万点取るぐらいまでの解説。([@showwin](https://twitter.com/showwin) の技術レベルが低い頃のスライドなのであまり価値ない)
* [ISHOCON1 〜個人参加のISUCON練習コンテスト〜](https://scouty.connpass.com/event/65322/) (by [@showwin](https://twitter.com/showwin))
  * ISUCON7のフォーミングアップとして個人でISHOCONを解くイベントをやりました。
* [ISUCON勉強会 ISHOCON1 を開催しました](https://www.wantedly.com/companies/scouty/post_articles/79778) (by [@showwin](https://twitter.com/showwin))
  * ISHOCON1の開催レポート。
* [ISHOCON1に参加した #scouty_ishocon](http://utgwkk.hateblo.jp/entry/2017/10/07/214659) (by [@utgwkk](https://twitter.com/utgwkk))
  * 上のイベントの参加レビュー。
* [ISUCON模試を開催して運営&参加してきた](http://saboyutaka.hatenablog.com/entry/2017/10/09/003257)
  * ISUCON模擬試験のイベントでISHOCON1を問題として採用して頂きました。
* [ISHOCON1をCrystalで書いたお気持ち](http://goryudyuma.hatenablog.jp/entry/2018/03/14/174935) (by [@Goryudyuma](https://twitter.com/Goryudyuma))
  * Crystal実装を追加してくださったお気持ちブログです。必読です！
* [今年もWantedlyの新卒研修で社内ISUCONを行いました！](https://www.wantedly.com/companies/wantedly/post_articles/117958) (by [@kobayang](https://github.com/kobayang))
  * Wantedlyさんの社内ISUCONでISHOCON1を使って頂きました。ベンチマーカーを外に出すのと複数台構成にすると戦略がかなり変わりますね。
* [ISHOCON1をScalaで書いたお気持ち](https://goryudyuma.hatenablog.jp/entry/2018/06/11/170711)(by [@Goryudyuma](https://twitter.com/Goryudyuma))
  * Scala実装を追加頂いた時の記事です。 [@Goryudyuma](https://twitter.com/Goryudyuma)さんの記事は毎回気付かされる部分があり、ありがたいです。
* [エンジニア新人研修の一環で株式会社はてな社内ISUCONを開催しました](https://developer.hatenastaff.com/entry/2021/03/12/103000)(by [@astj](https://github.com/astj))
  * はてなさんの新人研修でISHOCON1を使って頂きました。コンテスト後にきちんと振り返り時間を設けている点が良いですね！
* [『みんなで解く ISUCON勉強会』を開催しました！](https://zenn.dev/lovegraph/articles/e4e120b6d204fb)(by [@yokoe24](https://twitter.com/yokoe24))
  * ラブグラフさん主催のイベントでISHOCON1を使って頂きました。ISUCON未経験の方たちも楽しめたようで嬉しいです。


## ISHOCONシリーズ
* [ISHOCON1](https://github.com/showwin/ISHOCON1)
* [ISHOCON2](https://github.com/showwin/ISHOCON2)
