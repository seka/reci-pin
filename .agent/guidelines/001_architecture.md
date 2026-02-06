---
description: Architecture Guidelines
---
# アーキテクチャガイドライン

## Backend (Go)
**ヘキサゴナル/クリーンアーキテクチャ** のレイヤー構造に厳密に従ってください:

1.  **`domain/model`** (ドメインモデル)
    -   純粋なビジネスロジックの構造体。
    -   **外部依存なし**。

2.  **`domain/repository`** (ドメインリポジトリ)
    -   インターフェースのみ定義。
    -   **実装は含まない**。

3.  **`usecase`** (ユースケース)
    -   アプリケーションのビジネスルール。
    -   `domain` に依存する。

4.  **`infrastructure`** (インフラストラクチャ)
    -   データベース、外部APIなど。
    -   `domain/repository` インターフェースに依存（実装）する。

5.  **`server/handler`** (ハンドラー)
    -   HTTPトランスポート層。
    -   `usecase` に依存する。

## Frontend (Angular)
**Atomic Design** と **Feature-Sliced Design** のハイブリッド構成を採用しています:

### ディレクトリ構造
-   **`core/`**: アプリケーション全体で **1回だけインスタンス化** されるもの。
    -   Singleton Services (AuthService, ApiService), Guards, Interceptors, Global Styles/Assets。
    -   `AppModule` (または `main.ts`) でのみインポートされるべき。
-   **`shared/`**: 複数の機能で **再利用** されるもの。
    -   UI Components (Atoms, Molecules, Organisms), Pipes, Directives。
    -   各 Feature Module でインポートして使用する。
-   **`features/`**: 機能単位のモジュール (例: `auth`, `recipes`)。
    -   各機能内に `pages` (Smart Component) や `components` (Dumb Component) を配置。

### コンポーネント分類 (Atomic Design)
1.  **Atoms (原子)**
    -   最小単位のUIパーツ (例: Button, Input, Icon)。
    -   ロジックを持たない。
2.  **Molecules (分子)**
    -   Atomsを組み合わせた要素 (例: FormField)。
3.  **Organisms (生体)**
    -   独立して機能するUIブロック (例: Header, AuthCard, RecipeCard)。
4.  **Templates (テンプレート)**
    -   具体的なコンテンツを持たないレイアウト構造。
5.  **Pages (ページ)**
    -   ルーティングのエントリーポイント。データ取得などの副作用を管理。

### デザインシステム (DADS)
-   **Tokens**: 色、スペーシング、タイポグラフィなどの基本値は `src/app/core/assets/styles/` 以下のトークンファイルで管理します。
-   **Styles**: コンポーネント固有のスタイルは SCSS Modules ではなく、Angular の View Encapsulation を利用して各コンポーネントの `.scss` に記述します。
