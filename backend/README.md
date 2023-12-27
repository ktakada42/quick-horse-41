# バックエンド開発ドキュメント

## ディレクトリ構造

```
.
├── Dockerfile
├── README.md
├── configs
│   └── *_config.go
├── go.mod
├── go.sum
├── [パッケージ名]
│   ├── controller.go
│   ├── repository.go
│   ├── service.go
│   └── usecase.go
│   └── *_test.go
├── main.go
├── mock
│   ├── mock_[パッケージ名]
└── utils
    └── *_utils.go
```

### ルール

- パッケージごとにディレクトリを分ける
- パッケージ内では、層ごとにファイルを分割する

## アーキテクチャ

<img width="612" alt="スクリーンショット 2023-12-27 10 31 41" src="https://github.com/ktakada42/quick-horse-41/assets/93365024/e3ddfad9-4dc1-4d58-afdf-c588c7e7e033">

レイヤードアーキテクチャを採用  
各層の役割・責任は下記の通り  
依存関係は上から並んでいる通りになっており、下の層は上の層に依存している  
上の層は下の層のことは知らない

- Entity  
  ビジネスデータを定義。ビジネスにとって不可欠な概念であり、システムからは独立している。
- Repository  
  DB とのやりとりを担当。
- Service  
  UseCase のヘルパー的な層。詳細なロジックを記述。
- UseCase  
  アプリケーション固有のビジネスルール。ユーザーの存在確認など。
- Controller  
  入出力を担当。例えばリクエストの検証や、レスポンスの作成など。

## 動作確認

ローカル環境を立ち上げている場合、http://localhost:8000/ 以下へのアクセスがバックエンドコンテナにルーティングされます。
