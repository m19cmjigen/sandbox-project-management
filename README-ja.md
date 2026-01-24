# 全社プロジェクト進捗可視化プラットフォーム

JiraのプロジェクトとIssueを自動で取り込み、組織階層でまとめて遅延状況を可視化するシステム。

## 主な機能

### 🎯 プロジェクト管理
- Jira Cloudからのプロジェクト自動同期
- 組織階層によるプロジェクト分類
- プロジェクト別の進捗統計

### 📊 進捗可視化
- ダッシュボードでの全体サマリー
- プロジェクトヒートマップ（RED/YELLOW/GREEN）
- 組織別・プロジェクト別の詳細ビュー

### ⚡ 遅延検知
- 自動的な遅延ステータス計算
- 期日ベースの3段階評価（RED/YELLOW/GREEN）
- 遅延チケットのフィルタリング

### 🔄 Jira統合
- 定期的な自動同期
- 手動同期トリガー
- リトライ機能付きエラーハンドリング

## 技術スタック

### バックエンド
- **言語**: Go 1.21+
- **フレームワーク**: Gin
- **データベース**: PostgreSQL 15
- **アーキテクチャ**: Clean Architecture

### フロントエンド
- **言語**: TypeScript
- **フレームワーク**: React 18
- **UIライブラリ**: Material-UI (MUI)

## クイックスタート

### 1. 環境変数の設定

```bash
cp .env.example .env
# .envファイルを編集してJira APIトークンなどを設定
```

### 2. データベースのセットアップ

```bash
make db-up              # PostgreSQL起動
make db-migrate         # マイグレーション適用
```

### 3. バックエンドの起動

```bash
make backend-run        # APIサーバー起動 (http://localhost:8080)
```

### 4. Jira同期の実行

```bash
make sync-once          # 1回のみ同期実行
# または
make sync-scheduler     # 定期同期（1時間ごと）
```

### 5. フロントエンドの起動

```bash
cd frontend
npm install
npm run dev            # 開発サーバー起動 (http://localhost:5173)
```

## 詳細ドキュメント

- [セットアップガイド](docs/setup.md)
- [API仕様](docs/api.md)
- [アーキテクチャ](docs/architecture.md)

## トラブルシューティング

### Jira同期エラー

環境変数を確認：
```bash
echo $JIRA_BASE_URL
echo $JIRA_EMAIL
echo $JIRA_API_TOKEN
```

APIトークンは [Atlassian Account Settings](https://id.atlassian.com/manage-profile/security/api-tokens) で作成できます。

## ライセンス

MIT License
