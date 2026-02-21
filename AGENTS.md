# AIエージェント 作業ガイドライン

# AIエージェント 作業ガイドライン

## 基本原則 (General Principles)

### 🔧 安全第一 (Safety First)
-   **秘匿情報の保護**: APIキーやパスワードなどはログに出力しない。
-   **破壊的操作の禁止**: `rm -rf`, `git reset --hard` 等のコマンドは実行せず、ユーザーに提案するに留める。
-   **変更の最小化**: 必要以上のファイルを変更せず、常に元に戻せる状態を維持する。

### 🗣 コミュニケーション (Communication)
-   **言語**: AIの思考・回答は原則として **日本語** で行う。
-   **確認**: 指示が曖昧な場合は、作業前に質問して明確化する。

## プロジェクトドキュメント & ルール

具体的な標準については、`.agent/` 内の詳細ガイドラインを参照してください。

### 📜 コアガイドライン
-   **[アーキテクチャ](.agent/guidelines/001_architecture.md)**: ヘキサゴナル/クリーンアーキテクチャの定義。
-   **[コードスタイル & リファクタリング](.agent/guidelines/002_coding_standards.md)**: コミット基準、変数命名、リファクタリングルール。
-   **[テスト](.agent/guidelines/003_testing.md)**: カバレッジ目標、モック、実行方法。

### 🛠 ワークフロー
-   **[機能開発](.agent/workflows/001_feature_development.md)**: Git Worktree (`git gtr`) を使用した標準フロー。

#### Git Worktree Runner (`gtr`) の導入と利用
並行して作業を行う際は `git gtr` を利用してください。

##### インストール (macOS)
```bash
brew tap coderabbitai/tap
brew install git-gtr
```

##### リポジトリ初期設定
新しい環境で作業を開始する際は、以下の設定を実行してください。
```bash
git gtr config add gtr.copy.include ".env*"
git gtr config add gtr.copy.include "proxy/certs/*.pem"
git gtr config add gtr.copy.include "backend/vendor"
git gtr config add gtr.copy.include "frontend/node_modules"
```

##### 基本的な使い方

```bash
# ワークツリーを作成 (新しいブランチで)
git gtr new my-feature

# ワークツリー内でコマンドを実行 (適切なディレクトリを指定)
git gtr run my-feature npm --prefix frontend test
git gtr run my-feature make -C backend build

# マージ済みのワークツリーを一括削除 (gh CLIが必要)
git gtr clean --merged
```

---
**注記**: タスクを開始する前に必ずこれらのドキュメントを参照し、プロジェクトの標準に準拠してください。
