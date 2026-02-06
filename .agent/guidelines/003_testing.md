---
description: Testing Guidelines
---
# テストガイドライン

## Backend (Go)

### カバレッジ (Coverage)
-   ユースケース層とハンドラー層で **カバー率 80%以上** を目指してください。

### モック (Mocking)
-   ユニットテストには `go.uber.org/mock/gomock` を使用してください。

### テーブル駆動テスト (Table Driven Tests)
-   可読性と保守性のために **テーブル駆動テスト (Table Driven Tests)** を使用してください。

### 実行 (Execution)
-   論理的な変更セットごとに必ずテスト (`go test ./...`) を実行してください。

## Frontend (Angular)

### テスト戦略
-   **ユニットテスト (Jasmine/Karma)**: コンポーネントのロジック、サービスの振る舞い、Pipeの変換結果を検証します。
-   **コンポーネントテスト**: `TestBed` を使用しますが、深い子コンポーネントのレンダリングは避け (浅いレンダリング)、`NO_ERRORS_SCHEMA` やスタブを活用して独立性を保ってください。

### 記述ルール
-   `describe` ブロックでテスト対象を明確にし、`it` ブロックで「何がどうなるべきか」を自然言語で記述してください。
-   非同期処理のテストには `fakeAsync` / `tick` または `waitForAsync` を適切に使用してください。
