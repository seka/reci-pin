---
description: Code Style & Refactoring Guidelines
---
# コードスタイル & リファクタリングガイドライン

## Backend (Go)

### コミット粒度 (Commit Granularity)
コミットは粒度が細かく、論理的であるべきです。DB変更、ロジック、テストをすべて含むような巨大なコミットは避けてください。以下のように分割してください:
-   **Phase 1**: DBスキーマ & モデル
-   **Phase 2**: リポジトリ / インフラ
-   **Phase 3**: ドメイン / ユースケースロジック
-   **Phase 4**: インターフェース / ハンドラー / ルーティング
-   **Phase 5**: テスト

### 変数のクリーンさ (Variable Cleanliness)
-   未使用の変数やインポートは直ちに削除してください。
-   リンター (Lint) ツールでエラーがゼロになるようにしてください。

### リファクタリング (Refactoring)
-   リファクタリングを行う際は、必ず対応するテストも更新してください。
-   完了する前に必ず **テスト通過率 100%** を確認してください。

### 命名規則 (Naming Convention)
-   **ファイル名**: `snake_case` (例: `user_repository.go`)
-   **関数/構造体**: `PascalCase` (Public), `camelCase` (Private)
-   **インターフェース**: `er` サフィックスを推奨 (例: `Reader`, `Writer`)

## Frontend (Angular)

### TypeScript & Strict Mode
-   `strict: true` を前提とし、`any` 型の使用は原則禁止とします。
-   Null安全性 (`strictNullChecks`) を確保し、Optional Chaining (`?.`) や Nullish Coalescing (`??`) を活用してください。

### RxJS
-   **メモリリーク防止**: `Subscription` は適切に解除してください。Angular 16+ では `takeUntilDestroyed` オペレータの使用を推奨します。
-   **Async Pipe**: コンポーネントでの手動 `subscribe` は避け、テンプレート側で `async` パイプを使用してください。

### スタイル (SCSS)
-   グローバルなスタイル汚染を避けるため、可能な限りコンポーネント内 (`:host`) にスタイルを閉じ込めてください。
-   DADSデザイントークン (`src/app/core/assets/styles/tokens.scss`) の変数を積極的に使用し、マジックナンバーの使用は避けてください。

### 命名規則 (Naming Convention)
-   **ファイル名**: `kebab-case` (例: `user-profile.component.ts`)
-   **クラス名**: `PascalCase` (例: `UserProfileComponent`)
-   **変数/メソッド**: `camelCase` (例: `fetchData()`)
-   **セレクタ**: `kebab-case` で `app-` プレフィックス (例: `app-user-profile`)
