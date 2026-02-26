# Frontend (Reci-Pin)

Angular で構築されたレシピ管理アプリケーションのフロントエンドです。

## 技術構成 (Tech Stack)

- **Framework**: Angular 19.1.2
- **Language**: TypeScript
- **Styling**: SCSS, Angular Material, CSS Variables (Design Tokens)
- **Testing**: Vitest, Angular Testing Library
- **Documentation**: Storybook
- **Form Management**: Angular Reactive Forms

## アーキテクチャ

Angular の推奨される構成に従い、関心の分離を意識したディレクトリ構成を採用しています。

```
Components → Services → HttpClient → API
```

### ディレクトリ構造

- `src/app/core/`: シングルトンサービス、インターセプター、ガード、共通モデル
- `src/app/shared/`: 複数の機能で使用される再利用可能なコンポーネントやパイプ
- `src/app/features/`: 機能ごとのモジュール（認証、レシピ管理など）
  - 各機能ディレクトリ内に `components`, `services`, `pages` を配置
- `src/assets/`: 画像や静的ファイル
- `src/styles/`: グローバルスタイルとデザイントークン

## 開発ガイド

Docker Compose を使用せずに、直接フロントエンドを開発する場合の手順です。

### 前提条件

- **Node.js** (v20+)
- **Yarn**

### 1. 依存関係のインストール

```bash
yarn install
```

### 2. 開発サーバーの起動

```bash
yarn start
```

ブラウザで `http://localhost:4200/` を開いてください。
※ バックエンド API や Proxy と連携させて開発する場合は [Root README](../README.md) の Docker Compose 手順を推奨します。

### 3. その他のコマンド

- **コンポーネントの生成**: `ng generate component path/to/name`
- **ユニットテスト**: `yarn test`
- **Storybook**: `yarn storybook`

---

[Root README に戻る](../README.md)
