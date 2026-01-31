# AI エージェント開発ガイドライン

このドキュメントは、AI エージェントが本プロジェクト (Reci-Pin) で開発を行う際のガイドラインです。

## プロジェクト概要

レシピまとめサイト（Reci-Pin）は、Web上に点在するレシピのURLを保存・管理するための個人用アプリケーションです。

### 技術スタック

- **Backend**: Go + chi ルーター + PostgreSQL
- **Frontend**: Angular + TypeScript + Angular Material
- **Infrastructure**: Docker Compose + Nginx

## アーキテクチャ原則

### Backend (Go)

**クリーンアーキテクチャ**を採用しています。以下の層構造を遵守してください:

1. **Domain層** (`internal/domain/`)
   - ビジネスロジックの中核
   - 外部依存を持たない純粋なビジネスルール
   - エンティティとリポジトリインターフェースを定義

2. **Infrastructure層** (`internal/infrastructure/`)
   - 外部システムとの接続実装
   - リポジトリインターフェースの具体実装
   - データベース、外部API等

3. **UseCase層** (`internal/usecase/`)
   - アプリケーション固有のビジネスロジック
   - ドメインエンティティとリポジトリを組み合わせて処理

4. **Server層** (`internal/server/`)
   - HTTP通信の入出力処理
   - リクエストのバリデーション
   - レスポンスの整形

**依存関係のルール**:
```
Server → UseCase → Domain ← Infrastructure
```

- 上位層は下位層に依存可能
- 下位層(Domain)は上位層に依存してはいけない
- InfrastructureはDomainのインターフェースに依存

### Frontend (Angular)

**Angular推奨構成**に従います:

1. **Core Module** (`src/app/core/`)
   - アプリケーション全体で1度だけ読み込まれる
   - Services, Guards, Interceptors を配置
   - Singleton パターンで実装

2. **Shared Module** (`src/app/shared/`)
   - 複数の機能モジュールで共有されるコンポーネント
   - 純粋なプレゼンテーション層
   - ビジネスロジックを含まない

3. **Feature Modules** (`src/app/auth/`, `src/app/recipes/`)
   - 機能ごとにモジュール分割
   - 遅延ロード (Lazy Loading) を活用
   - 関心の分離を徹底

## 開発ルール

### コーディング規約

#### Backend (Go)

- **命名規則**:
  - パッケージ名: 小文字、単語区切りなし (`userservice` ではなく `user`)
  - 関数/メソッド: キャメルケース (`GetUserByID`)
  - 変数: キャメルケース (`userID`)
  
- **エラーハンドリング**:
  - カスタムエラーは `internal/domain/errors` で定義
  - エラーは常にラップして返す (`fmt.Errorf("failed to...: %w", err)`)
  
- **テスト**:
  - `_test.go` で単体テストを実装
  - テーブル駆動テストを推奨
  - モックは必要に応じて生成

#### Frontend (Angular)

- **命名規則**:
  - コンポーネント: ケバブケース (`recipe-list.component.ts`)
  - サービス: ケバブケース + `.service` (`recipe.service.ts`)
  - クラス: パスカルケース (`RecipeListComponent`)
  
- **ディレクトリ構造**:
  - 機能ごとにディレクトリを作成
  - コンポーネント、サービス、モデルを同一ディレクトリに配置
  
- **RxJS**:
  - Observable は `$` サフィックス (`users$`)
  - unsubscribe を適切に管理（AsyncPipe推奨）

### Git ワークフロー

1. **ブランチ戦略**:
   - `master`: 本番環境
   - `feature/*`: 新機能開発
   - `fix/*`: バグ修正

2. **コミットメッセージ**:
   - Conventional Commits 形式を使用
   - プレフィックス: `feat`, `fix`, `docs`, `refactor`, `test`, `chore`
   - 例: `feat(backend): add recipe search API`

3. **コミット粒度**:
   - 1コミット = 1つの論理的な変更
   - レイヤーごと、機能ごとに分けてコミット

### Docker 開発環境

- **ホットリロード**:
  - Backend: Air を使用
  - Frontend: Angular CLI の watch モード
  
- **データベースマイグレーション**:
  - SQLファイルでバージョン管理
  - `migrations/` ディレクトリに配置
  - ファイル名: `001_init.sql`, `002_add_tags.sql` (連番形式)

## 実装時の注意点

### セキュリティ

- **認証**: JWT を使用
- **パスワード**: bcrypt でハッシュ化
- **バリデーション**: フロントエンド・バックエンド両方で実施
- **CORS**: 必要最小限の設定

### パフォーマンス

- **N+1問題**: JOIN を適切に使用
- **インデックス**: 検索対象カラムにインデックス作成
- **画像**: リサイズ・圧縮を実装（将来的な拡張）

### テスト

- **Backend**:
  - 単体テスト（各層ごと）
  - 統合テスト（API エンドポイント）
  
- **Frontend**:
  - コンポーネントテスト
  - サービステスト

## タスク管理

- 大きな機能は小さなタスクに分割
- タスクごとにブランチを作成
- 各ステップでコミット

## 参考プロジェクト

本プロジェクトは [fish-auction](https://github.com/seka/fish-auction) の構成を参考にしています。

---

**重要**: このガイドラインに従って開発を進めることで、保守性が高く、拡張しやすいコードベースを維持できます。
