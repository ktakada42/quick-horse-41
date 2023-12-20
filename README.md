# 技術書管理アプリ

CA Tech Lounge で開発中の技術書管理アプリのリポジトリ

## 環境構築

- Docker コンテナで開発環境を構築します。
  デフォルトでは下記ポートを使います。

  - 80: nginx
  - 8000: バックエンド
  - 3306: DB

- DB との接続やアプリケーションで環境変数が必要となります。  
  プロジェクトのルートディレクトリに、`.env`ファイルを配置してください。  
  CA Tech Lounge 生は[こちら](https://drive.google.com/file/d/1eULGiXgCy3o73pstTGt08z_jTjCnQaKK/view?usp=drive_link)から`.env`ファイルをダウンロードできます。

- 上記手順が終わったら、下記コマンドを実行することでサーバーが起動します。

```
docker compose -f compose-local.yaml up -d
```
