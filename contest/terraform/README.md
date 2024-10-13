# terraform

## 概要

terraformでISHOCON1の環境が構築できます。  
社内ISUCON等のコンテストの開催準備として使うと便利です。  
開催者(以下admins)と参加者(以下players)のロールに分けて、サーバーを複数台準備できます。

## 必要なもの

- 支払情報が紐付けされたAWSアカウント
- IAMユーザのアクセスIDとシークレットアクセスキー

## 手順

### 1. AWSアカウントの登録

```shell
$ aws configure --profile ishocon1
```

### 2. adminsとplayersに、GitHubにて秘密鍵を登録してもらう

playerはコンテストで使用するインスタンスにログインするために、GitHubに登録されている秘密鍵を使用する。  
GitHubにて秘密鍵を登録後、各自のPCで以下のコマンドを実行し自身のGitHubのアカウントIDが表示されることを確認する。

```shell
$ ssh -T git@github.com
```

### 3. stateファイルを入れるS3バケットを作成
tfstateの管理をローカルで行う場合には、作成不要。

### 4. FIXMEを直す

- [ ] terraform.tfの中のbucketに、先ほど作成したS3のbucket nameを入れる
  - ローカルでtfstateを管理し、S3を使わない場合にはterraform.tfを削除する
- [ ] users.tfにadmins, playersのGitHubアカウントIDを入れる
  - 詳細は[ドキュメント](https://docs.github.com/ja/github/authenticating-to-github/connecting-to-github-with-ssh)を参照のこと
  - adminsは全インスタンスに入ることができる
    
### 5. terraform apply してリソースの作成

```shell
$ terraform apply
```

outputに競技で使用するインスタンスのIPアドレスが入っているので、playerの人に渡す。

## 注意点

- spot instanceで立てるようにしてコストの削減を図っています
  - コンテスト中に絶対に落ちて欲しくないなどの理由により spot instanceを使わない場合、main.tfの `aws_spot_instance_request` を `aws_instance` に変えてください
- デフォルトでは、ISHOCON1の推奨インスタンスタイプである `c5.xlarge` でインスタンスが起動します
- このterraformを使用することで発生した問題に対して責任は一切取りません
  - コードをよく読んだ上で、自己責任でご使用ください
