# レシピまとめサイト (Reci-Pin)

レシピの参照先URLを保存・管理するための個人用Webアプリケーションです。

## 概要

様々なWebサイトに点在するレシピを一箇所にまとめ、タグやメモで整理・検索できるようにします。

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

- **Frontend**: Angular (最新版) + TypeScript + Angular Material
- **Backend**: Go + chi ルーター
- **Database**: PostgreSQL
- **Web Server**: Nginx (Reverse Proxy)
- **Infrastructure**: Docker Compose

## アーキテクチャ

本プロジェクトは、保守性と拡張性を高めるために、関心の分離を意識したディレクトリ構成を採用しています。

### Backend (Go)

クリーンアーキテクチャの考え方を取り入れ、以下の3層構造で実装されています。

```
handler → usecase → repository → Database
```

- **handler**: HTTPリクエストのパース、バリデーション、レスポンスの返却
- **usecase**: 主要なビジネスロジック
- **repository**: データベースなど、外部ストレージへのデータアクセス

### Frontend (Angular)

Angular の推奨される構成に従い、機能ごとにモジュール・コンポーネントを分離しています。

```
Components → Services → HttpClient → API
```

- **Components**: UIコンポーネントとテンプレート
- **Services**: ビジネスロジックとAPI通信
- **Guards**: ルーティング保護
- **Interceptors**: HTTP通信の共通処理

## 開発環境のセットアップ

### Requirements

- Docker Desktop
- Git

### Setup

1. リポジトリをクローン

```bash
git clone https://github.com/seka/reci-pin.git
cd reci-pin
```

2. 環境変数の設定

```bash
cp .env.example .env
# 必要に応じて .env を編集
```

3. Docker コンテナを起動

```bash
docker-compose up -d
```

4. データベースマイグレーションを実行

```bash
docker-compose exec backend ./cmd/migrate
```

### Access

- **Frontend**: http://localhost
- **Backend API**: http://localhost/api
- **PostgreSQL**: localhost:5432

## Directory Structure

```
.
├── backend/          # Go API サーバー
│   ├── cmd/         # エントリーポイント
│   ├── internal/    # 内部パッケージ
│   │   ├── domain/      # ドメイン層（エンティティ、リポジトリIF）
│   │   ├── infrastructure/ # インフラ層（DB実装）
│   │   ├── usecase/     # ユースケース層
│   │   └── server/      # サーバー層（ハンドラー、ミドルウェア）
│   └── migrations/  # DBマイグレーション
├── frontend/        # Angular アプリケーション
│   └── src/
│       └── app/
│           ├── core/     # サービス、インターセプター、ガード
│           ├── shared/   # 共通コンポーネント
│           ├── auth/     # 認証機能
│           └── recipes/  # レシピ機能
├── nginx/           # Nginx 設定
├── docker-compose.yml
└── AGENTS.md        # AI エージェント用ガイドライン
```

## Development Rules

開発の進め方やAIエージェントの利用については [AGENTS.md](./AGENTS.md) を参照してください。

## License

MIT
