# Backend (Reci-Pin)

Go で構築されたレシピ管理アプリケーションのバックエンド API サーバーです。

## 技術構成 (Tech Stack)

- **Language**: Go 1.23+
- **Router**: chi
- **Database**: PostgreSQL
- **ORM/Query Builder**: Standard `database/sql` + `sqlc` (optional)
- **Search**: Elasticsearch
- **Storage**: AWS S3 (LocalStack in development)
- **Testing**: `testing` pkg, Uber Mock, golangci-lint

## アーキテクチャ

クリーンアーキテクチャの考え方を取り入れ、保守性とテスト容易性を高めるために以下の3層構造で実装されています。

```
handler → usecase → repository → Database
```

### ディレクトリ構造

- `cmd/api/`: アプリケーションのエントリーポイント
- `internal/domain/`: エンティティ、リポジトリ・ユースケースのインターフェース
- `internal/usecase/`: ビジネスロジックの実装
- `internal/server/`: HTTP ハンドラー、ルート定義、ミドルウェア
- `internal/infrastructure/`: データベース、外部 API、外部サービスの具体的な実装
- `migrations/`: データベースマイグレーションファイル

## 開発ガイド

Docker Compose を使用せずに、直接バックエンドを開発する場合の手順です。

### 前提条件

- **Go** (v1.23+)
- **PostgreSQL**, **Elasticsearch**, **LocalStack** (S3) が起動していること
  - ※ データベースなどのインフラのみを Docker で起動する場合は `docker-compose up postgres elasticsearch localstack` を実行してください。

### 1. 準備

```bash
make bootstrap
```

### 2. サーバーの起動

```bash
# デフォルトで localhost:8080 で起動
# 接続先の設定を .env や Makefile で調整が必要な場合があります
make app.run
```

### 3. その他のコマンド

- **ユニットテスト**: `make app.test`
- **リンター**: `make app.lint`
- **シード投入**: `make seed.run` (DBが起動している必要あり)

---

[Root README に戻る](../README.md)
