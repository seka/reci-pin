# 機能開発ワークフロー (Git Worktree / gtr)

このワークフローは、**git-worktree-runner (gtr)** を使用して新機能の実装や修正を行う標準手順を定義します。

1.  **ワークツリーの作成 (Create a Worktree)**
    -   `master` をベースにタスク用の新しいワークツリーを作成します。
    -   `.env` や `node_modules` などの必要なファイルは自動的にコピーされます。
    ```bash
    git gtr new feature/your-feature-name
    ```

2.  **実装と検証 (Implement & Verify)**
    -   作成されたディレクトリ内でコード変更を行います。
    -   `git gtr run` を使用して、メインのリポジトリディレクトリからコマンドを実行することも可能です。
    ```bash
    # ワークツリー内でビルド・テストを実行 (ディレクトリ指定が必要)
    git gtr run feature/your-feature-name make -C backend build
    git gtr run feature/your-feature-name npm --prefix frontend test
    ```

3.  **変更のコミット (Commit Changes)**
    -   ワークツリーディレクトリ内でコミットを行います。
    ```bash
    git add .
    git commit -m "feat: 変更内容の説明"
    ```

4.  **マスターへのマージ (Merge into Master)**
    -   検証後、機能ブランチを `master` にマージします。
    -   メインのリポジトリディレクトリに移動して実行します。
    ```bash
    git checkout master
    git merge feature/your-feature-name
    ```

5.  **クリーンアップ (Cleanup)**
    -   マージが終わったら、不要になったワークツリーを一括削除します（gh CLIが必要です）。
    ```bash
    git gtr clean --merged
    ```
    -   または個別に削除します。
    ```bash
    git worktree remove ../your-feature-name
    git branch -d feature/your-feature-name
    ```
