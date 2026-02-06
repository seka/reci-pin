---
description: Feature Development Workflow (Worktree)
---
# 機能開発ワークフロー (Worktree)

このワークフローは、**Git Worktree** を使用して新機能の実装や修正を行う手順を定義します。

1.  **ワークツリーの作成 (Create a Worktree)**
    -   `master` をベースにタスク用の新しいワークツリーを作成します。
    -   これにより、メインの作業ディレクトリをクリーンに保ちつつ、並行作業が可能になります。
    ```bash
    # リポジトリルートにいることを確認
    git worktree add -b feature/your-feature-name ../your-feature-name master
    cd ../your-feature-name
    ```

2.  **実装と検証 (Implement & Verify)**
    -   ワークツリーディレクトリ内で必要なコード変更を行います。
    -   ビルドやテストを実行して機能を検証します。

3.  **変更のコミット (Commit Changes)**
    -   Conventional Commits (例: `feat:`, `fix:`, `refactor:`) に従って、明確かつ粒度の細かいメッセージでコミットします。
    ```bash
    git add .
    git commit -m "feat: 変更内容の説明"
    ```

4.  **マスターへのマージ (Merge into Master)**
    -   検証後、機能ブランチを `master` にマージします。
    -   これはメインのリポジトリディレクトリから行います。
    ```bash
    # メインリポジトリに戻る
    cd /path/to/repo/root 
    git merge feature/your-feature-name
    ```

5.  **クリーンアップ (Cleanup)**
    -   ワークツリーを削除し、機能ブランチを削除します。
    -   （マージ済みのブランチは不要なため）
    ```bash
    git worktree remove ../your-feature-name
    git branch -d feature/your-feature-name
    ```
