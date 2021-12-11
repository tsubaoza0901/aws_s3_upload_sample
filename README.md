# aws_s3_upload_sample

# 使用方法

## 1．Docker イメージのビルド&コンテナの起動

```
$ docker-compose up -d --build
```

## 2．データベースの作成

① DB コンテナ内へ移動

```
$ docker exec -it aws-s3-upload-sample-db bash
```

② DB 接続

```
root@ec19d85976f4:/# mysql -u root -h db -p
Enter password:
```

③ DB 作成

```
mysql> CREATE DATABASE awss3uploadsample;
```

## 3．マイグレーションファイルの実行

① アプリケーションコンテナ内へ移動

```
$ docker exec -it aws-s3-upload-sample bash
```

② マイグレーションファイルの実行

```
root@fe385569a625:/go/src/app/server_side# goose up
```

## 4．アプリケーションの起動

```
root@fe385569a625:/go/src/app/server_side# go run main.go
```

## 5.使用可能なエンドポイント

### ① S3への画像保存
・Method：POST   
・Endpoint：http://127.0.0.1:9111/image   
・Request Body：

```
{
    "encoded_url": "Please set encoded Image URL",
    "file_name": "Please set Image FileName"
}
```

### ② S3に保存した画像の
・Method：GET   
・Endpoint：http://127.0.0.1:9111/image/:id   

# その他
使用する際にはルート直下に「.env」ファイルを作成し、以下の内容を追加する必要あり。

```
BUCKET_NAME=
AWS_ACCESS_KEY_ID=
AWS_SECRET_ACCESS_KEY=
```