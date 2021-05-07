# terraform

## 概要

terraformでISHOCON1の環境を構築できます。

## 必要なもの

- 支払情報が紐付けされたAWSアカウント
- IAMアカウントのアクセスIDとシークレットアクセスキー

## 手順

### AWSアカウントの登録

```shell
$ aws configure --profile ishocon1
```

### admins, playersの皆さんに、SSHの鍵を登録してもらう

```shell
$ ssh -T git@github.com
```
で、自分のGitHubのアカウントIDが表示されることを各自に確認してもらう

### stateファイルを入れるS3を作成

### FIXMEを直す

- [ ] terraform.tfの中のbucketに、先ほど作成したS3のbucket nameを入れる
  - ローカルでtfstateを管理し、S3を使わない場合、terraform.tfを削除する
- [ ] users.tfにadmins, playersを入れる
  - SSHの鍵が登録されたgithubアカウントIDを入れること
    - [ドキュメント](https://docs.github.com/ja/github/authenticating-to-github/connecting-to-github-with-ssh)を参照のこと
  - adminsは全インスタンスに入れるユーザー
    
### terraform apply

- outputにIPアドレスが入っているので、playerの人に渡す

## 注意点

- spot instanceで立てるようにしてコストの削減を図っています
  - main.tfの `aws_spot_instance_request` を `aws_instance` に変えた上で、適宜修正をお願いします
- デフォルトでは、 `c5.xlarge` でインスタンスが立ちます
  - 大きすぎる、小さすぎるなどがあれば、適宜直してください
- playerは増やしたり減らしたりできます
  - adminは増やすことも減らすこともできません
- 動かなくても保証はしませんし、使用したことによる責任は一切取りません
  - VPCから作るので、安全だとは思いますが、コードをよく読んだ上で、自己責任でご使用ください
- PRは大歓迎です