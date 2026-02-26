# レシピまとめサイト (Reci-Pin)

様々なWebサイトに点在するレシピを一箇所にまとめ、タグやメモで整理・検索できるようにするための個人用Webアプリケーションです。

## 主要機能

1. **ユーザー管理**
   - アカウント作成
   - ログイン/ログアウト

2. **レシピ管理**
   - レシピの保存（レシピ名、参照先URL、メモ、作成日）
   - レシピの編集・削除
   - レシピ一覧表示

3. **画像管理**
   - 自分で作った料理の画像をアップロード
   - レシピに画像を紐づけ

4. **タグ管理**
   - 材料や気分に応じたタグの作成
   - レシピに複数のタグを紐づけ

5. **検索機能**
   - タグによる絞り込み
   - メモ内容による検索

## 技術構成 (Tech Stack)

このプロジェクトは以下の技術スタックで構成されています。詳細は各ディレクトリの README を参照してください。

- **[Frontend (Angular)](./frontend/README.md)**: Angular (19系), TypeScript, Angular Material, Vitest
- **[Backend (Go)](./backend/README.md)**: Go, chi, PostgreSQL, Docker

---

## セットアップ

このプロジェクトは Docker Compose を使用して素早く開発環境を構築できます。各サービスを個別に実行したい場合は、それぞれのディレクトリの README を参照してください。

- **[Frontend 開発ガイド](./frontend/README.md)**
- **[Backend 開発ガイド](./backend/README.md)**

### クイックスタート (Docker Compose)

#### 1. 前提条件の確認

以下のツールがインストールされていることを確認してください。

- **Docker Desktop** / **Docker Compose**
- **Git**
- **mkcert**: ローカルで HTTPS (Proxy) を使用するために必要です。
  - macOS: `brew install mkcert nss && mkcert -install`

#### 2. リポジトリの準備

```bash
git clone https://github.com/seka/reci-pin.git
cd reci-pin

# 設定（ポート番号やJWTシークレット等）をカスタマイズしたい場合のみ実行
# docker-compose は自動的に .env ファイルを読み込み、デフォルト値を上書きします
# cp .env.example .env
```

#### 3. SSL 証明書の生成

Proxy (Nginx) で使用する証明書を生成します。

```bash
mkdir -p proxy/certs
mkcert -key-file proxy/certs/key.pem -cert-file proxy/certs/cert.pem localhost
```

#### 4. アプリケーションの起動

```bash
docker-compose up -d
```

起動後、以下のURLでアクセス可能です。

- **Frontend**: [https://localhost](https://localhost)
- **Backend API**: [https://localhost/api](https://localhost/api)
- **Mailhog (メール確認)**: [http://localhost:8025](http://localhost:8025)

#### 5. 初期データの投入 (Seed)

```bash
# データベースのシード、画像のアップロード、検索インデックスの同期を一括で実行します
docker-compose exec backend make seed.run
```

## Directory Structure

```
.
├── backend/          # Go API サーバー ([README](./backend/README.md))
├── frontend/         # Angular アプリケーション ([README](./frontend/README.md))
├── proxy/            # プロキシサーバー設定 (Nginx)
├── docker-compose.yml
└── AGENTS.md         # AI エージェント用ガイドライン
```

## Development Rules

開発の進め方やAIエージェントの利用については [AGENTS.md](./AGENTS.md) を参照してください。

## License

MIT
