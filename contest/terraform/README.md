# Terraform

## 概要

TerraformでISHOCON1の環境が構築できます。
社内ISUCON等のコンテストの開催準備として使うと便利です。

## 作成されるもの
* チーム毎に1台の競技に使用するEC2インスタンス
* 各チームのスコアがリアルタイムに見れるポータルサイト

## 必要なもの

- 支払情報が紐付けされたAWSアカウント
  - ポータルサイトは低コストになるように実装していますが、主に以下の費用がかかります。
    - DynamoDB, Lambda, CloudWatch, API Gateway, S3
- IAMユーザのアクセスIDとシークレットアクセスキー

## 手順

### 1. AWSアカウントの登録

```shell
$ aws configure
```

### 2. adminsとplayersに、GitHubにて秘密鍵を登録してもらう

playerはコンテストで使用するインスタンスにログインするために、GitHubに登録されている秘密鍵を使用する。
GitHubにて秘密鍵を登録後、各自のPCで以下のコマンドを実行し自身のGitHubのアカウントIDが表示されることを確認する。

```shell
$ ssh -T git@github.com
```

### 3. main/main.tf を編集

- [ ] (必須) `admins`と`teams`(チーム名と参加者名)に各々のGitHubアカウントIDを記述
  - 詳細は[ドキュメント](https://docs.github.com/ja/github/authenticating-to-github/connecting-to-github-with-ssh)を参照のこと
  - adminsは全インスタンスに入ることができる
- [ ] (任意) module/variables.tf を参考に変更したいパラメーターを main/main.tf に定義

### 4. terraform apply してリソースの作成

```shell
$ cd main
$ terraform apply

...

Outputs:

apigateway_url = "https://123456789.execute-api.ap-northeast-1.amazonaws.com/"
ip_addr = {
  "team1" = "11.111.111.111"
  "team2" = "11.111.111.112"
}
portal_url = "http://ishocon1-portal123456789.s3-website-ap-northeast-1.amazonaws.com"
```

Outputに競技で使用するインスタンスのIPアドレスとポータルサイトのアドレスが出力されるので、参加者に共有する。

```
$ ssh ishocon@<instance_ip>
```

でインスタンスにSSHできる。

## 注意点

- デフォルトでは、ISHOCON1の推奨インスタンスタイプである `c7i.xlarge` でインスタンスが起動します
- このTerraformを使用することで発生した問題に対して責任は一切取りません
  - コードをよく読んだ上で、自己責任でご使用ください
