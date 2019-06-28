# kenchan

KEN_ALL.CSVを郵便番号で検索するAPIサーバ.

## インストール

とりあえず試すには, Releaseページからバイナリをダウンロードし, ターミナルやコマンドプロンプトで以下のように実行してください.

実行すると日本郵政のサイトからKEN_ALL.zipをhttpsで取得して動作します.

```
# linux/macの場合
./kenchan -a "*"

# windowsの場合
kenchan.exe -a "*"
```

sampleディレクトリ内にあるsample.htmlを開くとAPIの動作サンプルページが動きます.

## API

文字コードはUTF-8です.

### リクエスト

以下のURLへGETリクエストを送ってください.

postcodeは郵便番号で前方一致します.  
ハイフンはあってもなくても構いません.  
レスポンス最大件数は起動時のオプションで指定した数です.

```
http://[hostname]/search?q=[zipcode]
```

### レスポンス

レスポンスは下記のフィールドの配列です.

| フィールド名     | 項目名           |
| ---------------- | ---------------- |
| zipcode          | 郵便番号         |
| address1         | 都道府県名       |
| address2         | 市区町村名       |
| address3         | 町域名           |
| kana1            | 都道府県名(カナ) |
| kana2            | 市区町村名(カナ) |
| kana3            | 町域名(カナ)     |

## オプションなど

ヘルプをどうぞ.

```
kenchan -h
Usage of kenchan:
  -a string
        レスポンスに付与するAccess-Control-Allow-Originヘッダを指定します
  -f string
        利用するcsvファイルパスを指定します (default "KEN_ALL.CSV")
  -l int
        検索結果件数のリミットを指定します (default 20)
  -p int
        受付ポートを指定します (default 8080)
  -r    csvファイルを更新する際に指定します
  -u string
        csvファイル取得に利用するURLを指定します (default "https://www.post.japanpost.jp/zipcode/dl/kogaki/zip/ken_all.zip")
```

